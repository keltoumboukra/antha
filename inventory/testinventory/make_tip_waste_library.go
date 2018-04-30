// Copyright (C) 2017 The Antha authors. All rights reserved.
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

package testinventory

import "github.com/antha-lang/antha/antha/anthalib/wtype"

func makeTipwastes() (tipwastes []*wtype.LHTipwaste) {
	tipwastes = append(tipwastes, makeGilsonTipWaste(), makeCyBioTipwaste(), makeManualTipwaste(), makeTecanTipwaste())
	return
}

func makeGilsonTipWaste() *wtype.LHTipwaste {
	shp := wtype.NewShape("box", "mm", 123.0, 80.0, 92.0)
	w := wtype.NewLHWell("ul", 800000.0, 800000.0, shp, 0, 123.0, 80.0, 92.0, 0.0, "mm")
	lht := wtype.NewLHTipwaste(6000, "Gilsontipwaste", "gilson", wtype.Coordinates{X: 127.76, Y: 85.48, Z: 92.0}, w, 49.5, 31.5, 0.0)
	return lht
}

// TODO figure out tip capacity
func makeCyBioTipwaste() *wtype.LHTipwaste {
	shp := wtype.NewShape("box", "mm", 90.5, 171.0, 90.0)
	w := wtype.NewLHWell("ul", 800000.0, 800000.0, shp, 0, 90.5, 171.0, 90.0, 0.0, "mm")
	lht := wtype.NewLHTipwaste(700, "CyBiotipwaste", "cybio", wtype.Coordinates{X: 127.76, Y: 85.48, Z: 90.5}, w, 85.5, 45.0, 0.0)
	return lht
}

// TODO figure out tip capacity
func makeManualTipwaste() *wtype.LHTipwaste {
	shp := wtype.NewShape("box", "mm", 90.5, 171.0, 90.0)
	w := wtype.NewLHWell("ul", 800000.0, 800000.0, shp, 0, 90.5, 171.0, 90.0, 0.0, "mm")
	lht := wtype.NewLHTipwaste(1000000, "Manualtipwaste", "ACMEBagsInc", wtype.Coordinates{X: 127.76, Y: 85.48, Z: 90.5}, w, 85.5, 45.0, 0.0)
	return lht
}

func makeTecanTipwaste() *wtype.LHTipwaste {
	shp := wtype.NewShape("box", "mm", 90.5, 171.0, 90.0)
	w := wtype.NewLHWell("ul", 800000.0, 800000.0, shp, 0, 90.5, 171.0, 90.0, 0.0, "mm")
	lht := wtype.NewLHTipwaste(2000, "Tecantipwaste", "Tecan", wtype.Coordinates{X: 127.76, Y: 85.48, Z: 90.5}, w, 85.5, 45.0, 0.0)
	return lht
}
