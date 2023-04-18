package app

type Providers interface {
	AuthenticateProvider()
	DescribeProviderSubnets()
	BuildSubnetData() []SubnetsData
}
