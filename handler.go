package main

import (
	"log"
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
		for {
			log.Println("waiting for 30 seconds ...")
			<-time.After(30 * time.Second)
			log.Println("reconnecting ...")
			err := h.wac.Restore()
			if err == nil {
				break
			}
		}
	} else {
		log.Printf("error occoured: %v\n", err)
	}
}

// HandleTextMessage handles a text message
func (h Handler) HandleTextMessage(message whatsapp.TextMessage) {
	if message.Info.Timestamp < uint64(h.startTime.Unix()) ||
		message.Info.Timestamp < uint64(time.Now().Unix()-30) ||
		message.Info.FromMe || message.Info.RemoteJid == "status@broadcast" {
		return
	}

	addSenderJid(&message)

	log.Printf("%v %v", message.Info.RemoteJid, message.Text)
	replyMsg, ok := getReply(h, &message)
	if !ok {
		return
	}

	defaultMessageInfo := whatsapp.MessageInfo{
		RemoteJid: message.Info.RemoteJid,
	}
	defaultMessageContextInfo := whatsapp.ContextInfo{
		QuotedMessage: &proto.Message{
			Conversation: &message.Text,
		},
		QuotedMessageID: message.Info.Id,
		Participant:     message.Info.SenderJid,
	}

	switch m := replyMsg.(type) {
	case whatsapp.TextMessage:
		m.Info = defaultMessageInfo
		m.ContextInfo = defaultMessageContextInfo
		h.wac.Send(m)
	case whatsapp.ImageMessage:
		m.Info = defaultMessageInfo
		m.ContextInfo = defaultMessageContextInfo
		h.wac.Send(m)
	}
}
