package model

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/ziyifast/log"
	"ziyi.game.com/config"
)

type Monster struct {
	GameObj
	img         *ebiten.Image
	speedFactor int
}

func NewMonster(cfg *config.Config) *Monster {
	image, _, err := ebitenutil.NewImageFromFile("/Users/ziyi2/GolandProjects/MyTest/demo_home/game_demo/images/monster.bmp")
	if err != nil {
		log.Fatal("%v", err)
	}
	width, height := image.Bounds().Dx(), image.Bounds().Dy()
	a := &Monster{
		img:         image,
		speedFactor: cfg.MonsterSpeedFactor,
	}
	a.GameObj.width = width
	a.GameObj.height = height
	a.GameObj.x = 0
	a.GameObj.y = 0
	return a
}

func (a *Monster) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(a.X()), float64(a.Y()))
	screen.DrawImage(a.img, op)
}

func (a *Monster) OutOfScreen(cfg *config.Config) bool {
	if a.Y()+a.Height() > cfg.ScreenHeight {
		return true
	}
	return false
}
