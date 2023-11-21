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
	"fmt"
	"gowabot/system/lib"
	"image/jpeg"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/amiruldev20/waSocket"
	"github.com/amiruldev20/waSocket/types/events"
	"github.com/nickalie/go-webpbin"
	// "github.com/joho/godotenv"
	// "github.com/nickalie/go-webpbin"
	// "github.com/chai2010/webp"
)

func Msg(sock *waSocket.Client, msg *events.Message) {

	// err := godotenv.Load()
	// if err != nil {
	// 	panic("Error load file .env")
	// }

	var (
		prefix  = os.Getenv("BOT_PREFIX")
		self, _ = strconv.ParseBool(strings.ToLower(os.Getenv("BOT_SELF")))
		owner   = os.Getenv("OWNER_NUMBER")
	)

	// botNumber := os.Getenv("BOT_NUMBER")

	/* my function */
	m := lib.NewSimp(sock, msg)
	//from := msg.Info.Chat
	sender := msg.Info.Sender.String()
	pushName := msg.Info.PushName
	isOwner := strings.Contains(sender, owner)
	//isAdmin := m.GetGroupAdmin(from, sender)
	//isBotAdm := m.GetGroupAdmin(from, botNumber + "@s.whatsapp.net")
	//isGroup := msg.Info.IsGroup
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

	// response command if chat with prefix
	if strings.HasPrefix(args[0], prefix) {
		command := strings.ToLower(args[0])
		command = strings.Split(command, prefix)[1]

		// Self
		if self && !isOwner {
			return
		}

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
			m.SendImg(m.Msg.Info.Chat, query)
			break
		}
	}

}
