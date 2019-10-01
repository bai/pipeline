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
	"encoding/json"
	"fmt"

	"emperror.dev/errors"
	"github.com/banzaicloud/pipeline/config"
	"github.com/banzaicloud/pipeline/internal/clusterfeature"
	"github.com/banzaicloud/pipeline/internal/clusterfeature/clusterfeatureadapter"
	"github.com/banzaicloud/pipeline/internal/clusterfeature/features"
	"github.com/banzaicloud/pipeline/internal/common"
)

// FeatureOperator implements the Logging feature operator
type FeatureOperator struct {
	clusterGetter  clusterfeatureadapter.ClusterGetter
	clusterService clusterfeature.ClusterService
	helmService    features.HelmService
	config         Configuration
	logger         common.Logger
}

// MakeFeatureOperator returns a Logging feature operator
func MakeFeatureOperator(
	clusterGetter clusterfeatureadapter.ClusterGetter,
	clusterService clusterfeature.ClusterService,
	helmService features.HelmService,
	config Configuration,
	logger common.Logger,
) FeatureOperator {
	return FeatureOperator{
		clusterGetter:  clusterGetter,
		clusterService: clusterService,
		helmService:    helmService,
		config:         config,
		logger:         logger,
	}
}

// Name returns the name of the Logging feature
func (op FeatureOperator) Name() string {
	return featureName
}

// Apply applies the provided specification to the cluster feature
func (op FeatureOperator) Apply(ctx context.Context, clusterID uint, spec clusterfeature.FeatureSpec) error {
	if err := op.clusterService.CheckClusterReady(ctx, clusterID); err != nil {
		return err
	}

	logger := op.logger.WithContext(ctx).WithFields(map[string]interface{}{"cluster": clusterID, "feature": featureName})

	//boundSpec, err := bindFeatureSpec(spec)
	_, err := bindFeatureSpec(spec)
	if err != nil {
		return clusterfeature.InvalidFeatureSpecError{
			FeatureName: featureName,
			Problem:     err.Error(),
		}
	}

	var chartValues = loggingOperatorValues{
		// todo (colin): extend me
	}

	valuesBytes, err := json.Marshal(chartValues)
	if err != nil {
		logger.Debug("failed to marshal chartValues")
		return errors.WrapIf(err, "failed to decode chartValues")
	}

	var chartName = op.config.operator.chartName

	if err := op.helmService.ApplyDeployment(
		ctx,
		clusterID,
		op.config.pipelineSystemNamespace,
		chartName,
		config.LoggingReleaseName,
		valuesBytes,
		op.config.operator.chartVersion,
	); err != nil {
		return errors.WrapIf(err, fmt.Sprintf("failed to install deployment: %q", chartName))
	}

	return nil
}

// Deactivate deactivates the cluster feature
func (op FeatureOperator) Deactivate(ctx context.Context, clusterID uint, spec clusterfeature.FeatureSpec) error {
	if err := op.clusterService.CheckClusterReady(ctx, clusterID); err != nil {
		return err
	}

	logger := op.logger.WithContext(ctx).WithFields(map[string]interface{}{"cluster": clusterID, "feature": featureName})

	// delete deployment
	if err := op.helmService.DeleteDeployment(ctx, clusterID, config.LoggingReleaseName); err != nil {
		logger.Info(fmt.Sprintf("failed to delete feature deployment: %q", config.LoggingReleaseName))

		return errors.WrapIf(err, fmt.Sprintf("failed to uninstall deployment: %q", config.LoggingReleaseName))
	}

	return nil
}
