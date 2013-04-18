/* mobile.go - basic logic for the PC and NPCs
 *
 * Copyright (c) 2013 by Mark Smith <mark@qq.is>.
 *
 * See the included LICENSE file for information. 
 */

package main

import (
	"github.com/nsf/termbox-go"
	"math/rand"
)

type SightType int
type BreatheType int
type AlignmentType int

const (
	// Sight constants.
	S_Eyes SightType = iota
	S_Telepathy
	S_Omniscient
	S_Blind

	// Breathing constants.
	B_Air BreatheType = iota
	B_Water
	B_Any
	B_None

	// Alignment constants.
	A_Chaotic AlignmentType = iota
	A_Neutral
	A_Lawful
)

// Attributes always go in a package, so we put them into a convenient struct
// to save us time and typing later.
type Attr struct {
	Str, Con, Dex, Wis, Int, Cha int
}

// Species contains basic information about what something is. This is kind of
// the meta information. I.e., it's easy for us to say that X applies to a
// species of things. Horses are a species, dwarves are a species.
type SAttr struct {
	Base  Attr // Start with Base ...
	Bonus Attr // ... and add 1-Bonus
	Max   Attr // But never over Max.
}
type Species struct {
	P          Behavior
	Name       string
	Render     rune
	Color      termbox.Attribute
	Alignment  AlignmentType
	MinLevel   int // Never generate at XPLvl < this.
	HitDice    int // 1dHitDice per level HP
	HitBonus   int // + HitBonus (per level)
	Attributes SAttr
}

// Prototype defines behaviors for something. I.e., does it have eyes, what
// kind of attacks does it have, etc. 
type Prototype struct {
	Sight    SightType
	Breathes BreatheType
	NumHands int
	Sessile  bool
	Speed    int
}

// Behavior should be implemented by everything that can be a mobile. This
// means the PC and all NPCs should adhere to this and implement these. This is
// how we determine what can be done.
type Behavior interface {
	Hands() int
	HasHands() bool
	BaseSpeed() int
}

// Mobile is an instance of a creature.
type Mobile struct {
	S         *Species
	X, Y      int
	HP, HPmax int
	Level     int
	XP        int
	Speed     int
}

// Hands returns how many hands this creature has.
func (self *Prototype) Hands() int {
	return self.NumHands
}

// HasHands returns whether or not this creature has any hands.
func (self *Prototype) HasHands() bool {
	return self.NumHands > 0
}

// BaseSpeed returns the raw, unmodified speed of this creature. Speed is
// how many ticks pass when they do something that takes a full turn.
func (self *Prototype) BaseSpeed() int {
	return self.Speed
}

// NewMobile builds a new mobile of a certain species. It starts at level 1. If
// you want to level it, you can do that with another function.
func newMobile(spec *Species) *Mobile {
	self := &Mobile{S: spec}
	self.HPmax = rand.Intn(spec.HitDice) + spec.HitBonus
	self.HP = self.HPmax
	self.Level = 1
	self.XP = 0
	self.Speed = self.S.P.BaseSpeed()
	self.X = -1
	self.Y = -1
	return self
}
