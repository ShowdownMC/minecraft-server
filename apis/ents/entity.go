package ents

import "github.com/ShowdownMC/minecraft-server/apis/base"

type Entity interface {
	Sender
	base.Unique

	EntityUUID() int64
}
