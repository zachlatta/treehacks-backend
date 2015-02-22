package main

import (
	"errors"
	"log"
)

type gameManager struct {
	currentGame *game
	register    chan *conn
	unregister  chan *conn
}

var gm = gameManager{
	register:   make(chan *conn),
	unregister: make(chan *conn),
}

func (gm *gameManager) run() {
	go gm.processConns()
	gm.currentGame = newGame()
}

func (gm *gameManager) processConns() {
	for {
		select {
		case c := <-gm.register:
			gm.currentGame.state.addPlayerToTeam(newPlayer(c))
		case c := <-gm.unregister:
			gm.currentGame.state.removePlayerWithConnFromTeam(c)
		}
	}
}

type game struct {
	state *gameState
}

func (g *game) start() {
	log.Println("Starting game")
}

func newGame() *game {
	game := game{
		state: newGameState(),
	}

	// Initialize players
	for conn := range h.conns {
		player := newPlayer(conn)
		game.state.addPlayerToTeam(player)
	}

	return &game
}

type gameState struct {
	redTeam  *team
	blueTeam *team
}

func newGameState() *gameState {
	return &gameState{
		redTeam:  newTeam(),
		blueTeam: newTeam(),
	}
}

// addPlayerToTeam adds a player to the game state's team with the least players.
func (s *gameState) addPlayerToTeam(p *player) {
	redLen := len(s.redTeam.players)
	blueLen := len(s.blueTeam.players)
	if redLen > blueLen {
		s.blueTeam.players[p] = true
	} else {
		s.redTeam.players[p] = true
	}
}

func (s *gameState) removePlayerWithConnFromTeam(c *conn) error {
	removed := false
	teams := []*team{s.redTeam, s.blueTeam}
	for _, team := range teams {
		for player := range team.players {
			if c == player.conn {
				if removed {
					log.Panic("multiple players were associated with a single connection. offending player and connection: ", *player, *c)
				}

				delete(team.players, player)
				removed = true
			}
		}
	}

	if !removed {
		return errors.New("no player found with connection")
	}

	return nil
}

type team struct {
	points  int
	players map[*player]bool
}

func newTeam() *team {
	return &team{
		players: make(map[*player]bool),
	}
}

type player struct {
	conn *conn
}

func newPlayer(c *conn) *player {
	return &player{
		conn: c,
	}
}

type miniGame interface {
	start(*gameState)
}

type tugOfWar struct {
}

func (t *tugOfWar) start(gm *gameState) {
}
