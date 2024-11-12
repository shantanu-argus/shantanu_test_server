package system

import (
	"fmt"
	"math"
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
	comp "shantanu-starter-game/component"
	"shantanu-starter-game/msg"
)

// MovementSystem updates the player's location based on its velocity and direction.
func MovementSystem(world cardinal.WorldContext) error {
	return cardinal.NewSearch().Entity(
		filter.Exact(filter.Component[comp.Player](), filter.Component[comp.Movement]())).
		Each(world, func(id types.EntityID) bool {
			movement, err := cardinal.GetComponent[comp.Movement](world, id)
			if err != nil {
				return true
			}

			// TODO: How to get delta time?
			movement.CurrentLocation.X += movement.Velocity * float64(movement.CurrentDirection.X)
			movement.CurrentLocation.Y += movement.Velocity * float64(movement.CurrentDirection.Y)

			if err := cardinal.SetComponent[comp.Movement](world, id, movement); err != nil {
				return true
			}
			return true
		})
}

func MovementValidationSystem(world cardinal.WorldContext) error {
	return cardinal.EachMessage[msg.MovementPlayerMsg, msg.MovementPlayerMsgReply](
		world,
		func(movement cardinal.TxData[msg.MovementPlayerMsg]) (msg.MovementPlayerMsgReply, error) {
			playerID, playerMovementData, err := queryTargetPlayerMovementData(world, movement.Msg.TargetNickname)
			if err != nil {
				return msg.MovementPlayerMsgReply{}, fmt.Errorf("failed to update movement: %w", err)
			}
			const errorTolerance float64 = 0.001
			if math.Abs(playerMovementData.CurrentLocation.X-movement.Msg.LocationX) > errorTolerance ||
				math.Abs(playerMovementData.CurrentLocation.Y-movement.Msg.LocationY) > errorTolerance {
				return msg.MovementPlayerMsgReply{IsValid: false, LocationX: playerMovementData.CurrentLocation.X, LocationY: playerMovementData.CurrentLocation.Y}, nil
			} else {
				playerMovementData.CurrentLocation.X = movement.Msg.LocationX
				playerMovementData.CurrentLocation.Y = movement.Msg.LocationY
				if err := cardinal.SetComponent[comp.Movement](world, playerID, playerMovementData); err != nil {
					return msg.MovementPlayerMsgReply{IsValid: false, LocationX: playerMovementData.CurrentLocation.X, LocationY: playerMovementData.CurrentLocation.Y}, fmt.Errorf("failed to update movement: %w", err)
				}
				return msg.MovementPlayerMsgReply{IsValid: true, LocationX: playerMovementData.CurrentLocation.X, LocationY: playerMovementData.CurrentLocation.Y}, nil
			}
		})
}
