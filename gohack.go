/* gohack - a nethack game in Go
 *
 * Copyright (c) 2013 by Mark Smith <mark@qq.is>.
 *
 * See the included LICENSE file for information. 
 */

package main

import (
	"github.com/nsf/termbox-go"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// The player starts out as human. Enjoy!
	PC = newMobile(Fauna["human"])

	startScreen()
	defer endScreen()

	World = generateLevel()

	// Drop the player on the upstairs.
	ds := findLocation(T_Floor, ST_StairsUp)
	if len(ds) < 1 {
		panic("no upstairs on world!")
	}
	if !teleport(PC, ds[0].X, ds[0].Y) {
		panic("failed to teleport PC onto level")
	}

MAIN:
	for {
		// Depending on what happens, we advance time. If the user does
		// something like move, then that happens and then time advances by the
		// user's speed, roughly. See the game timing document.
		advance := 0
		xMove, yMove := 0, 0

		// First, draw/redraw the screen.
		drawScreen()

		// Now get the user's input. We should be in a command cycle here: if
		// the user was typing somethign like a sentence, that'd be handled by
		// another part of the program.
		evt := termbox.PollEvent()
		switch evt.Ch {
		case 'q':
			break MAIN
		case 'h':
			xMove = -1
		case 'j':
			yMove = 1
		case 'k':
			yMove = -1
		case 'l':
			xMove = 1
		case 'y':
			xMove, yMove = -1, -1
		case 'u':
			xMove, yMove = 1, -1
		case 'b':
			xMove, yMove = -1, 1
		case 'n':
			xMove, yMove = 1, 1
		}

		if xMove != 0 || yMove != 0 {
			tryMove(PC, PC.X+xMove, PC.Y+yMove)
		}

		advanceTime(advance)
	}
}
