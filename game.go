package main

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type gameManager struct {
	currentGame *game
	register    chan *conn
	unregister  chan *conn
	incoming    chan *connMsg
}

var gm = gameManager{
	register:   make(chan *conn),
	unregister: make(chan *conn),
	incoming:   make(chan *connMsg),
}

func (gm *gameManager) run() {
	go gm.processConns()
	for {
		log.Println("Created new game")
		gm.currentGame = newGame()
		for player := range gm.currentGame.state.redTeam.players {
			player.conn.send <- NewNewGameEvent("red")
		}
		for player := range gm.currentGame.state.blueTeam.players {
			player.conn.send <- NewNewGameEvent("blue")
		}
		gm.currentGame.run()
		log.Println("Game finished")
	}
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
	state     *gameData
	miniGames []miniGame
}

func newGame() *game {
	game := game{
		state: newGameData(),
		miniGames: []miniGame{
			&tugOfWar{},
			&shipRace{},
			&shakeWar{},
		},
	}

	// Initialize players
	for conn := range h.conns {
		player := newPlayer(conn)
		game.state.addPlayerToTeam(player)
	}

	return &game
}

func (g *game) run() {
	for _, miniGame := range g.miniGames {
		h.broadcast <- NewSwitchStateEvent(miniGame.name(), miniGame.seconds())
		miniGame.run(g.state)
	}
}

type gameData struct {
	redTeam  *team
	blueTeam *team
}

func newGameData() *gameData {
	return &gameData{
		redTeam:  newTeam(),
		blueTeam: newTeam(),
	}
}

// addPlayerToTeam adds a player to the game state's team with the least players.
func (s *gameData) addPlayerToTeam(p *player) {
	redLen := len(s.redTeam.players)
	blueLen := len(s.blueTeam.players)
	if redLen > blueLen {
		s.blueTeam.players[p] = true
	} else {
		s.redTeam.players[p] = true
	}
}

func (s *gameData) removePlayerWithConnFromTeam(c *conn) error {
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
	wins    int
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
	name() string
	seconds() int
	run(*gameData)
}

type tugOfWar struct {
}

func (t *tugOfWar) name() string {
	return "tugWar"
}

func (t *tugOfWar) seconds() int {
	return 5
}

func (t *tugOfWar) run(*gameData) {
	log.Println("~LET'S TUG SOME ROPES~")

	ticker := time.NewTicker(time.Second)
	secondsRemaining := t.seconds()

loop:
	for {
		select {
		case <-ticker.C:
			secondsRemaining--
			if secondsRemaining == 0 {
				break loop
			}
			//h.broadcast <- NewTickEvent(secondsRemaining)
		case m := <-gm.incoming:
			fmt.Println(string(m.body))
			evt, err := decodeConnMsg(m)
			if err != nil {
				log.Println("Error decoding message", err)
				continue
			}

			fmt.Println(evt)
		}
	}

	log.Println("Done!")
}

type shipRace struct {
}

func (s *shipRace) name() string {
	return "shipRace"
}

func (s *shipRace) seconds() int {
	return 5
}

func (s *shipRace) run(*gameData) {
	log.Println("~LET'S RACE SOME SHIPS~")

	ticker := time.NewTicker(time.Second)
	secondsRemaining := s.seconds()

loop:
	for {
		select {
		case <-ticker.C:
			secondsRemaining--
			if secondsRemaining == 0 {
				break loop
			}
			//h.broadcast <- NewTickEvent(secondsRemaining)
		}
	}

	log.Println("Done!")
}

type shakeWar struct {
}

func (s *shakeWar) name() string {
	return "shakeWar"
}

func (s *shakeWar) seconds() int {
	return 5
}

func (s *shakeWar) run(*gameData) {
	log.Println("~SHAKE IT LIKE A POLAROID PICTURE~")

	ticker := time.NewTicker(time.Second)
	secondsRemaining := s.seconds()

loop:
	for {
		select {
		case <-ticker.C:
			secondsRemaining--
			if secondsRemaining == 0 {
				break loop
			}
			//h.broadcast <- NewTickEvent(secondsRemaining)
		}
	}

	log.Println("Done!")
}
