package model

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/ziyifast/log"
	"ziyi.game.com/config"
)

type Alien struct {
	GameObj
	img         *ebiten.Image
	speedFactor int
}

func NewAlien(cfg *config.Config) *Alien {
	image, _, err := ebitenutil.NewImageFromFile("/Users/ziyi2/GolandProjects/MyTest/demo_home/game_demo/images/alien.bmp")
	if err != nil {
		log.Fatal("%v", err)
	}
	width, height := image.Bounds().Dx(), image.Bounds().Dy()
	a := &Alien{
		img:         image,
		speedFactor: cfg.AlienSpeedFactor,
	}
	a.GameObj.width = width
	a.GameObj.height = height
	a.GameObj.x = 0
	a.GameObj.y = 0
	return a
}

func (a *Alien) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(a.X()), float64(a.Y()))
	screen.DrawImage(a.img, op)
}

func (a *Alien) OutOfScreen(cfg *config.Config) bool {
	if a.Y() > cfg.ScreenHeight {
		return true
	}
	return false
}
