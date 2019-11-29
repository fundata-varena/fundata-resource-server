package business

import "github.com/fundata-varena/fundata-resource-server/model"

func GetResource(resourceType, id string) (*model.ResourceLocal, error) {
	ops := new(model.ResourceOps)
	return ops.GetResource(resourceType, id)
}