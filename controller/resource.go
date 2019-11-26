package controller

import (
	"github.com/fundata-varena/fundata-resource-server/business"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetResource(ctx *gin.Context) {
	resourceType := ctx.Query("source_type")
	resourceId := ctx.Query("source_id")
	if resourceType == "" || resourceId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": "",
		})
	}

	response := business.GetResource(resourceType, resourceId)
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message": "",
		"data": response,
	})
}

func GetResources(ctx *gin.Context) {

}
