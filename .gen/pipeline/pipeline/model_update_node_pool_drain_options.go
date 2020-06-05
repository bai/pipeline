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

// UpdateNodePoolDrainOptions - Drain options for old nodes.
type UpdateNodePoolDrainOptions struct {

	// How long should drain wait for pod eviction (in seconds)
	Timeout int32 `json:"timeout,omitempty"`

	// Whether the process should fail if draining fails/times out.
	FailOnError bool `json:"failOnError,omitempty"`

	// Only evict those pods that matches this selector.
	PodSelector string `json:"podSelector,omitempty"`
}