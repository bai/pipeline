// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logging

import (
	"context"

	"github.com/banzaicloud/pipeline/internal/clusterfeature"
	"github.com/banzaicloud/pipeline/internal/clusterfeature/clusterfeatureadapter"
	"github.com/banzaicloud/pipeline/internal/common"
)

type FeatureManager struct {
	clusterGetter clusterfeatureadapter.ClusterGetter
	config        Configuration
	logger        common.Logger
}

// NewVaultFeatureManager builds a new feature manager component
func MakeFeatureManager(
	clusterGetter clusterfeatureadapter.ClusterGetter,
	config Configuration,
	logger common.Logger,
) FeatureManager {
	return FeatureManager{
		clusterGetter: clusterGetter,
		config:        config,
		logger:        logger,
	}
}

func (FeatureManager) Name() string {
	return featureName
}

func (m FeatureManager) GetOutput(ctx context.Context, clusterID uint, spec clusterfeature.FeatureSpec) (clusterfeature.FeatureOutput, error) {
	// TODO (colin): implement me
	return clusterfeature.FeatureOutput{}, nil
}

func (m FeatureManager) ValidateSpec(ctx context.Context, spec clusterfeature.FeatureSpec) error {
	vaultSpec, err := bindFeatureSpec(spec)
	if err != nil {
		return err
	}

	if err := vaultSpec.Validate(); err != nil {
		return clusterfeature.InvalidFeatureSpecError{
			FeatureName: featureName,
			Problem:     err.Error(),
		}
	}

	return nil
}

func (m FeatureManager) PrepareSpec(ctx context.Context, spec clusterfeature.FeatureSpec) (clusterfeature.FeatureSpec, error) {
	return spec, nil
}
