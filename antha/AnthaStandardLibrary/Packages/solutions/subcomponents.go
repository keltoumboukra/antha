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

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// List of the components and corresponding concentrations contained within an LHComponent
type ComponentList struct {
	Components map[string]wunit.Concentration `json:"Components"`
}

// add a single entry to a component list
func (c ComponentList) Add(component *wtype.LHComponent, conc wunit.Concentration) (newlist ComponentList) {
	complist := make(map[string]wunit.Concentration)
	for k, v := range c.Components {
		complist[k] = v
	}
	if _, found := complist[component.CName]; !found {
		complist[component.CName] = conc
	}

	newlist.Components = complist
	return
}

// Get a single concentration set point for a named component present in a component list.
// An error will be returned if the component is not present.
func (c ComponentList) Get(component *wtype.LHComponent) (conc wunit.Concentration, err error) {
	conc, found := c.Components[component.CName]

	if found {
		return conc, nil
	} else {
		return conc, &notFound{Name: component.CName}
	}
	return
}

// Get a single concentration set point using just the name of a component present in a component list.
// An error will be returned if the component is not present.
func (c ComponentList) GetByName(component string) (conc wunit.Concentration, err error) {
	conc, found := c.Components[component]

	if found {
		return conc, nil
	} else {
		return conc, &notFound{Name: component}
	}
	return
}

// List all Components and concentration set points presnet in a component list.
// if verbose is set to true the field annotations for each component and concentration will be included for each component.
func (c ComponentList) List(verbose bool) string {
	var s []string

	for k, v := range c.Components {

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
}

func (a *notFound) Error() string {
	return "component " + a.Name + " not found"
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
