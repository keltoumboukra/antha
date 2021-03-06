package liquidhandling

import (
	"context"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/testinventory"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

func TestInputSampleAutoAllocate(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	rbt := makeGilson(ctx)
	rq := NewLHRequest()

	cmp1, err := inventory.NewComponent(ctx, inventory.WaterType)
	if err != nil {
		t.Fatal(err)
	}
	cmp2, err := inventory.NewComponent(ctx, "dna_part")
	if err != nil {
		t.Fatal(err)
	}

	s1 := mixer.Sample(cmp1, wunit.NewVolume(50.0, "ul"))
	s2 := mixer.Sample(cmp2, wunit.NewVolume(25.0, "ul"))

	mo := mixer.MixOptions{
		Components: []*wtype.LHComponent{s1, s2},
		PlateType:  "pcrplate_skirted_riser20",
		Address:    "A1",
		PlateNum:   1,
	}

	ins := mixer.GenericMix(mo)

	rq.LHInstructions[ins.ID] = ins

	pl, err := inventory.NewPlate(ctx, "pcrplate_skirted_riser20")
	if err != nil {
		t.Fatal(err)
	}

	rq.Input_platetypes = append(rq.Input_platetypes, pl)

	rq.ConfigureYourself()

	lh := Init(rbt)

	lh.Plan(ctx, rq)

	expected := make(map[string]float64)

	expected["dna_part"] = 30.5
	expected["water"] = 55.5

	testSetup(rbt, expected, t)
}

func testSetup(rbt *liquidhandling.LHProperties, expected map[string]float64, t *testing.T) {
	for _, p := range rbt.Plates {
		for _, w := range p.Wellcoords {
			if !w.IsEmpty() {
				v, ok := expected[w.WContents.CName]

				if !ok {
					t.Errorf("unexpected component in plating area: %s", w.WContents.CName)
				}

				if v != w.WContents.Vol {
					t.Errorf("volume of component %s was %v should be %v", w.WContents.CName, w.WContents.Vol, v)
				}

				delete(expected, w.WContents.CName)
			}
		}
	}

	if len(expected) != 0 {
		t.Errorf("unexpected components remaining: %v", expected)
	}

}
func TestInPlaceAutoAllocate(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	rbt := makeGilson(ctx)
	rq := NewLHRequest()

	cmp1, err := inventory.NewComponent(ctx, inventory.WaterType)
	if err != nil {
		t.Fatal(err)
	}
	cmp2, err := inventory.NewComponent(ctx, "dna_part")
	if err != nil {
		t.Fatal(err)
	}

	cmp1.Vol = 100.0
	cmp2.Vol = 50.0

	mo := mixer.MixOptions{
		Components: []*wtype.LHComponent{cmp1, cmp2},
		PlateType:  "pcrplate_skirted_riser20",
		Address:    "A1",
		PlateNum:   1,
	}

	ins := mixer.GenericMix(mo)

	rq.LHInstructions[ins.ID] = ins

	pl, err := inventory.NewPlate(ctx, "pcrplate_skirted_riser20")
	if err != nil {
		t.Fatal(err)
	}

	rq.Input_platetypes = append(rq.Input_platetypes, pl)

	rq.ConfigureYourself()

	lh := Init(rbt)

	lh.Plan(ctx, rq)

	expected := make(map[string]float64)

	expected["dna_part"] = 55.5
	expected["water"] = 100.0

	testSetup(rbt, expected, t)

}
