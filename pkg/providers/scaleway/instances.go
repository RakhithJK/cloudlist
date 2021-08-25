package scaleway

import (
	"context"

	"github.com/projectdiscovery/cloudlist/pkg/schema"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// instanceProvider is an instance provider for scaleway API
type instanceProvider struct {
	id          string
	instanceAPI *instance.API
}

// GetInstances returns all the instances in the store for a provider.
func (d *instanceProvider) GetResource(ctx context.Context) (*schema.Resources, error) {
	list := schema.NewResources()

	for _, zone := range scw.AllZones {
		req := &instance.ListServersRequest{
			Zone: zone,
		}
		var totalResults uint32
		for {
			resp, err := d.instanceAPI.ListServers(req)
			if err != nil {
				return nil, err
			}

			for _, server := range resp.Servers {
				totalResults++

				var ip4, privateIP4 string
				if server.PublicIP != nil && server.PublicIP.Address != nil {
					ip4 = server.PublicIP.Address.String()
				}
				if server.PrivateIP != nil {
					privateIP4 = *server.PrivateIP
				}
				if privateIP4 != "" {
					list.Append(&schema.Resource{
						Provider:    providerName,
						ID:          d.id,
						PrivateIpv4: privateIP4,
					})
				}
				list.Append(&schema.Resource{
					Provider:   providerName,
					ID:         d.id,
					PublicIPv4: ip4,
					Public:     true,
				})
			}
			if resp.TotalCount == totalResults {
				break
			}
			*req.Page = *req.Page + 1
		}
	}
	return list, nil
}
