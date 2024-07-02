package mysql

import (
	"github.com/mqnoy/logistics-app/core/domain"
	"github.com/mqnoy/logistics-app/core/dto"
	"github.com/mqnoy/logistics-app/core/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type mysqlGoodRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) domain.GoodRepository {
	return &mysqlGoodRepository{
		db: db,
	}
}

func (m mysqlGoodRepository) InsertGood(row model.Good) (*model.Good, error) {
	err := m.db.Create(&row).Error
	return &row, err
}

func (m mysqlGoodRepository) SelectGoodByCode(code string) (row *model.Good, err error) {
	if err := m.db.
		Where("code=?", code).First(&row).
		Error; err != nil {
		return nil, err
	}

	return row, nil
}

func (m mysqlGoodRepository) InsertGoodStock(row model.GoodStock) (*model.GoodStock, error) {
	err := m.db.Create(&row).Error
	return &row, err
}

func (m mysqlGoodRepository) SelectGoodById(id string) (*model.Good, error) {
	var row model.Good
	if err := m.db.
		Joins("GoodStock").
		First(&row, "Good.id = ?", id).
		Error; err != nil {
		return nil, err
	}

	return &row, nil
}

func (m mysqlGoodRepository) SelectAndCountGood(param dto.ListParam[dto.FilterCommonParams]) (result dto.SelectAndCount[model.Good], err error) {
	var rows []*model.Good
	var count int64

	filters := param.Filters
	orders := param.Orders
	pagination := param.Pagination
	whereClause := clause.Where{}

	if filters.Keyword != "" {
		whereClause.Exprs = append(whereClause.Exprs, clause.Where{
			Exprs: []clause.Expression{
				clause.Or(
					clause.Like{
						Column: clause.Column{Name: "name"},
						Value:  "%" + filters.Keyword + "%",
					},
					clause.Like{
						Column: clause.Column{Name: "code"},
						Value:  "%" + filters.Keyword + "%",
					},
				),
			},
		})
	}

	if filters.IsActive != nil {
		whereClause.Exprs = append(whereClause.Exprs, clause.Eq{
			Column: "is_active",
			Value:  *filters.IsActive,
		})
	}

	mDB := m.db
	if len(whereClause.Exprs) > 0 {
		mDB = m.db.Clauses(whereClause)
	}

	mDB.Model(&model.Good{}).Count(&count)

	if err = mDB.
		Joins("GoodStock").
		Limit(pagination.Limit).Offset(pagination.Offset).
		Order(orders).
		Find(&rows).Error; err != nil {
		return result, err
	}

	return dto.SelectAndCount[model.Good]{
		Rows:  rows,
		Count: count,
	}, nil
}
