package ast

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// QPCRInstruction is a high-level instruction to perform a QPCR analysis.
type QPCRInstruction struct {
	ID           string
	ComponentIn  []*wtype.LHComponent
	ComponentOut []*wtype.LHComponent
	Definition   string
	Barcode      string
	Command      string
	TagAs        string
}

func (ins QPCRInstruction) String() string {
	return fmt.Sprint("QPCRInstruction")
}

// NewQPCRInstruction creates a new QPCRInstruction
func NewQPCRInstruction() *QPCRInstruction {
	var inst QPCRInstruction
	inst.ID = wtype.GetUUID()
	return &inst
}
