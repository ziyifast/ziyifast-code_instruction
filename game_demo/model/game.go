package model

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/ziyifast/log"
	"ziyi.game.com/config"
)

type Game struct {
	input   *Input
	ship    *Ship
	config  *config.Config
	bullets map[*Bullet]struct{}
	aliens  map[*Alien]struct{}
}

func NewGame() *Game {
	c := config.LoadConfig()
	//set window size & title
	ebiten.SetWindowSize(c.ScreenWidth, c.ScreenHeight)
	ebiten.SetWindowTitle(c.Title)
	g := &Game{
		input:   &Input{},
		ship:    NewShip(c.ScreenWidth, c.ScreenHeight),
		config:  c,
		bullets: make(map[*Bullet]struct{}),
		aliens:  map[*Alien]struct{}{},
	}
	//初始化外星人
	g.createAliens()
	return g
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
	//draw aliens
	for a := range g.aliens {
		a.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.config.ScreenWidth, g.config.ScreenHeight
}

func (g *Game) Update() error {
	log.Infof("update....")
	g.input.Update(g)
	//更新子弹位置
	for b := range g.bullets {
		if b.outOfScreen() {
			delete(g.bullets, b)
		}
		b.y -= b.speedFactor
	}
	//更新敌人位置
	for a := range g.aliens {
		a.y += a.speedFactor
	}
	return nil
}

func (g *Game) addBullet(bullet *Bullet) {
	g.bullets[bullet] = struct{}{}
}

func (g *Game) createAliens() {
	a := NewAlien(g.config)
	//外星人之间需要有间隔
	availableSpaceX := g.config.ScreenWidth - 2*a.width
	numAliens := availableSpaceX / (2 * a.width)
	for i := 0; i < numAliens; i++ {
		alien := NewAlien(g.config)
		alien.x = alien.width + 2*alien.width*i
		alien.y = alien.height
		g.addAliens(alien)
	}
}

func (g *Game) addAliens(alien *Alien) {
	g.aliens[alien] = struct{}{}
}
