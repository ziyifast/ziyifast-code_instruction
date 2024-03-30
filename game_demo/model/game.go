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
	//检查是否
	g.CheckCollision()
	return nil
}

func (g *Game) addBullet(bullet *Bullet) {
	g.bullets[bullet] = struct{}{}
}

func (g *Game) createAliens() {
	a := NewAlien(g.config)
	//外星人之间需要有间隔
	availableSpaceX := g.config.ScreenWidth - 2*a.Width()
	numAliens := availableSpaceX / (2 * a.Width())
	for i := 0; i < numAliens; i++ {
		alien := NewAlien(g.config)
		alien.x = alien.Width() + 2*alien.Width()*i
		alien.y = alien.Height()
		g.addAliens(alien)
	}
}

func (g *Game) addAliens(alien *Alien) {
	g.aliens[alien] = struct{}{}
}

func (g *Game) CheckCollision() {
	for alien := range g.aliens {
		for bullet := range g.bullets {
			if checkCollision(bullet, alien) {
				delete(g.aliens, alien)
				delete(g.bullets, bullet)
			}
		}
	}
}

// 检测子弹是否击中敌人
func checkCollision(entity1 Entity, entity2 Entity) bool {
	//只需要计算子弹顶点在敌人矩形之中，就认为击中敌人
	entity2Top := entity2.Y()
	entity2Left := entity2.X()
	entity2Bottom := entity2.Y() + entity2.Height()
	entity2Right := entity2.X() + entity2.Width()
	x, y := entity1.X(), entity1.Y()
	//击中敌人左上角
	if x > entity2Left && x < entity2Right && y > entity2Top && y < entity2Bottom {
		return true
	}
	//击中敌人右上角
	x, y = entity1.X(), entity1.Y()+entity1.Height()
	if x > entity2Left && x < entity2Right && y > entity2Bottom && y < entity2Top {
		return true
	}
	//左下角
	x, y = entity1.X()+entity1.Width(), entity1.Y()
	if y > entity2Top && y < entity2Bottom && x > entity2Left && x < entity2Right {
		return true
	}
	//右下角
	x, y = entity1.X()+entity1.Width(), entity1.Y()+entity1.Height()
	if y > entity2Top && y < entity2Bottom && x > entity2Left && x < entity2Right {
		return true
	}
	return false
}
