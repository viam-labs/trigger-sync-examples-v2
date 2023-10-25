# trigger-sync-examples
Module to manage syncing data based on a trigger from a component

## What does this accomplish
This is a module that takes in the datamanager, a color detector, and camera. If the color is detected by the vision service, syncing is triggered and the captured files are synced to the cloud. Otherwise, the data is captured but remains on the local directory until the next time the color is detected.

## How to run this
- Clone this repo
- `cd sync_module && go build && sudo chmod a+rx sync_module` to build the executable and have the correct permissions to run in RDK
- Use the config but change the module executable path to your local directory and camera to your own camera

Ensure that RDK is up-to-date, as we had to make changes to RDK (https://github.com/viamrobotics/rdk/pull/3135) to ensure this runs. This may require replacing the rdk import here to a local version of RDK that contains the necessary change and running the config with the local RDK, as in:
- `go mod edit -replace go.viam.com/rdk=../../viam/rdk && go mod tidy`
- `go run web/cmd/server/main.go -debug -config ~/Downloads/viam-sync-bot-main.json`
