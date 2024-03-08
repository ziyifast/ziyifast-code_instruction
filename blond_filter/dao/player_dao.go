package dao

import (
	"github.com/aobco/log"
	"myTest/demo_home/blond_filter/model"
	"myTest/demo_home/blond_filter/pg"
	"time"
)

type playerDao struct {
}

var PlayerDao = new(playerDao)

func (p *playerDao) InsertOne(player model.Player) (int64, error) {
	return pg.Engine.InsertOne(player)
}

func (p *playerDao) GetById(id int64) (*model.Player, error) {
	log.Infof("query postgres,time:%v", time.Now().String())
	player := new(model.Player)
	get, err := pg.Engine.Where("id=?", id).Get(player)
	if err != nil {
		log.Errorf("%v", err)
	}
	if !get {
		return nil, nil
	}
	return player, nil
}
