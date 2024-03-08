package service

import (
	"github.com/ziyifast/log"
	"myTest/demo_home/blond_filter/dao"
	"myTest/demo_home/blond_filter/model"
	"myTest/demo_home/blond_filter/util"
)

type playerService struct {
}

var PlayerService = new(playerService)

func (s *playerService) FindById(id int64) (*model.Player, error) {
	// query blond filter
	if !util.CheckExist(id) {
		log.Infof("the player does not exist in the blond filter,return it!!! ")
		return nil, nil
	}

	//query redis
	player, err := util.PlayerCache.GetById(id)
	if err != nil {
		return nil, err
	}
	if player != nil {
		return player, nil
	}
	//query db and cache result
	p, err := dao.PlayerDao.GetById(id)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	if p != nil {
		err = util.PlayerCache.Put(p)
		if err != nil {
			log.Errorf("%v", err)
		}
		return p, nil
	}
	return p, nil
}
