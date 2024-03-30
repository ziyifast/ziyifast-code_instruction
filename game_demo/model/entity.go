package model

type Entity interface {
	Width() int
	Height() int
	X() int
	Y() int
}
