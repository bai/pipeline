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
	pkgCluster "github.com/banzaicloud/pipeline/cluster"
	"github.com/banzaicloud/pipeline/config"
	"github.com/banzaicloud/pipeline/internal/clusterfeature"
	"github.com/banzaicloud/pipeline/internal/clusterfeature/clusterfeatureadapter"
	"github.com/banzaicloud/pipeline/internal/clusterfeature/features"
	"github.com/banzaicloud/pipeline/internal/common"
	pkgSecret "github.com/banzaicloud/pipeline/pkg/secret"
	"github.com/banzaicloud/pipeline/secret"
)

// FeatureOperator implements the Logging feature operator
type FeatureOperator struct {
	clusterGetter  clusterfeatureadapter.ClusterGetter
	clusterService clusterfeature.ClusterService
	helmService    features.HelmService
	secretStore    features.SecretStore
	config         Configuration
	logger         common.Logger
}

// MakeFeatureOperator returns a Logging feature operator
func MakeFeatureOperator(
	clusterGetter clusterfeatureadapter.ClusterGetter,
	clusterService clusterfeature.ClusterService,
	helmService features.HelmService,
	secretStore features.SecretStore,
	config Configuration,
	logger common.Logger,
) FeatureOperator {
	return FeatureOperator{
		clusterGetter:  clusterGetter,
		clusterService: clusterService,
		helmService:    helmService,
		secretStore:    secretStore,
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

	// generate TLS secret
	if err := op.generateTLSSecret(ctx, clusterID); err != nil {
		return errors.WrapIf(err, "failed to generate TLS secret for logging feature")
	}

	// install TLS secret to cluster
	if err := op.installTLSSecretToCluster(ctx, clusterID); err != nil {
		return errors.WrapIf(err, "failed to install TLS secret")
	}

	// install logging-operator
	if err := op.installLoggingOperator(ctx, clusterID, logger); err != nil {
		return errors.WrapIf(err, fmt.Sprintf("failed to install deployment: %q", op.config.operator.chartName))
	}

	// install logging-operator-logging
	if err := op.installLoggingOperatorLogging(ctx, clusterID, logger); err != nil {
		return errors.WrapIf(err, fmt.Sprintf("failed to install deployment: %q", op.config.logging.chartName))
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
	if err := op.helmService.DeleteDeployment(ctx, clusterID, config.LoggingOperatorReleaseName); err != nil {
		logger.Info(fmt.Sprintf("failed to delete feature deployment: %q", config.LoggingOperatorReleaseName))

		return errors.WrapIf(err, fmt.Sprintf("failed to uninstall deployment: %q", config.LoggingOperatorReleaseName))
	}
	// TODO (colin): remove TLS secret from Vault

	return nil
}

func (op FeatureOperator) installLoggingOperator(ctx context.Context, clusterID uint, logger common.Logger) error {
	var chartValues = loggingOperatorValues{
		// todo (colin): extend me
	}

	valuesBytes, err := json.Marshal(chartValues)
	if err != nil {
		logger.Debug("failed to marshal chartValues")
		return errors.WrapIf(err, "failed to decode chartValues")
	}

	return op.helmService.ApplyDeployment(
		ctx,
		clusterID,
		op.config.pipelineSystemNamespace,
		op.config.operator.chartName,
		config.LoggingOperatorReleaseName,
		valuesBytes,
		op.config.operator.chartVersion,
	)
}

func (op FeatureOperator) installLoggingOperatorLogging(ctx context.Context, clusterID uint, logger common.Logger) error {
	var chartValues = loggingOperatorLoggingValues{
		// todo (colin): extend me
	}

	valuesBytes, err := json.Marshal(chartValues)
	if err != nil {
		logger.Debug("failed to marshal chartValues")
		return errors.WrapIf(err, "failed to decode chartValues")
	}

	return op.helmService.ApplyDeployment(
		ctx,
		clusterID,
		op.config.pipelineSystemNamespace,
		op.config.logging.chartName,
		config.LoggingReleaseName,
		valuesBytes,
		op.config.logging.chartVersion,
	)
}

func (op FeatureOperator) generateTLSSecret(ctx context.Context, clusterID uint) error {
	var namespace = op.config.pipelineSystemNamespace
	var tlsHost = "fluentd." + namespace + ".svc.cluster.local"
	var secretName = getTLSSecretName(clusterID)

	req := &secret.CreateSecretRequest{
		Name: secretName,
		Type: pkgSecret.TLSSecretType,
		Tags: []string{
			pkgSecret.TagBanzaiReadonly,
			secretTag,
		},
		Values: map[string]string{
			pkgSecret.TLSHosts: tlsHost,
		},
	}

	// store TLS secret
	_, err := op.secretStore.Store(ctx, req)
	if err != nil {
		return errors.WrapIf(err, "failed generate TLS secrets to logging operator")
	}

	return nil
}

func (op FeatureOperator) installTLSSecretToCluster(ctx context.Context, clusterID uint) error {
	cl, err := op.clusterGetter.GetClusterByIDOnly(ctx, clusterID)
	if err != nil {
		return errors.WrapIf(err, "failed to get cluster")
	}

	var namespace = op.config.pipelineSystemNamespace
	var secretName = getTLSSecretName(cl.GetID())
	secretRequest := pkgCluster.InstallSecretRequest{
		SourceSecretName: secretName,
		Namespace:        namespace,
		Update:           true,
	}

	_, err = pkgCluster.InstallSecret(cl, secretName, secretRequest)
	if err != nil {
		return errors.WrapIfWithDetails(err, "failed to install secret to the cluster", "clusterID", cl.GetID())
	}

	return nil
}

func (op FeatureOperator) deleteTLSSecret(ctx context.Context) error {
	// 	op.secretStore.Delete(ctx, )
	return nil
}
