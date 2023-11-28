/*
###################################
# Name: Mywa BOT                  #
# Version: 1.0.1                  #
# Developer: Amirul Dev           #
# Library: waSocket               #
# Contact: 085157489446           #
###################################
# Thanks to:
# Vnia
*/
package lib

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"strings"
	"whatsapp-bot-go/system/dto"

	"github.com/amiruldev20/waSocket"
	waProto "github.com/amiruldev20/waSocket/binary/proto"
	"github.com/amiruldev20/waSocket/types"
	"github.com/amiruldev20/waSocket/types/events"

	"google.golang.org/protobuf/proto"
)

type renz struct {
	sock *waSocket.Client
	Msg  *events.Message
}

func NewSimp(Cli *waSocket.Client, m *events.Message) *renz {
	return &renz{
		sock: Cli,
		Msg:  m,
	}
}

/* parse jid */
func (m *renz) parseJID(arg string) (types.JID, bool) {
	if arg[0] == '+' {
		arg = arg[1:]
	}
	if !strings.ContainsRune(arg, '@') {
		return types.NewJID(arg, types.DefaultUserServer), true
	} else {
		recipient,
			err := types.ParseJID(arg)
		if err != nil {
			fmt.Printf("Invalid JID %s: %v\n", arg, err)
			return recipient, false
		} else if recipient.User == "" {
			fmt.Printf("Invalid JID %s: no server specified\n", arg)
			return recipient, false
		}
		return recipient,
			true
	}
}

/* send react */
func (m *renz) React(react string) {
	_,
		err := m.sock.SendMessage(context.Background(), m.Msg.Info.Chat, m.sock.BuildReaction(m.Msg.Info.Chat, m.Msg.Info.Sender, m.Msg.Info.ID, react))
	if err != nil {
		return
	}
}

/* send message */
func (m *renz) SendMsg(jid types.JID, teks string) {
	_,
		err := m.sock.SendMessage(context.Background(), jid, &waProto.Message{Conversation: proto.String(teks)})
	if err != nil {
		return
	}
}

/* send sticker */
func (m *renz) SendSticker(jid types.JID, data []byte, extra ...dto.ExtraSend) {
	var contextInfo *waProto.ContextInfo
	var req dto.ExtraSend
	if len(extra) > 1 {
		log.Println("only one extra parameter may be provided to SendMessage")
		return
	} else if len(extra) == 1 {
		req = extra[0]
	}

	if req.Reply {
		// Isi contextInfo jika Reply adalah true
		contextInfo = &waProto.ContextInfo{
			Expiration:    proto.Uint32(86400),
			StanzaId:      &m.Msg.Info.ID,
			Participant:   proto.String(m.Msg.Info.Sender.String()),
			QuotedMessage: m.Msg.Message,
		}
	}

	uploadImg, err := m.sock.Upload(context.Background(), data, waSocket.MediaImage)

	if err != nil {
		log.Println(err)
		return
	}

	_, err = m.sock.SendMessage(context.Background(), m.Msg.Info.Chat, &waProto.Message{
		StickerMessage: &waProto.StickerMessage{
			Url:           proto.String(uploadImg.URL),
			FileSha256:    uploadImg.FileSHA256,
			FileEncSha256: uploadImg.FileEncSHA256,
			MediaKey:      uploadImg.MediaKey,
			Mimetype:      proto.String(http.DetectContentType(data)),
			DirectPath:    proto.String(uploadImg.DirectPath),
			FileLength:    proto.Uint64(uint64(len(data))),
			ContextInfo:   contextInfo,
			// FirstFrameSidecar: data,
			// PngThumbnail:      data,
		},
	})

	if err != nil {
		log.Println(err)
		return
	}
}

/* send image */
func (m *renz) SendImg(jid types.JID, data []byte) {

	uploadedImg, err := m.sock.Upload(context.Background(), data, waSocket.MediaImage)

	if err != nil {
		log.Println("Failed to upload file:", err)
		return
	}

	hs := sha256.New()
	// hs.Write(imgByte)
	_, err = m.sock.SendMessage(context.Background(), m.Msg.Info.Chat, &waProto.Message{
		ImageMessage: &waProto.ImageMessage{
			JpegThumbnail: hs.Sum(data), // blm work
			Url:           proto.String(uploadedImg.URL),
			DirectPath:    proto.String(uploadedImg.DirectPath),
			MediaKey:      uploadedImg.MediaKey,
			Mimetype:      proto.String(http.DetectContentType(data)),
			FileEncSha256: uploadedImg.FileEncSHA256,
			FileSha256:    uploadedImg.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
		},
	})

	if err != nil {
		log.Println(err)
		return
	}

}

/* delete message */

// To delete someone else's message in the group, fill the sender parameter with the participant ID.
//
// If you want to delete your own messages, fill in the sender parameter with types.EmptyJID.
//
// MessageID parameter can be filled with stanzaId
func (m *renz) DeleteMsg(chat types.JID, sender types.JID, messageID string) {
	_, err := m.sock.SendMessage(context.Background(), chat, m.sock.BuildRevoke(chat, sender, messageID))
	if err != nil {
		log.Println("Error deleting message:", err)
		return
	}
}

/* send reply */
func (m *renz) Reply(teks string) {
	_,
		err := m.sock.SendMessage(context.Background(), m.Msg.Info.Chat, &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: proto.String(teks),
			ContextInfo: &waProto.ContextInfo{
				Expiration:    proto.Uint32(86400),
				StanzaId:      &m.Msg.Info.ID,
				Participant:   proto.String(m.Msg.Info.Sender.String()),
				QuotedMessage: m.Msg.Message,
			},
		},
	})
	if err != nil {
		return
	}
}

/* send replyAsSticker */
func (m *renz) ReplyAsSticker(data []byte) {
	m.SendSticker(m.Msg.Info.Chat, data, dto.ExtraSend{Reply: true})
}

/* send adReply */
func (m *renz) ReplyAd(teks string) {
	var isImage = waProto.ContextInfo_ExternalAdReplyInfo_IMAGE
	_, err := m.sock.SendMessage(context.Background(), m.Msg.Info.Chat, &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: proto.String(teks),
			ContextInfo: &waProto.ContextInfo{
				ExternalAdReply: &waProto.ContextInfo_ExternalAdReplyInfo{
					Title:                 proto.String("MywaBOT 2023"),
					Body:                  proto.String("Made with waSocket by Amirul Dev"),
					MediaType:             &isImage,
					ThumbnailUrl:          proto.String("https://telegra.ph/file/eb7261ee8de82f8f48142.jpg"),
					MediaUrl:              proto.String("https://wa.me/stickerpack/amirul.dev"),
					SourceUrl:             proto.String("https://chat.whatsapp.com/ByQt0u0bz4NJfNPEUfDHps"),
					ShowAdAttribution:     proto.Bool(true),
					RenderLargerThumbnail: proto.Bool(true),
				},
				Expiration:    proto.Uint32(86400),
				StanzaId:      &m.Msg.Info.ID,
				Participant:   proto.String(m.Msg.Info.Sender.String()),
				QuotedMessage: m.Msg.Message,
			},
		},
	})
	if err != nil {
		return
	}
}

/* send contact */
func (m *renz) SendContact(jid types.JID, number string, nama string) {
	_,
		err := m.sock.SendMessage(context.Background(), jid, &waProto.Message{
		ContactMessage: &waProto.ContactMessage{
			DisplayName: proto.String(nama),
			Vcard:       proto.String(fmt.Sprintf("BEGIN:VCARD\nVERSION:3.0\nN:%s;;;\nFN:%s\nitem1.TEL;waid=%s:+%s\nitem1.X-ABLabel:Mobile\nEND:VCARD", nama, nama, number, number)),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:      &m.Msg.Info.ID,
				Participant:   proto.String(m.Msg.Info.Sender.String()),
				QuotedMessage: m.Msg.Message,
			},
		},
	})
	if err != nil {
		return
	}
}

/* create channel */
func (m *renz) CreateChannel(title, description string) {
	metadata,
		err := m.sock.CreateNewsletter(waSocket.CreateNewsletterParams{
		Name:        title,
		Description: description,
		// Picture: profilePicture,
	})
	if err != nil {
		m.Reply("Error creating channel:" + err.Error())
		return
	}
	jid := metadata.ID
	m.Reply(fmt.Sprintf("Success create channel\nJID: %s\nName: %s\nDescription: %s\nLink: https://whatsapp.com/channel/%s", jid, metadata.ThreadMeta.Name.Text, metadata.ThreadMeta.Description.Text, metadata.ThreadMeta.InviteCode))
}

/* fetch group admin */
func (m *renz) FetchGroupAdmin(Jid types.JID) ([]string, error) {
	var Admin []string
	resp, err := m.sock.GetGroupInfo(Jid)
	if err != nil {
		return Admin, err
	} else {
		for _, group := range resp.Participants {
			if group.IsAdmin || group.IsSuperAdmin {
				Admin = append(Admin, group.JID.String())
			}
		}
	}
	return Admin, nil
}

/* get group admin */
func (m *renz) GetGroupAdmin(jid types.JID, sender string) bool {
	if !m.Msg.Info.IsGroup {
		return false
	}
	admin, err := m.FetchGroupAdmin(jid)
	if err != nil {
		return false
	}
	for _, v := range admin {
		if v == sender {
			return true
		}
	}
	return false
}

/* get link group */
func (m *renz) LinkGc(Jid types.JID, reset bool) string {
	link,
		err := m.sock.GetGroupInviteLink(Jid, reset)

	if err != nil {
		panic(err)
	}
	return link
}

func (m *renz) GetCMD() string {
	extended := m.Msg.Message.GetExtendedTextMessage().GetText()
	text := m.Msg.Message.GetConversation()
	imageMatch := m.Msg.Message.GetImageMessage().GetCaption()
	videoMatch := m.Msg.Message.GetVideoMessage().GetCaption()
	//pollVote := m.Msg.Message.GetPollUpdateMessage().GetVote()
	tempBtnId := m.Msg.Message.GetTemplateButtonReplyMessage().GetSelectedId()
	btnId := m.Msg.Message.GetButtonsResponseMessage().GetSelectedButtonId()
	listId := m.Msg.Message.GetListResponseMessage().GetSingleSelectReply().GetSelectedRowId()
	var command string
	if text != "" {
		command = text
	} else if imageMatch != "" {
		command = imageMatch
	} else if videoMatch != "" {
		command = videoMatch
	} else if extended != "" {
		command = extended
		/*
		   } else if pollVote != "" {
		   command = pollVote
		*/
	} else if tempBtnId != "" {
		command = tempBtnId
	} else if btnId != "" {
		command = btnId
	} else if listId != "" {
		command = listId
	}
	return command
}
