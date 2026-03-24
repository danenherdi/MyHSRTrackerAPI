package services

import (
	"fmt"
	"log"
	"time"

	"MyHSRTrackerAPI/database"
	"MyHSRTrackerAPI/models"
)

var GachaTypes = []string{"11", "1", "2", "12"}

func SyncWarpLogs(authURL string) (int, error) {
	q, err := ExtractQueryParams(authURL)
	if err != nil {
		return 0, err
	}

	totalImported := 0

	for _, gType := range GachaTypes {
		endId := "0"
		for {
			res, err := FetchGachaLog(q, gType, endId)
			if err != nil {
				return totalImported, err
			}
			if res.Retcode != 0 {
				log.Printf("API Error: %s", res.Message)
				return totalImported, fmt.Errorf("API error: %s", res.Message)
			}

			if res.Data == nil || len(res.Data.List) == 0 {
				break
			}

			for _, logItem := range res.Data.List {
				// Prevent duplicate insert if ID exists
				var existing models.WarpLog
				if err := database.DB.Where("id = ?", logItem.ID).First(&existing).Error; err == nil {
					// We reached records that are already imported for this gacha type
					goto NextGachaType
				}

				// Parse time
				logTime, _ := time.Parse("2006-01-02 15:04:05", logItem.Time)
				rankType := 3 // fallback
				fmt.Sscanf(logItem.RankType, "%d", &rankType)

				newLog := models.WarpLog{
					ID:        logItem.ID,
					UID:       logItem.UID,
					ItemID:    logItem.ItemID,
					Name:      logItem.Name,
					ItemType:  logItem.ItemType,
					RankType:  rankType,
					GachaType: logItem.GachaType,
					Time:      logTime,
				}
				database.DB.Create(&newLog)
				totalImported++
			}

			// Pagination
			lastItem := res.Data.List[len(res.Data.List)-1]
			endId = lastItem.ID

			// Sleep to avoid rate limit (500ms should be safe)
			time.Sleep(500 * time.Millisecond)
		}
	NextGachaType:
	}

	return totalImported, nil
}
