package model

type ResourceUpdated struct {
	ResourceID   string `json:"resource_id"`
	ResourceType string `json:"resource_type"`
	Size         string `json:"size"`
	UpdatedTime  string `json:"updated_time"`
}
