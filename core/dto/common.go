package dto

import (
	"github.com/mqnoy/logistics-app/core/model"
	"github.com/mqnoy/logistics-app/core/util"
	"gorm.io/datatypes"
)

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
	TotalItems int64 `json:"total_items"`
	Offset     int   `json:",omitempty"`
}

type ListResponse[T any] struct {
	Rows     []T        `json:"rows"`
	MetaData Pagination `json:"metadata"`
}

type SelectAndCount[M any] struct {
	Rows  []*M
	Count int64
}

type DetailParam struct {
	ID      string
	Session AuthorizedUser
}

// list param
type ListParam[T any] struct {
	Filters    T
	Orders     string
	Pagination Pagination
	Session    AuthorizedUser
}

type CreateParam[T any] struct {
	CreateValue T
	Session     AuthorizedUser
}

// update param
type UpdateParam[T any] struct {
	ID          string
	UpdateValue T
	Session     AuthorizedUser
}

type Timestamp struct {
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

type FilterCommonParams struct {
	Keyword  string
	Name     string
	IsActive *bool
	MemberId string
	IsDone   *bool
}

func ComposeTimestamp(m model.TimestampColumn) Timestamp {
	return Timestamp{
		CreatedAt: util.DateToEpoch(m.CreatedAt),
		UpdatedAt: util.DateToEpoch(m.UpdatedAt),
	}
}

type EntitySnapshot struct {
	Snapshot datatypes.JSON
}
