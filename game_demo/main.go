package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ziyifast/log"
	"ziyi.game.com/model"
)

func main() {
	err := ebiten.RunGame(model.NewGame())
	if err != nil {
		log.Fatal("%v", err)
	}
}
