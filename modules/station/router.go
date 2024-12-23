package station

import (
	"net/http"

	"github.com/arjunstein/mrt-schedules/common/response"
	"github.com/gin-gonic/gin"
)

func Initiate(router *gin.RouterGroup) {
	stationService := NewService()
	station := router.Group("/stations")
	station.GET("", func(c *gin.Context) {
		GetAllStations(c, stationService)
	})

	station.GET("/:id", func(c *gin.Context) {
		CheckSchedulesByStation(c, stationService)
	})
}

func GetAllStations(c *gin.Context, service Service) {
	// code service here
	datas, err := service.GetAllStations()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Success get all stations",
		Data:    datas,
		})
}

func CheckSchedulesByStation(c *gin.Context, service Service) {
	// code service here
	id := c.Param("id")
	datas, err := service.CheckSchedulesByStation(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Successfully get schedules by station",
		Data:    datas,
	})
}