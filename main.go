package main

import (
	"flag"
	"io"
	"mqttcli/logger"
	"mqttcli/mqtt"
	"os"
)

func main() {
	con, stop := mqtt.NewConnection()
	defer stop()

	publishCmd := flag.NewFlagSet("publish", flag.ExitOnError)
	publishPayload := publishCmd.String("payload", "", "payload")
	stdinPayload := publishCmd.Bool("stdin", false, "read payload from stdin (ignores -payload)")

	if len(os.Args) < 2 {
		logger.Fail("expected 'publish' or 'subscribe' subcommands")
	}

	switch os.Args[1] {
	case "publish":
		err := publishCmd.Parse(os.Args[2:])
		if err != nil {
			logger.Fail("Error parsing flags:", err)
		}

		if *stdinPayload {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				logger.Fail("Error reading payload from stdin:", err)
			}
			con.Publish(data)

		} else {
			if !wasFlagPassed(publishCmd, "payload") {
				logger.Fail("expected -payload or -stdin")
			}
			con.Publish([]byte(*publishPayload))
		}

	case "subscribe":
		con.Subscribe()

	default:
		logger.Fail("expected 'publish' or 'subscribe' subcommands")
	}
}

func wasFlagPassed(fs *flag.FlagSet, name string) bool {
	found := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
