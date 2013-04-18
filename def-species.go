/* def-species.go - defines the species in the game
 *
 * Copyright (c) 2013 by Mark Smith <mark@qq.is>.
 *
 * See the included LICENSE file for information. 
 */

package main

import (
	"github.com/nsf/termbox-go"
)

// 	                    Str, Con, Dex, Wis, Int, Cha
var DefMaxAttr Attr = Attr{18, 18, 18, 18, 18, 18}

// FaunaList contains the list of possible creatures in the game. This includes
// things that could be PCs (base races).
var FaunaList []*Species = []*Species{
	&Species{Humanoid, "human", '@', termbox.ColorWhite,
		A_Neutral, 1, 6, 2, SAttr{
			Attr{10, 10, 10, 10, 10, 10}, Attr{6, 6, 6, 6, 6, 6}, DefMaxAttr,
		}},
	&Species{Humanoid, "elf", '@', termbox.ColorWhite,
		A_Neutral, 1, 6, 0, SAttr{
			Attr{8, 8, 10, 12, 12, 10}, Attr{6, 6, 6, 6, 6, 6}, DefMaxAttr,
		}},
	&Species{Jelly, "blue jelly", 'j', termbox.ColorBlue + termbox.AttrBold,
		A_Neutral, 3, 3, 0, SAttr{
			Attr{2, 4, 0, 0, 0, 0}, Attr{0, 0, 0, 0, 0, 0}, DefMaxAttr,
		}},
}
var Fauna map[string]*Species = make(map[string]*Species)

func init() {
	// Simply makes the Fauna map. This saves redundancy in the data structure
	// at the cost of a little startup time.
	for _, sp := range FaunaList {
		Fauna[sp.Name] = sp
	}
}
