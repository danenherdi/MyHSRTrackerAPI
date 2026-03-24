package handlers

import (
	"net/http"
	"strconv"

	"MyHSRTrackerAPI/database"
	"MyHSRTrackerAPI/models"
	"MyHSRTrackerAPI/services"
	"MyHSRTrackerAPI/utils"

	"github.com/gin-gonic/gin"
)

func ImportWarp(c *gin.Context) {
	var input struct {
		AuthURL string `json:"auth_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format or missing auth_url")
		return
	}

	count, err := services.SyncWarpLogs(input.AuthURL)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Import successful", gin.H{
		"imported": count,
	})
}

func GetWarpList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	gachaType := c.Query("gacha_type")
	uid := c.Query("uid")

	var records []models.WarpLog
	query := database.DB.Model(&models.WarpLog{})

	if gachaType != "" {
		query = query.Where("gacha_type = ?", gachaType)
	}
	if uid != "" {
		query = query.Where("uid = ?", uid)
	}

	offset := (page - 1) * size
	query.Order("time DESC").Limit(size).Offset(offset).Find(&records)

	var total int64
	// Reset limits for total count
	queryCount := database.DB.Model(&models.WarpLog{})
	if gachaType != "" {
		queryCount = queryCount.Where("gacha_type = ?", gachaType)
	}
	if uid != "" {
		queryCount = queryCount.Where("uid = ?", uid)
	}
	queryCount.Count(&total)

	utils.SuccessResponse(c, http.StatusOK, "Records retrieved", gin.H{
		"list":  records,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

func GetWarpStats(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "uid parameter is required")
		return
	}

	// Calculate pity for each gacha type for the given UID
	var types = []string{"11", "1", "2", "12"}
	stats := make(map[string]map[string]interface{})

	for _, gType := range types {
		var total int64
		database.DB.Model(&models.WarpLog{}).Where("uid = ? AND gacha_type = ?", uid, gType).Count(&total)

		var last5Star models.WarpLog
		res := database.DB.Where("uid = ? AND gacha_type = ? AND rank_type = ?", uid, gType, 5).Order("time DESC").First(&last5Star)

		var pity int64
		if res.Error == nil {
			database.DB.Model(&models.WarpLog{}).Where("uid = ? AND gacha_type = ? AND time > ?", uid, gType, last5Star.Time).Count(&pity)
		} else {
			pity = total
		}

		stats[gType] = map[string]interface{}{
			"total": total,
			"pity":  pity,
		}
	}

	utils.SuccessResponse(c, http.StatusOK, "Stats retrieved", stats)
}
