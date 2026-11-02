/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package providerclient

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	armresources "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

func (client *Client) ListProviders(ctx context.Context) (result []*armresources.Provider, rerr error) {
	pager := client.NewListPager(nil)
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		result = append(result, nextResult.Value...)
	}
	return result, nil
}
func (client *Client) GetProvider(ctx context.Context, resourceProviderNamespace string) (*armresources.Provider, error) {
	resp, err := client.Get(ctx, resourceProviderNamespace, nil)
	if err != nil {
		return nil, err
	}
	return &resp.Provider, nil
}

func (client *Client) GetVirtualMachineSupportedZones(ctx context.Context) (map[string][]*string, error) {
	// TEST ONLY: keep the invalid VM zone available without faking Network capabilities.
	return map[string][]*string{
		"eastus2": to.SliceOfPtrs("1", "2", "3", "4", "5"),
	}, nil
}
