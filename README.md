# Feels Good Station

A little Raspberry Pi-powered hardware station to collect data that may influence a good night sleep.

The station constantly measures:

- Humidity
- Temperature
- Ambient light
- Noise level
- Movement 

The data is published hourly as CSV files to a linked Dropbox folder.

## Running the software

- Create a `.env` file on your code folder. Configure it with your dropbox credentials and a temp folder:

```
DROPBOX_CLIENT_ID=<your client id>
DROPBOX_ACCESS_TOKEN=<your auth token>
DROPBOX_FOLDER=<a folder on dropbox to store your csvs>
TMP_FOLDER=<local folder where your csvs will be stored until they're uploaded>
```

- Build the binary for the ARM platform running `./build.sh`

- `./feelsgood` will start running it!

## Schematics


