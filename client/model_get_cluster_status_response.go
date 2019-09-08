/*
 * Pipeline API
 *
 * Pipeline is a feature rich application platform, built for containers on top of Kubernetes to automate the DevOps experience, continuous application development and the lifecycle of deployments.
 *
 * API version: latest
 * Contact: info@banzaicloud.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

type GetClusterStatusResponse struct {
	Status        string `json:"status,omitempty"`
	StatusMessage string `json:"statusMessage,omitempty"`
	Name          string `json:"name,omitempty"`
	Cloud         string `json:"cloud,omitempty"`
	Distribution  string `json:"distribution,omitempty"`
	Version       string `json:"version,omitempty"`
	Spot          bool   `json:"spot,omitempty"`
	Location      string `json:"location,omitempty"`
	Id            int32  `json:"id,omitempty"`
	Logging       bool   `json:"logging,omitempty"`
	Monitoring    bool   `json:"monitoring,omitempty"`
	Servicemesh   bool   `json:"servicemesh,omitempty"`
	Securityscan  bool   `json:"securityscan,omitempty"`
	CreatedAt     string `json:"createdAt,omitempty"`
	StartedAt     string `json:"startedAt,omitempty"`
	CreatorName   string `json:"creatorName,omitempty"`
	CreatorId     int32  `json:"creatorId,omitempty"`
	Region        string `json:"region,omitempty"`
	// The lifespan of the cluster expressed in minutes after which it is automatically deleted. Zero value means the cluster is never automatically deleted.
	TtlMinutes   int32                     `json:"ttlMinutes,omitempty"`
	NodePools    map[string]NodePoolStatus `json:"nodePools,omitempty"`
	TotalSummary ResourceSummary           `json:"totalSummary,omitempty"`
}
