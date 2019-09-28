/*
 * Pipeline API
 *
 * Pipeline is a feature rich application platform, built for containers on top of Kubernetes to automate the DevOps experience, continuous application development and the lifecycle of deployments. 
 *
 * API version: latest
 * Contact: info@banzaicloud.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package pipeline

type CreateClusterRequestV2 struct {

	Name string `json:"name"`

	Features []Feature `json:"features,omitempty"`

	SecretId string `json:"secretId,omitempty"`

	SecretName string `json:"secretName,omitempty"`

	SshSecretId string `json:"sshSecretId,omitempty"`

	ScaleOptions ScaleOptions `json:"scaleOptions,omitempty"`

	Type string `json:"type"`

	Kubernetes CreatePkeClusterKubernetes `json:"kubernetes"`

	// Non-existent resources will be created in this location. Existing resources that must have the same location as the cluster will be validated against this.
	Location string `json:"location,omitempty"`

	// Required resources will be created in this resource group.
	ResourceGroup string `json:"resourceGroup"`

	Network PkeOnAzureClusterNetwork `json:"network,omitempty"`

	Nodepools []PkeOnAzureNodePool `json:"nodepools,omitempty"`

	// List of access points (i.e. load balancers, floating IPs) to be created for the cluster. Access points are implemented using cloud provider specific resources.
	AccessPoints []string `json:"accessPoints,omitempty"`

	// List of access point references for the API server; currently, public and private are the only valid values
	ApiServerAccessPoints []string `json:"apiServerAccessPoints,omitempty"`
}
