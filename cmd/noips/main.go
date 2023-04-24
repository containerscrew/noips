package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nanih98/noips/internal/logging"
	"github.com/nanih98/noips/internal/parseargs"
	"github.com/nanih98/noips/pkg/app"
	"github.com/nanih98/noips/pkg/infra"
	"go.uber.org/zap"
	"log"
)

var (
	customProvider app.Providers
	metrics        app.Metrics
)

func main() {
	// Argument parser
	parser, err := parseargs.ParseArgs()

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
