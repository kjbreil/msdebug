package msdebug

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/xjtdy888/mailslot"
	"golang.org/x/text/encoding/charmap"
)

// Client is the "sending" proxy mailslot
type Client struct {
	Mailslot  string // mailslot is the proxy mailslot, created locally that forwards messages to the server
	Send      chan *Message
	Receive   chan *Message
	frequency time.Duration
}

// NewClient creates a client from a mailslot proxy address given
func NewClient(mailslot string) *Client {
	return &Client{
		Mailslot:  mailslot,
		frequency: time.Second,
	}
}

// Start the server, something needs to sit on the send channel otherwise it will block
func (c *Client) Start() error {
	ms, err := mailslot.New(c.Mailslot, 0, mailslot.MAILSLOT_WAIT_FOREVER)
	if err != nil {
		return fmt.Errorf("could not create mailslot with: %v", err)
	}
	// make the send channel, limited to 1 message since the messages will queue up in the mailslot
	c.Send = make(chan *Message, 1)
	// make the receive channel, 100 message right now as a queue
	c.Receive = make(chan *Message, 100)

	// spawn sender goroutine
	go c.sender(ms)
	// spawn receier goroutine
	go c.receiver()
	return nil
}

func (c *Client) sender(ms *mailslot.MailSlot) {
	for {
		info, _ := ms.Info()

		// the count will be greater than 0 if there is a message waiting
		if info.Count > 0 {
			// read the message
			m, err := readMailslot(ms, info.NextSize)
			if err != nil {
				panic(err)
			}
			c.Send <- m
		}
		time.Sleep(c.frequency)
	}
}

func (c *Client) receiver() {
	for {
		// wait on and receive the message
		m := <-c.Receive
		// there was no callback on the message meaning it wasn't supposed to be sent
		if m.Callback == "" {
			continue
		}
		enc, _ := charmap.Windows1252.NewEncoder().String(m.Data)

		// make a reader of the string to be copied into the mailslot
		r := strings.NewReader(enc)
		// open the mailslot
		ms, err := mailslot.Open(m.Callback)
		if err != nil {
			// TODO: mailslot down so should put message back on queue and sleep
			log.Panicf("there was an error opening receiver mailslot: %v", err)
		}
		_, err = io.Copy(ms, r)
		ms.Close()
	}
}
