package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name:                 "mapr",
		Usage:                "Apply a command to values of structured data, with optional key filtering.",
		UsageText:            "[global options] command",
		Version:              "0.0.1",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "input-format",
				Value: "json",
				Usage: "Input data format",
			},
			cli.StringFlag{
				Name:  "input-file",
				Value: "",
				Usage: "Input file (or pipe to stdin)",
			},
			cli.StringFlag{
				Name:  "key-filter-type",
				Value: "",
				Usage: "[prefix|suffix]",
			},
			cli.StringFlag{
				Name:  "key-filter",
				Value: "",
				Usage: "String to filter on",
			},
			cli.BoolFlag{
				Name:  "key-filter-strip",
				Usage: "Remove the filter from the resulting key name",
			},
			cli.StringFlag{
				Name:  "command-reference",
				Value: "{{value}}",
				Usage: "The reference to interpolate",
			},
			cli.BoolFlag{
				Name:  "command-no-trim",
				Usage: "Do not trim whitespace and newlines from command output",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				log.Fatal("Must pass a command to apply to values")
			}

			var bytes []byte
			var fileErr error

			if c.String("input-file") == "" {
				bytes, fileErr = ioutil.ReadAll(os.Stdin)
			} else {
				bytes, fileErr = ioutil.ReadFile(c.String("input-file"))

			}

			if fileErr != nil {
				log.Fatal(fileErr)
			}

			switch strings.ToLower(c.String("input-format")) {
			case "json":
				certainlyJSON := parseJSON(bytes)

				outputJSON(mapValues(c, certainlyJSON))
			}

			return nil
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func filterMatch(key string, keyFilterType string, keyFilter string) bool {
	if keyFilter == "" || keyFilterType == "" {
		return false
	}

	switch keyFilterType {
	case "prefix":
		return strings.HasPrefix(key, keyFilter)
	case "suffix":
		return strings.HasSuffix(key, keyFilter)
	}

	return false
}

func mapValues(c *cli.Context, data map[string]interface{}) map[string]interface{} {
	keyFilter := c.String("key-filter")
	keyFilterType := c.String("key-filter-type")
	keyFilterStrip := c.Bool("key-filter-strip")

	for k, v := range data {
		if value, ok := v.(string); ok {
			if keyFilterType != "" && keyFilter != "" {
				if filterMatch(k, keyFilterType, keyFilter) {
					value = applyCommand(c, value)

					if keyFilterStrip {
						delete(data, k)

						switch keyFilterType {
						case "prefix":
							k = strings.TrimPrefix(k, keyFilter)
						case "suffix":
							k = strings.TrimSuffix(k, keyFilter)
						}
					}
				}
			} else {
				value = applyCommand(c, value)
			}

			if !c.Bool("command-no-trim") {
				data[k] = strings.TrimSpace(value)
			} else {
				data[k] = value
			}
		} else if value, ok := v.(map[string]interface{}); ok {
			mapValues(c, value)
		}
	}

	return data
}

func parseJSON(stdIn []byte) map[string]interface{} {
	var maybeJSON interface{}
	err := json.Unmarshal(stdIn, &maybeJSON)

	if err != nil {
		println(err)
		os.Exit(1)
	}

	return maybeJSON.(map[string]interface{})
}

func outputJSON(certainlyJSON map[string]interface{}) {
	jsonOutput, err := json.Marshal(certainlyJSON)

	if err != nil {
		println(err)
		os.Exit(1)
	} else {
		os.Stdout.Write(jsonOutput)
		os.Exit(0)
	}
}

func applyCommand(c *cli.Context, dataValue string) string {
	shellCommand := strings.ReplaceAll(c.Args().Get(0), c.String("command-reference"), dataValue)
	out, err := exec.Command("sh", "-c", shellCommand).Output()

	if err != nil {
		log.Fatal(err)
	}

	return string(out)
}
