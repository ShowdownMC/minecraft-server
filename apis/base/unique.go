package base

import "github.com/ShowdownMC/minecraft-server/apis/uuid"

type Unique interface {
	UUID() uuid.UUID
}
