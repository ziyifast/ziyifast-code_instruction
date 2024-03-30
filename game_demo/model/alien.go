package model

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/ziyifast/log"
	"ziyi.game.com/config"
)

type Alien struct {
	img         *ebiten.Image
	width       int
	height      int
	x           int
	y           int
	speedFactor int
}

func NewAlien(cfg *config.Config) *Alien {
	image, _, err := ebitenutil.NewImageFromFile("/Users/ziyi2/GolandProjects/MyTest/demo_home/game_demo/images/alien.bmp")
	if err != nil {
		log.Fatal("%v", err)
	}
	width, height := image.Bounds().Dx(), image.Bounds().Dy()
	return &Alien{
		img:         image,
		width:       width,
		height:      height,
		x:           0,
		y:           0,
		speedFactor: cfg.AlienSpeedFactor,
	}
}

func (a *Alien) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(a.x), float64(a.y))
	screen.DrawImage(a.img, op)
}
