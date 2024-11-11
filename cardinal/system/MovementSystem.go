package system

import (
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
	comp "shantanu-starter-game/component"
)

// RegenSystem replenishes the player's HP at every tick.
// This provides an example of a system that doesn't rely on a transaction to update a component.
func MovementSystem(world cardinal.WorldContext) error {
	return cardinal.NewSearch().Entity(
		filter.Exact(filter.Component[comp.Player](), filter.Component[comp.Movement]())).
		Each(world, func(id types.EntityID) bool {
			movement, err := cardinal.GetComponent[comp.Movement](world, id)
			if err != nil {
				return true
			}

			// TODO: How to get delta time?
			movement.Location.X += movement.Velocity * float32(movement.Direction.X)
			movement.Location.Y += movement.Velocity * float32(movement.Direction.Y)

			if err := cardinal.SetComponent[comp.Movement](world, id, movement); err != nil {
				return true
			}
			return true
		})
}
