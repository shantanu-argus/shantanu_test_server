package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"

	comp "shantanu-starter-game/component"
	"shantanu-starter-game/msg"
)

const (
	InitialHP = 100
)

const (
	InitialXLocation = 0
	InitialYLocation = 0
	InitialVelocity  = 0
	InitialDirection = 0
)

// PlayerSpawnerSystem spawns players based on `CreatePlayer` transactions.
// This provides an example of a system that creates a new entity.
func PlayerSpawnerSystem(world cardinal.WorldContext) error {
	return cardinal.EachMessage[msg.CreatePlayerMsg, msg.CreatePlayerResult](
		world,
		func(create cardinal.TxData[msg.CreatePlayerMsg]) (msg.CreatePlayerResult, error) {
			id, err := cardinal.Create(world,
				comp.Player{Nickname: create.Msg.Nickname},
				comp.Health{HP: InitialHP},
				comp.Movement{
					CurrentLocation: comp.Location{X: InitialXLocation, Y: InitialYLocation},
					Velocity:        InitialVelocity,
					CurrentDirection: comp.Direction{
						comp.Directions[InitialDirection].X,
						comp.Directions[InitialDirection].Y,
					},
				},
			)
			if err != nil {
				return msg.CreatePlayerResult{}, fmt.Errorf("error creating player: %w", err)
			}

			err = world.EmitEvent(map[string]any{
				"event": "new_player",
				"id":    id,
			})
			if err != nil {
				return msg.CreatePlayerResult{}, err
			}
			return msg.CreatePlayerResult{Success: true}, nil
		})
}
