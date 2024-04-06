package config

import (
	"encoding/json"
	"github.com/ziyifast/log"
	"image/color"
	"os"
)

type Config struct {
	ScreenWidth        int        `json:"screen_width"`
	ScreenHeight       int        `json:"screen_height"`
	Title              string     `json:"title"`
	BgColor            color.RGBA `json:"bg_color"`
	MoveSpeed          int        `json:"move_speed"`
	BulletWidth        int        `json:"bullet_width"`
	BulletHeight       int        `json:"bullet_height"`
	BulletSpeed        int        `json:"bullet_speed"`
	BulletColor        color.RGBA `json:"bullet_color"`
	MaxBulletNum       int        `json:"max_bullet_num"`  //页面中最多子弹数量
	BulletInterval     int64      `json:"bullet_interval"` //发射子弹间隔时间
	MonsterSpeedFactor int        `json:"monster_speed_factor"`
	TitleFontSize      int        `json:"title_font_size"`
	FontSize           int        `json:"font_size"`
	SmallFontSize      int        `json:"small_font_size"`
	FailedCountLimit   int        `json:"failed_count_limit"` //最多能遗漏多少怪物
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
