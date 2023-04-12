package providers

import "github.com/nanih98/noips/internal/providers/app"

type Providers interface {
	AuthenticateProvider()
	DescribeProviderSubnets()
	BuildSubnetData() []app.SubnetsData
	// En mi lógica de app, monto mi custom struct con la data que necesito de lo que me devuelve el provider en el DescribeProviderSubnets()
}
