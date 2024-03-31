package model

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"ziyi.game.com/config"
)

type Bullet struct {
	GameObj
	image       *ebiten.Image
	speedFactor int
}

// NewBullet 添加子弹
func NewBullet(cfg *config.Config, ship *Ship) *Bullet {
	rect := image.Rect(0, 0, cfg.BulletWidth, cfg.BulletHeight)
	img := ebiten.NewImageWithOptions(rect, nil)
	img.Fill(cfg.BulletColor)
	b := &Bullet{
		image:       img,
		speedFactor: cfg.BulletSpeed,
	}
	b.GameObj.width = cfg.BulletWidth
	b.GameObj.height = cfg.BulletHeight
	b.GameObj.y = ship.Y() + (ship.Height()-cfg.BulletHeight)/2
	b.GameObj.x = ship.X() + (ship.Width()-cfg.BulletWidth)/2
	return b
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.X()), float64(b.Y()))
	screen.DrawImage(b.image, op)
}

func (b *Bullet) outOfScreen() bool {
	return b.Y() < -b.Height()
}
