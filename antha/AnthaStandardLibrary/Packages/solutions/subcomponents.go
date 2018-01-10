// Part of the Antha language
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

// solutions is a utility package for working with solutions of LHComponents
package solutions

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/pubchem"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// ComponentListSample is a sample of a Component list at a specified volume
// when two ComponentListSamples are mixed a new diluted ComponentList is generated
type ComponentListSample struct {
	ComponentList
	Volume wunit.Volume
}

// mixComponentLists merges two componentListSamples.
// When two ComponentListSamples are mixed a new diluted ComponentList is generated.
// An error may be generated if two components with the same name exist within the two lists with incompatible concentration units.
// In this instance, the molecular weight for that component will be looked up in pubchem in order to change the units in both lists to g/l,
// which will be able to be added.
func mixComponentLists(sample1, sample2 ComponentListSample) (newList ComponentList, err error) {

	var errs []string

	complist := make(map[string]wunit.Concentration)

	sample1DilutionRatio := sample1.Volume.SIValue() / (sample1.Volume.SIValue() + sample2.Volume.SIValue())

	sample2DilutionRatio := sample2.Volume.SIValue() / (sample1.Volume.SIValue() + sample2.Volume.SIValue())

	for key, conc := range sample1.Components {

		newConc := wunit.MultiplyConcentration(conc, sample1DilutionRatio)

		if existingConc, found := complist[key]; found {
			sumOfConcs, newerr := wunit.AddConcentrations(newConc, existingConc)
			if newerr != nil {
				// attempt unifying base units
				molecule, newerr := pubchem.MakeMolecule(key)
				if newerr != nil {
					errs = append(errs, newerr.Error())
				} else {
					newConcG := molecule.GramPerL(newConc)
					existingConcG := molecule.GramPerL(existingConc)

					sumOfConcs, newerr = wunit.AddConcentrations(newConcG, existingConcG)
					if newerr != nil {
						errs = append(errs, newerr.Error())
					}
				}
			}
			complist[key] = sumOfConcs
		} else {
			complist[key] = newConc
		}
	}

	for key, conc := range sample2.Components {
		newConc := wunit.MultiplyConcentration(conc, sample2DilutionRatio)

		if existingConc, found := complist[key]; found {
			sumOfConcs, newerr := wunit.AddConcentrations(newConc, existingConc)
			if newerr != nil {
				// attempt unifying base units
				molecule, newerr := pubchem.MakeMolecule(key)
				if newerr != nil {
					errs = append(errs, newerr.Error())
				} else {
					newConcG := molecule.GramPerL(newConc)
					existingConcG := molecule.GramPerL(existingConc)

					sumOfConcs, newerr = wunit.AddConcentrations(newConcG, existingConcG)
					if newerr != nil {
						errs = append(errs, newerr.Error())
					}
				}
			}
			complist[key] = sumOfConcs
		} else {
			complist[key] = newConc
		}
	}
	newList.Components = complist

	if len(errs) > 0 {
		err = fmt.Errorf(strings.Join(errs, "; "))
	}

	return
}

// SimulateMix simulates the resulting list of components and concentrations
// which would be generated by mixing the samples together.
// This will only add the component name itself to the new component list if the sample has no components
// this is to prevent potential duplication since if a component has a list of sub components the name
// is considered to be an alias and the component list the true meaning of what the component is.
// If any sample concentration of zero is found the component list will be made but an error returned.
func SimulateMix(samples ...*wtype.LHComponent) (newComponentList ComponentList, mixSteps []ComponentListSample, err error) {

	var errs []string
	var nonZeroVols []wunit.Volume
	var forTotalVol wunit.Volume
	// top up volume will only be used if a SampleForTotalVolume command is used
	var bufferIndex int = -1

	for i, sample := range samples {
		if sample.Volume().RawValue() == 0.0 && sample.Tvol > 0 {
			forTotalVol = wunit.NewVolume(sample.Tvol, sample.Vunit)
			bufferIndex = i
		}
		nonZeroVols = append(nonZeroVols, sample.Volume())
	}

	topUpVolume := wunit.SubtractVolumes(forTotalVol, nonZeroVols...)

	var volsSoFar []wunit.Volume

	for i, sample := range samples {

		var volToAdd wunit.Volume

		if i == bufferIndex {
			volToAdd = topUpVolume
		} else {
			volToAdd = sample.Volume()
		}

		if i == 0 {

			newComponentList, err = GetSubComponents(sample)
			if err != nil {
				newComponentList.Components = make(map[string]wunit.Concentration)
				if sample.Conc == 0 {
					errs = append(errs, "zero concentration found for sample "+sample.CName)
					newComponentList.Components[sample.Name()] = wunit.NewConcentration(0.0, "g/L")
				} else {
					newComponentList.Components[sample.Name()] = sample.Concentration()
				}
				mixSteps = append(mixSteps, ComponentListSample{newComponentList, volToAdd})
			}
			volsSoFar = append(volsSoFar, volToAdd)
		}

		if i < len(samples)-1 {

			nextSample := samples[i+1]
			nextList, err := GetSubComponents(nextSample)
			if err != nil {
				nextList.Components = make(map[string]wunit.Concentration)
				if nextSample.Conc == 0 {
					errs = append(errs, "zero concentration found for sample "+sample.CName)
					nextList.Components[nextSample.Name()] = wunit.NewConcentration(0.0, "g/L")
				} else {
					nextList.Components[nextSample.Name()] = nextSample.Concentration()
				}
			}

			if i != 0 {
				volsSoFar = append(volsSoFar, volToAdd)
			}

			volOfPreviousSamples := wunit.AddVolumes(volsSoFar...)

			previousMixStep := ComponentListSample{newComponentList, volOfPreviousSamples}

			var nexSampleVolToAdd wunit.Volume

			if i+1 == bufferIndex {
				nexSampleVolToAdd = topUpVolume
			} else {
				nexSampleVolToAdd = nextSample.Volume()
			}
			nextMixStep := ComponentListSample{nextList, nexSampleVolToAdd}
			newComponentList, err = mixComponentLists(previousMixStep, nextMixStep)

			if err != nil {
				errs = append(errs, err.Error())
			}

			mixSteps = append(mixSteps, nextMixStep)

		}

	}

	if len(errs) > 0 {
		err = fmt.Errorf(strings.Join(errs, "; "))
	}
	return newComponentList, mixSteps, nil
}

// List of the components and corresponding concentrations contained within an LHComponent
type ComponentList struct {
	Components map[string]wunit.Concentration `json:"Components"`
}

// add a single entry to a component list
func (c ComponentList) Add(component *wtype.LHComponent, conc wunit.Concentration) (newlist ComponentList) {

	componentName := NormaliseName(component.Name())

	complist := make(map[string]wunit.Concentration)
	for k, v := range c.Components {
		complist[k] = v
	}
	if _, found := complist[componentName]; !found {
		complist[componentName] = conc
	}

	newlist.Components = complist
	return
}

// Get a single concentration set point for a named component present in a component list.
// An error will be returned if the component is not present.
func (c ComponentList) Get(component *wtype.LHComponent) (conc wunit.Concentration, err error) {

	componentName := NormaliseName(component.Name())

	conc, found := c.Components[componentName]

	if found {
		return conc, nil
	} else {
		return conc, &notFound{Name: component.CName, All: c.AllComponents()}
	}
	return
}

// Get a single concentration set point using just the name of a component present in a component list.
// An error will be returned if the component is not present.
func (c ComponentList) GetByName(component string) (conc wunit.Concentration, err error) {

	component = NormaliseName(component)

	conc, found := c.Components[component]

	if found {
		return conc, nil
	} else {
		return conc, &notFound{Name: component, All: c.AllComponents()}
	}
	return
}

// List all Components and concentration set points presnet in a component list.
// if verbose is set to true the field annotations for each component and concentration will be included for each component.
func (c ComponentList) List(verbose bool) string {
	var s []string

	var sortedKeys []string

	for key, _ := range c.Components {
		sortedKeys = append(sortedKeys, key)
	}

	sort.Strings(sortedKeys)

	for _, k := range sortedKeys {

		v := c.Components[k]

		var message string
		if verbose {
			message = fmt.Sprintln("Component: ", k, "Conc: ", v)
		} else {
			message = v.ToString() + " " + k
		}

		s = append(s, message)
	}
	var list string
	if verbose {
		list = strings.Join(s, ";")
	} else {
		list = strings.Join(s, "---")
	}
	return list
}

// Returns all component names present in component list, sorted in alphabetical order.
func (c ComponentList) AllComponents() []string {
	var s []string

	for k, _ := range c.Components {
		s = append(s, k)
	}

	sort.Strings(s)

	return s
}

type alreadyAdded struct {
	Name string
}

func (a *alreadyAdded) Error() string {
	return "component " + a.Name + " already added"
}

type notFound struct {
	Name string
	All  []string
}

func (a *notFound) Error() string {
	return "component " + a.Name + " not found. Found: " + strings.Join(a.All, ";")
}

// returns error if component already found
func AddSubComponent(component *wtype.LHComponent, subcomponent *wtype.LHComponent, conc wunit.Concentration) (*wtype.LHComponent, error) {
	var err error

	if component == nil {
		return nil, fmt.Errorf("No component specified so cannot add subcomponent")
	}
	if subcomponent == nil {
		return nil, fmt.Errorf("No subcomponent specified so cannot add subcomponent")
	}
	if _, found := component.Extra["History"]; !found {
		complist := make(map[string]wunit.Concentration)

		complist[subcomponent.CName] = conc

		var newlist ComponentList

		newlist = newlist.Add(subcomponent, conc)

		if len(newlist.Components) == 0 {

			return component, fmt.Errorf("No subcomponent added! list still empty")
		}

		if _, err := newlist.Get(subcomponent); err != nil {
			return component, fmt.Errorf("No subcomponent added, no subcomponent to get: %s!", err.Error())

		}

		component, err = setHistory(component, newlist)

		if err != nil {
			return component, err
		}

		history, err := getHistory(component)

		if err != nil {
			return component, fmt.Errorf("Error getting History for %s: %s", component.CName, err.Error())
		}

		if len(history.Components) == 0 {
			return component, fmt.Errorf("No history added!")
		}
		return component, nil
	} else {

		history, err := getHistory(component)

		if err != nil {
			return component, err
		}

		if _, found := history.Components[subcomponent.CName]; !found {
			history = history.Add(subcomponent, conc)
			component, err = setHistory(component, history)
			return component, err
		} else {
			return component, &alreadyAdded{Name: subcomponent.CName}
		}
	}
}

// Add a component list to a component.
// Any existing component list will be overwritten
func AddSubComponents(component *wtype.LHComponent, allsubComponents ComponentList) (*wtype.LHComponent, error) {

	for _, compName := range allsubComponents.AllComponents() {
		var comp wtype.LHComponent

		comp.CName = compName

		conc, err := allsubComponents.Get(&comp)

		if err != nil {
			return component, err
		}

		component, err = AddSubComponent(component, &comp, conc)

		if err != nil {
			return component, err
		}
	}

	return component, nil
}

// return a component list from a component
func GetSubComponents(component *wtype.LHComponent) (componentMap ComponentList, err error) {

	components, err := getHistory(component)

	if err != nil {
		return componentMap, fmt.Errorf("Error getting componentList for %s: %s", component.CName, err.Error())
	}

	if len(components.Components) == 0 {
		return components, fmt.Errorf("No sub components found for %s", component.CName)
	}

	return components, nil
}

// utility function to allow the object properties to be retained when serialised
func serialise(compList ComponentList) ([]byte, error) {

	return json.Marshal(compList)
}

// utility function to allow the object properties to be retained when serialised
func deserialise(data []byte) (compList ComponentList, err error) {
	compList = ComponentList{}
	err = json.Unmarshal(data, &compList)
	return
}

// Return a component list from a component.
// Users should use getSubComponents function.
func getHistory(comp *wtype.LHComponent) (compList ComponentList, err error) {

	history, found := comp.Extra["History"]

	if !found {
		return compList, fmt.Errorf("No component list found")
	}

	var bts []byte

	bts, err = json.Marshal(history)
	if err != nil {
		return
	}

	err = json.Unmarshal(bts, &compList)

	if err != nil {
		err = fmt.Errorf("Problem getting %s history. History found: %+v; error: %s", comp.Name(), history, err.Error())
	}

	return
}

// Add a component list to a component.
// Any existing component list will be overwritten.
// Users should use add SubComponents function
func setHistory(comp *wtype.LHComponent, compList ComponentList) (*wtype.LHComponent, error) {

	comp.Extra["History"] = compList // serialisedList

	return comp, nil
}
