// Package selectivesync implements a datasync manager.
package selectivesync

import (
	"context"
	"time"

	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/services/datamanager"
)

// Model is the full model definition.
var Model = resource.NewModel("selectivesync", "demo", "time")

func init() {
	resource.RegisterComponent(sensor.API, Model,
		resource.Registration[sensor.Sensor, *Config]{
			Constructor: newSelectiveSyncer,
		})
}

type Config struct{}

// Validate validates the config and returns implicit dependencies.
func (cfg *Config) Validate(path string) ([]string, error) {
	return []string{}, nil
}

type timeSyncer struct {
	resource.Named

	cancelCtx  context.Context
	cancelFunc func()
	logger     logging.Logger
}

func newSelectiveSyncer(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger,
) (sensor.Sensor, error) {
	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	v := &timeSyncer{
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

func (s *timeSyncer) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
	return nil
}

// DoCommand simply echos whatever was sent.
func (s *timeSyncer) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return cmd, nil
}

func (s *timeSyncer) Readings(context.Context, map[string]interface{}) (map[string]interface{}, error) {
	currentTime := time.Now()
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()) // midnight
	endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 1, 0, 0, 0, currentTime.Location())   // 1AM

	// If it is between midnight and 1AM, sync.
	if currentTime.After(startTime) && currentTime.Before(endTime) {
		return datamanager.CreateShouldSyncReading(true), nil
	}

	// Otherwise, do not sync.
	return datamanager.CreateShouldSyncReading(false), nil
}

// Close closes the underlying generic.
func (s *timeSyncer) Close(ctx context.Context) error {
	s.cancelFunc()
	return nil
}
