package sse

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// SSEBroker est une structure pour gérer les événements envoyés par le serveur (SSE) pour notifier les clients.
type SSEBroker struct {
	Clients        map[uuid.UUID]chan string
	NewClients     chan (chan string)
	DefunctClients chan uuid.UUID
	Messages       chan string
}

// NewSSEBroker initialise un nouveau courtier SSE.
func NewSSEBroker() *SSEBroker {
	return &SSEBroker{
		Clients:        make(map[uuid.UUID]chan string),
		NewClients:     make(chan (chan string)),
		DefunctClients: make(chan uuid.UUID),
		Messages:       make(chan string),
	}
}

// NotifyClients envoie un message à tous les clients connectés.
func (b *SSEBroker) NotifyClients(msg string) {
	b.Messages <- msg
}

// Start exécute la boucle principale de l'événement pour le courtier SSE.
func (b *SSEBroker) Start() {
	for {
		select {
		case s := <-b.NewClients:
			id := uuid.New()
			b.Clients[id] = s
		case id := <-b.DefunctClients:
			delete(b.Clients, id)
		case msg := <-b.Messages:
			for _, s := range b.Clients {
				s <- msg
			}
		}
	}
}

// @Summary Server-Sent Events (SSE) notifications
// @Description Listens for real-time notifications using SSE
// @Produce  text/event-stream
// @Success 200 {string} string "SSE connection established"
// @Router /events [get]
func SSEHandler(b *SSEBroker) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/event-stream")
		c.Response().Header().Set("Cache-Control", "no-cache")
		c.Response().Header().Set("Connection", "keep-alive")
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")

		messageChan := make(chan string)
		b.NewClients <- messageChan

		defer func() {
			for id, ch := range b.Clients {
				if ch == messageChan {
					b.DefunctClients <- id
					break
				}
			}
		}()

		// Check for closed connection
		ctx := c.Request().Context()

		for {
			select {
			case <-ctx.Done():
				return nil
			case msg := <-messageChan:
				_, err := fmt.Fprintf(c.Response(), "data: %s\n\n", msg)
				c.Response().Flush()
				if err != nil {
					return err
				}
			}
		}
	}
}
