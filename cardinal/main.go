package main

import (
	"errors"

	"github.com/rs/zerolog/log"
	"pkg.world.dev/world-engine/cardinal"

	"shantanu-starter-game/component"
	"shantanu-starter-game/msg"
	"shantanu-starter-game/query"
	"shantanu-starter-game/system"
)

/*
 TODO:
1. setup a new component and system for movement.
2. should take in direction and have an onTick flow to lerp from current position to new position.
3. whenever you get a new position, you should validate if its close enought to lerped position.
4. if its close enough, you should update to that position continue lerping.
5. if its not close enough, you should throw a an error. Client should handle this error and try again.
*/

func main() {
	w, err := cardinal.NewWorld(cardinal.WithDisableSignatureVerification())
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	MustInitWorld(w)

	Must(w.StartGame())
}

// MustInitWorld registers all components, messages, queries, and systems. This initialization happens in a helper
// function so that this can be used directly in tests.
func MustInitWorld(w *cardinal.World) {
	// Register components
	// NOTE: You must register your components here for it to be accessible.
	Must(
		cardinal.RegisterComponent[component.Player](w),
		cardinal.RegisterComponent[component.Health](w),
		cardinal.RegisterComponent[component.Movement](w),
	)

	// Register messages (user action)
	// NOTE: You must register your transactions here for it to be executed.
	Must(
		cardinal.RegisterMessage[msg.CreatePlayerMsg, msg.CreatePlayerResult](w, "create-player"),
		cardinal.RegisterMessage[msg.AttackPlayerMsg, msg.AttackPlayerMsgReply](w, "attack-player"),
		cardinal.RegisterMessage[msg.MovementPlayerMsg, msg.MovementPlayerMsgReply](w, "movement-player"),
	)

	// Register queries
	// NOTE: You must register your queries here for it to be accessible.
	Must(
		cardinal.RegisterQuery[query.PlayerHealthRequest, query.PlayerHealthResponse](w, "player-health", query.PlayerHealth),
	)

	// Each system executes deterministically in the order they are added.
	// This is a neat feature that can be strategically used for systems that depends on the order of execution.
	// For example, you may want to run the attack system before the regen system
	// so that the player's HP is subtracted (and player killed if it reaches 0) before HP is regenerated.
	Must(cardinal.RegisterSystems(w,
		system.AttackSystem,
		system.RegenSystem,
		system.PlayerSpawnerSystem,
		system.MovementValidationSystem,
		system.MovementSystem,
	))

	Must(cardinal.RegisterInitSystems(w,
		system.SpawnDefaultPlayersSystem,
	))
}

func Must(err ...error) {
	e := errors.Join(err...)
	if e != nil {
		log.Fatal().Err(e).Msg("")
	}
}
