package files

import (
	//"github.com/dmitryDevGoMid/gokeeper/grpcapi/gateway/internal/auth"

	"net/http"

	"github.com/dmitryDevGoMid/gokeeper/grpcapi/gateway/internal/config"
	"github.com/dmitryDevGoMid/gokeeper/grpcapi/gateway/internal/files/routes"

	"github.com/gin-gonic/gin"
)

var cancelDownloadFile = make(map[string]chan struct{})

func RegisterRoutes(r *gin.Engine, cfg *config.Config) { //, authSvc *auth.ServiceClient) {

	svc := &ServiceClient{
		Client: InitServiceClient(cfg),
	}

	r.Use(CheckAuthByToken(cfg))

	routes := r.Group("/files")
	routes.POST("/upload", svc.SendFiles)
	routes.POST("/download", svc.GetFile)

	routes.POST("/download/cancel", func(c *gin.Context) {
		udiFile := c.GetHeader("File-ID")
		if value, ok := cancelDownloadFile[udiFile]; ok {
			value <- struct{}{}
			c.JSON(http.StatusOK, gin.H{"cancel": "success"})
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"cancel": "not found"})

	})
}

func (svc *ServiceClient) SendFiles(ctx *gin.Context) {
	routes.SendFiles(ctx, svc.Client)
}

func (svc *ServiceClient) GetFile(ctx *gin.Context) {
	routes.GetFile(ctx, svc.Client, &cancelDownloadFile)
}
