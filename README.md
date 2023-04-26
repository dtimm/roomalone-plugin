# RoomAlone-plugin

## Setup

To install the required packages for this plugin, run the following command:

```bash
go mod install
```

To run the plugin, enter the following command:

```bash
go run ./cmd/server
```

Once the local server is running:

1. Navigate to https://chat.openai.com. 
2. In the Model drop down, select "Plugins" (note, if you don't see it there, you don't have access yet).
3. Select "Plugin store"
4. Select "Develop your own plugin"
5. Enter in `localhost:8080` since this is the URL the server is running on locally, then select "Find manifest file".

The plugin should now be installed and enabled! You can start by asking to play a game of RoomAlone!
