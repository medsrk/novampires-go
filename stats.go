package main

import "fmt"

type Stat interface {
	Value() float32
	AddBonus(bonus float32)
	RemoveBonus(bonus float32)
	SetBonus(bonus float32)
	ResetBonus()
	SetBase(base float32)
	String() string
}

type Speed struct {
	base  float32
	bonus float32
}

func NewSpeed(base, bonus float32) Speed {
	return Speed{
		base:  base,
		bonus: bonus,
	}
}

func (s *Speed) Value() float32 {
	return s.base * (1 + (s.bonus / 100))
}

func (s *Speed) AddBonus(bonus float32) {
	s.bonus += bonus
}

func (s *Speed) RemoveBonus(bonus float32) {
	s.bonus -= bonus
}

func (s *Speed) SetBonus(bonus float32) {
	s.bonus = bonus
}

func (s *Speed) ResetBonus() {
	s.bonus = 0
}

func (s *Speed) SetBase(base float32) {
	s.base = base
}

func (s *Speed) String() string {
	return fmt.Sprint(s.Value() / 100)
}
