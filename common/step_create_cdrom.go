package common

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/packer/tmp"

	"github.com/kdomanski/iso9660"
)

// StepCreateCDROM will create a CD-ROM ISO-9660 disk from the given paths.
type StepCreateCDROM struct {
	Paths []string
	Label string

	cdromPath string
}

func (s *StepCreateCDROM) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	if len(s.Paths) == 0 {
		log.Println("No files specified. CD-ROM will not be made.")
		return multistep.ActionContinue
	}

	if s.Label == "" {
		s.Label = "packer"
	} else {
		log.Printf("CD-ROM label is set to %s", s.Label)
	}

	ui := state.Get("ui").(packer.Ui)
	ui.Say("Creating CD-ROM...")

	writer, err := iso9660.NewWriter()
	if err != nil {
		state.Put("error", fmt.Errorf("Error creating CD-ROM writer: %s", err))
		return multistep.ActionHalt
	}
	defer writer.Cleanup()

	for _, p := range s.Paths {
		fh, err := os.Open(p)
		if err != nil {
			state.Put("error", fmt.Errorf("Error opening %#v for copying into CD-ROM image: %w", p, err))
			return multistep.ActionHalt
		}
		defer fh.Close()

		err = writer.AddFile(fh, p)
		if err != nil {
			state.Put("error", fmt.Errorf("Error adding file %#v into CD-ROM image: %w", p, err))
			return multistep.ActionHalt
		}
	}

	outfh, err := tmp.File("packer")
	if err != nil {
		state.Put("error", fmt.Errorf("Error creating temporary file for CD-ROM: %w", err))
		return multistep.ActionHalt
	}
	defer outfh.Close()

	s.cdromPath = outfh.Name()

	if err := writer.WriteTo(outfh, s.Label); err != nil {
		state.Put("error", fmt.Errorf("Error writing CD-ROM image to %#v: %w", outfh.Name(), err))
		return multistep.ActionHalt
	}
	if err := outfh.Close(); err != nil {
		state.Put("error", fmt.Errorf("Error writing CD-ROM image to %#v: %w", outfh.Name(), err))
		return multistep.ActionHalt
	}
	return multistep.ActionContinue
}

func (s *StepCreateCDROM) Cleanup(multistep.StateBag) {
	if s.cdromPath != "" {
		log.Printf("Deleting CD-ROM ISO: %#v", s.cdromPath)
		os.Remove(s.cdromPath)
	}
}
