# trigger-sync-examples
Module to manage syncing data based on a trigger from a component

## What does this accomplish
This is a module that takes in the datamanager, a color detector, and camera. If the color is detected by the vision service, syncing is triggered and the captured files are synced to the cloud. Otherwise, the data is captured but remains on the local directory until the next time the color is detected. The color detector may need to be modified for your environment. If sync is enabled and using trigger sync module, data may be synced even before encountering a trigger. As such, ensure that you disable sync on the data manager service in your config to constrain syncing to be trigger-based.

## How to run this
- Clone this repo
- `cd sync_module && go build && sudo chmod a+rx sync_module` to build the executable and have the correct permissions to run in RDK
- Use the config but change the module executable path to your local directory and camera to your own camera

Ensure that RDK is up-to-date, as we had to make changes to RDK (https://github.com/viamrobotics/rdk/pull/3135) to ensure this runs. This may require running the config with the local RDK, as in `go run web/cmd/server/main.go -debug -config ~/Downloads/viam-sync-bot-main.json` or updating viam-server to use head `brew install --HEAD viam-server` or `curl https://storage.googleapis.com/packages.viam.com/apps/viam-server/viam-server-latest-aarch64.AppImage -o viam-server && chmod 755 viam-server && sudo ./viam-server --aix-install`
