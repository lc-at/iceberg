package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Rhymen/go-whatsapp"
	"github.com/disintegration/imaging"
	"github.com/p4kl0nc4t/iceberg/integration/wolfram"
)

func getReply(h Handler, message *whatsapp.TextMessage) (interface{}, bool) {
	defaultMessage := whatsapp.TextMessage{}
	if !isGroupChat(message) {
		defaultMessage.Text = cnf.getMessageTemplate("private_chat")
	} else if cond, _ := (&groupModel{JID: message.Info.RemoteJid}).isExist(); !cond {
		if message.Text == "@register" {
			groupName := h.wac.Store.Contacts[message.Info.RemoteJid].Name
			checkError((&groupModel{message.Info.RemoteJid, groupName}).add())
			defaultMessage.Text = fmt.Sprintf(cnf.getMessageTemplate("register_success"), groupName)
		} else {
			defaultMessage.Text = cnf.getMessageTemplate("not_registered")
		}
	} else if message.Text == "@unregister" {
		checkError((&groupModel{JID: message.Info.RemoteJid}).delete())
		defaultMessage.Text = cnf.getMessageTemplate("unregister")
	} else if text := getTextReply(h, message); len(text) > 0 {
		defaultMessage.Text = text
	} else {
		if customReply, ok := getCustomReply(h, message); ok {
			if text, ok := customReply.(string); ok {
				defaultMessage.Text = text
			} else {
				return customReply, true
			}
		}
	}

	ok := true
	if len(defaultMessage.Text) == 0 {
		ok = false
	}

	return defaultMessage, ok
}

func getTextReply(h Handler, message *whatsapp.TextMessage) string {
	switch {
	case message.Text == "@ping":
		return "Pong!"
	case message.Text == "@menu":
		return cnf.getMessageTemplate("menu")
	case strings.HasPrefix(message.Text, "@tambah"):
		args := strings.SplitN(message.Text, " ", 3)
		if message.ContextInfo.QuotedMessageID == "" {
			return cnf.getMessageTemplate("no_assignment_description")
		} else if len(args) != 3 {
			return cnf.getMessageTemplate("invalid_add_assignment_args")
		} else if len(args[2]) > 30 || len(args[1]) > 10 {
			return cnf.getMessageTemplate("assignment_too_long")
		}
		checkError((&assignmentModel{
			GroupJID:    message.Info.RemoteJid,
			Subject:     args[1],
			Description: message.ContextInfo.QuotedMessage.GetConversation(),
			Deadline:    args[2],
		}).add())
		return fmt.Sprintf(cnf.getMessageTemplate("assignment_added"))
	case strings.HasPrefix(message.Text, "@hapus"):
		args := strings.SplitN(message.Text, " ", 2)
		if len(args) != 2 {
			return cnf.getMessageTemplate("invalid_args")
		}
		var assignmentID []int
		for _, v := range strings.Split(args[1], ",") {
			ID, err := strconv.Atoi(strings.TrimSpace(v))
			assignment := &assignmentModel{
				ID:       ID,
				GroupJID: message.Info.RemoteJid,
			}
			if err != nil {
				return cnf.getMessageTemplate("invalid_args")
			} else if cond, _ := assignment.isExist(); !cond {
				return cnf.getMessageTemplate("invalid_assignment_id")
			}
			assignmentID = append(assignmentID, ID)
		}
		for _, ID := range assignmentID {
			checkError((&assignmentModel{ID: ID,
				GroupJID: message.Info.RemoteJid}).delete())
		}
		return cnf.getMessageTemplate("assignment_deleted")
	case message.Text == "@tugas":
		assignmentRows, err := (&assignmentModel{
			GroupJID: message.Info.RemoteJid}).query()
		checkError(err)
		var assignments []string
		for _, assignment := range assignmentRows {
			assignment.humanReadableValues()
			assignments = append(assignments, fmt.Sprintf(
				cnf.getMessageTemplate("assignment_item"),
				assignment.Subject, assignment.ID,
				assignment.Description, assignment.Deadline,
			))
		}
		var formattedAssignments string
		if len(assignments) == 0 {
			formattedAssignments = cnf.getMessageTemplate("empty_assignment_list")
		} else {
			formattedAssignments = strings.Join(assignments, "\n")
		}
		date := time.Now().Format("2006/01/02 15:04:05")
		dayname, _ := cnf.getNameByDay(int(time.Now().Weekday()))
		return fmt.Sprintf(
			cnf.getMessageTemplate("assignment_list"),
			fmt.Sprintf("%s, %s", strings.Title(dayname), date),
			formattedAssignments)
	case message.Text == "@tentang":
		return cnf.getMessageTemplate("about")
	default:
		return ""
	}

}

func getCustomReply(h Handler, message *whatsapp.TextMessage) (interface{}, bool) {
	if !strings.HasPrefix(message.Text, "@wolfram") || len(cnf.WolframAlphaAppID) == 0 {
		return nil, false
	}

	args := strings.SplitN(message.Text, " ", 2)
	if len(args) != 2 || len(strings.TrimSpace(args[1])) == 0 {
		return cnf.getMessageTemplate("invalid_args"), true
	}

	client := wolfram.Client{AppID: cnf.WolframAlphaAppID}
	result, err := client.Simple(args[1])

	if err != nil {
		if err == wolfram.ErrInvalidInput {
			return cnf.getMessageTemplate("wolfram_bad_input"), true
		}
		log.Printf("error: %v\n", err)
		return cnf.getMessageTemplate("wolfram_error"), true

	}

	thumb, err := getThumbnail(bytes.NewReader(result))

	if err != nil {
		fmt.Println(err.Error())
		log.Printf("error: %v\n", err)
	}

	return whatsapp.ImageMessage{
		Type:      "image/jpeg",
		Content:   bytes.NewReader(result),
		Thumbnail: thumb,
	}, true
}

func isGroupChat(message *whatsapp.TextMessage) bool {
	return strings.HasSuffix(message.Info.RemoteJid, "g.us")
}

func addSenderJid(message *whatsapp.TextMessage) {
	message.Info.SenderJid = message.Info.RemoteJid
	if len(message.Info.Source.GetParticipant()) != 0 {
		message.Info.SenderJid = message.Info.Source.GetParticipant()
	}
}

func getThumbnail(image io.Reader) ([]byte, error) {
	img, err := imaging.Decode(image)
	if err != nil {
		return nil, err
	}

	b := img.Bounds()
	imgWidth := b.Max.X
	imgHeight := b.Max.Y

	thumbWidth := 100
	thumbHeight := 100

	if imgWidth > imgHeight {
		thumbHeight = 56
	} else {
		thumbWidth = 56
	}

	thumb := imaging.Thumbnail(img, thumbWidth, thumbHeight, imaging.CatmullRom)

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, thumb, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
