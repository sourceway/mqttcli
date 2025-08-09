# MQTT CLI

This is a small command-line tool designed to publish a single message to MQTT v5 or subscribe to a topic to receive a
single message.

## Usage

Ensure all required environment variables are set and then run the tool with one of the two supported commands.

### Publish

You can publish a single message to a topic by using the `publish` command.

The payload can be provided via the `-payload` flag.

```bash
mqttcli publish -payload=<message>
```

Or via stdin.

```bash
cat payload | mqttcli publish -stdin
```

### Subscribe (untested)

You can subscribe to a topic to receive a single message by using the `subscribe` command.
The received message will be printed to stdout. Binary payloads should be supported.

```bash
mqttcli subscribe
```

## Configuration

Configuration is primarily done via environment variables.

### `MQTT_SERVER` (and `MQTT_USERNAME` and `MQTT_PASSWORD`)

MQTT Server to connect to, e.g., `mqtt://mqtt.your.server:1883`; supports `ws://` and `wss://` (and whatever
else [paho.golang](https://github.com/eclipse-paho/paho.golang/) supports).

Credentials can be provided via `MQTT_USERNAME` and `MQTT_PASSWORD` or by encoding them into the URL.
For parts of the credentials encoded in the URL, their corresponding environment variable will be ignored.

#### Examples

Username and password in environment variables:

```bash
MQTT_SERVER=mqtt://mqtt.your.server:1883
MQTT_USERNAME=username
MQTT_PASSWORD=password
```

Username encoded in the URL, password in environment variables:

```bash
MQTT_SERVER=mqtt://username@mqtt.your.server:1883
# MQTT_USERNAME= # would be ignored anyway
MQTT_PASSWORD=password
```

Username and password encoded in the URL:

```bash
MQTT_SERVER=mqtt://username:password@mqtt.your.server:1883
# MQTT_USERNAME= # would be ignored anyway
# MQTT_PASSWORD= # would be ignored anyway
```

### `MQTT_TOPIC`

Topic to publish to or subscribe to.

### `MQTT_CLIENT_ID`

Client ID to use. If not provided, a random one will be generated.

### `MQTT_QOS`

QoS to use for publishing, defaults to `1`.

### Disclaimer

I'm not a Go developer, so this is probably not the best way to do things. Feel free to contribute and improve this
tool.
