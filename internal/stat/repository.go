package stat

import (
	"go/adv-demo/pkg/db"
	"time"

	"gorm.io/datatypes"
)

type StatRepository struct {
	Database *db.Db
}

func NewStatRepository(db *db.Db) *StatRepository{
	return &StatRepository{
		Database: db,
	}
}

func (repo *StatRepository) AddClick(linkId uint){
	var stat Stat
	curr_date := datatypes.Date(time.Now())

	repo.Database.DB.Find(&stat, "link_id = ? and date = ?", linkId, curr_date)

	if stat.ID == 0{
		repo.Database.DB.Create(&Stat{
			LinkId: linkId,
			Clicks: 1,
			Date: curr_date,
		})
	} else {
		stat.Clicks += 1
		repo.Database.DB.Save(&stat)
	}
}

func (repo *StatRepository) GetStats(by string, from, to time.Time) []GetStatResponse{
	var stats []GetStatResponse
	var selectQuery string

	switch by {
	case GroupByDay:
		selectQuery = "to_char(date, 'YYYY-MM-DD') as period, sum(clicks)"
	case GroupByMonth:
		selectQuery = "to_char(date, 'YYYY-MM') as period, sum(clicks)"
	}

	repo.Database.DB.Table("stats").
		Select(selectQuery).
		Where("date BETWEEN ? AND ?", from, to).
		Group("period").
		Order("period").
		Scan(&stats)
	
	return stats
}