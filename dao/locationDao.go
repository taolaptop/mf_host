package dao

import "github.com/user/entity"

type Location struct {
}

type ILocation interface {
	Save(location entity.LocationSnapshot)
	Get(mobile string) entity.LocationSnapshot
}

func (l Location) Save(location entity.LocationSnapshot) {
	effect, _ := db.Where("id=?", location.Id).Update(location)
	if effect == 0 {
		db.Insert(location)
	}
}

func (l Location) Get(mobile string) entity.LocationSnapshot {
	return entity.LocationSnapshot{}
}
