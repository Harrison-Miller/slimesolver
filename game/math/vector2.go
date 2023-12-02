package math

import "fmt"

var NegVec = Vector2{-1, -1}

type Vector2 struct {
	X, Y int
}

func (v Vector2) String() string {
	return fmt.Sprintf("(%d, %d)", v.X, v.Y)
}

func (v Vector2) Equals(other Vector2) bool {
	return v.X == other.X && v.Y == other.Y
}
