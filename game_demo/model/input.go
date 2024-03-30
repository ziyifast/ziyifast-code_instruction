package model

import (
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

type Input struct {
	lastBulletTime time.Time //上次子弹发射时间，避免用户一直按着连续发子弹
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
		if len(g.bullets) < cfg.MaxBulletNum && time.Since(i.lastBulletTime).Milliseconds() > cfg.BulletInterval {
			//发射子弹
			bullet := NewBullet(cfg, s)
			g.addBullet(bullet)
			i.lastBulletTime = time.Now()
		}
	}
}
