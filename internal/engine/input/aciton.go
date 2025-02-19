package input

type Action uint8

const (
	MoveUp Action = iota
	MoveDown
	MoveLeft
	MoveRight
	Attack
	UseAbility1
	UseAbility2
	UseAbility3
	Pause
	ToggleDebug
)

type ActionState struct {
	Active       bool
	JustPressed  bool
	JustReleased bool
}

func (a *Action) String() string {
	return [...]string{"MoveUp", "MoveDown", "MoveLeft", "MoveRight", "Attack", "UseAbility1", "UseAbility2", "UseAbility3", "Pause", "Toggle Debug"}[*a]
}
