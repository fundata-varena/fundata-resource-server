package model

import "time"

type ResourceUpdated struct {
	ResourceID   string `json:"resource_id"`
	ResourceType string `json:"resource_type"`
	Size         string `json:"size"`
	UpdatedTime  time.Time `json:"updated_time"`
}
