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