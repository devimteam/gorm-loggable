package loggable_test

import (
	"testing"

	"github.com/devimteam/gorm-loggable"
)

type Client struct {
	Passports  []Passport
	ClientInfo ClientInfo
	loggable.LoggableModel
}

type Passport struct {
	loggable.LoggableModel
}

type ClientInfo struct {
	Socials []Social
	loggable.LoggableModel
}

type Social struct {
	loggable.LoggableModel
}

// Test RecursiveSetLogableEnabled func
// which sets loggableModel.Disabled property for any provided struct.
func TestRecursiveSetLoggableDisabled(t *testing.T) {

	x := Client{
		Passports: []Passport{
			{},
			{},
		},
		ClientInfo: ClientInfo{
			Socials: []Social{
				{},
				{},
			},
		},
	}
	loggable.RecursiveSetLoggableEnabled(&x, false)

	if x.Enabled() {
		t.Error("loggable for client enabled")
	}

	if x.Passports[0].Enabled() {
		t.Error("loggable for passport enabled")
	}

	if x.ClientInfo.Enabled() {
		t.Error("loggable for client info")
	}

	if x.ClientInfo.Socials[0].Enabled() {
		t.Error("loggable for client info socials")
	}

	if x.ClientInfo.Socials[1].Enabled() {
		t.Error("loggable for client info socials")
	}
}
