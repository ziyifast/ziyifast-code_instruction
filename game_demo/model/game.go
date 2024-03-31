package model

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/ziyifast/log"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"ziyi.game.com/config"
)

type Mode int

const (
	ModeTitle Mode = iota
	ModeGame
	ModeOver
)

type Game struct {
	input       *Input
	ship        *Ship
	config      *config.Config
	bullets     map[*Bullet]struct{}
	aliens      map[*Alien]struct{}
	mode        Mode
	failedCount int //记录失败次数（未击中的次数）
}

func (g *Game) init() {
	fmt.Println("恢复初始状态...")
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
	g.CreateFonts()
	return g
}

func (g *Game) Draw(screen *ebiten.Image) {
	var titleTexts []string
	var texts []string
	switch g.mode {
	case ModeTitle:
		titleTexts = []string{"ALIEN INVASION"}
		texts = []string{"", "", "", "", "", "", "", "PRESS SPACE KEY", "", "OR LEFT MOUSE"}
	case ModeGame:
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
	case ModeOver:
		texts = []string{"", "GAME OVER!"}
	}
	for i, l := range titleTexts {
		x := (g.config.ScreenWidth - len(l)*g.config.TitleFontSize) / 2
		text.Draw(screen, l, titleArcadeFont, x, (i+4)*g.config.TitleFontSize, color.White)
	}
	for i, l := range texts {
		x := (g.config.ScreenWidth - len(l)*g.config.FontSize) / 2
		text.Draw(screen, l, arcadeFont, x, (i+4)*g.config.FontSize, color.White)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.config.ScreenWidth, g.config.ScreenHeight
}

func (g *Game) Update() error {
	switch g.mode {
	case ModeTitle:
		if g.input.IsKeyPressed() {
			g.mode = ModeGame
		}
	case ModeGame:

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
		//检查是否击相撞（击中敌人）
		g.CheckKillAlien()
		//检查是否飞机碰到外星人

	case ModeOver:
		//游戏结束，恢复初始状态
		if g.input.IsKeyPressed() {
			g.init()
			g.mode = ModeTitle
		}
	}
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

func (g *Game) CheckKillAlien() {
	for alien := range g.aliens {
		for bullet := range g.bullets {
			if checkCollision(bullet, alien) {
				delete(g.aliens, alien)
				delete(g.bullets, bullet)
			}
		}
	}
}

func (g *Game) CheckShipCrashed() {
	for alien := range g.aliens {
		if checkCollision(g.ship, alien) {
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

//加载页面字体

var (
	titleArcadeFont font.Face
	arcadeFont      font.Face
	smallArcadeFont font.Face
)

// CreateFonts 初始化页面字体信息
func (g *Game) CreateFonts() {
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatalf("%v", err)
	}
	const dpi = 72
	titleArcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(g.config.TitleFontSize),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	arcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(g.config.FontSize),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	smallArcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(g.config.SmallFontSize),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}
