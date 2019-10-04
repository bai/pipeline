/*
 * Anchore Engine API Server
 *
 * This is the Anchore Engine API. Provides the primary external API for users of the service.
 *
 * API version: 0.1.12
 * Contact: nurmi@anchore.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package anchore
// GroupSyncResult struct for GroupSyncResult
type GroupSyncResult struct {
	// The name of the group
	Group string `json:"group,omitempty"`
	Status string `json:"status,omitempty"`
	// The number of images updated by the this group sync, across all accounts. This is typically only non-zero for vulnerability feeds which update images' vulnerability results during the sync.
	UpdatedImageCount int32 `json:"updated_image_count,omitempty"`
	// The number of feed data records synced down as either updates or new records
	UpdatedRecordCount int32 `json:"updated_record_count,omitempty"`
	// The duration of the group sync in seconds
	TotalTimeSeconds float32 `json:"total_time_seconds,omitempty"`
}