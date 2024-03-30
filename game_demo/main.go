package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/ziyifast/log"
	"ziyi.game.com/model"
)

type Game struct {
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "hi~")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 300, 240
}

func main() {
	err := ebiten.RunGame(model.NewGame())
	if err != nil {
		log.Fatal("%v", err)
	}
}
