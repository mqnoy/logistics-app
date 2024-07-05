package mysql

import (
	"log"

	"github.com/mqnoy/logistics-app/core/domain"
	"github.com/mqnoy/logistics-app/core/dto"
	"github.com/mqnoy/logistics-app/core/model"
	"github.com/mqnoy/logistics-app/core/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type mysqlOrderRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) domain.OrderRepository {
	return &mysqlOrderRepository{
		db: db,
	}
}

func (m mysqlOrderRepository) WithTrx(trxHandle *gorm.DB) domain.OrderRepository {
	if trxHandle == nil {
		log.Println("transaction not found")
		return m
	}
	m.db = trxHandle

	return m
}

func (m mysqlOrderRepository) InsertOrder(row model.Order) (*model.Order, error) {
	err := m.db.Create(&row).Error
	return &row, err
}

func (m mysqlOrderRepository) SelectAndCountOrder(param dto.ListParam[dto.FilterOrderParams]) (result dto.SelectAndCount[model.Order], err error) {
	var rows []*model.Order
	var count int64

	filters := param.Filters
	orders := param.Orders
	pagination := param.Pagination
	whereClause := clause.Where{}

	if len(filters.RequestAt) > 1 {
		from, err := util.NumberToEpoch(filters.RequestAt[0])
		if err != nil {
			return result, err
		}

		to, err := util.NumberToEpoch(filters.RequestAt[1])
		if err != nil {
			return result, err
		}
		expressions := []clause.Expression{
			clause.Expr{
				SQL: "request_at BETWEEN ? AND ?", Vars: []interface{}{
					from, to,
				},
			},
		}
		whereClause.Exprs = append(whereClause.Exprs, clause.Where{Exprs: expressions})

	}

	if filters.OrderType > 0 {
		whereClause.Exprs = append(whereClause.Exprs, clause.Eq{
			Column: clause.Column{Name: "type"},
			Value:  filters.OrderType,
		})
	}

	mDB := m.db
	if len(whereClause.Exprs) > 0 {
		mDB = m.db.Clauses(whereClause)
	}

	if filters.GoodId != "" {
		mDB = m.db.Joins("INNER JOIN OrderItem ON OrderItem.order_id = Order.id").
			Where("OrderItem.good_id = ?", filters.GoodId).
			Clauses(whereClause)
	}

	mDB.Model(&model.Order{}).Count(&count)

	if err = mDB.
		Limit(pagination.Limit).Offset(pagination.Offset).
		Order(orders).
		Find(&rows).Error; err != nil {
		return result, err
	}

	return dto.SelectAndCount[model.Order]{
		Rows:  rows,
		Count: count,
	}, nil
}

func (m mysqlOrderRepository) SelectOrderById(id string) (row *model.Order, err error) {
	if err := m.db.
		Preload("OrderItem").
		First(&row, "id = ?", id).
		Error; err != nil {
		return nil, err
	}

	return row, nil
}

func (m mysqlOrderRepository) InsertOrderItem(row model.OrderItem) (*model.OrderItem, error) {
	err := m.db.Create(&row).Error
	return &row, err
}

func (m mysqlOrderRepository) UpdateOrderById(id string, values interface{}) error {
	return m.db.Model(model.Order{}).Where("id = ?", id).Updates(values).Error
}
