package event

import (
	"github.com/ShowdownMC/minecraft-server/apis/game/level"
)

type BlockEvent struct {
	level.Block
}

type BlockBreakEvent struct {
	BlockEvent
	PlayerEvent
	Cancellable
}
