package model

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/ziyifast/log"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"math/rand"
	"time"
	"ziyi.game.com/config"
)

type Mode int

const (
	ModeTitle Mode = iota
	ModeGame
	ModeOver
)

var r *rand.Rand

func init() {
	source := rand.NewSource(time.Now().UnixMicro())
	r = rand.New(source)
}

type Game struct {
	input            *Input
	ship             *Ship
	config           *config.Config
	bullets          map[*Bullet]struct{}
	monsters         map[*Monster]struct{}
	mode             Mode
	failedCountLimit int
	failedCount      int
}

func (g *Game) init() {
	g.mode = ModeTitle
	g.failedCount = 0
	g.bullets = make(map[*Bullet]struct{})
	g.monsters = make(map[*Monster]struct{})
	g.ship = NewShip(g.config.ScreenWidth, g.config.ScreenHeight)
	g.createMonsters()
}

func NewGame() *Game {
	c := config.LoadConfig()
	//set window size & title
	ebiten.SetWindowSize(c.ScreenWidth, c.ScreenHeight)
	ebiten.SetWindowTitle(c.Title)
	g := &Game{
		input:            &Input{},
		ship:             NewShip(c.ScreenWidth, c.ScreenHeight),
		config:           c,
		bullets:          make(map[*Bullet]struct{}),
		monsters:         make(map[*Monster]struct{}),
		failedCount:      0,
		failedCountLimit: c.FailedCountLimit,
	}
	//初始化外星人
	g.createMonsters()
	g.CreateFonts()
	return g
}

func (g *Game) Draw(screen *ebiten.Image) {
	var titleTexts []string
	var texts []string
	switch g.mode {
	case ModeTitle:
		titleTexts = []string{"RUN GOPHER"}
		texts = []string{"", "", "", "", "", "", "", "PRESS SPACE KEY", "", "OR LEFT MOUSE"}
	case ModeGame:
		//set screen color
		screen.Fill(g.config.BgColor)
		//draw gopher
		g.ship.Draw(screen, g.config)
		//draw bullet
		for b := range g.bullets {
			b.Draw(screen)
		}
		//draw monsters
		for a := range g.monsters {
			a.Draw(screen)
		}
	case ModeOver:
		screen.Fill(color.Black)
		g.Update()
		texts = []string{"", "GAME OVER!"}
	}
	for i, l := range titleTexts {
		x := (g.config.ScreenWidth - len(l)*g.config.TitleFontSize) / 2
		text.Draw(screen, l, titleArcadeFont, x, (i+4)*g.config.TitleFontSize, color.RGBA{
			R: 0,
			G: 100,
			B: 0,
			A: 0,
		})
	}
	for i, l := range texts {
		x := (g.config.ScreenWidth - len(l)*g.config.FontSize) / 2
		text.Draw(screen, l, arcadeFont, x, (i+4)*g.config.FontSize, color.RGBA{
			R: 0,
			G: 100,
			B: 0,
			A: 0,
		})
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
		g.input.Update(g)
		//更新子弹位置
		for b := range g.bullets {
			if b.outOfScreen() {
				delete(g.bullets, b)
			}
			b.y -= b.speedFactor
		}
		//更新敌人位置
		for a := range g.monsters {
			a.y += a.speedFactor
		}
		//检查是否击相撞（击中敌人）
		g.CheckKillMonster()
		//外星人溜走 或者 是否飞机碰到外星人
		if g.failedCount >= g.failedCountLimit || g.CheckShipCrashed() {
			g.mode = ModeOver
			log.Warnf("over..........")
		}
		go func() {
			if len(g.monsters) < 0 {
				//下一波怪物
				g.createMonsters()
			}
		}()
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

func (g *Game) createMonsters() {
	a := NewMonster(g.config)
	//怪物之间需要有间隔
	availableSpaceX := g.config.ScreenWidth - 2*a.Width()
	numMonsters := availableSpaceX / (2 * a.Width())
	//预设怪物数量
	for i := 0; i < numMonsters; i++ {
		monster := NewMonster(g.config)
		monster.x = monster.Width() + 2*monster.Width()*i
		monster.y = monster.Height() + r.Intn(g.config.ScreenHeight/10)
		g.addMonsters(monster)
	}
}

func (g *Game) addMonsters(monster *Monster) {
	g.monsters[monster] = struct{}{}
}

func (g *Game) CheckKillMonster() {
	for monster := range g.monsters {
		for bullet := range g.bullets {
			if checkCollision(bullet, monster) {
				delete(g.monsters, monster)
				delete(g.bullets, bullet)
			}
		}
		if monster.OutOfScreen(g.config) {
			g.failedCount++
			delete(g.monsters, monster)
		}
	}
}

func (g *Game) CheckShipCrashed() bool {
	for monster := range g.monsters {
		if checkCollision(g.ship, monster) {
			return true
		}
	}
	return false
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
