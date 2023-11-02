# trigger-sync-examples

Module to manage syncing data based on a trigger from a component

## What does this accomplish

1. time-interval-trigger is a module that triggers sync to the cloud only between midnight and 1AM. Otherwise, the data is captured but remains in the local directory until the next day at midnight.

2. color-trigger is a module that takes a color detector and camera and wraps it into a sensor. If the color is detected by the vision service, syncing is triggered and the captured files are synced to the cloud. Otherwise, the data is captured but remains in the local directory until the next time the color is detected. The color detector may need to be modified for your environment.

The selective sync sensor checks for the sync condition at the interval configured by the data manager (e.g. `sync_interval_mins`) and will only sync if sync is enabled.

## How to run this

- Clone this repo
- `cd time-interval-trigger && go build && sudo chmod a+rx time-interval-trigger` or `cd color-trigger && go build && sudo chmod a+rx color-trigger` to build the executable and have the correct permissions to run in RDK
- Use the config but change the module executable path to your local directory

Ensure that RDK is up-to-date, as we had to make changes to RDK (https://github.com/viamrobotics/rdk/pull/3174) to ensure this runs. This may require running the config with the local RDK, as in go run web/cmd/server/main.go -debug -config ~/Downloads/viam-sync-bot-main.json or updating viam-server to use head brew install --HEAD viam-server or curl https://storage.googleapis.com/packages.viam.com/apps/viam-server/viam-server-latest-aarch64.AppImage -o viam-server && chmod 755 viam-server && sudo ./viam-server --aix-install
