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
package message

import (
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"whatsapp-bot-go/system/helpers"
	"whatsapp-bot-go/system/lib"

	"github.com/amiruldev20/waSocket"
	waProto "github.com/amiruldev20/waSocket/binary/proto"
	"github.com/amiruldev20/waSocket/types"
	"github.com/amiruldev20/waSocket/types/events"
	"google.golang.org/protobuf/proto"

	"github.com/joho/godotenv"
	"github.com/nickalie/go-webpbin"
)

func Msg(sock *waSocket.Client, msg *events.Message) {

	err := godotenv.Load()
	if err != nil {
		panic("Error load file .env")
	}

	var (
		prefix    = os.Getenv("BOT_PREFIX")
		self, _   = strconv.ParseBool(strings.ToLower(os.Getenv("BOT_SELF")))
		owner     = os.Getenv("OWNER_NUMBER")
		botNumber = os.Getenv("BOT_NUMBER")
	)

	/* my function */
	m := lib.NewSimp(sock, msg)
	from := msg.Info.Chat
	sender := msg.Info.Sender.String()
	pushName := msg.Info.PushName
	isOwner := strings.Contains(sender, owner)
	isAdmin := m.GetGroupAdmin(from, sender)
	isBotAdm := m.GetGroupAdmin(from, botNumber+"@s.whatsapp.net")
	isGroup := msg.Info.IsGroup
	args := strings.Split(m.GetCMD(), " ")
	query := strings.Join(args[1:], ` `)
	//extended := msg.Message.GetExtendedTextMessage()
	//quotedMsg := extended.GetContextInfo().GetQuotedMessage()
	//quotedImage := quotedMsg.GetImageMessage()
	//quotedVideo := quotedMsg.GetVideoMessage()
	//quotedSticker := quotedMsg.GetStickerMessage()

	//-- CONSOLE LOG
	// fmt.Println(msg)
	fmt.Println("\n===============================\nNAME: " + pushName + "\nJID: " + sender + "\nTYPE: " + msg.Info.Type + "\nMessage: " + m.GetCMD() + "")
	//fmt.Println(m.Msg.Message.GetPollUpdateMessage().GetMetadata())

	if self && !isOwner {
		return
	}

	if !strings.HasPrefix(args[0], prefix) {
		command := strings.ToLower(args[0])

		switch command {
		//
		}
	}

	// response command if chat with prefix
	if strings.HasPrefix(args[0], prefix) {
		command := strings.ToLower(args[0])
		command = strings.Split(command, prefix)[1]

		// Self

		switch command {
		case "bot":
			m.Reply("Bot Active")
			break
		case "ping":
			m.Reply("Pong!")
			break
		case "st":
			if msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage != nil {
				quoted := msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage
				if quoted.ImageMessage == nil {
					m.Reply("Reply gambar lalu ketik /st")
					return
				}

				imgmsg := quoted.GetImageMessage()
				data, err := sock.Download(imgmsg)

				if err != nil {
					log.Println(err)
					return
				}

				randomJpgImg := "./temp/" + lib.GenerateRandomString(5) + ".jpg"
				randomWebpImg := "./temp/" + lib.GenerateRandomString(5) + ".webp"
				if err := os.WriteFile(randomJpgImg, data, 0600); err != nil {
					log.Printf("Failed to save image: %v", err)
					return
				}

				log.Printf("Saved image in %s", randomJpgImg)

				imgbyte, err := os.ReadFile(randomJpgImg)
				if err != nil {
					fmt.Println("Error reading file:", err)
					return
				}

				decodeImg, err := jpeg.Decode(bytes.NewReader(imgbyte))
				if err != nil {
					fmt.Println("Error decoding file:", err)
					return
				}

				fmt.Println("convert jpg to webp...")
				f, err := os.Create(randomWebpImg)

				if err != nil {
					log.Println(err)
					return
				}

				if err := webpbin.Encode(f, decodeImg); err != nil {
					f.Close()
					log.Println(err)
					return
				}

				if err := f.Close(); err != nil {
					log.Println(err)
					return
				}

				fmt.Println("Success convert to webp")
				webpByte, err := os.ReadFile(randomWebpImg)
				if err != nil {
					fmt.Println("Error reading file:", err)
					return
				}

				fmt.Println("Sending webp as sticker...")

				m.ReplyAsSticker(webpByte)

				// delete file image
				err = os.Remove(randomJpgImg)
				err = os.Remove(randomWebpImg)
			}
			break
		case "sendimg":
			if query == "" {
				m.Reply("Url cannot empty")
				return
			}

			m.SendImg(from, query)
			break

			//command create channel
			// jangan brutal ntar turu nangid :'(
		case "nc":
			if !isOwner {
				m.Reply("Hanya owner bot")
				return
			}
			split := strings.Split(query, "|")
			title := split[0]
			desc := strings.Join(split[1:], " ")
			m.CreateChannel(title, desc)
			break
		case "add":
			if !isGroup {
				m.Reply(helpers.NotGroup)
				return
			}
			if query == "" {
				m.Reply(fmt.Sprintf("Contoh penggunaan:\n%sadd 628xxxxx", prefix))
				return
			}
			if !isBotAdm {
				m.Reply(helpers.BotNotAdmin)
				return
			}
			if !isAdmin {
				m.Reply(helpers.NotAdmin)
				return
			}

			m.React("⏱️")

			ok, err := sock.IsOnWhatsApp([]string{query})

			if err != nil {
				log.Println("Error:", err)
				return
			}

			if len(ok) == 0 {
				return
			}

			if !ok[0].IsIn {
				m.React("❌")
				m.Reply("Nomor tidak terdaftar di WhatsApp")
				return
			}

			res, err := sock.UpdateGroupParticipants(from, []types.JID{ok[0].JID}, waSocket.ParticipantChangeAdd)

			if err != nil {
				log.Println("Error adding participant:", err)
				return
			}

			for _, item := range res {
				if item.Status == "403" {
					info, _ := sock.GetGroupInfo(from)
					pp, _ := sock.GetProfilePictureInfo(from, &waSocket.GetProfilePictureParams{
						Preview: true,
					})
					getimg, err := http.Get(pp.URL)
					if err != nil {
						log.Println("Error getting the image:", err)
						return
					}
					defer getimg.Body.Close()
					imgByte, _ := io.ReadAll(getimg.Body)
					exp, _ := strconv.ParseInt(item.Content.Attrs["expiration"].(string), 10, 64)
					log.Printf("\nParticipant is private: %s %s %s %d", item.Status, item.JID, item.Content.Attrs["code"].(string), exp)
					sock.SendMessage(context.TODO(), item.JID, &waProto.Message{
						GroupInviteMessage: &waProto.GroupInviteMessage{
							InviteCode:       proto.String(item.Content.Attrs["code"].(string)),
							InviteExpiration: proto.Int64(exp),
							GroupJid:         proto.String(info.JID.String()),
							GroupName:        proto.String(info.Name),
							Caption:          proto.String(info.Topic),
							JpegThumbnail:    imgByte,
						},
					})
					m.React("⚠️")
				} else if item.Status == "409" {
					log.Printf("\nParticipant already in group: %s %s %+v", item.Status, item.JID, item.Content)
					m.React("❌")
				} else if item.Status == "200" {
					log.Printf("\nAdded participant: %s %s %+v", item.Status, item.JID, item.Content)
					m.React("✅")
				} else {
					log.Printf("\nUnknown status: %s %s %+v", item.Status, item.JID, item.Content)
					m.React("❌")
				}
			}

			break

		case "kick":
			if !isGroup {
				m.Reply(helpers.NotGroup)
				return
			}
			if query == "" {
				m.Reply(fmt.Sprintf("Contoh penggunaan:\n%skick @mention", prefix))
				return
			}
			if !isBotAdm {
				m.Reply(helpers.BotNotAdmin)
				return
			}
			if !isAdmin {
				m.Reply(helpers.NotAdmin)
				return
			}
			m.React("⏱️")
			if m.Msg.Message.ExtendedTextMessage.ContextInfo.MentionedJid != nil {
				participant := m.Msg.Message.ExtendedTextMessage.ContextInfo.MentionedJid[0]
				parse_participant, _ := types.ParseJID(participant)

				_, err := sock.UpdateGroupParticipants(from, []types.JID{parse_participant}, waSocket.ParticipantChangeRemove)

				if err != nil {
					log.Println("Error removing participant:", err)
					return
				}
				m.React("✅")

				// m.Reply("Sayonara")
			}
			break
		case "pp":
			sock.GetProfilePictureInfo(from, &waSocket.GetProfilePictureParams{
				Preview: true,
			})
			break
		}
	}

}
