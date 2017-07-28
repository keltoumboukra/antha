package image

import (
	"context"
	"testing"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/download"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/testinventory"
)

func TestSelectLibrary(t *testing.T) {
	SelectLibrary("UV")
}

func TestSelectColors(t *testing.T) {
	SelectColor("JuniperGFP")
}

func TestMakeAnthaImg(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	//downloading image for the test
	imgFile, err := download.File("http://orig08.deviantart.net/a19f/f/2008/117/6/7/8_bit_mario_by_superjerk.jpg", "Downloaded file")
	if err != nil {
		t.Error(err)
	}

	//opening image
	imgBase, err := OpenFile(imgFile)
	if err != nil {
		t.Fatal(err)
	}

	palette := SelectLibrary("UV")

	//initiating components
	var components []*wtype.LHComponent
	component, err := inventory.NewComponent(ctx, "Gluc")
	if err != nil {
		t.Fatal(err)
	}

	//making the array to make palette. It's the same length than the array from the "UV" library
	for i := 1; i < 48; i++ {
		components = append(components, component.Dup())
	}
	//getting palette
	anthaPalette := MakeAnthaPalette(palette, components)

	//getting plate
	plate, err := inventory.NewPlate(ctx, "greiner384")
	if err != nil {
		t.Fatal(err)
	}

	//testing function
	MakeAnthaImg(imgBase, anthaPalette, plate)
}

func TestMakeAnthaPalette(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	//getting palette
	palette := SelectLibrary("UV")

	//initiating component
	var components []*wtype.LHComponent
	component, err := inventory.NewComponent(ctx, "Gluc")
	if err != nil {
		t.Fatal(err)
	}

	//making the array to test. It's the same length than the array from the "UV" library
	for i := 1; i < 48; i++ {
		components = append(components, component.Dup())
	}

	//running the function
	MakeAnthaPalette(palette, components)
}

func TestSelectLivingColorLibrary(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	SelectLivingColorLibrary(ctx, "ProteinPaintBox")
}

func TestSelectLivingColor(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())
	SelectLivingColor(ctx, "UVDasherGFP")
}

func TestMakeLivingImg(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	//downloading image for the test
	imgFile, err := download.File("http://orig08.deviantart.net/a19f/f/2008/117/6/7/8_bit_mario_by_superjerk.jpg", "Downloaded file")
	if err != nil {
		t.Error(err)
	}

	//opening image
	imgBase, err := OpenFile(imgFile)
	if err != nil {
		t.Error(err)
	}

	//initiating components
	var components []*wtype.LHComponent
	component, err := inventory.NewComponent(ctx, "Gluc")
	if err != nil {
		t.Fatal(err)
	}

	//making the array to make palette. It's the same length than the array from the "ProteinPaintbox" library
	for i := 1; i < 3; i++ {
		components = append(components, component.Dup())
	}

	//Selecting livingPalette
	selectedPalette := SelectLivingColorLibrary(ctx, "ProteinPaintBox")

	//Making palette
	livingPalette := MakeLivingPalette(selectedPalette, components)

	//getting plate
	plate, err := inventory.NewPlate(ctx, "greiner384")
	if err != nil {
		t.Fatal(err)
	}

	//testing function
	MakeLivingImg(imgBase, livingPalette, plate)
}

func TestMakeLivingGIF(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	//------------------------------------------------
	//Making antha image
	//------------------------------------------------

	//downloading image for the test
	imgFile, err := download.File("http://orig08.deviantart.net/a19f/f/2008/117/6/7/8_bit_mario_by_superjerk.jpg", "Downloaded file")
	if err != nil {
		t.Error(err)
	}

	//opening image
	imgBase, err := OpenFile(imgFile)
	if err != nil {
		t.Error(err)
	}

	//initiating components
	var components []*wtype.LHComponent
	component, err := inventory.NewComponent(ctx, "Gluc")
	if err != nil {
		t.Fatal(err)
	}

	//making the array to make palette. It's the same length than the array from the "ProteinPaintbox" library
	for i := 1; i < 3; i++ {
		components = append(components, component.Dup())
	}

	//Selecting livingPalette
	selectedPalette := SelectLivingColorLibrary(ctx, "ProteinPaintBox")

	//Making palette
	livingPalette := MakeLivingPalette(selectedPalette, components)

	//getting plate
	plate, err := inventory.NewPlate(ctx, "greiner384")
	if err != nil {
		t.Fatal(err)
	}

	//generating images. We only use 2 since that's what we use for our construct
	anthaImg1, _ := MakeLivingImg(imgBase, livingPalette, plate)
	anthaImg2, _ := MakeLivingImg(imgBase, livingPalette, plate)

	//Merge them
	var anthaImgs []LivingImg
	anthaImgs = append(anthaImgs, *anthaImg1, *anthaImg2)

	//------------------------------------------------
	//Testing GIF functions
	//------------------------------------------------

	MakeLivingGIF(anthaImgs)
}