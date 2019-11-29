package controller

import (
	"github.com/fundata-varena/fundata-resource-server/business"
	"github.com/fundata-varena/fundata-resource-server/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func GetResource(ctx *gin.Context) {
	resourceType := ctx.Query("resource_type")
	resourceId := ctx.Query("resource_id")

	response := response{}
	response["code"] = http.StatusOK
	response["message"] = "success"

	if resourceType == "" || resourceId == "" {
		response["code"] = http.StatusBadRequest
		response["message"] = "illegal parameters"
		ctx.JSON(http.StatusBadRequest, gin.H(response))
		return
	}

	row, err := business.GetResource(resourceType, resourceId)
	if err != nil {
		response["code"] = http.StatusInternalServerError
		response["message"] = "internal server error"
		ctx.JSON(http.StatusInternalServerError, gin.H(response))
		return
	}

	response["data"] = row

	ctx.JSON(http.StatusOK, gin.H(response))
}

func GetResources(ctx *gin.Context) {
	response := response{}
	response["code"] = http.StatusOK
	response["message"] = "success"

	resources := ctx.QueryArray("resources")
	allHasSeparator := true
	for _, resource := range resources {
		// 要求每个参数里都有半角逗号分隔符
		if strings.Contains(resource, RESOURCES_SEPARATOR) {
			continue
		}
		allHasSeparator = false
		break
	}
	if len(resources) == 0 || !allHasSeparator {
		response["code"] = http.StatusBadRequest
		response["message"] = "illegal parameters"
		ctx.JSON(http.StatusBadRequest, gin.H(response))
		return
	}

	var rows []*model.ResourceLocal

	for _, resource := range resources {
		arrStr := strings.Split(resource, RESOURCES_SEPARATOR)
		row, err := business.GetResource(arrStr[0], arrStr[1])
		if err != nil {
			continue
		}
		rows = append(rows, row)
	}

	response["data"] = rows

	ctx.JSON(http.StatusOK, gin.H(response))
}
