package liquidhandling

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type SplitBlockInstruction struct {
	GenericRobotInstruction
	Inss []*wtype.LHInstruction
}

func NewSplitBlockInstruction(inss []*wtype.LHInstruction) SplitBlockInstruction {
	sb := SplitBlockInstruction{}
	sb.Inss = inss
	sb.GenericRobotInstruction.Ins = RobotInstruction(&sb)
	return sb
}

func (sp SplitBlockInstruction) InstructionType() int {
	return SPB
}

func (sp SplitBlockInstruction) GetParameter(p string) interface{} {
	return nil
}

// this instruction does not generate anything
// it just modifies the components in the robot
func (sp SplitBlockInstruction) Generate(ctx context.Context, policy *wtype.LHPolicyRuleSet, robot *LHProperties) ([]RobotInstruction, error) {
	// this may need more work

	for _, ins := range sp.Inss {
		if ins.Type != wtype.LHISPL {
			return []RobotInstruction{}, fmt.Errorf("Splitblock fed non-split instruction, type %s", ins.InsType())
		}

		// if Components is a sample we'll probably want to change ParentID instead
		// that may not work
		robot.UpdateComponentID(ins.Components[0].ID, ins.Results[1])
	}

	return []RobotInstruction{}, nil
}
