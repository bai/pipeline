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
	"fmt"

	"github.com/banzaicloud/pipeline/config"
	"github.com/banzaicloud/pipeline/internal/clusterfeature/clusterfeatureadapter"
)

const (
	featureName        = "logging"
	providerAmazonS3   = "s3"
	providerGoogleGCS  = "gcs"
	providerAlibabaOSS = "oss"
	providerAzure      = "azure"
	secretTag          = "feature:logging"
)

type obj = map[string]interface{}

func getTLSSecretName(clusterID uint) string {
	return fmt.Sprintf("logging-tls-%d", clusterID)
}

func getReleaseSecretTag() string {
	return fmt.Sprintf("release:%s", config.LoggingOperatorReleaseName)
}

func getSecretClusterTags(cluster clusterfeatureadapter.Cluster) []string {
	clusterNameSecretTag := fmt.Sprintf("cluster:%s", cluster.GetName())
	clusterUidSecretTag := fmt.Sprintf("clusterUID:%s", cluster.GetUID())
	return []string{clusterNameSecretTag, clusterUidSecretTag}
}
