package common

import (
	"testing"

	"github.com/hashicorp/packer/helper/multistep"
)

func TestStepCreateCDROM_Impl(t *testing.T) {
	var raw interface{}
	raw = new(StepCreateCDROM)
	if _, ok := raw.(multistep.Step); !ok {
		t.Fatalf("StepCreateCDROM should be a step")
	}
}
