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
	"testing"

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
