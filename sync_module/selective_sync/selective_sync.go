// Package selectivesync implements a datasync manager.
package selectivesync

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.viam.com/utils"

	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/services/datamanager"

	"go.viam.com/rdk/services/vision"
)

// Model is the full model definition.
var Model = resource.NewModel("selectivesync", "demo", "vision")

func init() {
	resource.RegisterComponent(sensor.API, Model,
		resource.Registration[sensor.Sensor, *Config]{
			Constructor: newSelectiveSyncer,
		})
}

// Config contains the name to the underlying camera and the name of the vision service to be used.
type Config struct {
	Camera        string `json:"camera"`
	VisionService string `json:"vision_service"`
}

// Validate validates the config and returns implicit dependencies.
func (cfg *Config) Validate(path string) ([]string, error) {
	if cfg.Camera == "" {
		return nil, fmt.Errorf(`expected "camera" attribute in %q`, path)
	}
	if cfg.VisionService == "" {
		return nil, fmt.Errorf(`expected "vision_service" attribute in %q`, path)
	}

	return []string{cfg.Camera, cfg.VisionService}, nil
}

type visionSyncer struct {
	resource.Named
	camera        camera.Camera
	visionService vision.Service

	cancelCtx  context.Context
	cancelFunc func()
	logger     logging.Logger
}

func newSelectiveSyncer(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger,
) (sensor.Sensor, error) {
	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	v := &visionSyncer{
		Named:      conf.ResourceName().AsNamed(),
		logger:     logger,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
	}
	if err := v.Reconfigure(ctx, deps, conf); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *visionSyncer) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
	dmConfig, err := resource.NativeConfig[*Config](conf)
	if err != nil {
		return err
	}

	s.camera, err = camera.FromDependencies(deps, dmConfig.Camera)
	if err != nil {
		return errors.Wrapf(err, "unable to get camera %v for visionSyncer", dmConfig.VisionService)
	}
	s.visionService, err = vision.FromDependencies(deps, dmConfig.VisionService)
	if err != nil {
		return errors.Wrapf(err, "unable to get vision service %v for visionSyncer", dmConfig.VisionService)
	}

	return nil
}

// DoCommand simply echos whatever was sent.
func (s *visionSyncer) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return cmd, nil
}

func (s *visionSyncer) Readings(context.Context, map[string]interface{}) (map[string]interface{}, error) {
	readings := datamanager.CreateShouldSyncReading(s.ToSync())
	return readings, nil
}

func (s *visionSyncer) ToSync() bool {
	stream, err := s.camera.Stream(s.cancelCtx)
	defer utils.UncheckedErrorFunc(func() error {
		return stream.Close(s.cancelCtx)
	})
	if err != nil {
		s.logger.Error("could not get camera stream")
		return false
	}
	// Check for stuff, if true Sync
	img, release, err := stream.Next(s.cancelCtx)
	defer release()
	if err != nil {
		s.logger.Error("could not get next image")
		return false
	}
	detections, err := s.visionService.Detections(s.cancelCtx, img, map[string]interface{}{})
	if err != nil {
		s.logger.Error("could not get detections")
		return false
	}
	if len(detections) != 0 {
		s.logger.Info("time to sync")
		return true
	}
	return false
}

// Close closes the underlying generic.
func (s *visionSyncer) Close(ctx context.Context) error {
	s.cancelFunc()
	return nil
}
