package model

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/ziyifast/log"
	_ "golang.org/x/image/bmp"
	"ziyi.game.com/config"
)

type Ship struct {
	image  *ebiten.Image
	width  int
	height int
}

func NewShip() *Ship {
	image, _, err := ebitenutil.NewImageFromFile("/Users/ziyi2/GolandProjects/MyTest/demo_home/game_demo/images/ship.bmp")
	if err != nil {
		log.Fatalf("%v", err)
	}
	width, height := image.Bounds().Dx(), image.Bounds().Dy()
	return &Ship{
		image:  image,
		width:  width,
		height: height,
	}
}

func (ship *Ship) Draw(screen *ebiten.Image, cfg *config.Config) {
	// draw by self
	op := &ebiten.DrawImageOptions{}
	//init ship at the screen center
	op.GeoM.Translate(float64(cfg.ScreenWidth-ship.width)/2, float64(cfg.ScreenHeight-ship.height))
	screen.DrawImage(ship.image, op)
}
