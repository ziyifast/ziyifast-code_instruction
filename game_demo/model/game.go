package model

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"ziyi.game.com/config"
)

type Game struct {
	input   *Input
	ship    *Ship
	config  *config.Config
	bullets map[*Bullet]struct{}
}

func NewGame() *Game {
	c := config.LoadConfig()
	//set window size & title
	ebiten.SetWindowSize(c.ScreenWidth, c.ScreenHeight)
	ebiten.SetWindowTitle(c.Title)
	return &Game{
		input:   &Input{},
		ship:    NewShip(c.ScreenWidth, c.ScreenHeight),
		config:  c,
		bullets: make(map[*Bullet]struct{}),
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "hello world")
	//set screen color
	screen.Fill(g.config.BgColor)
	//draw ship
	g.ship.Draw(screen, g.config)
	//draw bullet
	for b := range g.bullets {
		b.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.config.ScreenWidth, g.config.ScreenHeight
}

func (g *Game) Update() error {
	g.input.Update(g)
	//更新子弹位置
	for b := range g.bullets {
		b.y -= b.speedFactor
	}
	return nil
}

func (g *Game) addBullet(bullet *Bullet) {
	g.bullets[bullet] = struct{}{}
}
