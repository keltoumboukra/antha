package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

type VolumeSet []wunit.Volume

func NewVolumeSet(n int) VolumeSet {
	var vs VolumeSet
	vs = make([]wunit.Volume, n)
	for i := 0; i < n; i++ {
		vs[i] = (wunit.NewVolume(0.0, "ul"))
	}
	return vs
}

// Add behaves inconsistently with Sub... this is a design error
func (vs VolumeSet) Add(v wunit.Volume) {
	for i := 0; i < len(vs); i++ {
		vs[i].Add(v)
	}
}

// add two volume sets

func (vs VolumeSet) AddA(vs2 VolumeSet) {
	s := len(vs2)

	if len(vs) < s {
		s = len(vs)
	}

	for i := 0; i < s; i++ {
		vs[i].Add(vs2[i])
	}
}

func (vs VolumeSet) Sub(v wunit.Volume) VolumeSet {
	ret := make(VolumeSet, len(vs))
	for i := 0; i < len(vs); i++ {
		ret[i] = wunit.CopyVolume(vs[i])
		ret[i].Subtract(v)
	}
	return ret
}

func (vs VolumeSet) SubA(vs2 VolumeSet) {
	// maintain consistency with the above but one or the other must change
	ret := make(VolumeSet, len(vs))

	for i := 0; i < len(vs); i++ {
		ret[i] = wunit.CopyVolume(vs[i])

		if i < len(vs2) {
			v := vs2[i]
			ret[i].Subtract(v)
		}
	}
}

func (vs VolumeSet) SetEqualTo(v wunit.Volume, multi int) {
	for i := 0; i < multi; i++ {
		vs[i] = wunit.CopyVolume(v)
	}
}

func (vs VolumeSet) Dup() VolumeSet {
	return vs.GetACopy()
}

func (vs VolumeSet) GetACopy() VolumeSet {
	r := make([]wunit.Volume, len(vs))
	for i := 0; i < len(vs); i++ {
		r[i] = wunit.CopyVolume(vs[i])
	}
	return r
}

func (vs VolumeSet) NonZeros() VolumeSet {
	vols := make(VolumeSet, 0, len(vs))

	for _, v := range vs {
		if !v.IsZero() {
			vols = append(vols, v)
		}
	}

	return vols
}

func (vs VolumeSet) IsZero() bool {
	return len(vs.NonZeros()) == 0
}

func (vs VolumeSet) Min() wunit.Volume {
	if len(vs) == 0 {
		return wunit.ZeroVolume()
	}
	v := vs[0]

	for i := 1; i < len(vs); i++ {
		if vs[i].LessThan(v) {
			v = vs[i].Dup()
		}
	}

	return v
}

func countSetSize(set []int) int {
	c := 0
	for _, v := range set {
		if v != -1 {
			c += 1
		}
	}

	return c
}
