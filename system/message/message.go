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
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
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
	bot := botNumber + "@s.whatsapp.net"
	isOwner := strings.Contains(sender, owner)
	isAdmin := m.GetGroupAdmin(from, sender)
	isBotAdm := m.GetGroupAdmin(from, bot)
	isGroup := msg.Info.IsGroup
	args := strings.Split(m.GetCMD(), " ")
	query := strings.Join(args[1:], ` `)
	// extended := msg.Message.GetExtendedTextMessage()
	// contextInfo := extended.GetContextInfo()
	// quotedMsg := extended.GetContextInfo().GetQuotedMessage()
	//quotedImage := quotedMsg.GetImageMessage()
	//quotedVideo := quotedMsg.GetVideoMessage()
	//quotedSticker := quotedMsg.GetStickerMessage()

	//-- CONSOLE LOG
	// fmt.Println(msg)
	fmt.Println("\n===============================\nNAME: " + pushName + "\nJID: " + sender + "\nTYPE: " + msg.Info.Type + "\nMessage: " + m.GetCMD() + "")
	//fmt.Println(m.Msg.Message.GetPollUpdateMessage().GetMetadata())

	// Self

	if self && !isOwner {
		return
	}

	if !strings.HasPrefix(args[0], prefix) {
		command := strings.ToLower(args[0])

		switch command {
		case "bot":
			m.Reply("Bot Active " + m.Msg.Info.PushName)
			break
		}
		return
	}

	// response command if chat with prefix
	if strings.HasPrefix(args[0], prefix) {
		command := strings.ToLower(args[0])
		command = strings.Split(command, prefix)[1]

		switch command {
		case "ping":
			now := time.Now()
			mdate := time.Unix(m.Msg.Info.Timestamp.Unix(), 0)
			mtime := now.Sub(mdate)
			ms := mtime.Seconds()
			pingStr := fmt.Sprintf("Pong!\n%.3f seconds", ms)
			m.Reply(pingStr)
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

			// delete message
		case "del":
			if isGroup {
				if !isBotAdm {
					m.Reply(helpers.BotNotAdmin)
					return
				}
				if !isAdmin {
					m.Reply(helpers.NotAdmin)
					return
				}

				if m.Msg.Message.ExtendedTextMessage.ContextInfo.Participant != nil {
					var target types.JID
					ctx := m.Msg.Message.ExtendedTextMessage.ContextInfo
					stanza := ctx.StanzaId
					participant, _ := types.ParseJID(*ctx.Participant)
					parse_bot_number, _ := types.ParseJID(bot)
					if participant == parse_bot_number {
						target = types.EmptyJID
					} else {
						target = participant
					}

					_, err := sock.SendMessage(context.Background(), from, sock.BuildRevoke(from, target, *stanza))
					if err != nil {
						log.Println(err)
						return
					}
					// log.Println(res)
				}
			}

			break
			// group command
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
					m.Reply("Nomor di private")
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

			}
			break
		case "ai":
			if query == "" {
				m.Reply("Masukkan pertanyaanmu")
				return
			}

			type Data struct {
				Status bool   `json:"status"`
				Data   string `json:"data"`
			}

			data := Data{}

			m.React("⏱️")

			apiUrl := "https://vihangayt.me/tools/chatgptv4"
			params := url.Values{}
			params.Add("q", query)

			// Membuat URL dengan query parameters
			fullURL, err := url.ParseRequestURI(apiUrl)
			if err != nil {
				log.Println("Error parsing URL:", err)
				m.React("❌")
				return
			}
			fullURL.RawQuery = params.Encode()

			err = lib.ReqGet(fullURL.String(), &data)
			if err != nil {
				m.Reply("Error: " + err.Error())
				m.React("❌")
				return
			}

			m.Reply(data.Data)
			m.React("✅")

			break
		}
	}
	return

}
