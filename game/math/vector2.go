package math

import "fmt"

type Vector2 struct {
	X, Y int
}

func (v Vector2) String() string {
	return fmt.Sprintf("(%d, %d)", v.X, v.Y)
}
