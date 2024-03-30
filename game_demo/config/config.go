package config

import (
	"encoding/json"
	"github.com/ziyifast/log"
	"image/color"
	"os"
)

type Config struct {
	ScreenWidth  int        `json:"screen_width"`
	ScreenHeight int        `json:"screen_height"`
	Title        string     `json:"title"`
	BgColor      color.RGBA `json:"bg_color"`
	MoveSpeed    int        `json:"move_speed"`
}

func LoadConfig() *Config {
	file, err := os.Open("./config.json")
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer file.Close()
	config := new(Config)
	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return config
}
