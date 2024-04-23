package app_engine

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHttpApiRoutes(r *gin.Engine) {
	r.GET("/pods", func(ctx *gin.Context) {
		query := struct {
			Namespace string `form:"namespace"`
		}{}
		if err := ctx.BindQuery(&query); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		list, err := ListPod(ctx.Request.Context(), query.Namespace)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, list)
	})
}
