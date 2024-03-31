package model

import (
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

type Input struct {
	lastBulletTime time.Time //上次子弹发射时间，避免用户一直按着连续发子弹
}

func (i *Input) IsKeyPressed() bool {
	//按下空格或者鼠标左键，游戏开始
	if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		return true
	}
	return false
}

func (i *Input) Update(g *Game) {
	cfg := g.config
	s := g.ship
	//listen the key event
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		s.GameObj.x -= cfg.MoveSpeed
		//防止飞船跑出页面 prevents movement out of the page
		if s.X() < -s.Width()/2 {
			s.x = -s.Width() / 2
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		s.GameObj.x += cfg.MoveSpeed
		if s.X() > cfg.ScreenWidth-s.Width()/2 {
			s.GameObj.x = cfg.ScreenWidth - s.Width()/2
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
