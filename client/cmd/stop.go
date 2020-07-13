package cmd

import (
	"fmt"
	"os"

	"alda.io/client/emitter"
	log "alda.io/client/logging"
	"alda.io/client/system"
	"github.com/spf13/cobra"
)

func init() {
	stopCmd.Flags().StringVarP(
		&playerID, "player-id", "i", "", "The ID of the player process to tell to stop playback",
	)

	stopCmd.Flags().IntVarP(
		&port, "port", "p", -1, "The port of the player process to tell to stop playback",
	)
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop playback",
	Run: func(_ *cobra.Command, args []string) {
		players := []system.PlayerState{}

		// Determine the players to which to send a "stop" message based on the
		// provided CLI options.
		switch {
		// Port is explicitly specified, so use that port.
		case port != -1:
			players = append(players, system.PlayerState{
				ID: "unknown", State: "unknown", Port: port,
			})
		// Player ID is specified; look up the player by ID and use its port.
		case playerID != "":
			player, err := system.FindPlayerByID(playerID)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			players = append(players, player)
		// Send a "stop" message to all known players.
		default:
			knownPlayers, err := system.ReadPlayerStates()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			players = append(players, knownPlayers...)
		}

		for _, player := range players {
			err := (emitter.OSCEmitter{Port: player.Port}).EmitStopMessage()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			log.Info().
				Interface("player", player).
				Msg("Sent \"stop\" message to player process.")
		}
	},
}
