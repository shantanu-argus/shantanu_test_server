package system

import (
	"fmt"
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
			movement.CurrentLocation.X += movement.Velocity * float32(movement.CurrentDirection.X)
			movement.CurrentLocation.Y += movement.Velocity * float32(movement.CurrentDirection.Y)

			if err := cardinal.SetComponent[comp.Movement](world, id, movement); err != nil {
				return true
			}
			return true
		})
}

func MovementValidationSystem(world cardinal.WorldContext) error {
	return cardinal.EachMessage[msg.AttackPlayerMsg, msg.AttackPlayerMsgReply](
		world,
		func(attack cardinal.TxData[msg.AttackPlayerMsg]) (msg.AttackPlayerMsgReply, error) {
			playerID, playerHealth, err := queryTargetPlayer(world, attack.Msg.TargetNickname)
			if err != nil {
				return msg.AttackPlayerMsgReply{}, fmt.Errorf("failed to inflict damage: %w", err)
			}

			playerHealth.HP -= AttackDamage
			if err := cardinal.SetComponent[comp.Health](world, playerID, playerHealth); err != nil {
				return msg.AttackPlayerMsgReply{}, fmt.Errorf("failed to inflict damage: %w", err)
			}

			return msg.AttackPlayerMsgReply{Damage: AttackDamage}, nil
		})
}
