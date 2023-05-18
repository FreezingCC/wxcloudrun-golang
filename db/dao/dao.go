package dao

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

const tableName = "Counters"

// ClearCounter 清除Counter
func (imp *CounterInterfaceImp) ClearCounter(id int32) error {
	cli := db.Get()
	return cli.Table(tableName).Delete(&model.CounterModel{Id: id}).Error
}

// UpsertCounter 更新/写入counter
func (imp *CounterInterfaceImp) UpsertCounter(counter *model.CounterModel) error {
	cli := db.Get()
	return cli.Table(tableName).Save(counter).Error
}

// GetCounter 查询Counter
func (imp *CounterInterfaceImp) GetCounter(id int32) (*model.CounterModel, error) {
	var err error
	var counter = new(model.CounterModel)

	cli := db.Get()
	err = cli.Table(tableName).Where("id = ?", id).First(counter).Error

	return counter, err
}

func (sImp UserMaxScoreImp) GetScoreByUserId(userId string) (*model.UserMaxScore, error) {
	cli := db.Get()
	var score = new(model.UserMaxScore)
	err := cli.Table("user_max_score").Where("user_id = ?", userId).First(score).Error
	return score, err
}

func (sImp UserMaxScoreImp) UpdateScoreByUserId(userId string, score int32) error {
	scoreM := model.UserMaxScore{
		UserId: userId,
		Score:  score,
	}
	cli := db.Get()
	table := cli.Table("user_max_score")
	_ = table.Clauses(
		clause.OnConflict{Columns: []clause.Column{{Table: "user_max_score", Name: "user_id"}},
			DoUpdates: clause.Set{{Column: clause.Column{Table: "user_max_score", Name: "score"},
				Value: gorm.Expr("IF(score < ?, ?, score)", scoreM.Score, scoreM.Score)}}},
	).Create(&scoreM)
	return nil
}
