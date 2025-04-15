package mode

import (
	"github.com/ShowdownMC/minecraft-server/apis/util"
	"github.com/ShowdownMC/minecraft-server/impl/base"
	"github.com/ShowdownMC/minecraft-server/impl/prot/server"
)

/**
 * handshake
 */

func HandleState0(watcher util.Watcher) {

	watcher.SubAs(func(packet *server.PacketIHandshake, conn base.Connection) {
		conn.SetState(packet.State)
	})

}
