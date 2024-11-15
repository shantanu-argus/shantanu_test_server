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
	fmt.Println("Hello from movement system")

	return cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[comp.Movement]())).
		Each(world, func(id types.EntityID) bool {
			movement, err := cardinal.GetComponent[comp.Movement](world, id)
			if err != nil {
				fmt.Println("Error in movement system")
				return true
			}
			// Simulate the motion of the player
			movement.CurrentLocation.X += movement.Velocity * float64(movement.CurrentDirection.X) * 0.2
			movement.CurrentLocation.Y += movement.Velocity * float64(movement.CurrentDirection.Y) * 0.2

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
			fmt.Println("Hello from movement validation system", movement)
			playerID, playerMovementData, err := queryTargetPlayerMovementData(world, movement.Msg.TargetNickname)
			if err != nil {
				return msg.MovementPlayerMsgReply{}, fmt.Errorf("failed to update movement: %w", err)
			}
			const errorTolerance float64 = 5
			if math.Abs(playerMovementData.CurrentLocation.X-movement.Msg.LocationX) > errorTolerance ||
				math.Abs(playerMovementData.CurrentLocation.Y-movement.Msg.LocationY) > errorTolerance {
				fmt.Println(fmt.Errorf("player moved too far from server location"))
				// force player to reset to servers location
				return msg.MovementPlayerMsgReply{IsValid: false,
						LocationX: playerMovementData.CurrentLocation.X,
						LocationY: playerMovementData.CurrentLocation.Y},
					nil
			}

			// update the movement if its error margin is within the tolerance
			playerMovementData.Velocity = movement.Msg.Velocity
			playerMovementData.CurrentDirection = comp.Directions[movement.Msg.Direction]
			playerMovementData.CurrentLocation.X = movement.Msg.LocationX
			playerMovementData.CurrentLocation.Y = movement.Msg.LocationY
			fmt.Println(playerMovementData)
			if err := cardinal.SetComponent[comp.Movement](world, playerID, playerMovementData); err != nil {
				return msg.MovementPlayerMsgReply{IsValid: false,
						LocationX: playerMovementData.CurrentLocation.X,
						LocationY: playerMovementData.CurrentLocation.Y},
					fmt.Errorf("failed to update movement: %w", err)
			}
			fmt.Println("Updated player movement ", err)
			return msg.MovementPlayerMsgReply{IsValid: true,
					LocationX: playerMovementData.CurrentLocation.X,
					LocationY: playerMovementData.CurrentLocation.Y},
				nil
		})
}
