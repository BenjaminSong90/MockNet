package middleware

import (
	"github.com/gin-gonic/gin"
	"mock_net/model"
)

func EnhanceMDW(mockConfig model.ProjectConfig) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set("project_config", mockConfig)
	}
}
