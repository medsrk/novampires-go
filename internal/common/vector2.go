package common

import (
	"fmt"
	"math"
)

type Vector2 struct {
	X, Y float64
}

func (v Vector2) String() string {
	return fmt.Sprintf("(%.2f, %.2f)", v.X, v.Y)
}

func (v Vector2) Add(v2 Vector2) Vector2 {
	return Vector2{v.X + v2.X, v.Y + v2.Y}
}

func (v Vector2) Sub(v2 Vector2) Vector2 {
	return Vector2{v.X - v2.X, v.Y - v2.Y}
}

func (v Vector2) Scale(s float64) Vector2 {
	return Vector2{v.X * s, v.Y * s}
}

func (v Vector2) Div(s float64) Vector2 {
	return Vector2{v.X / s, v.Y / s}
}

func (v Vector2) Length() float64 {
	return v.LengthSquared() * 0.5
}

func (v Vector2) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v Vector2) Normalize() Vector2 {
	return v.Div(v.Length())
}

func (v Vector2) Normalized() Vector2 {
	mag := v.Magnitude()
	if mag == 0 {
		return Vector2{}
	}
	return Vector2{X: v.X / mag, Y: v.Y / mag}
}

func (v Vector2) Magnitude() float64 {
	return math.Sqrt(v.MagnitudeSquared())
}

func (v Vector2) MagnitudeSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v Vector2) Dot(v2 Vector2) float64 {
	return v.X*v2.X + v.Y*v2.Y
}

func (v Vector2) Lerp(v2 Vector2, t float64) Vector2 {
	return v.Add(v2.Sub(v).Scale(t))
}

func (v Vector2) Distance(v2 Vector2) float64 {
	return v.Sub(v2).Length()
}

func (v Vector2) DistanceSquared(v2 Vector2) float64 {
	return v.Sub(v2).LengthSquared()
}

func (v Vector2) Angle() float64 {
	return math.Atan2(v.Y, v.X)
}

func (v Vector2) AngleBetween(v2 Vector2) float64 {
	return math.Acos(v.Dot(v2) / (v.Length() * v2.Length()))
}

func (v Vector2) Rotate(angle float64) Vector2 {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	return Vector2{
		X: v.X*cos - v.Y*sin,
		Y: v.X*sin + v.Y*cos,
	}
}

func (v Vector2) Reflect(normal Vector2) Vector2 {
	return v.Sub(normal.Scale(2 * v.Dot(normal)))
}
