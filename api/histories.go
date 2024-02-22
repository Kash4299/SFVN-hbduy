package api

import (
	"sfvn-hbduy/common/response"
	"sfvn-hbduy/common/util"
	"sfvn-hbduy/service"
	"time"

	"github.com/gin-gonic/gin"
)

type Histories struct {
	historiesService service.IHistoriesService
}

func APIHistoriesHandler(r *gin.Engine, historiesService service.IHistoriesService) {
	handler := &Histories{
		historiesService: historiesService,
	}
	Group := r.Group("v1/get_histories")
	{
		Group.GET("", handler.GetHistories)
	}
}

func (h *Histories) GetHistories(c *gin.Context) {
	start_date := c.Query("start_date")
	if start_date == "" {
		c.JSON(response.BadRequestMsg("start_date is required"))
		return
	}
	end_date := c.Query("end_date")
	if end_date == "" {
		c.JSON(response.BadRequestMsg("end_date is required"))
		return
	}
	symbol := c.Query("symbol")
	if symbol == "" {
		c.JSON(response.BadRequestMsg("symbol is required"))
		return
	}

	period := c.Query("period")

	startTime, err := time.Parse("2006-01-02", start_date)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	endTime, err := time.Parse("2006-01-02", end_date)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	// validate start_date < end_date
	if startTime.After(endTime) {
		c.JSON(response.BadRequestMsg("start_date must before end_date"))
	}

	days := int(endTime.Sub(startTime).Hours() / 24)
	if !util.IsValidDays(days) {
		c.JSON(response.BadRequestMsg("Out of range"))
		return
	}

	code, result := h.historiesService.GetHistories(c, symbol, days, period)
	c.JSON(code, result)
}
