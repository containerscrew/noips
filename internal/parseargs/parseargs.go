package parseargs

import (
	"github.com/akamensky/argparse"
	"os"
)

type Arguments struct {
	Profile  *string
	Region   *string
	Provider *string
	LogLevel *string
	Usage    string
}

func ParseArgs() (Arguments, error) {
	var arguments Arguments
	parser := argparse.NewParser("noips", "Prometheus exporter to check the number of ips in your subnet of your cloud provider")
	arguments.Profile = parser.String("e", "e", &argparse.Options{Required: false, Help: "AWS PROFILE", Default: os.Getenv("AWS_PROFILE")})
	arguments.Region = parser.String("r", "region", &argparse.Options{Required: false, Help: "AWS region", Default: "eu-west-1"})
	arguments.Provider = parser.Selector("p", "provider", []string{"aws", "gcp", "azure"}, &argparse.Options{Required: false, Default: "aws", Help: "Set provider"})
	arguments.LogLevel = parser.Selector("l", "log-level", []string{"info", "debug"}, &argparse.Options{Required: false, Default: "info", Help: "Log level of the application"})

	if err := parser.Parse(os.Args); err != nil {
		return Arguments{
			Usage: parser.Usage(err),
		}, err
	}

	return arguments, nil
}
