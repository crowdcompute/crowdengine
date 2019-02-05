// Copyright 2018 The crowdcompute:crowdengine Authors
// This file is part of the crowdcompute:crowdengine library.
//
// The crowdcompute:crowdengine library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The crowdcompute:crowdengine library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the crowdcompute:crowdengine library. If not, see <http://www.gnu.org/licenses/>.

package terminal

import (
	"github.com/crowdcompute/crowdengine/log"

	"github.com/peterh/liner"
)

// Stdin represents the terminal
var Stdin = newTerminal()

// Terminal represents the liner obj
type Terminal struct {
	*liner.State
	warned     bool
	supported  bool
	normalMode liner.ModeApplier
	rawMode    liner.ModeApplier
}

// GetPassphrase gets the password from stdin
func (t *Terminal) GetPassphrase(prompt string, confirmation bool) (passwd string, err error) {
	if prompt != "" {
		t.State.Prompt(prompt)
	}
	pass, err := t.getPassword("Passphrase:")
	if err != nil {
		log.Fatalf("Error while reading passphrase: %v", err)
	}

	if confirmation {
		confirm, err := t.getPassword("Repeat passphrase: ")
		if err != nil {
			log.Fatalf("Error while reading passphrase confirmation: %v", err)
		}
		if pass != confirm {
			log.Fatalf("Passphrases do not match")
		}
	}
	return pass, nil
}

// getPassword gets the password from stdin
func (t *Terminal) getPassword(prompt string) (passwd string, err error) {
	if t.supported {
		t.rawMode.ApplyMode()
		defer t.normalMode.ApplyMode()
		return t.State.PasswordPrompt(prompt)
	}
	if !t.warned {
		log.Println("Terminal is unsupported and password will be shown!")
		t.warned = true
	}

	log.Print(prompt)
	passwd, err = t.State.Prompt("")
	log.Println()
	return passwd, err
}

// newTerminal returns a terminal instance
func newTerminal() *Terminal {
	t := new(Terminal)
	normalMode, _ := liner.TerminalMode()
	t.State = liner.NewLiner()
	rawMode, err := liner.TerminalMode()
	if err != nil || !liner.TerminalSupported() {
		t.supported = false
	} else {
		t.supported = true
		t.normalMode = normalMode
		t.rawMode = rawMode
		normalMode.ApplyMode()
	}
	t.SetCtrlCAborts(true)
	t.SetTabCompletionStyle(liner.TabPrints)
	t.SetMultiLineMode(true)
	return t
}
