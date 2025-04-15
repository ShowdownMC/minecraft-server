package event

import (
	"github.com/ShowdownMC/minecraft-server/impl/base"
	"github.com/ShowdownMC/minecraft-server/impl/data/plugin"
)

type PlayerConnJoinEvent struct {
	Conn base.PlayerAndConnection
}

type PlayerConnQuitEvent struct {
	Conn base.PlayerAndConnection
}

type PlayerPluginMessagePullEvent struct {
	Conn base.PlayerAndConnection

	Channel string
	Message plugin.Message
}
