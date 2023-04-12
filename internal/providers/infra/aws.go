package infra

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/nanih98/noips/internal/providers/app"
)

type AWSData struct {
	Profile        string
	Region         string
	Config         aws.Config
	SubnetResponse []types.Subnet
}

func NewAWSProvider(profile, region string) *AWSData {
	return &AWSData{
		Profile:        profile,
		Region:         region,
		Config:         aws.Config{},
		SubnetResponse: []types.Subnet{},
	}
}

func (a *AWSData) AuthenticateProvider() {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithSharedConfigProfile(a.Profile),
		config.WithRegion(a.Region),
	)

	if err != nil {
		errors.New(err.Error())
	}

	a.Config = cfg
}

func (a *AWSData) DescribeProviderSubnets() {
	client := ec2.NewFromConfig(a.Config)

	input := &ec2.DescribeSubnetsInput{
		MaxResults: aws.Int32(100),
	}

	paginator := ec2.NewDescribeSubnetsPaginator(client, input, func(o *ec2.DescribeSubnetsPaginatorOptions) {
		o.Limit = 1000
	})

	var subnetsResponseItems []types.Subnet

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())

		if err != nil {
			errors.New(err.Error())
		}
		subnetsResponseItems = append(subnetsResponseItems, output.Subnets...)
	}

	a.SubnetResponse = subnetsResponseItems
}

func (a *AWSData) BuildSubnetData() []app.SubnetsData {
	var subnets []app.SubnetsData
	for _, item := range a.SubnetResponse {
		subnets = append(subnets, app.SubnetsData{
			ID:            *item.SubnetId,
			CIDR:          *item.CidrBlock,
			AvailableIPV4: *item.AvailableIpAddressCount,
		})
	}
	return subnets
}
