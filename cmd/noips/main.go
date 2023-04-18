package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/gin-gonic/gin"
	"github.com/nanih98/noips/internal/logging"
	"github.com/nanih98/noips/pkg/app"
	"github.com/nanih98/noips/pkg/infra"
	"go.uber.org/zap"
	"log"
	"os"
)

var (
	customProvider app.Providers
	metrics        app.Metrics
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

func main() {
	parser, err := ParseArgs()

	if err != nil {
		log.Fatal(fmt.Println(parser.Usage))
	}

	// Setup application logger
	logger := logging.NewLogger(&logging.LoggerOptions{
		LogLevel: *parser.LogLevel,
	})

	logger.Log.Info("Starting prometheus noips exporter",
		zap.String("provider", *parser.Provider),
		zap.String("author", "github.com/nanih98"),
	)

	switch *parser.Provider {
	case "aws":
		customProvider = infra.NewAWSProvider(*parser.Profile, *parser.Region)
	case "gcp":
	case "azure":
	}

	// Authenticate
	customProvider.AuthenticateProvider()

	// Describe subnets
	customProvider.DescribeProviderSubnets()

	// Start metrics
	metrics = infra.NewMetrics()
	metrics.RegisterMetrics()

	// Api config
	router := app.APIConfiguration()

	// Metrics handler
	router.GET("/metrics", gin.WrapH(metrics.RegisterHandler()), func(context *gin.Context) {
		// If someone ask for the metrics, the metrics are refreshed, if not, you don't need to perform api calls to the provider :)
		subnets := customProvider.BuildSubnetData()
		for _, subnet := range subnets {
			metrics.UpdateMetrics(subnet)
		}
	})

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
