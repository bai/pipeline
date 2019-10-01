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

	"emperror.dev/errors"
	"github.com/banzaicloud/pipeline/internal/clusterfeature"
	"github.com/banzaicloud/pipeline/internal/clusterfeature/clusterfeatureadapter"
	"github.com/banzaicloud/pipeline/internal/clusterfeature/features"
	"github.com/banzaicloud/pipeline/internal/common"
)

type FeatureManager struct {
	clusterGetter clusterfeatureadapter.ClusterGetter
	config        Configuration
	secretStore   features.SecretStore
	logger        common.Logger
}

// NewVaultFeatureManager builds a new feature manager component
func MakeFeatureManager(
	clusterGetter clusterfeatureadapter.ClusterGetter,
	config Configuration,
	secretStore features.SecretStore,
	logger common.Logger,
) FeatureManager {
	return FeatureManager{
		clusterGetter: clusterGetter,
		config:        config,
		secretStore:   secretStore,
		logger:        logger,
	}
}

func (FeatureManager) Name() string {
	return featureName
}

func (m FeatureManager) GetOutput(ctx context.Context, clusterID uint, spec clusterfeature.FeatureSpec) (clusterfeature.FeatureOutput, error) {
	// TODO (colin): extends me

	var tlsSecretName = getTLSSecretName(clusterID)
	tlsSecretID, err := m.secretStore.GetIDByName(ctx, tlsSecretName)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to get TLS secret")
	}

	var output = clusterfeature.FeatureOutput{
		"loggingOperator": map[string]interface{}{
			"version": m.config.operator.chartVersion,
		},
		"logging": map[string]interface{}{
			"version": m.config.logging.chartVersion,
		},
		"tls": map[string]interface{}{
			"secretId": tlsSecretID,
		},
	}
	return output, nil
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
