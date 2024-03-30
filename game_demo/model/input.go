package model

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Input struct {
}

func (i *Input) Update(g *Game) {
	cfg := g.config
	s := g.ship
	//listen the key event
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		s.x -= cfg.MoveSpeed
		//防止飞船跑出页面 prevents movement out of the page
		if s.x < -s.width/2 {
			s.x = -s.width / 2
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		s.x += cfg.MoveSpeed
		if s.x > cfg.ScreenWidth-s.width/2 {
			s.x = cfg.ScreenWidth - s.width/2
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		//发射子弹
		bullet := NewBullet(cfg, s)
		g.addBullet(bullet)
	}
}
