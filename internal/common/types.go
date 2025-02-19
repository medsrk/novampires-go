package common

type TargetInfo struct {
	ID       uint64
	Pos      Vector2
	Vel      Vector2
	Radius   float64
	Priority float64 // higher is more important
}

type Rectangle struct {
	Pos  Vector2
	Size Vector2
}

func (r Rectangle) Contains(p Vector2) bool {
	return p.X >= r.Pos.X &&
		p.X <= r.Pos.X+r.Size.X &&
		p.Y >= r.Pos.Y &&
		p.Y <= r.Pos.Y+r.Size.Y
}

func (r Rectangle) Intersects(other Rectangle) bool {
	return r.Pos.X < other.Pos.X+other.Size.X &&
		r.Pos.X+r.Size.X > other.Pos.X &&
		r.Pos.Y < other.Pos.Y+other.Size.Y &&
		r.Pos.Y+r.Size.Y > other.Pos.Y
}

func (r Rectangle) Center() Vector2 {
	return Vector2{
		X: r.Pos.X + r.Size.X/2,
		Y: r.Pos.Y + r.Size.Y/2,
	}
}
