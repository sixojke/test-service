package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/sixojke/test-service/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()

	good := router.Group("/good")
	{
		good.POST("/create", h.createGood)
		good.PATCH("/update", h.updateGood)
		good.DELETE("/remove", h.deleteGood)
		good.PATCH("/reprioritiize", h.reprioritiizeGood)
	}

	goods := router.Group("/goods")
	{
		goods.GET("/list", h.getListGoods)
	}

	return router
}
