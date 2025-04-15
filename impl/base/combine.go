package base

import "github.com/ShowdownMC/minecraft-server/apis/ents"

type PlayerAndConnection struct {
	Connection
	ents.Player
}
