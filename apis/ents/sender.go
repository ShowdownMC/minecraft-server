package ents

import (
	"github.com/ShowdownMC/minecraft-server/apis/base"
)

type Sender interface {
	base.Named

	SendMessage(message ...interface{})
}
