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
		message.Info.FromMe {
		return
	}
	addSenderJid(&message)
	log.Printf("%v %v", message.Info.RemoteJid, message.Text)
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
		Text: getTextReply(h, &message),
	}
	h.wac.Send(text)
}

func getTextReply(h Handler, message *whatsapp.TextMessage) string {
	if isPrivateChat(message) {
		return "Hi! Namaku *Iceburg*. Undang aku ke grup kelas" +
			"mu supaya aku bisa mengelola tugas kalian."
	} else if cond, _ := groupExists(message.Info.RemoteJid); !cond {
		if message.Text == "@register" {
			groupName := h.wac.Store.Contacts[message.Info.RemoteJid].Name
			checkError(addGroup(message.Info.RemoteJid, groupName))
			return fmt.Sprintf("Sukses! Grup ```%v``` sukses terdaftar."+
				" Gunakan perintah _@unregister_ untuk mereset data.", groupName)
		}
		return "Hi! Namaku *Iceburg*. Daftarkan grup ini dengan mengirim pe" +
			"rintah _@register_."
	} else if message.Text == "@unregister" {
		checkError(deleteGroup(message.Info.RemoteJid))
		return "Selamat tinggal! Terima kasih telah menggunakan *Iceburg*."
	}
	return "woy"
}

func isPrivateChat(message *whatsapp.TextMessage) bool {
	return !strings.HasSuffix(message.Info.RemoteJid, "g.us")
}

func addSenderJid(message *whatsapp.TextMessage) {
	message.Info.SenderJid = message.Info.RemoteJid
	if len(message.Info.Source.GetParticipant()) != 0 {
		message.Info.SenderJid = message.Info.Source.GetParticipant()
	}
}
