package model

type Player struct {
	Id   int64  `xorm:"id" json:"id"`
	Name string `xorm:"name" json:"name"`
	Age  int    `xorm:"age" json:"age"`
}

func (p *Player) TableName() string {
	return "player"
}
