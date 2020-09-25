package msdebug

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/xjtdy888/mailslot"
	"golang.org/x/text/encoding/charmap"
)

// Message contains the message and the callback mailslot
type Message struct {
	Data     string
	Callback string
}

func (m *Message) String() string {
	return fmt.Sprintf("Callback: %s, Message: %s", m.Callback, m.Data)
}

// Format the message for a mailslot
func (m *Message) Format() string {
	enc, _ := charmap.Windows1252.NewEncoder().String(fmt.Sprintf("%s,%s", m.Callback, m.Data))
	return enc
}

func readMailslot(ms *mailslot.MailSlot, size int32) (*Message, error) {
	buf := make([]byte, size)
	n, err := ms.Read(buf)
	if err == io.EOF {
		return nil, fmt.Errorf("mailslot reached eof")
	}

	dec, _ := charmap.Windows1252.NewDecoder().Bytes(buf[:n])
	s := string(dec)

	sa := strings.Split(s, ",")
	if len(sa) < 2 {
		return &Message{
			Data:     s,
			Callback: "",
		}, nil
	}
	// TODO: This will not work if there is a comma in the message, need to join on all but 0 oridinal
	c, m := sa[0], strings.Join(sa[1:], ",")
	return &Message{
		Data:     m,
		Callback: c,
	}, nil
}

// Send a message with callback to specific mailslot
func (m *Message) Send(msAddress string) {
	r := strings.NewReader(m.Format())
	// open the mailslot
	ms, err := mailslot.Open(msAddress)
	if err != nil {
		// TODO: mailslot down so should put message back on queue and sleep
		log.Panicf("there was an error opening receiver mailslot: %v", err)
	}
	_, err = io.Copy(ms, r)
	ms.Close()

}
