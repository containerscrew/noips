package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/gin-gonic/gin"
	"github.com/nanih98/noips/internal/providers"
	"github.com/nanih98/noips/internal/providers/infra"
	"github.com/nanih98/noips/pkg/app"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"os"
)

func main() {
	// Argparser
	parser := argparse.NewParser("noips", "Prometheus exporter to check the number of ips in your subnet of your cloud provider")
	profile := parser.String("p", "profile", &argparse.Options{Required: false, Help: "Custom nameserver", Default: os.Getenv("AWS_PROFILE")})
	region := parser.String("r", "region", &argparse.Options{Required: false, Help: "AWS region", Default: "eu-west-1"})
	provider := parser.Selector("p", "provider", []string{"aws", "gcp", "azure"}, &argparse.Options{Required: false, Default: "aws", Help: "Set provider"})
	//logLevel := parser.Selector("l", "log-level", []string{"info", "debug"}, &argparse.Options{Required: false, Default: "info", Help: "Log level of the application"})
	if err := parser.Parse(os.Args); err != nil {
		log.Fatal(fmt.Println(parser.Usage(err)))
	}

	// Setup application logger
	//logger := logging.NewLogger(&logging.LoggerOptions{
	//	LogLevel: *logLevel,
	//})

	var customProvider providers.Providers
	var metrics providers.Metrics

	switch *provider {
	case "aws":
		customProvider = infra.NewAWSProvider(*profile, *region)
	case "gcp":
	case "azure":
	default:
		fmt.Println("error") // change this
	}

	// Authenticate
	customProvider.AuthenticateProvider()

	// Describe subnets
	customProvider.DescribeProviderSubnets()

	// Build our custom data for subnets
	subnets := customProvider.BuildSubnetData()

	// Start metrics
	metrics = infra.NewMetrics()
	metrics.RegisterMetrics()
	
	// Api config
	router := app.APIConfiguration()

	// Metrics handler
	router.GET("/metrics", gin.WrapH(metrics.RegisterHandler()), func(context *gin.Context) {
		// If someone ask for the metrics, the metrics are refreshed :)
		for _, subnet := range subnets {
			log.Println("Refreshing data for", subnet.ID, subnet.CIDR, subnet.AvailableIPV4)
			m.info.With(prometheus.Labels{"subnet_id": subnet.ID, "subnet_cidr": subnet.CIDR}).Set(float64(subnet.AvailableIPV4))
		}
	})

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
