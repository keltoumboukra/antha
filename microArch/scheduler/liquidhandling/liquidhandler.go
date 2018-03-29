// liquidhandling/Liquidhandler.go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
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

package liquidhandling

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/microArch/driver"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/logger"
)

// the liquid handler structure defines the interface to a particular liquid handling
// platform. The structure holds the following items:
// - an LHRequest structure defining the characteristics of the platform
// - a channel for communicating with the liquid handler
// additionally three functions are defined to implement platform-specific
// implementation requirements
// in each case the LHRequest structure passed in has some additional information
// added and is then passed out. Features which are already defined (e.g. by the
// scheduler or the user) are respected as constraints and will be left unchanged.
// The three functions define
// - setup (SetupAgent): How sources are assigned to plates and plates to positions
// - layout (LayoutAgent): how experiments are assigned to outputs
// - execution (ExecutionPlanner): generates instructions to implement the required plan
//
// The general mechanism by which requests which refer to specific items as opposed to
// those which only state that an item of a particular kind is required is by the definition
// of an 'inst' tag in the request structure with a guid. If this is defined and valid
// it indicates that this item in the request (e.g. a plate, stock etc.) is a specific
// instance. If this is absent then the GUID will either be created or requested
//

// NB for flexibility we should not make the properties object part of this but rather
// send it in as an argument

type Liquidhandler struct {
	Properties       *liquidhandling.LHProperties
	FinalProperties  *liquidhandling.LHProperties
	SetupAgent       func(context.Context, *LHRequest, *liquidhandling.LHProperties) (*LHRequest, error)
	LayoutAgent      func(context.Context, *LHRequest, *liquidhandling.LHProperties) (*LHRequest, error)
	ExecutionPlanner func(context.Context, *LHRequest, *liquidhandling.LHProperties) (*LHRequest, error)
	PolicyManager    *LHPolicyManager
	plateIDMap       map[string]string // which plates are before / after versions
}

// initialize the liquid handling structure
func Init(properties *liquidhandling.LHProperties) *Liquidhandler {
	lh := Liquidhandler{}
	lh.SetupAgent = BasicSetupAgent
	lh.LayoutAgent = ImprovedLayoutAgent
	//lh.ExecutionPlanner = ImprovedExecutionPlanner
	lh.ExecutionPlanner = ExecutionPlanner3
	lh.Properties = properties
	lh.FinalProperties = properties
	lh.plateIDMap = make(map[string]string)
	return &lh
}

func (this *Liquidhandler) PlateIDMap() map[string]string {
	ret := make(map[string]string, len(this.plateIDMap))

	for k, v := range this.plateIDMap {
		ret[k] = v
	}

	return ret
}

// catch errors early
func ValidateRequest(request *LHRequest) error {
	if len(request.LHInstructions) == 0 {
		return wtype.LHError(wtype.LH_ERR_OTHER, "Nil plan requested: no Mix Instructions present")
	}

	// no component can have all three of Conc, Vol and TVol set to 0:

	for _, ins := range request.LHInstructions {
		for i, cmp := range ins.Components {
			if cmp.Vol == 0.0 && cmp.Conc == 0.0 && cmp.Tvol == 0.0 {
				errstr := fmt.Sprintf("Nil mix (no volume, concentration or total volume) requested: %d : ", i)

				for j := 0; j < len(ins.Components); j++ {
					ss := ins.Components[i].CName
					if j == i {
						ss = strings.ToUpper(ss)
					}

					if j != len(ins.Components)-1 {
						ss += ", "
					}

					errstr += ss
				}
				return wtype.LHError(wtype.LH_ERR_OTHER, errstr)
			}
		}
	}
	return nil
}

// high-level function which requests planning and execution for an incoming set of
// solutions
func (this *Liquidhandler) MakeSolutions(ctx context.Context, request *LHRequest) error {
	err := ValidateRequest(request)

	if err != nil {
		return err
	}

	request.ConfigureYourself()

	//f := func() {
	err = this.Plan(ctx, request)
	if err != nil {
		return err
	}

	err = this.Execute(request)

	if err != nil {
		return err
	}

	// output some info on the final setup

	OutputSetup(this.Properties)

	// and afer

	fmt.Println("SETUP AFTER: ")
	OutputSetup(this.FinalProperties)

	return nil
}

// run the request via the driver
func (this *Liquidhandler) Execute(request *LHRequest) error {

	if (*request).Instructions == nil {
		return wtype.LHError(wtype.LH_ERR_OTHER, "Cannot execute request: no instructions")
	}
	// add setup instructions to the request instruction stream

	this.add_setup_instructions(request)

	// set up the robot with extra calls not included in instructions

	err := this.do_setup(request)

	if err != nil {
		return err
	}

	instructions := (*request).Instructions

	// some timing info for the log (only) for now

	timer := this.Properties.GetTimer()
	var d time.Duration

	for _, ins := range instructions {

		if (*request).Options.PrintInstructions {
			fmt.Println(liquidhandling.InsToString(ins))
		}
		_, ok := ins.(liquidhandling.TerminalRobotInstruction)

		if !ok {
			fmt.Printf("ERROR: Got instruction ", liquidhandling.InsToString(ins), "which is wrong type")
			continue
		}

		err := ins.(liquidhandling.TerminalRobotInstruction).OutputTo(this.Properties.Driver)

		if err != nil {
			return wtype.LHError(wtype.LH_ERR_DRIV, err.Error())
		}

		// The graph view depends on the string generated in this step
		str := ""
		if ins.InstructionType() == liquidhandling.TFR {
			mocks := liquidhandling.MockAspDsp(ins)
			for _, ii := range mocks {
				str += liquidhandling.InsToString2(ii) + "\n"
			}
		} else {
			str = liquidhandling.InsToString2(ins) + "\n"
		}

		request.InstructionText += str

		//fmt.Println(liquidhandling.InsToString(ins))

		if timer != nil {
			d += timer.TimeFor(ins)
		}
	}

	logger.Debug(fmt.Sprintf("Total time estimate: %s", d.String()))
	request.TimeEstimate = d.Seconds()

	return nil
}

func (this *Liquidhandler) revise_volumes(rq *LHRequest) error {
	// XXX -- HARD CODE 8 here
	lastPlate := make([]string, 8)
	lastWell := make([]string, 8)

	vols := make(map[string]map[string]wunit.Volume)

	for _, ins := range rq.Instructions {
		if ins.InstructionType() == liquidhandling.MOV {
			lastPlate = make([]string, 8)
			lastPos := ins.GetParameter("POSTO").([]string)

			for i, p := range lastPos {
				lastPlate[i] = this.Properties.PosLookup[p]
			}

			lastWell = ins.GetParameter("WELLTO").([]string)
		} else if ins.InstructionType() == liquidhandling.ASP {
			for i, _ := range lastPlate {
				if i >= len(lastWell) {
					break
				}
				lp := lastPlate[i]
				lw := lastWell[i]

				if lp == "" {
					continue
				}

				ppp := this.Properties.PlateLookup[lp].(*wtype.LHPlate)

				lwl := ppp.Wellcoords[lw]

				if !lwl.IsAutoallocated() {
					continue
				}

				_, ok := vols[lp]

				if !ok {
					vols[lp] = make(map[string]wunit.Volume)
				}

				v, ok := vols[lp][lw]

				if !ok {
					v = wunit.NewVolume(0.0, "ul")
					vols[lp][lw] = v
				}
				//v.Add(ins.Volume[i])

				insvols := ins.GetParameter("VOLUME").([]wunit.Volume)
				v.Add(insvols[i])
				// double add of carry volume here?
				v.Add(rq.CarryVolume)

			}
		} else if ins.InstructionType() == liquidhandling.TFR {
			tfr := ins.(*liquidhandling.TransferInstruction)
			for _, mtf := range tfr.Transfers {
				for _, tf := range mtf.Transfers {
					lpos, lw := tf.PltFrom, tf.WellFrom

					lp := this.Properties.PosLookup[lpos]
					ppp := this.Properties.PlateLookup[lp].(*wtype.LHPlate)
					lwl := ppp.Wellcoords[lw]

					if !lwl.IsAutoallocated() {
						continue
					}

					_, ok := vols[lp]

					if !ok {
						vols[lp] = make(map[string]wunit.Volume)
					}

					v, ok := vols[lp][lw]

					if !ok {
						v = wunit.NewVolume(0.0, "ul")
						vols[lp][lw] = v
					}
					//v.Add(ins.Volume[i])

					v.Add(tf.Volume)
				}
			}
		}
	}

	// apply evaporation
	for _, vc := range rq.Evaps {
		loctox := strings.Split(vc.Location, ":")

		// ignore anything where the location isn't properly set

		if len(loctox) < 2 {
			continue
		}

		plateID := loctox[0]
		wellcrds := loctox[1]

		wellmap, ok := vols[plateID]

		if !ok {
			continue
		}

		vol := wellmap[wellcrds]
		vol.Add(vc.Volume)
	}

	// now go through and set the plates up appropriately

	for plateID, wellmap := range vols {
		plate, ok := this.FinalProperties.Plates[this.Properties.PlateIDLookup[plateID]]
		plate2, _ := this.Properties.Plates[this.Properties.PlateIDLookup[plateID]]

		if !ok {
			err := wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprint("NO SUCH PLATE: ", plateID))
			return err
		}

		// what's it like here?

		for crd, unroundedvol := range wellmap {
			rv, _ := wutil.Roundto(unroundedvol.RawValue(), 1)
			vol := wunit.NewVolume(rv, unroundedvol.Unit().PrefixedSymbol())
			well := plate.Wellcoords[crd]
			well2 := plate2.Wellcoords[crd]

			if well.IsAutoallocated() {
				vol.Add(well.ResidualVolume())
				well2.WContents.SetVolume(vol)
				well.WContents.SetVolume(well.ResidualVolume())
				well.WContents.ID = wtype.GetUUID()
				well.DeclareNotTemporary()
				well2.DeclareNotTemporary()
			}
		}
	}

	// finally get rid of any temporary stuff

	this.Properties.RemoveUnusedAutoallocatedComponents()
	this.FinalProperties.RemoveUnusedAutoallocatedComponents()

	pidm := make(map[string]string, len(this.Properties.Plates))
	for pos, _ := range this.Properties.Plates {
		p1, ok1 := this.Properties.Plates[pos]
		p2, ok2 := this.FinalProperties.Plates[pos]

		if (!ok1 && ok2) || (ok1 && !ok2) {

			if ok1 {
				fmt.Println("BEFORE HAS: ", p1)
			}

			if ok2 {
				fmt.Println("AFTER  HAS: ", p2)
			}

			return (wtype.LHError(8, fmt.Sprintf("Plate disappeared from position %s", pos)))
		}

		if !(ok1 && ok2) {
			continue
		}

		this.plateIDMap[p1.ID] = p2.ID
		pidm[p2.ID] = p1.ID
	}

	// this is many shades of wrong but likely to save us a lot of time
	for _, pos := range this.Properties.InputSearchPreferences() {
		p1, ok1 := this.Properties.Plates[pos]
		p2, ok2 := this.FinalProperties.Plates[pos]

		if ok1 && ok2 {
			for _, wa := range p1.Cols {
				for _, w := range wa {
					// copy the outputs to the correct side
					// and remove the outputs from the initial state
					if !w.Empty() {
						w2, ok := p2.Wellcoords[w.Crds]
						if ok {
							// there's no strict separation between outputs and
							// inputs here
							if w.IsAutoallocated() {
								continue
							} else if w.IsUserAllocated() {
								// swap old and new
								c := w.WContents.Dup()
								w.Clear()
								c2 := w2.WContents.Dup()
								w2.Clear()
								w.Add(c2)
								w2.Add(c)
							} else {
								// replace
								w2.Clear()
								w2.Add(w.WContents)
								w.Clear()
							}
						}
					}
				}

			}

			//fmt.Println(p2, " ", p1)
			//fmt.Println("Plate ID Map: ", p2.ID, " --> ", p1.ID)

			//	this.plateIDMap[p2.ID] = p1.ID
		}
	}

	// all done

	return nil
}

func (this *Liquidhandler) add_setup_instructions(rq *LHRequest) {
	instructions := make([]liquidhandling.TerminalRobotInstruction, 0, len(rq.Instructions)+10)

	instructions = append(instructions, liquidhandling.NewRemoveAllPlatesInstruction())

	for position, plateid := range this.Properties.PosLookup {
		if plateid == "" {
			continue
		}
		plate := this.Properties.PlateLookup[plateid]
		name := plate.(wtype.Named).GetName()

		ins := liquidhandling.NewAddPlateToInstruction(position, name, plate)

		instructions = append(instructions, ins)
	}
	instructions = append(instructions, rq.Instructions...)
	rq.Instructions = instructions
}

func (this *Liquidhandler) do_setup(rq *LHRequest) error {
	/*
		stat := this.Properties.Driver.RemoveAllPlates()

		if stat.Errorcode == driver.ERR {
			return wtype.LHError(wtype.LH_ERR_DRIV, stat.Msg)
		}

		for position, plateid := range this.Properties.PosLookup {
			if plateid == "" {
				continue
			}
			plate := this.Properties.PlateLookup[plateid]
			name := plate.(wtype.Named).GetName()

			stat = this.Properties.Driver.AddPlateTo(position, plate, name)
			if stat.Errorcode == driver.ERR {
				return wtype.LHError(wtype.LH_ERR_DRIV, stat.Msg)
			}
		}
	*/

	if _, ok := this.Properties.Driver.(liquidhandling.LowLevelLiquidhandlingDriver); ok {
		stat := this.Properties.Driver.(liquidhandling.LowLevelLiquidhandlingDriver).UpdateMetaData(this.Properties)
		if stat.Errorcode == driver.ERR {
			return wtype.LHError(wtype.LH_ERR_DRIV, stat.Msg)
		}
	}

	return nil
}

// This runs the following steps in order:
// - determine required inputs
// - request inputs	--- should be moved out
// - define robot setup
// - define output layout
// - generate the robot instructions
// - request consumables and other device setups e.g. heater setting
//
// as described above, steps only have an effect if the required inputs are
// not defined beforehand
//
// so essentially the idea is to parameterise all requests to liquid handlers
// using a Command structure called LHRequest
//
// Depending on its state of completeness, the request structure may be executable
// immediately or may need some additional definition. The purpose of the liquid
// handling service is to provide methods to invoke when parts of the request need
// further definition.
//
// when running a request we should be able to provide mechanisms for pushing requests
// back into the queue to allow them to be cached
//
// this should be OK since the LHRequest parameterises all state including instructions
// for asynchronous drivers we have to determine how far the program got before it was
// paused, which should be tricky but possible.
//

func checkSanityIns(request *LHRequest) {
	// check instructions for basic sanity

	good := true
	for _, ins := range request.LHInstructions {
		if ins.Type == wtype.LHIMIX {
			v := wunit.NewVolume(0.0, "ul")
			tv := wunit.NewVolume(0.0, "ul")
			for _, c := range ins.Components {
				// need to be a bit careful but...

				if c.Vol < 0.0 {
					fmt.Println("NEGATIVE VOLUME!!!! ", c.CName, " ", c.Vol)
					good = false
					continue
				}

				if c.Vol != 0.0 {
					v.Add(c.Volume())
				} else if c.Tvol != 0.0 {
					if !tv.IsZero() && !tv.EqualTo(c.TotalVolume()) {
						fmt.Println("ERROR: MULTIPLE DISTINCT TOTAL VOLUMES SPECIFIED FOR ", ins.ID, " ", ins.Results[0].CName, " COMPONENT ", c)
						good = false
					}

					tv = c.TotalVolume()
				}
			}

			if tv.IsZero() && !v.EqualTo(ins.Results[0].Volume()) {
				fmt.Println("OH DEAR DEAR DEAR: VOLUME INCONSISTENCY FOR ", ins.ID, " ", ins.Results[0].CName, " COMP: ", v, " PROD: ", ins.Results[0].Volume())
				good = false
			} else if !tv.IsZero() && !tv.EqualTo(ins.Results[0].Volume()) {
				fmt.Println("ERROR: VOLUME INCONSISTENCY FOR ", ins.ID, " ", ins.Results[0].CName, " COMP: ", tv, " PROD: ", ins.Results[0].Volume())
				good = false
			} else if ins.PlateID != "" {
				// compare result volume to the well volume

				plat, ok := request.GetPlate(ins.PlateID)

				if !ok {
					// possibly an issue
				} else if plat.Welltype.MaxVolume().LessThan(ins.Results[0].Volume()) {
					fmt.Println("WARNING: EXCESS VOLUME REQUIRED FOR ", ins.ID, " ", ins.Results[0].CName, " WANT: ", ins.Results[0].Volume(), " MAX FOR PLATE OF TYPE ", plat.Type, ": ", plat.Welltype.MaxVolume())
					//good = false
				}
			}
		}
	}

	if !good {
		panic("URGH - volume issues here")
	}

}

func checkInstructionOrdering(request *LHRequest) {
	ch := request.InstructionChain

	for {
		if ch == nil {
			break
		}

		onlyAllowOneInstructionType(ch)

		ch = ch.Child
	}
}

func onlyAllowOneInstructionType(c *IChain) {
	m := make(map[string]bool)
	inss := c.Values

	for _, i := range inss {
		m[i.InsType()] = true
	}

	if len(m) != 1 {
		panic(fmt.Errorf("Only one instruction type per stage is allowed, found %v at stage %d", m, c.Depth))
	}
}

func checkDestinationSanity(request *LHRequest) {
	for _, ins := range request.LHInstructions {
		// non-mix instructions are fine
		if ins.Type != wtype.LHIMIX {
			continue
		}

		if ins.PlateID == "" || ins.Platetype == "" || ins.Welladdress == "" {
			fmt.Println("INS ", ins, " NOT WELL FORMED: HAS PlateID ", ins.PlateID != "", " HAS platetype ", ins.Platetype != "", " HAS WELLADDRESS ", ins.Welladdress != "")
			panic(fmt.Errorf("After layout all mix instructions must have plate IDs, plate types and well addresses"))
		}
	}
}

func anotherSanityCheck(request *LHRequest) {
	p := map[*wtype.LHComponent]*wtype.LHInstruction{}

	for _, ins := range request.LHInstructions {
		// we must not share pointers

		for _, c := range ins.Components {
			ins2, ok := p[c]
			if ok {
				panic(fmt.Sprintf("POINTER REUSE: Instructions %s %s for component %s %s", ins.ID, ins2.ID, c.ID, c.CName))
			}

			p[c] = ins
		}

		ins2, ok := p[ins.Results[0]]

		if ok {
			panic(fmt.Sprintf("POINTER REUSE: Instructions %s %s for component %s %s", ins.ID, ins2.ID, ins.Results[0].ID, ins.Results[0].CName))
		}

		p[ins.Results[0]] = ins
	}
}

func forceSanity(request *LHRequest) {
	for _, ins := range request.LHInstructions {
		for i := 0; i < len(ins.Components); i++ {
			ins.Components[i] = ins.Components[i].Dup()
		}

		ins.Results[0] = ins.Results[0].Dup()
	}
}

func (this *Liquidhandler) Plan(ctx context.Context, request *LHRequest) error {
	// figure out the output order

	err := set_output_order(request)

	if err != nil {
		return err
	}

	if request.Options.PrintInstructions {
		for _, insID := range request.Output_order {
			ins := request.LHInstructions[insID]
			fmt.Print(ins.InsType(), " G:", ins.Generation(), " ", ins.ID, " ", wtype.ComponentVector(ins.Components), " ", ins.PlateName, " ID(", ins.PlateID, ") ", ins.Welladdress, ": ", ins.ProductIDs())

			if ins.IsMixInPlace() {
				fmt.Print(" INPLACE")
			}

			fmt.Println()
		}
		request.InstructionChain.Print()
	}

	// assert we should have some instruction ordering

	if len(request.Output_order) == 0 {
		return fmt.Errorf("Error with instruction sorting: Have %d want %d instructions", len(request.Output_order), len(request.LHInstructions))
	}

	// assert that we must keep prompts separate from mixes

	checkInstructionOrdering(request)

	forceSanity(request)
	// convert requests to volumes and determine required stock concentrations
	checkSanityIns(request)
	instructions, stockconcs, err := solution_setup(request, this.Properties)

	if err != nil {
		return err
	}

	checkSanityIns(request)

	request.LHInstructions = instructions
	request.Stockconcs = stockconcs

	// set up the mapping of the outputs
	// tried moving here to see if we can use results in fixVolumes
	request, err = this.Layout(ctx, request)

	if err != nil {
		return err
	}
	forceSanity(request)
	anotherSanityCheck(request)

	// assert: all instructions should now be assigned specific plate IDs, types and wells
	checkDestinationSanity(request)

	if request.Options.FixVolumes {
		// see if volumes can be corrected
		request, err = FixVolumes(request)

		if err != nil {
			return err
		}
		if request.Options.PrintInstructions {
			fmt.Println("")
			fmt.Println("POST VOLUME FIX")
			fmt.Println("")
			for _, insID := range request.Output_order {
				ins := request.LHInstructions[insID]
				fmt.Print(ins.InsType(), " G:", ins.Generation(), " ", ins.ID, " ", wtype.ComponentVector(ins.Components), " ", ins.PlateName, " ID(", ins.PlateID, ") ", ins.Welladdress, ": ", ins.ProductIDs())

				if ins.IsMixInPlace() {
					fmt.Print(" INPLACE")
				}

				fmt.Println()
			}
		}
	}
	checkSanityIns(request)

	// looks at components, determines what inputs are required
	request, err = this.GetInputs(request)

	if err != nil {
		return err
	}
	// define the input plates
	// should be merged with the above
	request, err = input_plate_setup(ctx, request)

	if err != nil {
		return err
	}

	// next we need to determine the liquid handler setup
	request, err = this.Setup(ctx, request)
	if err != nil {
		return err
	}

	// remove dummy mix-in-place instructions

	request = removeDummyInstructions(request)

	// now make instructions
	request, err = this.ExecutionPlan(ctx, request)

	if err != nil {
		return err
	}

	// counts tips used in this run -- reads instructions generated above so must happen
	// after execution planning
	request, err = this.countTipsUsed(request)

	if err != nil {
		return err
	}

	// Ensures tip boxes and wastes are correct for initial and final robot states
	this.Refresh_tipboxes_tipwastes(request)

	// revise the volumes - this makes sure the volumes requested are correct
	err = this.revise_volumes(request)

	if err != nil {
		return err
	}
	// ensure the after state is correct
	this.fix_post_ids()
	err = this.fix_post_names(request)

	if err != nil {
		return err
	}

	return nil
}

// resolve question of where something is requested to go
const NoID = "NOID"
const NoName = "NONAME"
const NoWell = "NOWELL"

func assembleLoc(ins *wtype.LHInstruction) string {
	id := NoID
	if ins.PlateID != "" {
		id = ins.PlateID
	}

	name := NoName

	if ins.PlateName != "" {
		name = ins.PlateName
	}

	well := NoWell

	if ins.Welladdress != "" {
		well = ins.Welladdress
	}

	return strings.Join([]string{id, name, well}, ":")
}

// sort out inputs
func (this *Liquidhandler) GetInputs(request *LHRequest) (*LHRequest, error) {
	instructions := (*request).LHInstructions

	inputs := make(map[string][]*wtype.LHComponent, 3)
	vmap := make(map[string]wunit.Volume)

	allinputs := make([]string, 0, 10)

	ordH := make(map[string]int, len(instructions))

	inPlaceLocations := make(map[string]string, len(instructions))

	//	for _, instruction := range instructions {
	for _, insID := range request.Output_order {
		// ignore non-mixes

		instruction := instructions[insID]

		if instruction.InsType() != "MIX" {
			continue
		}

		components := instruction.Components

		for ix, component := range components {
			// Ignore components which already exist

			if component.IsInstance() {
				continue
			}

			// what if this is a mix in place?
			if ix == 0 && !component.IsSample() {
				// these components come in as instances -- hence 1 per well
				// but if not allocated we need to do so
				inputs[component.CNID()] = make([]*wtype.LHComponent, 0, 1)
				inputs[component.CNID()] = append(inputs[component.CNID()], component)
				allinputs = append(allinputs, component.CNID())
				vmap[component.CNID()] = component.Volume()
				component.DeclareInstance()

				// if this already exists do nothing
				_, ok := ordH[component.CNID()]

				if !ok {
					ordH[component.CNID()] = len(ordH)
					// assign like this: ID:NAME:WELL
					// if ID is blank we call it NOID
					loc := assembleLoc(instruction)
					inPlaceLocations[component.CNID()] = loc
				}
			} else {
				cmps, ok := inputs[component.Kind()]
				if !ok {
					cmps = make([]*wtype.LHComponent, 0, 3)
					allinputs = append(allinputs, component.Kind())
				}

				_, ok = ordH[component.Kind()]

				if !ok {
					ordH[component.Kind()] = len(ordH)
				}

				cmps = append(cmps, component)
				inputs[component.Kind()] = cmps

				// similarly add the volumes up

				vol := vmap[component.Kind()]

				if vol.IsNil() {
					vol = wunit.NewVolume(0.0, "ul")
				}

				v2a := wunit.NewVolume(component.Vol, component.Vunit)

				// we have to add the carry volume here
				// this is roughly per transfer so should be OK
				v2a.Add(request.CarryVolume)
				vol.Add(v2a)

				vmap[component.Kind()] = vol
			}
		}
	}

	// work out how much we have and how much we need
	// need to consider what to do with IDs

	// invert the Hash

	var err error
	(*request).Input_order, err = OrdinalFromHash(ordH)

	if err != nil {
		return request, err
	}

	var requestinputs map[string][]*wtype.LHComponent
	requestinputs = request.Input_solutions

	if len(requestinputs) == 0 {
		requestinputs = make(map[string][]*wtype.LHComponent, 5)
	}

	vmap2 := make(map[string]wunit.Volume, len(vmap))
	vmap3 := make(map[string]wunit.Volume, len(vmap))

	for _, k := range allinputs {
		// vola: how much comes in
		ar := requestinputs[k]
		vola := wunit.NewVolume(0.00, "ul")
		for _, cmp := range ar {
			vold := wunit.NewVolume(cmp.Vol, cmp.Vunit)
			vola.Add(vold)
		}
		// volb: how much we asked for
		volb := vmap[k].Dup()
		volb.Subtract(vola)
		vmap2[k] = vola

		if volb.GreaterThanFloat(0.0001) {
			vmap3[k] = volb
		}
		// toggle HERE for DEBUG
		if false {
			volc := vmap[k]
			logger.Debug(fmt.Sprint("COMPONENT ", k, " HAVE : ", vola.ToString(), " WANT: ", volc.ToString(), " DIFF: ", volb.ToString()))
		}
	}

	(*request).Input_vols_required = vmap
	(*request).Input_vols_supplied = vmap2
	(*request).Input_vols_wanting = vmap3

	// add any new inputs

	for k, v := range inputs {
		if requestinputs[k] == nil {
			requestinputs[k] = v
		}
	}

	(*request).Input_solutions = requestinputs

	return request, nil
}

func OrdinalFromHash(m map[string]int) ([]string, error) {
	s := make([]string, len(m))

	// no collisions allowed!

	for k, v := range m {
		if s[v] != "" {
			return nil, fmt.Errorf("Error: ordinal %d appears twice!", v)
		}

		s[v] = k
	}

	return s, nil
}

// define which labware to use
func (this *Liquidhandler) GetPlates(ctx context.Context, plates map[string]*wtype.LHPlate, major_layouts map[int][]string, ptype *wtype.LHPlate) (map[string]*wtype.LHPlate, error) {
	if plates == nil {
		plates = make(map[string]*wtype.LHPlate, len(major_layouts))

		// assign new plates
		for i := 0; i < len(major_layouts); i++ {
			//newplate := wtype.New_Plate(ptype)
			newplate, err := inventory.NewPlate(ctx, ptype.Type)
			if err != nil {
				return nil, err
			}
			plates[newplate.ID] = newplate
		}
	}

	// we should know how many plates we need
	for k, plate := range plates {
		if plate.Inst == "" {
			//stockrequest := execution.GetContext().StockMgr.RequestStock(makePlateStockRequest(plate))
			//plate.Inst = stockrequest["inst"].(string)
		}

		plates[k] = plate
	}

	return plates, nil
}

// generate setup for the robot
func (this *Liquidhandler) Setup(ctx context.Context, request *LHRequest) (*LHRequest, error) {
	// assign the plates to positions
	// this needs to be parameterizable
	return this.SetupAgent(ctx, request, this.Properties)
}

// generate the output layout
func (this *Liquidhandler) Layout(ctx context.Context, request *LHRequest) (*LHRequest, error) {
	// assign the results to destinations
	// again needs to be parameterized

	return this.LayoutAgent(ctx, request, this.Properties)
}

// make the instructions for executing this request
func (this *Liquidhandler) ExecutionPlan(ctx context.Context, request *LHRequest) (*LHRequest, error) {
	// necessary??
	this.FinalProperties = this.Properties.Dup()
	temprobot := this.Properties.Dup()
	//saved_plates := this.Properties.SaveUserPlates()

	var rq *LHRequest
	var err error

	if request.Options.ExecutionPlannerVersion == "ep3" {
		rq, err = ExecutionPlanner3(ctx, request, this.Properties)
	} else {
		rq, err = this.ExecutionPlanner(ctx, request, this.Properties)
	}

	this.FinalProperties = temprobot

	//this.Properties.RestoreUserPlates(saved_plates)

	return rq, err
}

func OutputSetup(robot *liquidhandling.LHProperties) {
	logger.Debug("DECK SETUP INFO")
	logger.Debug("Tipboxes: ")

	for k, v := range robot.Tipboxes {
		logger.Debug(fmt.Sprintf("%s %s: %s", k, robot.PlateIDLookup[k], v.Type))
	}

	logger.Debug("Plates:")

	for k, v := range robot.Plates {

		logger.Debug(fmt.Sprintf("%s %s: %s %s", k, robot.PlateIDLookup[k], v.PlateName, v.Type))

		//TODO Deprecate
		if strings.Contains(v.GetName(), "Input") {
			wtype.AutoExportPlateCSV(v.GetName()+".csv", v)
		}

		v.OutputLayout()
	}

	logger.Debug("Tipwastes: ")

	for k, v := range robot.Tipwastes {
		logger.Debug(fmt.Sprintf("%s %s: %s capacity %d", k, robot.PlateIDLookup[k], v.Type, v.Capacity))
	}

}

//ugly
func (lh *Liquidhandler) fix_post_ids() {
	for _, p := range lh.FinalProperties.Plates {
		for _, w := range p.Wellcoords {
			if w.IsUserAllocated() {
				w.WContents.ID = wtype.GetUUID()
			}
		}
	}
}

func (lh *Liquidhandler) fix_post_names(rq *LHRequest) error {
	// Instructions updating a well
	assignment := make(map[*wtype.LHWell]*wtype.LHInstruction)
	for _, inst := range rq.LHInstructions {
		// ignore non -mix instructions
		if inst.Type != wtype.LHIMIX {
			continue
		}

		tx := strings.Split(inst.Results[0].Loc, ":")
		newid, ok := lh.plateIDMap[tx[0]]
		if !ok {
			return wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprintf("No output plate mapped to %s", tx[0]))
		}

		ip, ok := lh.FinalProperties.PlateLookup[newid]
		if !ok {
			return wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprintf("No output plate %s", newid))
		}

		p, ok := ip.(*wtype.LHPlate)
		if !ok {
			return wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprintf("Got %s, should have *wtype.LHPlate", reflect.TypeOf(ip)))
		}

		well, ok := p.Wellcoords[tx[1]]
		if !ok {
			return wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprintf("No well %s on plate %s", tx[1], tx[0]))
		}

		oldInst := assignment[well]
		if oldInst == nil {
			assignment[well] = inst
		} else if prev, cur := oldInst.Results[0].Generation(), inst.Results[0].Generation(); prev < cur {
			assignment[well] = inst
		}
	}

	for well, inst := range assignment {
		well.WContents.CName = inst.Results[0].CName
	}

	return nil
}

func dummy(ins *wtype.LHInstruction) bool {
	if wtype.InsType(ins.Type) == "MIX" && ins.IsMixInPlace() && len(ins.Components) == 1 {
		// instructions of this form generally mean "do nothing"
		// but have very useful side-effects
		return true
	}

	return false
}

func removeDummyInstructions(rq *LHRequest) *LHRequest {
	toRemove := make(map[string]bool, len(rq.LHInstructions))
	for _, ins := range rq.LHInstructions {
		if dummy(ins) {
			toRemove[ins.ID] = true
		}
	}

	if len(toRemove) == 0 {
		//no dummies
		return rq
	}

	oo := make([]string, 0, len(rq.Output_order)-len(toRemove))

	for _, ins := range rq.Output_order {
		if toRemove[ins] {
			continue
		} else {
			oo = append(oo, ins)
		}
	}

	if len(oo) != len(rq.Output_order)-len(toRemove) {
		panic(fmt.Sprintf("Dummy instruction prune failed: before %d dummies %d after %d", len(rq.Output_order), len(toRemove), len(oo)))
	}

	rq.Output_order = oo

	// prune instructionChain

	rq.InstructionChain.PruneOut(toRemove)

	return rq
}

func (req *LHRequest) MergedInputOutputPlates() map[string]*wtype.LHPlate {
	m := make(map[string]*wtype.LHPlate, len(req.Input_plates)+len(req.Output_plates))
	addToMap(m, req.Input_plates)
	addToMap(m, req.Output_plates)
	return m
}

func addToMap(m, a map[string]*wtype.LHPlate) {
	for k, v := range a {
		m[k] = v
	}
}
