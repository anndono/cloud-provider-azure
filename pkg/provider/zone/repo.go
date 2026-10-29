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

// Generate mocks for the repository interface
//go:generate mockgen -destination=./mock_repo.go -package=zone -copyright_file ../../../hack/boilerplate/boilerplate.generatego.txt -source=repo.go Repository

import (
	"context"

	"sigs.k8s.io/cloud-provider-azure/pkg/azclient/providerclient"
	"sigs.k8s.io/cloud-provider-azure/pkg/log"
)

type Repository interface {
	ListZones(ctx context.Context) (map[string][]string, error)
}

type repo struct {
	client providerclient.Interface
}

func NewRepo(client providerclient.Interface) (Repository, error) {
	return &repo{client: client}, nil
}

func (r *repo) ListZones(ctx context.Context) (map[string][]string, error) {
	// TEST ONLY: bypass Compute discovery to reproduce a public IP rejecting VM zone 5.
	zones := map[string][]string{
		"eastus2": {"1", "2", "3", "4", "5"},
	}
	log.FromContextOrBackground(ctx).WithName("ListZones").V(0).Info(
		"TEST ONLY: Using fake VM zones without ARM discovery; do not deploy to production",
		"region", "eastus2", "zones", zones["eastus2"])

	return zones, nil
}
