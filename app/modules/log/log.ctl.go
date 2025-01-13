package log

import (
	logdto "app/app/modules/log/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LogController struct {
	logSvc *LogService
}

func newController(logService *LogService) *LogController {
	return &LogController{
		logSvc: logService,
	}
}

func (ctl *LogController) CreateLog(c *gin.Context) {
	ctx := c.Request.Context()
	var req logdto.LogDTORequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logs, err := ctl.logSvc.CreateLogs(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

func (ctl *LogController) GetAllLog(c *gin.Context) {
	ctx := c.Request.Context()

	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	logs, err := ctl.logSvc.GetAllLog(ctx, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)

}

func (ctl *LogController) DeleteLog(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	err := ctl.logSvc.DeleteLog(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Log deleted successfully"})
}
