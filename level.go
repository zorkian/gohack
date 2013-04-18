/* level.go - generates levels
 *
 * Copyright (c) 2013 by Mark Smith <mark@qq.is>.
 *
 * See the included LICENSE file for information. 
 */

package main

import (
	"math/rand"
)

const (
	ROOM_DENSITY   int = 5 // Higher = denser, but slower to generate.
	TUNNEL_DENSITY int = 5 // Higher = denser.
)

type Room struct {
	X, Y, W, H int
	Connected  bool
}

// randomRoom returns a randomly sized and configured room.
func randomRoom() *Room {
	w, h := 4+rand.Intn(12), 4+rand.Intn(4)
	return &Room{
		W: w, H: h,
		X: rand.Intn(WWIDTH - w),
		Y: rand.Intn(WHEIGHT - h),
	}
}

func (self *Room) Draw(level *Level) {
	for y := self.Y; y <= self.Y+self.H; y++ {
		for x := self.X; x <= self.X+self.W; x++ {
			level[y][x].Type = T_Wall

			if y == self.Y && x == self.X {
				level[y][x].Render = '┌'
			} else if y == self.Y && x == self.X+self.W {
				level[y][x].Render = '┐'
			} else if y == self.Y+self.H && x == self.X {
				level[y][x].Render = '└'
			} else if y == self.Y+self.H && x == self.X+self.W {
				level[y][x].Render = '┘'
			} else if y == self.Y || y == self.Y+self.H {
				level[y][x].Render = '─'
			} else if x == self.X || x == self.X+self.W {
				level[y][x].Render = '│'
			} else {
				level[y][x].Type = T_Floor
				level[y][x].Subtype = ST_Dirt
			}
			level[y][x].Lit = true
		}
	}
}

// Place puts a random type of something in the room in a random location that
// doesn't already have something else.
func (self *Room) Place(level *Level, ptype LocationType,
	psubtype LocationSubtype) bool {
	for i := 0; i < 20; i++ {
		x, y := self.X+rand.Intn(self.W-2), self.Y+rand.Intn(self.H-2)
		if level[y+1][x+1].Type == T_Floor && level[y+1][x+1].Subtype == ST_Dirt {
			level[y+1][x+1].Type = ptype
			level[y+1][x+1].Subtype = psubtype
			return true
		}
	}
	return false
}

func placeDoor(level *Level, x, y int) {
	loc := &level[y][x]
	loc.Type, loc.Subtype = T_Wall, ST_Door
	if rand.Intn(7) == 0 {
		loc.Open = true
	} else {
		loc.Open = false
	}
	if y < WHEIGHT && level[y+1][x].Type == T_Wall {
		loc.Render = '│'
	} else {
		loc.Render = '─'
	}
}

func doesIntersect(rooms []*Room, newroom *Room) bool {
	for _, room := range rooms {
		if !(newroom.X+newroom.W < room.X-1) &&
			!(newroom.Y+newroom.H < room.Y-1) &&
			!(newroom.X > room.X+room.W+1) &&
			!(newroom.Y > room.Y+room.H+1) {
			return true
		}
	}
	return false
}

func randomTunnelPoint(room *Room, level *Level) (x, y int) {
	pair, side := rand.Intn(2), rand.Intn(2)
	if pair == 0 {
		y = room.Y + rand.Intn(room.H-2) + 1
		if side == 0 {
			x = room.X - 1
			placeDoor(level, room.X, y)
		} else {
			x = room.X + room.W + 1
			placeDoor(level, room.X+room.W, y)
		}
	} else {
		x = room.X + rand.Intn(room.W-2) + 1
		if side == 0 {
			y = room.Y - 1
			placeDoor(level, x, room.Y)
		} else {
			y = room.Y + room.H + 1
			placeDoor(level, x, room.Y+room.H)
		}
	}
	return
}

func tryTunnel(from, to *Room, level *Level) bool {
	cx, cy := randomTunnelPoint(from, level)
	if cx < 0 || cy < 0 || cx >= WWIDTH || cy >= WHEIGHT ||
		(level[cy][cx].Type != T_Rock && level[cy][cx].Type != T_Tunnel) {
		return false
	}
	level[cy][cx].Type = T_Tunnel

	dx, dy := randomTunnelPoint(to, level)
	if dx < 0 || dy < 0 || dx >= WWIDTH || dy >= WHEIGHT ||
		(level[dy][dx].Type != T_Rock && level[dy][dx].Type != T_Tunnel) {
		return false
	}
	level[dy][dx].Type = T_Tunnel

	// Great, looks good ... start drawing a point from point to point, but do
	// not cross any walls. If we hit an edge, abort.
	for {
		// Done if cx/cy are hit.
		if cx == dx && cy == dy {
			break
		}

		// Try to move left/right.
		moved := false
		if dx < cx && (level[cy][cx-1].Type == T_Rock ||
			level[cy][cx-1].Type == T_Tunnel) {
			moved = true
			cx--
		} else if dx > cx && (level[cy][cx+1].Type == T_Rock ||
			level[cy][cx+1].Type == T_Tunnel) {
			moved = true
			cx++
		}
		level[cy][cx].Type = T_Tunnel

		// Try to move up/down.
		if dy < cy && (level[cy-1][cx].Type == T_Rock ||
			level[cy-1][cx].Type == T_Tunnel) {
			moved = true
			cy--
		} else if dy > cy && (level[cy+1][cx].Type == T_Rock ||
			level[cy+1][cx].Type == T_Tunnel) {
			moved = true
			cy++
		}
		level[cy][cx].Type = T_Tunnel

		if !moved {
			return false
		}
	}

	// If we get here, these rooms are connected OK now.
	return true
}

func generateLevel() (level *Level) {
	level = new(Level)

	// Make everything ROCK. (We salute you, etc.)
	for y := 0; y < WHEIGHT; y++ {
		for x := 0; x < WWIDTH; x++ {
			level[y][x].Type = T_Rock
			level[y][x].Diggable = true
		}
	}

	// There are types of levels, or will be. For now, we just generate levels
	// based on randomly dropping rooms. Each level is filled with rooms until
	// we run out of spaces to put them.
	rooms := make([]*Room, 0)
	for ct := 0; len(rooms) < 3 || ct < ROOM_DENSITY; {
		room := randomRoom()
		if doesIntersect(rooms, room) {
			ct++
			continue
		}
		rooms = append(rooms, room)
		room.Draw(level)
	}

	// Draw some tunnels, every room gets a chance at being connected to
	// another room. Doesn't always happen, but we try...
	from, roomct, lastroom := rooms[0], len(rooms), 0
	for i := 1; i < roomct-1; i++ {
		if tryTunnel(from, rooms[i], level) {
			from = rooms[i]
			lastroom = i
		}
	}

	// We put the downstairs in the first room, and the upstairs in the last
	// room we link. This ends up looking random to the user.
	ok := rooms[0].Place(level, T_Floor, ST_StairsUp)
	ok = ok && rooms[lastroom].Place(level, T_Floor, ST_StairsDown)
	if !ok {
		// Failed sanity check. This level is completely fucked without stairs,
		// so we have to bail out here. Just ask for another level.
		return generateLevel()
	}

	// Now we can add a few more tunnels randomly, just for fun.
	for i := 0; i < TUNNEL_DENSITY; i++ {
		idx1, idx2 := rand.Intn(roomct), rand.Intn(roomct)
		if idx1 != idx2 {
			tryTunnel(rooms[idx1], rooms[idx2], level)
		}
	}
	return
}

// findLocation returns an array of locations for a given type/subtype in he
// current World.
func findLocation(ptype LocationType, psubtype LocationSubtype) (res []LocXY) {
	res = make([]LocXY, 0)
	for y := 0; y < WHEIGHT; y++ {
		for x := 0; x < WWIDTH; x++ {
			if World[y][x].Type == ptype && World[y][x].Subtype == psubtype {
				res = append(res, LocXY{x, y})
			}
		}
	}
	return
}

// isPassable returns a true/false if a space is generally passible: i.e., it's
// not a wall, door, or other thing that can't be passed through. Note that
// this does not consider whether or not the space is already occupied.
func isPassable(x, y int) bool {
	loc := &World[y][x]
	switch loc.Type {
	case T_Floor, T_Tunnel, T_Water:
		return true
	case T_Wall:
		if loc.Subtype == ST_Door && loc.Open {
			return true
		}
		return false
	case T_Rock:
		return false
	}
	panic("isPassable had a fall-through type/subtype!")
}

// tryMove takes a thing and coordinates and attempts to move it there. If it
// fails, returns false.
func tryMove(thing *Mobile, x, y int) bool {
	if World[y][x].Occupant != nil {
		return false
	} else if !isPassable(x, y) {
		return false
	}

	if thing.X >= 0 && thing.Y >= 0 {
		if World[thing.Y][thing.X].Occupant != thing {
			panic("teleport: thing disagrees with World.Occupant!")
		}
		World[thing.Y][thing.X].Occupant = nil
		World[thing.Y][thing.X].Redraw = true
	}

	World[y][x].Occupant = thing
	World[y][x].Redraw = true
	thing.X, thing.Y = x, y
	return true
}

// teleport is a forgiving method of moving a thing from wherever it is to some
// new location. If the location given isn't valid, we try to find one nearby.
func teleport(thing *Mobile, x, y int) bool {
	return tryMove(thing, x, y)

	for i := 0; i < WWIDTH; i++ {
		for j := 0; j <= i; j++ {

		}
	}
	return true
}
