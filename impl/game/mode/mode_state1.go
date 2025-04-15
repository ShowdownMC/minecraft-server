package mode

import (
	"github.com/ShowdownMC/minecraft-server/apis/util"
	"github.com/ShowdownMC/minecraft-server/impl/base"
	"github.com/ShowdownMC/minecraft-server/impl/data/status"
	"github.com/ShowdownMC/minecraft-server/impl/prot/client"
	"github.com/ShowdownMC/minecraft-server/impl/prot/server"
)

/**
 * status
 */

func HandleState1(watcher util.Watcher) {

	watcher.SubAs(func(packet *server.PacketIRequest, conn base.Connection) {
		response := client.PacketOResponse{Status: status.DefaultResponse()}
		conn.SendPacket(&response)
	})

	watcher.SubAs(func(packet *server.PacketIPing, conn base.Connection) {
		response := client.PacketOPong{Ping: packet.Ping}
		conn.SendPacket(&response)
	})

}
