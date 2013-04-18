/* display.go - the gohack display library
 *
 * Copyright (c) 2013 by Mark Smith <mark@qq.is>.
 *
 * See the included LICENSE file for information. 
 */

package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"time"
)

var sWidth, sHeight int

func startScreen() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	// Ensure basic size is fine. Don't run if size is too small.
	w, h := termbox.Size()
	if w < WWIDTH || h < WHEIGHT {
		panic("screen is too small to play gohack")
	}
}

func drawScreen() {
	// Detect if the screen size has changed so we do a full clear+redraw.
	redraw := false
	w, h := termbox.Size()
	if w != sWidth || h != sHeight {
		if w < WWIDTH || h < WHEIGHT {
			fmt.Print("SCREEN TOO SMALL. Please make larger to continue.\n")
			time.Sleep(1 * time.Second)
			return
		}
		sWidth, sHeight = w, h
		termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
		redraw = true
	}

	// Draw things that either need drawing, or everything if we've been asked
	// to do a complete redraw.
	for y := 0; y < WHEIGHT; y++ {
		for x := 0; x < WWIDTH; x++ {
			if redraw || World[y][x].Redraw {
				World[y][x].Draw(x, y)
			}
		}
	}
	termbox.Flush()
}

func endScreen() {
	termbox.Close()
}

// Draw instructs a location to blit itself to the screen.
func (self *Location) Draw(x, y int) {
	self.Redraw = false
	char, fg, bg, bold := ' ', termbox.ColorWhite, termbox.ColorBlack, true

	// If occupied, that's easy... do it
	if self.Occupant != nil {
		if self.Occupant == PC {
			termbox.SetCell(x, y, self.Occupant.S.Render,
				self.Occupant.S.Color+termbox.AttrBold, bg)
		} else {
			termbox.SetCell(x, y, self.Occupant.S.Render,
				self.Occupant.S.Color, bg)
		}
		return
	}

	// Default rendering for other types of things.
	switch self.Type {
	case T_Floor:
		switch self.Subtype {
		case ST_Dirt:
			if self.Lit {
				char = '.'
			}
		case ST_Fountain:
			char, fg = '{', termbox.ColorBlue
		case ST_Grave:
			char, fg = '|', termbox.ColorRed
		case ST_Sink:
			char, fg = '#', termbox.ColorGreen
		case ST_Altar:
			char, fg = '_', termbox.ColorMagenta
		case ST_StairsDown:
			char, fg = '>', termbox.ColorRed
		case ST_StairsUp:
			char, fg = '<', termbox.ColorCyan
		}
	case T_Tunnel:
		char, bold = '▒', self.Lit
	case T_Water:
		char, fg = '~', termbox.ColorBlue
	case T_Wall:
		char = self.Render
		if self.Subtype == ST_Door {
			fg = termbox.ColorMagenta
			if !self.Open {
				char = '+'
			}
		}
	case T_Rock:
		// default
	default:
		char, bg = rune('a'+self.Type), termbox.ColorRed
	}
	if bold {
		fg = fg + termbox.AttrBold
	}
	termbox.SetCell(x, y, char, fg, bg)
}
