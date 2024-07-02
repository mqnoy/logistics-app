package model

type GoodStock struct {
	ID       int64 `gorm:"column:id"`
	Total    int   `gorm:"column:total;default:0"`
	GoodID   string
	Currency Good `gorm:"foreignKey:GoodID"`
	TimestampColumn
}
