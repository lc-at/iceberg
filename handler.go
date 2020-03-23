package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Rhymen/go-whatsapp"
	"github.com/Rhymen/go-whatsapp/binary/proto"
)

// Handler whatsapp connection handler
type Handler struct {
	wac       *whatsapp.Conn
	startTime time.Time
}

// HandleError handles an error
func (h Handler) HandleError(err error) {

	if e, ok := err.(*whatsapp.ErrConnectionFailed); ok {
		log.Printf("connection failed, underlying error: %v\n", e.Err)
		log.Println("waiting for 30 seconds ...")
		<-time.After(30 * time.Second)
		log.Println("reconnecting ...")
		err := h.wac.Restore()
		if err != nil {
			log.Fatalf("restore failed: %v", err)
		}
	} else {
		log.Printf("error occoured: %v\n", err)
	}
}

// HandleTextMessage handles a text message
func (h Handler) HandleTextMessage(message whatsapp.TextMessage) {
	if message.Info.Timestamp < uint64(h.startTime.Unix()) ||
		message.Info.Timestamp < uint64(time.Now().Unix()-30) ||
		message.Info.FromMe || isPrivateChat(message) {
		return
	}
	addSenderJid(&message)
	log.Printf("%v %v", message.Info.RemoteJid, message.Text)
	h.wac.Read(message.Info.SenderJid, message.Info.Id)
	text := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: message.Info.RemoteJid,
		},
		ContextInfo: whatsapp.ContextInfo{
			QuotedMessage: &proto.Message{
				Conversation: &message.Text,
			},
			QuotedMessageID: message.Info.Id,
			Participant:     message.Info.SenderJid,
		},
		Text: fmt.Sprintf("Kamu bilang: %v", message.Text),
	}
	h.wac.Send(text)
}

func isPrivateChat(textMsg whatsapp.TextMessage) bool {
	return !strings.HasSuffix(textMsg.Info.RemoteJid, "g.us")
}

func addSenderJid(textMsg *whatsapp.TextMessage) {
	textMsg.Info.SenderJid = textMsg.Info.RemoteJid
	if len(textMsg.Info.Source.GetParticipant()) != 0 {
		textMsg.Info.SenderJid = textMsg.Info.Source.GetParticipant()
	}
}
