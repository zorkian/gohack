/* state.go - tracks the status of things, particularly the map/world
 *
 * Copyright (c) 2013 by Mark Smith <mark@qq.is>.
 *
 * See the included LICENSE file for information. 
 */

package main

type LocationType int
type LocationSubtype int

const (
	// Some constants we use everywhere to define the playing field.
	WWIDTH  int = 76
	WHEIGHT int = 19

	// Types. A square can be one of these structural types.
	T_Floor LocationType = iota
	T_Tunnel
	T_Water
	T_Wall
	T_Rock

	// FLOOR subtypes.
	ST_Dirt LocationSubtype = iota
	ST_Fountain
	ST_Sink
	ST_Grave
	ST_Altar
	ST_StairsUp
	ST_StairsDown
	ST_Pit
	ST_Hole

	// TUNNEL subtypes.

	// WALL subtypes.
	ST_Door

	// DOOR subtypes.

	// WATER subtypes.
	ST_River LocationSubtype = iota
	ST_Moat
	ST_Lake

	// ROCK subtypes.
)

// Location is a single spot in the world. Each location has a type (what it is)
// but also can contain a monster (player/other) as well as items.
type LocXY struct {
	X, Y int
}
type Location struct {
	Type     LocationType
	Subtype  LocationSubtype
	Lit      bool
	Open     bool
	Locked   bool
	Diggable bool
	Redraw   bool
	Render   rune // Used for wall characters mostly.
	Items    []*Item
	Occupant *Mobile
}

// Level represents a layer of the dungeon. This contains a bunch of locations,
// which can contain many other things of course. It's an array of scanlines.
type Level [WHEIGHT][WWIDTH]Location

// TODO: This should be more than the present level.
var World *Level
var PC *Mobile
var Ticks int = 0
