package main

import "encoding/json"

type Event struct {
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}

// Client events

type TapEvent struct{}

// Server events

type SwitchStateEvent struct {
	StateName string `json:"stateName"`
	Seconds   int    `json:"seconds"`
}

func NewSwitchStateEvent(stateName string, seconds int) Event {
	return Event{
		Type: "switchState",
		Body: SwitchStateEvent{
			StateName: stateName,
			Seconds:   seconds,
		},
	}
}

type NewGameEvent struct {
	Team string `json:"team"`
}

func NewNewGameEvent(team string) Event {
	return Event{
		Type: "newGame",
		Body: NewGameEvent{
			Team: team,
		},
	}
}

type TickEvent struct {
	TimeRemaining int `json:"timeRemaining"`
}

func NewTickEvent(timeRemaining int) Event {
	return Event{
		Type: "tick",
		Body: TickEvent{
			TimeRemaining: timeRemaining,
		},
	}
}

var eventDecoder = map[string]func(interface{}) interface{}{
	"tap": func(body interface{}) interface{} {
		return TapEvent{}
	},
}

func decodeConnMsg(c *connMsg) (*Event, error) {
	var event Event
	if err := json.Unmarshal(c.body, &event); err != nil {
		return nil, err
	}
	event.Body = eventDecoder[event.Type](&event.Body)
	return &event, nil
}
