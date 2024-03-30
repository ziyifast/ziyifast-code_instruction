package model

// GameObj 后续除了普通敌人还可能有其他小boss，因此我们直接将所有物体抽象出来
type GameObj struct {
	width  int
	height int
	x      int
	y      int
}

func (o *GameObj) Width() int {
	return o.width
}

func (o *GameObj) Height() int {
	return o.height
}

func (o *GameObj) X() int {
	return o.x
}

func (o *GameObj) Y() int {
	return o.y
}
