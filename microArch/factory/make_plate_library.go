// anthalib/factory/make_plate_library.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
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
// 2 Royal College St, London NW1 0NH UK

package factory

import (
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/devices"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

//var commonwelltypes

func makePlateLibrary() map[string]*wtype.LHPlate {
	plates := make(map[string]*wtype.LHPlate)

	// deep square well 96
	swshp := wtype.NewShape("box", "mm", 8.2, 8.2, 41.3)
	welltype := wtype.NewLHWell("DSW96", "", "", "ul", 2000, 25, swshp, 3, 8.2, 8.2, 41.3, 4.7, "mm")
	plate := wtype.NewLHPlate("DSW96", "Unknown", 8, 12, 44.1, "mm", welltype, 9, 9, 0.0, 0.0, 0.0)
	plates[plate.Type] = plate

	// shallow round well flat bottom 96
	rwshp := wtype.NewShape("cylinder", "mm", 8.2, 8.2, 11)
	welltype = wtype.NewLHWell("SRWFB96", "", "", "ul", 500, 10, rwshp, 0, 8.2, 8.2, 11, 1.0, "mm")
	plate = wtype.NewLHPlate("SRWFB96", "Unknown", 8, 12, 15, "mm", welltype, 9, 9, 0.0, 0.0, 0.0)
	plates[plate.Type] = plate

	// deep well strip trough 12
	stshp := wtype.NewShape("box", "mm", 8.2, 72, 41.3)
	welltype = wtype.NewLHWell("DWST12", "", "", "ul", 15000, 1000, stshp, 3, 8.2, 72, 41.3, 4.7, "mm")
	plate = wtype.NewLHPlate("DWST12", "Unknown", 1, 12, 44.1, "mm", welltype, 9, 9, 0, 0, 0.0)
	plates[plate.Type] = plate

	// deep well strip trough 8
	stshp = wtype.NewShape("box", "mm", 115.0, 8.2, 41.3)
	welltype = wtype.NewLHWell("DWST8", "", "", "ul", 24000, 1000, stshp, 3, 115, 8.2, 41.3, 4.7, "mm")
	plate = wtype.NewLHPlate("DWST8", "Unknown", 8, 1, 44.1, "mm", welltype, 9, 9, 49.5, 0.0, 0.0)
	plates[plate.Type] = plate

	// deep well reservoir
	rshp := wtype.NewShape("box", "mm", 115.0, 72.0, 41.3)
	welltype = wtype.NewLHWell("DWR1", "", "", "ul", 300000, 20000, rshp, 3, 115, 72, 41.3, 4.7, "mm")
	plate = wtype.NewLHPlate("DWR1", "Unknown", 1, 1, 44.1, "mm", welltype, 9, 9, 49.5, 0.0, 0.0)
	plates[plate.Type] = plate

	// pcr plate with cooler
	cone := wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
	welltype = wtype.NewLHWell("pcrplate", "", "", "ul", 250, 5, cone, 0, 5.5, 5.5, 20.4, 1.4, "mm")
	//plate = wtype.NewLHPlate("pcrplate", "Unknown", 8, 12, 25.7, "mm", welltype, 9, 9, 0.0, 0.0, 6.5)
	//plates[plate.Type] = plate
	plate = wtype.NewLHPlate("pcrplate_with_cooler", "Unknown", 8, 12, 25.7, "mm", welltype, 9, 9, 0.0, 0.0, 15.5)
	plates[plate.Type] = plate

	// pcr plate with incubator
	cone = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
	welltype = wtype.NewLHWell("pcrplate", "", "", "ul", 250, 5, cone, 0, 5.5, 5.5, 20.4, 1.4, "mm")
	plate = wtype.NewLHPlate("pcrplate_with_incubater", "Unknown", 8, 12, 25.7, "mm", welltype, 9, 9, 0.0, 0.0, (15.5 + 44.0))
	plates[plate.Type] = plate

	// Block Kombi 2ml
	eppy := wtype.NewShape("cylinder", "mm", 8.2, 8.2, 45)

	wellxoffset := 18.0 // centre of well to centre of neightbouring well in x direction
	wellyoffset := 18.0 //centre of well to centre of neightbouring well in y direction
	xstart := 5.0       // distance from top left side of plate to first well
	ystart := 5.0       // distance from top left side of plate to first well
	zstart := 6.0       // offset of bottom of deck to bottom of well

	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype = wtype.NewLHWell("2mlEpp", "", "", "ul", 2000, 25, eppy, 3, 8.2, 8.2, 45, 4.7, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	plate = wtype.NewLHPlate("Kombi2mlEpp", "Unknown", 4, 2, 45, "mm", welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates[plate.Type] = plate

	// greiner 384 well plate flat bottom
	square := wtype.NewShape("box", "mm", 4, 4, 14)

	wellxoffset = 4.5 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 4.5 //centre of well to centre of neighbouring well in y direction
	xstart = 9.0      // distance from top left side of plate to first well
	ystart = 6.0      // distance from top left side of plate to first well
	zstart = 3.0      // offset of bottom of deck to bottom of well

	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype = wtype.NewLHWell("384flat", "", "", "ul", 100, 10, square, 0, 4, 4, 14, 1, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	plate = wtype.NewLHPlate("greiner384", "Unknown", 16, 24, 14, "mm", welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates[plate.Type] = plate

	// 250ml box reservoir (working vol estimated to be 100ml to prevent spillage on moving decks)
	reservoirbox := wtype.NewShape("box", "mm", 71, 107, 38) // 39?
	welltype = wtype.NewLHWell("Reservoir", "", "", "ul", 100000, 10000, reservoirbox, 0, 107, 71, 38, 3, "mm")
	plate = wtype.NewLHPlate("reservoir", "unknown", 1, 1, 45, "mm", welltype, 58, 13, 0, 0, 10)
	plates[plate.Type] = plate
	/*
		rwshp = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
		welltype = wtype.NewLHWell("pcrplate", "", "", "ul", 250, 5, rwshp, 0, 5.5, 5.5, 20.4, 1.4, "mm")
		//plate = wtype.NewLHPlate("pcrplate", "Unknown", 8, 12, 25.7, "mm", welltype, 9, 9, 0.0, 0.0, 6.5)
		//plates[plate.Type] = plate
		plate = wtype.NewLHPlate("pcrplate_with_skirt", "Unknown", 8, 12, 25.7, "mm", welltype, 9, 9, 0.0, 0.0, 15.5)
		plates[plate.Type] = plate
	*/
	/// placeholder for non plate container for testing
	rwshp = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
	welltype = wtype.NewLHWell("pcrplate", "", "", "ul", 250, 5, rwshp, 0, 5.5, 5.5, 20.4, 1.4, "mm")
	//plate = wtype.NewLHPlate("pcrplate", "Unknown", 8, 12, 25.7, "mm", welltype, 9, 9, 0.0, 0.0, 6.5)
	//plates[plate.Type] = plate
	plate = wtype.NewLHPlate("1L_DuranBottle", "Unknown", 8, 12, 25.7, "mm", welltype, 9, 9, 0.0, 0.0, 15.5)
	plates[plate.Type] = plate

	//NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64)
	return plates
}

func GetPlateByType(typ string) *wtype.LHPlate {
	plates := makePlateLibrary()
	p := plates[typ]
	return p.Dup()
}

func GetPlateList() []string {
	plates := makePlateLibrary()

	kz := make([]string, len(plates))
	x := 0
	for name, _ := range plates {
		kz[x] = name
		x += 1
	}
	return kz
}
