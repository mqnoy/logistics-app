package model

import "gorm.io/datatypes"

type GoodSnapShotColumn struct {
	GoodSnapShot datatypes.JSON `gorm:"column:good_snapshot;default:'{}'"`
}

func (m *GoodSnapShotColumn) ParseGoodSnapshot() (b []byte, err error) {
	b, err = m.GoodSnapShot.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return b, nil
}
