// antha/doc/testdata/b.go: Part of the Antha language
// Copyright (C) 2014 The Antha authors. All rights reserved.
// 
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
// 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
// 
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o 
// Synthace Ltd. The London Bioscience Innovation Centre
// 1 Royal College St, London NW1 0NH UK


package b

import "a"

// ----------------------------------------------------------------------------
// Basic declarations

const Pi = 3.14   // Pi
var MaxInt int    // MaxInt
type T struct{}   // T
var V T           // v
func F(x int) int {} // F
func (x *T) M()   {} // M

// Corner cases: association with (presumed) predeclared types

// Always under the package functions list.
func NotAFactory() int {}

// Associated with uint type if AllDecls is set.
func UintFactory() uint {}

// Associated with uint type if AllDecls is set.
func uintFactory() uint {}

// Should only appear if AllDecls is set.
type uint struct{} // overrides a predeclared type uint

// ----------------------------------------------------------------------------
// Exported declarations associated with non-exported types must always be shown.

type notExported int

const C notExported = 0

const (
	C1 notExported = iota
	C2
	c3
	C4
	C5
)

var V notExported
var V1, V2, v3, V4, V5 notExported

var (
	U1, U2, u3, U4, U5 notExported
	u6                 notExported
	U7                 notExported = 7
)

func F1() notExported {}
func f2() notExported {}