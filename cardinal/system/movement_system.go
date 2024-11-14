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
		filter.Exact(filter.Component[comp.Player](), filter.Component[comp.Movement]())).
		Each(world, func(id types.EntityID) bool {
			movement, err := cardinal.GetComponent[comp.Movement](world, id)
			if err != nil {
				fmt.Println("Error in movement system", err)
				return true
			}
			fmt.Println(movement)
			if err := cardinal.SetComponent[comp.Movement](world, id, movement); err != nil {
				return true
			}
			return true
		})
	//return cardinal.NewSearch().Entity(
	//	filter.Exact(filter.Component[comp.Movement]())).
	//	Each(world, func(id types.EntityID) bool {
	//		fmt.Println("Hello from movement system")
	//		movement, err := cardinal.GetComponent[comp.Movement](world, id)
	//		if err != nil {
	//			fmt.Println("Error in movement system")
	//			return true
	//		}
	//		velocity := 5.0
	//		movement.CurrentLocation.X += velocity * float64(movement.CurrentDirection.X) * 0.2
	//		movement.CurrentLocation.Y += velocity * float64(movement.CurrentDirection.Y) * 0.2
	//
	//		if err := cardinal.SetComponent[comp.Movement](world, id, movement); err != nil {
	//			return true
	//		}
	//		return true
	//	})
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
			const errorTolerance float64 = 0.2
			if math.Abs(playerMovementData.CurrentLocation.X-movement.Msg.LocationX) > errorTolerance ||
				math.Abs(playerMovementData.CurrentLocation.Y-movement.Msg.LocationY) > errorTolerance {
				return msg.MovementPlayerMsgReply{IsValid: false,
						LocationX: playerMovementData.CurrentLocation.X,
						LocationY: playerMovementData.CurrentLocation.Y},
					nil
			}
			playerMovementData.CurrentLocation.X = movement.Msg.LocationX
			playerMovementData.CurrentLocation.Y = movement.Msg.LocationY
			if err := cardinal.SetComponent[comp.Movement](world, playerID, playerMovementData); err != nil {
				return msg.MovementPlayerMsgReply{IsValid: false,
						LocationX: playerMovementData.CurrentLocation.X,
						LocationY: playerMovementData.CurrentLocation.Y},
					fmt.Errorf("failed to update movement: %w", err)
			}
			return msg.MovementPlayerMsgReply{IsValid: true,
					LocationX: playerMovementData.CurrentLocation.X,
					LocationY: playerMovementData.CurrentLocation.Y},
				nil
		})
}
