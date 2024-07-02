package model

type GoodStock struct {
	ID     uint   `gorm:"column:id;primaryKey"`
	Total  int    `gorm:"column:total;default:0"`
	GoodID string `gorm:"column:good_id"`
	TimestampColumn
}
