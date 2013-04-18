/* def-prototypes.go - defines the prototypes/behaviors
 *
 * Copyright (c) 2013 by Mark Smith <mark@qq.is>.
 *
 * See the included LICENSE file for information. 
 */

package main

// First, we need to define classes of prototypes. These are instantiated
// later, but we do it this way to define behavior for each of the types and
// allow behavioral inheritance to work.
type Animal struct {
	Prototype
}

// Represents things that can reason and think.
type Intelligent struct {
	Animal
}

// Definitions of the prototypes.
var Humanoid *Intelligent = &Intelligent{Animal{Prototype{
	S_Eyes, B_Air, 2, false, 100,
}}}
var Quadruped *Animal = &Animal{Prototype{
	S_Eyes, B_Air, 0, false, 80,
}}
var Jelly *Animal = &Animal{Prototype{
	S_Blind, B_None, 0, true, 100,
}}
var Octopus *Animal = &Animal{Prototype{
	S_Eyes, B_Water, 8, false, 80,
}}
