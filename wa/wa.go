package wa

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"main/ai"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	"google.golang.org/protobuf/proto"
)

var Client *whatsmeow.Client
var BotStartTime time.Time

//----------------------------------------------------
// Mengambil Text Message
//----------------------------------------------------

func GetTextMessage(msg *waProto.Message) string {

	if msg == nil {
		return ""
	}

	if msg.GetConversation() != "" {
		return msg.GetConversation()
	}

	if msg.GetExtendedTextMessage() != nil {
		return msg.GetExtendedTextMessage().GetText()
	}

	return ""
}

//----------------------------------------------------
// Mengirim Pesan
//----------------------------------------------------

func SendText(chat types.JID, text string) error {

	_, err := Client.SendMessage(
		context.Background(),
		chat,
		&waProto.Message{
			Conversation: proto.String(text),
		},
	)

	return err
}

//----------------------------------------------------
// Event Handler
//----------------------------------------------------

func eventHandler(evt any) {

	switch v := evt.(type) {

	case *events.Connected:

		fmt.Println("===================================")
		fmt.Println("WhatsApp Connected")
		fmt.Println("===================================")

	case *events.Disconnected:

		fmt.Println("===================================")
		fmt.Println("WhatsApp Disconnected")
		fmt.Println("===================================")

	case *events.Message:

		// Jangan balas pesan sendiri
		if v.Info.IsFromMe {
			return
		}

		// Jangan balas grup
		if v.Info.IsGroup {
			return
		}

		// Abaikan pesan lama
		if v.Info.Timestamp.Before(BotStartTime) {
			return
		}

		pesan := strings.TrimSpace(GetTextMessage(v.Message))

		if pesan == "" {
			return
		}

		fmt.Println("===================================")
		fmt.Println("Pesan Baru")
		fmt.Println("Dari :", v.Info.Sender)
		fmt.Println("Isi  :", pesan)
		fmt.Println("===================================")

		cmd := strings.ToLower(pesan)

		var balasan string

		switch cmd {

		case "menu":

			balasan = `===== MENU BOT =====

1. menu
2. info
3. prodi
4. universitas
5. hai

Untuk menggunakan AI :

!ai pertanyaan

Contoh :

!ai Apa itu Golang?
!ai Jelaskan Cloud Computing`

		case "info":

			balasan = `===== INFORMASI BOT =====

Nama :
Cloud Computing Bot

Versi :
1.0

Teknologi :

- Golang
- Gin
- WhatsMeow
- MySQL
- AI Qwen
- llama.cpp`

		case "hai":

			balasan = `Halo 👋

Selamat datang.

Silakan ketik menu
untuk melihat daftar perintah.`

		case "prodi":

			balasan = `===== PROGRAM STUDI =====

1. Sistem Informasi

2. Informatika

3. Manajemen

4. Akuntansi

5. Kebidanan

6. Keperawatan`

		case "universitas":

			balasan = `Universitas Duta Bangsa Surakarta

Alamat :

Jl. Bhayangkara No.55,
Tipes,
Serengan,
Surakarta

Website :

https://udb.ac.id`

		default:

			// hanya jika memakai !ai

			if strings.HasPrefix(cmd, "!ai") {

				prompt := strings.TrimSpace(
					strings.TrimPrefix(pesan, "!ai"),
				)

				if prompt == "" {

					balasan = `Contoh :

!ai Apa itu Golang?

!ai Jelaskan Docker

!ai Buat Pantun`

				} else {

					fmt.Println("Mengirim ke AI...")

					var err error

					balasan, err = ai.AskAI(prompt)

					if err != nil {

						fmt.Println(err)

						balasan = "Maaf, AI sedang offline."
					}
				}

			} else {

				balasan = `Perintah tidak dikenali.

Silakan ketik :

menu

untuk melihat daftar perintah.`
			}
		}

		err := SendText(v.Info.Chat, balasan)

		if err != nil {

			fmt.Println("Gagal mengirim :", err)

		} else {

			fmt.Println("Balasan berhasil dikirim")
		}
	}
}

//----------------------------------------------------
// Koneksi WhatsApp
//----------------------------------------------------

func KonekWa() {

	ctx := context.Background()

	dbLog := waLog.Stdout("Database", "INFO", true)

	container, err := sqlstore.New(
		ctx,
		"sqlite3",
		"file:examplestore.db?_foreign_keys=on",
		dbLog,
	)

	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice(ctx)

	if err != nil {
		panic(err)
	}

	clientLog := waLog.Stdout("Client", "INFO", true)

	Client = whatsmeow.NewClient(deviceStore, clientLog)

	Client.AddEventHandler(eventHandler)

	BotStartTime = time.Now()

	if Client.Store.ID == nil {

		fmt.Println("Belum ada session WhatsApp")

		qrChan, err := Client.GetQRChannel(ctx)

		if err != nil {
			panic(err)
		}

		err = Client.Connect()

		if err != nil {
			panic(err)
		}

		for evt := range qrChan {

			switch evt.Event {

			case "code":

				fmt.Println("==============================")
				fmt.Println("Scan QR WhatsApp")
				fmt.Println("==============================")

				qrterminal.GenerateHalfBlock(
					evt.Code,
					qrterminal.L,
					os.Stdout,
				)

			case "success":

				fmt.Println("Login Berhasil")

			default:

				fmt.Println("QR Event :", evt.Event)
			}
		}

	} else {

		fmt.Println("Session ditemukan")

		err = Client.Connect()

		if err != nil {
			panic(err)
		}
	}

}

func StopWA() {
	if Client != nil {
		Client.Disconnect()
		fmt.Println("===================================")
		fmt.Println("WhatsApp Disconnected")
		fmt.Println("===================================")
	}
}
