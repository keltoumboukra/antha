package mixer

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/inventory"
	driver "github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/sampletracker"
	planner "github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/human"
)

var (
	_ target.Device = &Mixer{}
)

// A Mixer is a device plugin for mixer devices
type Mixer struct {
	driver     driver.LiquidhandlingDriver
	properties *driver.LHProperties // Prototype to create fresh properties
	opt        Opt
}

func (a *Mixer) String() string {
	return "Mixer"
}

// CanCompile implements a Device
func (a *Mixer) CanCompile(req ast.Request) bool {
	// TODO: Add specific volume constraints
	can := ast.Request{
		Selector: []ast.NameValue{
			target.DriverSelectorV1Mixer,
			target.DriverSelectorV1Prompter,
		},
	}
	if a.properties.CanPrompt() {
		can.Selector = append(can.Selector, target.DriverSelectorV1Prompter)
	}
	return can.Contains(req)
}

// MoveCost implements a Device
func (a *Mixer) MoveCost(from target.Device) int {
	if from == a {
		return 0
	}
	return human.HumanByXCost + 1
}

// FileType returns the file type for generated files
func (a *Mixer) FileType() (ftype string) {
	if m := a.properties.Mnfr; len(m) != 0 {
		ftype = fmt.Sprintf("application/%s", strings.ToLower(m))
	}
	return
}

type lhreq struct {
	*planner.LHRequest     // A request
	*driver.LHProperties   // ... its state
	*planner.Liquidhandler // ... and its associated planner
}

func (a *Mixer) makeLhreq(ctx context.Context) (*lhreq, error) {
	// MIS -- this might be a hole. We probably need to invoke the sample tracker here
	addPlate := func(req *planner.LHRequest, ip *wtype.LHPlate) error {
		if _, seen := req.Input_plates[ip.ID]; seen {
			return fmt.Errorf("plate %q already added", ip.ID)
		}
		//req.Input_plates[ip.ID] = ip
		req.AddUserPlate(ip)
		return nil
	}

	req := planner.NewLHRequest()

	/// TODO --> a.opt.Destination isn't being passed through, this makes MixInto redundant

	if err := req.Policies.SetOption("USE_DRIVER_TIP_TRACKING", a.opt.UseDriverTipTracking); err != nil {
		return nil, err
	}
	if err := req.Policies.SetOption("USE_LLF", a.opt.UseLLF); err != nil {
		return nil, err
	}

	prop := a.properties.Dup()
	prop.Driver = a.properties.Driver
	plan := planner.Init(prop)

	if p := a.opt.MaxPlates; p != nil {
		req.Input_setup_weights["MAX_N_PLATES"] = *p
	}

	if p := a.opt.MaxWells; p != nil {
		req.Input_setup_weights["MAX_N_WELLS"] = *p
	}

	if p := a.opt.ResidualVolumeWeight; p != nil {
		req.Input_setup_weights["RESIDUAL_VOLUME_WEIGHT"] = *p
	}

	// TODO -- error check here to prevent nil values

	if p := a.opt.InputPlateTypes; len(p) != 0 {
		for _, v := range p {
			p, err := inventory.NewPlate(ctx, v)
			if err != nil {
				return nil, err
			}

			req.Input_platetypes = append(req.Input_platetypes, p)
		}
	}

	if p := a.opt.OutputPlateTypes; len(p) != 0 {
		for _, v := range p {
			p, err := inventory.NewPlate(ctx, v)
			if err != nil {
				return nil, err
			}
			req.Output_platetypes = append(req.Output_platetypes, p)
		}
	}

	if p := a.opt.TipTypes; len(p) != 0 {
		for _, v := range p {
			t, err := inventory.NewTipbox(ctx, v)
			if err != nil {
				return nil, err
			}
			req.Tips = append(req.Tips, t)
		}
	}

	if p := a.opt.InputPlateData; len(p) != 0 {
		for idx, bs := range p {
			buf := bytes.NewBuffer(bs)
			r, err := ParsePlateCSV(ctx, buf)
			if err != nil {
				return nil, fmt.Errorf("cannot parse data at idx %d: %s", idx, err)
			}

			if len(r.Warnings) != 0 {
				return nil, fmt.Errorf("cannot parse data at idx %d: %s", idx, strings.Join(r.Warnings, " "))
			}

			if err := addPlate(req, r.Plate); err != nil {
				return nil, err
			}
		}
	}

	if ips := a.opt.InputPlates; len(ips) != 0 {
		for _, ip := range ips {
			if err := addPlate(req, ip); err != nil {
				return nil, err
			}
		}
	}

	// add plates requested via protocol

	st := sampletracker.GetSampleTracker()

	parr := st.GetInputPlates()

	for _, p := range parr {
		if err := addPlate(req, p); err != nil {
			return nil, err
		}

	}

	// try to do better multichannel execution planning?

	req.Options.ExecutionPlannerVersion = a.opt.PlanningVersion

	// print instructions?

	req.Options.PrintInstructions = a.opt.PrintInstructions

	// model evaporation?

	req.Options.ModelEvaporation = a.opt.ModelEvaporation

	// deal with output sorting

	req.Options.OutputSort = a.opt.OutputSort

	// LiquidLevelFollowing

	req.Options.UseLLF = a.opt.UseLLF

	// legacy volume use

	req.Options.LegacyVolume = a.opt.LegacyVolume

	// volume fix

	req.Options.FixVolumes = a.opt.FixVolumes

	return &lhreq{
		LHRequest:     req,
		LHProperties:  prop,
		Liquidhandler: plan,
	}, nil
}

// Compile implements a Device
func (a *Mixer) Compile(ctx context.Context, nodes []ast.Node) ([]target.Inst, error) {
	var mixes []*wtype.LHInstruction
	for _, node := range nodes {
		if c, ok := node.(*ast.Command); !ok {
			return nil, fmt.Errorf("cannot compile %T", node)
		} else if m, ok := c.Inst.(*wtype.LHInstruction); !ok {
			return nil, fmt.Errorf("cannot compile %T", c.Inst)
		} else {
			mixes = append(mixes, m)
		}
	}

	mix, err := a.makeMix(ctx, mixes)
	if err != nil {
		return nil, err
	}

	return target.SequentialOrder(mix), nil
}

func (a *Mixer) saveFile(name string) ([]byte, error) {
	data, status := a.driver.GetOutputFile()
	if !status.OK {
		return nil, fmt.Errorf("%d: %s", status.Errorcode, status.Msg)
	} else if len(data) == 0 {
		return nil, nil
	}

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)
	bs := []byte(data)

	if err := tw.WriteHeader(&tar.Header{
		Name:    name,
		Mode:    0644,
		Size:    int64(len(bs)),
		ModTime: time.Now(),
	}); err != nil {
		return nil, err
	} else if _, err := tw.Write(bs); err != nil {
		return nil, err
	} else if err := tw.Close(); err != nil {
		return nil, err
	} else if err := gw.Close(); err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

func (a *Mixer) makeMix(ctx context.Context, mixes []*wtype.LHInstruction) (*target.Mix, error) {
	hasPlate := func(plates []*wtype.LHPlate, typ, id string) bool {
		for _, p := range plates {
			if p.Type == typ && (id == "" || p.ID == id) {
				return true
			}
		}
		return false
	}

	getID := func(mixes []*wtype.LHInstruction) (r wtype.BlockID) {
		m := make(map[wtype.BlockID]bool)
		for _, mix := range mixes {
			m[mix.BlockID] = true
		}
		for k := range m {
			r = k
			break
		}
		return
	}

	r, err := a.makeLhreq(ctx)
	if err != nil {
		return nil, err
	}

	for _, m := range mixes {
		if m.OutPlate != nil {
			p, ok := r.LHRequest.Output_plates[m.OutPlate.ID]
			if ok && p != m.OutPlate {
				return nil, fmt.Errorf("Mix setup error: Plate %s already requested in different state", p.ID)
			}
			r.LHRequest.Output_plates[m.OutPlate.ID] = m.OutPlate
		}
	}

	r.LHRequest.BlockID = getID(mixes)

	for _, mix := range mixes {
		if len(mix.Platetype) != 0 && !hasPlate(r.LHRequest.Output_platetypes, mix.Platetype, mix.PlateID) {
			p, err := inventory.NewPlate(ctx, mix.Platetype)
			if err != nil {
				return nil, err
			}
			p.ID = mix.PlateID
			r.LHRequest.Output_platetypes = append(r.LHRequest.Output_platetypes, p)
		}
		r.LHRequest.Add_instruction(mix)
	}

	err = r.Liquidhandler.MakeSolutions(ctx, r.LHRequest)
	// TODO: MIS unfortunately we need to make sure this stays up to date would
	// be better to remove this and just use the ones the liquid handler holds
	r.LHProperties = r.Liquidhandler.Properties

	if err != nil {
		return nil, err
	}

	name := a.opt.DriverOutputFileName
	if len(name) == 0 {
		// TODO: Desired filename not exposed in current driver interface, so pick
		// a name. So far, at least Gilson software cares what the filename is, so
		// use .sqlite for compatibility
		name = strings.Replace(fmt.Sprintf("%s.sqlite", time.Now().Format(time.RFC3339)), ":", "_", -1)
	}

	tarball, err := a.saveFile(name)
	if err != nil {
		return nil, err
	}

	return &target.Mix{
		Dev:             a,
		Request:         r.LHRequest,
		Properties:      r.LHProperties,
		FinalProperties: r.Liquidhandler.FinalProperties,
		Final:           r.Liquidhandler.PlateIDMap(),
		Files: target.Files{
			Tarball: tarball,
			Type:    a.FileType(),
		},
	}, nil
}

// New creates a new Mixer
func New(opt Opt, d driver.LiquidhandlingDriver) (*Mixer, error) {
	p, status := d.GetCapabilities()
	if !status.OK {
		return nil, fmt.Errorf("cannot get capabilities: %s", status.Msg)
	}

	update := func(addr *[]string, v []string) {
		if len(v) != 0 {
			*addr = v
		}
	}

	if len(opt.DriverSpecificInputPreferences) != 0 && p.CheckPreferenceCompatibility(opt.DriverSpecificInputPreferences) {
		update(&p.Input_preferences, opt.DriverSpecificInputPreferences)
	}
	if len(opt.DriverSpecificOutputPreferences) != 0 && p.CheckPreferenceCompatibility(opt.DriverSpecificOutputPreferences) {
		update(&p.Output_preferences, opt.DriverSpecificOutputPreferences)
	}

	if len(opt.DriverSpecificTipPreferences) != 0 && p.CheckTipPrefCompatibility(opt.DriverSpecificTipPreferences) {
		update(&p.Tip_preferences, opt.DriverSpecificTipPreferences)
	}

	if len(opt.DriverSpecificTipWastePreferences) != 0 && p.CheckPreferenceCompatibility(opt.DriverSpecificTipWastePreferences) {
		update(&p.Tipwaste_preferences, opt.DriverSpecificTipWastePreferences)
	}
	if len(opt.DriverSpecificWashPreferences) != 0 && p.CheckPreferenceCompatibility(opt.DriverSpecificWashPreferences) {
		update(&p.Wash_preferences, opt.DriverSpecificWashPreferences)
	}
	p.Driver = d
	return &Mixer{driver: d, properties: &p, opt: opt}, nil
}
