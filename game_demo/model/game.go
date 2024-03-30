package model

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"ziyi.game.com/config"
)

type Input struct {
}

type Game struct {
	input  *Input
	ship   *Ship
	config *config.Config
}

func NewGame() *Game {
	c := config.LoadConfig()
	//set window size & title
	ebiten.SetWindowSize(c.ScreenWidth, c.ScreenHeight)
	ebiten.SetWindowTitle(c.Title)
	return &Game{
		input:  &Input{},
		ship:   NewShip(),
		config: c,
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "hello world")
	//set screen color
	screen.Fill(g.config.BgColor)
	//draw ship
	g.ship.Draw(screen, g.config)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.config.ScreenWidth, g.config.ScreenHeight
}

func (g *Game) Update() error {
	return nil
}
