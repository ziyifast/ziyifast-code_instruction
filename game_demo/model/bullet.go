package model

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"ziyi.game.com/config"
)

type Bullet struct {
	image       *ebiten.Image
	width       int
	height      int
	x           int
	y           int
	speedFactor int
}

// NewBullet 添加子弹
func NewBullet(cfg *config.Config, ship *Ship) *Bullet {
	rect := image.Rect(0, 0, cfg.BulletWidth, cfg.BulletHeight)
	img := ebiten.NewImageWithOptions(rect, nil)
	img.Fill(cfg.BulletColor)
	return &Bullet{
		image:       img,
		width:       cfg.BulletWidth,
		height:      cfg.BulletHeight,
		x:           ship.x + (ship.width-cfg.BulletWidth)/2,
		y:           cfg.ScreenHeight - ship.height - cfg.BulletHeight,
		speedFactor: cfg.BulletSpeed,
	}
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.x), float64(b.y))
	screen.DrawImage(b.image, op)
}

func (b *Bullet) outOfScreen() bool {
	return b.y < -b.height
}
