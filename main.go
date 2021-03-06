package main

import (
	"github.com/termoose/irccloud/config"
	"github.com/termoose/irccloud/events"
	"github.com/termoose/irccloud/requests"
	"github.com/termoose/irccloud/ui"
	"log"
)

func main() {
	conf := config.Parse()

	session, err := requests.GetSessionToken(conf.Username, conf.Password)

	if err != nil {
		log.Print(err)
		return
	}

	wsConn := requests.NewConnection(session)
	view := ui.NewView(wsConn, conf.Triggers, conf.LastChan)

	defer func() {
		current := view.GetCurrentChannel()
		config.WriteLatestChannel(conf, current)
		view.Stop()
	}()

	eventHandler := events.NewHandler(session, view)

	go func() {
		for {
			msg, err := wsConn.ReadMessage()

			if err != nil {
				view.Stop()
				log.Print(err)

				return
			}

			eventHandler.Enqueue(msg)
		}
	}()

	view.Start()
}
