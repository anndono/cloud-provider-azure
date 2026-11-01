/*
Copyright 2024 The Kubernetes Authors.

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

package zone

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	armresources "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"sigs.k8s.io/cloud-provider-azure/pkg/azclient/providerclient/mock_providerclient"
)

func TestRepo_ListZones(t *testing.T) {
	t.Run("return fake eastus2 zones without calling ARM", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := mock_providerclient.NewMockInterface(ctrl)
		repo, err := NewRepo(client)
		assert.NoError(t, err)
		client.EXPECT().GetVirtualMachineSupportedZones(gomock.Any()).Times(0)

		zones, err := repo.ListZones(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, map[string][]string{
			"eastus2": {"1", "2", "3", "4", "5"},
		}, zones)
	})

	t.Run("refresh returns independent fake zones without calling ARM", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := mock_providerclient.NewMockInterface(ctrl)
		repo, err := NewRepo(client)
		assert.NoError(t, err)
		client.EXPECT().GetVirtualMachineSupportedZones(gomock.Any()).Times(0)

		zones, err := repo.ListZones(context.Background())
		assert.NoError(t, err)
		if !assert.Equal(t, []string{"1", "2", "3", "4", "5"}, zones["eastus2"]) {
			t.FailNow()
		}
		zones["eastus2"][4] = "4"
		delete(zones, "eastus2")

		refreshedZones, err := repo.ListZones(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, map[string][]string{
			"eastus2": {"1", "2", "3", "4", "5"},
		}, refreshedZones)
	})
}

func TestRepo_ListNetworkResourceZones(t *testing.T) {
	networkProvider := &armresources.Provider{
		ResourceTypes: []*armresources.ProviderResourceType{
			{
				ResourceType: to.Ptr("loadBalancers"),
				ZoneMappings: []*armresources.ZoneMapping{
					{Location: to.Ptr("West Europe"), Zones: []*string{to.Ptr("1"), to.Ptr("2"), to.Ptr("3"), to.Ptr("4")}},
				},
			},
			{
				ResourceType: to.Ptr("publicIPAddresses"),
				ZoneMappings: []*armresources.ZoneMapping{
					{Location: to.Ptr("East US 2"), Zones: []*string{to.Ptr("2"), to.Ptr("1"), to.Ptr("3")}},
					{Location: to.Ptr("westus"), Zones: []*string{to.Ptr("1")}},
				},
			},
		},
	}

	t.Run("return the zones of the requested resource type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := mock_providerclient.NewMockInterface(ctrl)
		repo, err := NewRepo(client)
		assert.NoError(t, err)
		client.EXPECT().GetProvider(gomock.Any(), "Microsoft.Network").Return(networkProvider, nil).Times(1)

		zones, err := repo.ListNetworkResourceZones(context.Background(), PublicIPResourceType)
		assert.NoError(t, err)
		assert.Equal(t, map[string][]string{
			"eastus2": {"2", "1", "3"},
			"westus":  {"1"},
		}, zones)
	})

	t.Run("match the resource type case-insensitively", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := mock_providerclient.NewMockInterface(ctrl)
		repo, err := NewRepo(client)
		assert.NoError(t, err)
		client.EXPECT().GetProvider(gomock.Any(), "Microsoft.Network").Return(networkProvider, nil).Times(1)

		zones, err := repo.ListNetworkResourceZones(context.Background(), "PublicIPAddresses")
		assert.NoError(t, err)
		assert.Equal(t, []string{"2", "1", "3"}, zones["eastus2"])
	})

	t.Run("return an empty map for an unlisted resource type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := mock_providerclient.NewMockInterface(ctrl)
		repo, err := NewRepo(client)
		assert.NoError(t, err)
		client.EXPECT().GetProvider(gomock.Any(), "Microsoft.Network").Return(networkProvider, nil).Times(1)

		zones, err := repo.ListNetworkResourceZones(context.Background(), "applicationGateways")
		assert.NoError(t, err)
		assert.Equal(t, map[string][]string{}, zones)
	})

	t.Run("return error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := mock_providerclient.NewMockInterface(ctrl)
		repo, err := NewRepo(client)
		assert.NoError(t, err)
		expectedErr := fmt.Errorf("error")
		client.EXPECT().GetProvider(gomock.Any(), "Microsoft.Network").Return(nil, expectedErr).Times(1)

		_, err = repo.ListNetworkResourceZones(context.Background(), PublicIPResourceType)
		assert.ErrorIs(t, err, expectedErr)
	})
}
