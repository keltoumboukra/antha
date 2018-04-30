package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"testing"
)

func TestUniquePlateName(t *testing.T) {
	mswl := func(s string) map[string]*wtype.LHPlate {
		return map[string]*wtype.LHPlate{s: {}}
	}

	type testData struct {
		Name         string
		InputPlates  map[string]*wtype.LHPlate
		OutputPlates map[string]*wtype.LHPlate
	}

	tests := []testData{
		{Name: "InputPlates", InputPlates: mswl("input_plate_1"), OutputPlates: mswl("")},
		{Name: "OutputPlates", InputPlates: mswl(""), OutputPlates: mswl("output_plate_1")},
		{Name: "Both", InputPlates: mswl("input_plate_1"), OutputPlates: mswl("output_plate_1")},
	}

	for i := range tests {
		dat := tests[i]

		doTheTest := func(t *testing.T) {
			rq := NewLHRequest()
			rq.Input_plates = dat.InputPlates
			rq.Output_plates = dat.OutputPlates

			for v := 0; v < 100; v++ {
				nom := getSafeInputPlateName(rq, 1)

				if rq.HasPlateNamed(nom) {
					t.Errorf("Plate named %s returned by getSafePlateName - already defined by request", nom)
				}

				rq.AddUserPlate(&wtype.LHPlate{PlateName: nom})

				if !rq.HasPlateNamed(nom) {
					t.Errorf("Plate named %s not recognised by request after addition", nom)
				}
			}
		}

		t.Run(dat.Name, doTheTest)
	}

}
