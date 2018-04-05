package liquidhandling

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/devices"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

//var commonwelltypes

func makePlateLibrary() map[string]*wtype.LHPlate {
	plates := make(map[string]*wtype.LHPlate)

	offset := 0.25
	riserheightinmm := 40.0 - offset
	incubatorheightinmm := devices.Shaker["3000 T-elm"]["Height"] * 1000

	inhecoincubatorinmm := devices.Shaker["InhecoStaticOnDeck"]["Height"] * 1000

	valueformaxheadtonotintoDSWplatewithp20tips := 4.5
	// deep square well 96
	swshp := wtype.NewShape("box", "mm", 8.2, 8.2, 41.3)
	welltype := wtype.NewLHWell("ul", 1000, 100, swshp, wtype.VWellBottom, 8.2, 8.2, 41.3, 4.7, "mm")
	plate := wtype.NewLHPlate("DSW96", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 44.1}, welltype, 9, 9, 0.0, 0.0, valueformaxheadtonotintoDSWplatewithp20tips)
	plates[plate.Type] = plate

	// deep square well 96 on riser
	swshp = wtype.NewShape("box", "mm", 8.2, 8.2, 41.3)
	welltype = wtype.NewLHWell("ul", 1000, 100, swshp, wtype.VWellBottom, 8.2, 8.2, 41.3, 4.7, "mm")
	plate = wtype.NewLHPlate("DSW96_riser", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 44.1}, welltype, 9, 9, 0.0, 0.0, riserheightinmm)
	plates[plate.Type] = plate

	// deep square well 96 on q instruments incubator
	swshp = wtype.NewShape("box", "mm", 8.2, 8.2, 41.3)
	welltype = wtype.NewLHWell("ul", 1000, 100, swshp, wtype.VWellBottom, 8.2, 8.2, 41.3, 4.7, "mm")
	plate = wtype.NewLHPlate("DSW96_incubator", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 44.1}, welltype, 9, 9, 0.0, 0.0, incubatorheightinmm)
	plates[plate.Type] = plate

	// deep square well 96 on inheco incubator
	swshp = wtype.NewShape("box", "mm", 8.2, 8.2, 41.3)
	welltype = wtype.NewLHWell("ul", 1000, 100, swshp, wtype.VWellBottom, 8.2, 8.2, 41.3, 4.7, "mm")
	plate = wtype.NewLHPlate("DSW96_inheco", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 44.1}, welltype, 9, 9, 0.0, 0.0, inhecoincubatorinmm)
	plates[plate.Type] = plate

	// 24 well deep square well plate on riser

	bottomtype := wtype.VWellBottom // 0 = flat, 2 = v shaped
	xdim := 16.8
	ydim := 16.8
	zdim := 41.3
	bottomh := 4.7

	wellcapacityinwelltypeunit := 11000.0
	welltypeunit := "ul"
	wellsperrow := 6
	wellspercolumn := 4
	residualvol := 650.0 // assume in ul

	wellxoffset := 18.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset := 18.0 //centre of well to centre of neighbouring well in y direction
	xstart := 4.5       // distance from top left side of plate to first well
	ystart := 4.5       // distance from top left side of plate to first well
	zstart := -1.0      // offset of bottom of deck to bottom of well (this includes agar estimate)

	zstart = zstart + riserheightinmm

	heightinmm := 44.1

	squarewell := wtype.NewShape("box", "mm", xdim, ydim, zdim)
	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype = wtype.NewLHWell(welltypeunit, wellcapacityinwelltypeunit, residualvol, squarewell, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	plate = wtype.NewLHPlate("DSW24_riser", "Unknown", wellspercolumn, wellsperrow, wtype.Coordinates{127.76, 85.48, heightinmm}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates[plate.Type] = plate

	// shallow round well flat bottom 96
	rwshp := wtype.NewShape("cylinder", "mm", 8.2, 8.2, 11)
	welltype = wtype.NewLHWell("ul", 500, 10, rwshp, 0, 8.2, 8.2, 11, 1.0, "mm")
	plate = wtype.NewLHPlate("SRWFB96", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 15}, welltype, 9, 9, 0.0, 0.0, 1.0)
	plates[plate.Type] = plate

	// shallow round well flat bottom 96 on riser
	rwshp = wtype.NewShape("cylinder", "mm", 8.2, 8.2, 11)
	welltype = wtype.NewLHWell("ul", 500, 10, rwshp, 0, 8.2, 8.2, 11, 1.0, "mm")
	plate = wtype.NewLHPlate("SRWFB96_riser", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 15}, welltype, 9, 9, 0.0, 0.0, 40.0)
	plates[plate.Type] = plate

	// deep well strip trough 12
	stshp := wtype.NewShape("box", "mm", 8.2, 72, 41.3)
	welltype = wtype.NewLHWell("ul", 15000, 1000, stshp, wtype.VWellBottom, 8.2, 72, 41.3, 4.7, "mm")
	plate = wtype.NewLHPlate("DWST12", "Unknown", 1, 12, wtype.Coordinates{127.76, 85.48, 44.1}, welltype, 9, 9, 0, 30.0, 0.0)
	plates[plate.Type] = plate

	// deep well strip trough 12 on riser
	stshp = wtype.NewShape("box", "mm", 8.2, 72, 41.3)
	welltype = wtype.NewLHWell("ul", 15000, 1000, stshp, wtype.VWellBottom, 8.2, 72, 41.3, 4.7, "mm")
	plate = wtype.NewLHPlate("DWST12_riser", "Unknown", 1, 12, wtype.Coordinates{127.76, 85.48, 44.1}, welltype, 9, 9, 0, 30.0, riserheightinmm+2.0)
	plates[plate.Type] = plate

	// deep well strip trough 8
	stshp = wtype.NewShape("box", "mm", 115.0, 8.2, 41.3)
	welltype = wtype.NewLHWell("ul", 24000, 1000, stshp, wtype.VWellBottom, 115, 8.2, 41.3, 4.7, "mm")
	plate = wtype.NewLHPlate("DWST8", "Unknown", 8, 1, wtype.Coordinates{127.76, 85.48, 44.1}, welltype, 9, 9, 49.5, 0.0, 0.0)
	plates[plate.Type] = plate

	// deep well reservoir
	rshp := wtype.NewShape("box", "mm", 115.0, 72.0, 41.3)
	welltype = wtype.NewLHWell("ul", 300000, 20000, rshp, wtype.VWellBottom, 115, 72, 41.3, 4.7, "mm")
	plate = wtype.NewLHPlate("DWR1", "Unknown", 1, 1, wtype.Coordinates{127.76, 85.48, 44.1}, welltype, 9, 9, 49.5, 0.0, 0.0)
	plates[plate.Type] = plate

	// pcr plate with cooler
	cone := wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
	welltype = wtype.NewLHWell("ul", 250, 5, cone, wtype.UWellBottom, 5.5, 5.5, 20.4, 1.4, "mm")
	//plate = wtype.NewLHPlate("pcrplate", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, 6.5)
	//plates[plate.Type] = plate
	plate = wtype.NewLHPlate("pcrplate_with_cooler", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, 15.5)
	plates[plate.Type] = plate

	// pcr plate skirted (on riser)
	cone = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
	welltype = wtype.NewLHWell("ul", 200, 5, cone, wtype.UWellBottom, 5.5, 5.5, 20.4, 1.4, "mm")

	plate = wtype.NewLHPlate("pcrplate_skirted_riser", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, riserheightinmm-1.25)

	plates[plate.Type] = plate

	// pcr plate skirted
	cone = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
	welltype = wtype.NewLHWell("ul", 200, 5, cone, wtype.UWellBottom, 5.5, 5.5, 20.4, 1.4, "mm")
	plate = wtype.NewLHPlate("pcrplate_skirted", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, 0.636)
	plates[plate.Type] = plate

	// pcr plate with incubator
	cone = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
	welltype = wtype.NewLHWell("ul", 250, 5, cone, wtype.UWellBottom, 5.5, 5.5, 20.4, 1.4, "mm")
	plate = wtype.NewLHPlate("pcrplate_with_incubater", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, (15.5 + 44.0))

	consar := []string{"position_1"}
	plate.SetConstrained("Pipetmax", consar)

	plates[plate.Type] = plate

	// Block Kombi 2ml
	eppy := wtype.NewShape("cylinder", "mm", 8.2, 8.2, 45)

	wellxoffset = 18.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 18.0 //centre of well to centre of neighbouring well in y direction
	xstart = 5.0       // distance from top left side of plate to first well
	ystart = 5.0       // distance from top left side of plate to first well
	zstart = 6.0       // offset of bottom of deck to bottom of well

	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype = wtype.NewLHWell("ul", 2000, 25, eppy, wtype.VWellBottom, 8.2, 8.2, 45, 4.7, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	plate = wtype.NewLHPlate("Kombi2mlEpp", "Unknown", 4, 2, wtype.Coordinates{127.76, 85.48, 45}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates[plate.Type] = plate

	// greiner 384 well plate flat bottom

	bottomtype = wtype.FlatWellBottom
	xdim = 4.0
	ydim = 4.0
	zdim = 14.0
	bottomh = 1.0

	wellxoffset = 4.5 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 4.5 //centre of well to centre of neighbouring well in y direction
	xstart = -2.5     // distance from top left side of plate to first well
	ystart = -2.5     // distance from top left side of plate to first well
	zstart = 2.5      // offset of bottom of deck to bottom of well

	square := wtype.NewShape("box", "mm", 4, 4, 14)
	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype = wtype.NewLHWell("ul", 125, 10, square, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	plate = wtype.NewLHPlate("greiner384", "Unknown", 16, 24, wtype.Coordinates{127.76, 85.48, 14}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates[plate.Type] = plate

	// greiner 384 well plate flat bottom on riser

	bottomtype = wtype.FlatWellBottom
	xdim = 4.0
	ydim = 4.0
	zdim = 12.0 // modified from 14
	bottomh = 1.0

	wellxoffset = 4.5               // centre of well to centre of neighbouring well in x direction
	wellyoffset = 4.5               //centre of well to centre of neighbouring well in y direction
	xstart = -2.5                   // distance from top left side of plate to first well
	ystart = -2.5                   // distance from top left side of plate to first well
	zstart = riserheightinmm + 0.25 // offset of bottom of deck to bottom of well

	square = wtype.NewShape("box", "mm", 4, 4, 14)
	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype = wtype.NewLHWell("ul", 125, 10, square, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	plate = wtype.NewLHPlate("greiner384_riser", "Unknown", 16, 24, wtype.Coordinates{127.76, 85.48, 14}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates[plate.Type] = plate

	// NUNC 1536 well plate flat bottom on riser

	bottomtype = wtype.FlatWellBottom
	xdim = 2.0 // of well
	ydim = 2.0
	zdim = 7.0
	bottomh = 0.5

	wellxoffset = 2.25           // centre of well to centre of neighbouring well in x direction
	wellyoffset = 2.25           //centre of well to centre of neighbouring well in y direction
	xstart = -2.5                // distance from top left side of plate to first well
	ystart = -2.5                // distance from top left side of plate to first well
	zstart = riserheightinmm + 2 // offset of bottom of deck to bottom of well

	square = wtype.NewShape("box", "mm", 2, 2, 7)
	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype = wtype.NewLHWell("ul", 13, 2, square, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	plate = wtype.NewLHPlate("nunc1536_riser", "Unknown", 32, 48, wtype.Coordinates{127.76, 85.48, 7}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates[plate.Type] = plate

	// 250ml box reservoir (working vol estimated to be 100ml to prevent spillage on moving decks)
	reservoirbox := wtype.NewShape("box", "mm", 71, 107, 38) // 39?
	welltype = wtype.NewLHWell("ul", 100000, 10000, reservoirbox, 0, 107, 71, 38, 3, "mm")
	plate = wtype.NewLHPlate("reservoir", "unknown", 1, 1, wtype.Coordinates{127.76, 85.48, 45}, welltype, 58, 13, 0, 0, 10)
	plates[plate.Type] = plate

	// Onewell SBS format Agarplate with colonies on riser (50ml agar) high res

	bottomtype = wtype.FlatWellBottom
	xdim = 2.0 // of well
	ydim = 2.0
	zdim = 7.0
	bottomh = 0.5

	wellxoffset = 2.25           // centre of well to centre of neighbouring well in x direction
	wellyoffset = 2.250          //centre of well to centre of neighbouring well in y direction
	xstart = -2.5                // distance from top left side of plate to first well
	ystart = -2.5                // distance from top left side of plate to first well
	zstart = riserheightinmm + 3 // offset of bottom of deck to bottom of well

	square = wtype.NewShape("box", "mm", 2, 2, 7)
	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype = wtype.NewLHWell("ul", 13, 2, square, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	// greiner one well with 50ml of agar in
	plate = wtype.NewLHPlate("Agarplateforpicking1536_riser", "Unknown", 32, 48, wtype.Coordinates{127.76, 85.48, 7}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates[plate.Type] = plate

	// Onewell SBS format Agarplate with colonies on riser (50ml agar) low res

	bottomtype = wtype.FlatWellBottom
	xdim = 4.0
	ydim = 4.0
	zdim = 14.0
	bottomh = 1.0

	wellxoffset = 4.5            // centre of well to centre of neighbouring well in x direction
	wellyoffset = 4.5            //centre of well to centre of neighbouring well in y direction
	xstart = -2.5                // distance from top left side of plate to first well
	ystart = -2.5                // distance from top left side of plate to first well
	zstart = riserheightinmm + 3 // offset of bottom of deck to bottom of well

	square = wtype.NewShape("box", "mm", 4, 4, 14)
	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype = wtype.NewLHWell("ul", 125, 10, square, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	// greiner one well with 50ml of agar in
	plate = wtype.NewLHPlate("Agarplateforpicking384_riser", "Unknown", 16, 24, wtype.Coordinates{127.76, 85.48, 14}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates[plate.Type] = plate

	// Onewell SBS format Agarplate with colonies on riser (30ml agar) low res

	zstart = 41 // offset of bottom of deck to bottom of well

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	// greiner one well with 50ml of agar in
	plate = wtype.NewLHPlate("30mlAgarplateforpicking384_riser", "Unknown", 16, 24, wtype.Coordinates{127.76, 85.48, 14}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates[plate.Type] = plate

	/*
		rwshp = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
		welltype = wtype.NewLHWell("pcrplate", "", "", "ul", 250, 5, rwshp, 0, 5.5, 5.5, 20.4, 1.4, "mm")
		//plate = wtype.NewLHPlate("pcrplate", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, 6.5)
		//plates[plate.Type] = plate
		plate = wtype.NewLHPlate("pcrplate_with_skirt", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, 15.5)
		plates[plate.Type] = plate
	*/
	/// placeholder for non plate container for testing
	rwshp = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
	welltype = wtype.NewLHWell("ul", 250, 5, rwshp, wtype.FlatWellBottom, 5.5, 5.5, 20.4, 1.4, "mm")
	//plate = wtype.NewLHPlate("pcrplate", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, 6.5)
	//plates[plate.Type] = plate
	plate = wtype.NewLHPlate("1L_DuranBottle", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, 15.5)
	plates[plate.Type] = plate

	//forward position

	//	ep48g := wtype.NewShape("trap", "mm", 2, 4, 2)
	//	welltype = wtype.NewLHWell("EPAGE48", "", "", "ul", 15, 0, ep48g, 0, 2, 4, 2, 48, "mm")
	//	plate = wtype.NewLHPlate("EPAGE48", "Invitrogen", 2, 26, wtype.Coordinates{127.76, 85.48, 50}, welltype, 4.5, 34, 0.0, 0.0, 2.0)
	//	plates[plate.Type] = plate

	//refactored for reverse position

	ep48g := wtype.NewShape("trap", "mm", 2, 4, 2)
	//can't reach all wells; change to 24 wells per row?
	welltype = wtype.NewLHWell("ul", 25, 0, ep48g, wtype.FlatWellBottom, 2, 4, 2, 2, "mm")
	//welltype = wtype.NewLHWell("384flat", "", "", "ul", 100, 10, square, bottomtype, xdim, ydim, zdim, bottomh, "mm")
	//plate = wtype.NewLHPlate("EPAGE48", "Invitrogen", 2, 26, wtype.Coordinates{127.76, 85.48, 50}, welltype, 4.5, 34, -1.0, 17.25, 49.5)
	plate = wtype.NewLHPlate("EPAGE48", "Invitrogen", 2, 26, wtype.Coordinates{127.76, 85.48, 48.5}, welltype, 4.5, 33.75, -1.0, 18.0, riserheightinmm+4.5)
	//plate = wtype.NewLHPlate("greiner384", "Unknown", 16, 24, wtype.Coordinates{127.76, 85.48, 14}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)

	plates[plate.Type] = plate

	// E-GEL 96 definition

	//same welltype as EPAGE

	// due to staggering of wells: 1 96well gel is set up as two well types

	// 1st type
	//can't reach all wells; change to 12 wells per row?
	plate = wtype.NewLHPlate("EGEL96_1", "Invitrogen", 4, 13, wtype.Coordinates{127.76, 85.48, 48.5}, welltype, 9, 18.0, 0, -1.0, riserheightinmm+5.5)
	//plate = wtype.NewLHPlate("greiner384", "Unknown", 16, 24, wtype.Coordinates{127.76, 85.48, 14}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates[plate.Type] = plate

	// 2nd type
	plate = wtype.NewLHPlate("EGEL96_2", "Invitrogen", 4, 13, wtype.Coordinates{127.76, 85.48, 48.5}, welltype, 9, 18.0, 4.0, 7.5, riserheightinmm+5.5)
	//plate = wtype.NewLHPlate("greiner384", "Unknown", 16, 24, wtype.Coordinates{127.76, 85.48, 14}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)

	plates[plate.Type] = plate

	// falcon 6 well plate with Agar flat bottom with 4ml per well

	bottomtype = wtype.FlatWellBottom
	xdim = 37.0
	ydim = 37.0
	zdim = 20.0
	bottomh = 9.0 //(this includes agar estimate)

	wellxoffset = 39.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 39.0 //centre of well to centre of neighbouring well in y direction
	xstart = 5.0       // distance from top left side of plate to first well
	ystart = 5.0       // distance from top left side of plate to first well
	zstart = 9.0       // offset of bottom of deck to bottom of well (this includes agar estimate)

	wellsperrow = 3
	wellspercolumn = 2
	heightinmm = 20.0

	circle := wtype.NewShape("cylinder", "mm", 37, 37, 20)
	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype = wtype.NewLHWell("ul", 100, 10, circle, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	plate = wtype.NewLHPlate("falcon6wellAgar", "Unknown", wellspercolumn, wellsperrow, wtype.Coordinates{127.76, 85.48, heightinmm}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates[plate.Type] = plate

	// Nunclon 12 well plate with Agar flat bottom 2ml per well

	bottomtype = wtype.FlatWellBottom
	xdim = 22.5 // diameter
	ydim = 22.5 // diameter
	zdim = 20.0
	bottomh = 9.0 //(this includes agar estimate)

	wellxoffset = 27.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 27.0 //centre of well to centre of neighbouring well in y direction
	xstart = 11.0      // distance from top left side of plate to first well
	ystart = 4.0       // distance from top left side of plate to first well
	zstart = 9.0       // offset of bottom of deck to bottom of well (this includes agar estimate)

	wellsperrow = 4
	wellspercolumn = 3
	heightinmm = 22.0

	circle = wtype.NewShape("cylinder", "mm", xdim, ydim, zdim)
	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype = wtype.NewLHWell("ul", 100, 10, circle, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	plate = wtype.NewLHPlate("Nuncon12wellAgar", "Unknown", wellspercolumn, wellsperrow, wtype.Coordinates{127.76, 85.48, heightinmm}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates[plate.Type] = plate

	//	WellXOffset float64
	//	WellYOffset float64
	//	WellXStart  float64
	//	WellYStart  float64
	//	WellZStart  float64

	zstart = 9.0 + incubatorheightinmm // offset of bottom of deck to bottom of well (this includes agar estimate)
	welltype = wtype.NewLHWell("ul", 100, 10, circle, bottomtype, xdim, ydim, zdim, bottomh, "mm")
	plate = wtype.NewLHPlate("Nuncon12wellAgar_incubator", "Unknown", wellspercolumn, wellsperrow, wtype.Coordinates{127.76, 85.48, heightinmm}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)

	consar = []string{"position_1"}
	plate.SetConstrained("Pipetmax", consar)

	plates[plate.Type] = plate
	/*
		rwshp = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
		welltype = wtype.NewLHWell("pcrplate", "", "", "ul", 250, 5, rwshp, 0, 5.5, 5.5, 20.4, 1.4, "mm")
		//plate = wtype.NewLHPlate("pcrplate", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, 6.5)
		//plates[plate.Type] = plate
		plate = wtype.NewLHPlate("pcrplate_with_skirt", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, 15.5)
		plates[plate.Type] = plate
	*/

	/// placeholder for non plate container for testing
	rwshp = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
	welltype = wtype.NewLHWell("ul", 250, 5, rwshp, wtype.FlatWellBottom, 5.5, 5.5, 20.4, 1.4, "mm")
	//plate = wtype.NewLHPlate("pcrplate", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, 6.5)
	//plates[plate.Type] = plate
	plate = wtype.NewLHPlate("1L_DuranBottle", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, 15.5)
	plates[plate.Type] = plate

	return plates
}

func getPlateByType(typ string) *wtype.LHPlate {
	plates := makePlateLibrary()
	p := plates[typ]
	return p.Dup()
}

func getPlateList() []string {
	plates := makePlateLibrary()

	kz := make([]string, len(plates))
	x := 0
	for name, _ := range plates {
		kz[x] = name
		x += 1
	}
	return kz
}
