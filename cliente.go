package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time" // Agregar esta importación
)

const serverURL = "http://0.0.0.0:8080"

type Message struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Solicitar apodo al usuario
	fmt.Print("Ingresar tu nombre: ")
	nick, _ := reader.ReadString('\n')
	nick = strings.TrimSpace(nick)

	if nick == "" {
		fmt.Println("El apodo no puede estar vacío. Por favor, reinicia el programa.")
		return
	}

	go receiveMessages(nick)

	time.Sleep(1 * time.Second)

	sendMessage(Message{
		Sender:  nick,
		Content: fmt.Sprintf("%s ha ingresado al chat", nick),
	})

	// Bucle para enviar mensajes
	for {

		msgContent, _ := reader.ReadString('\n')
		msgContent = strings.TrimSpace(msgContent)

		if msgContent == "salir" {
			// Notificar salida
			sendMessage(Message{
				Sender:  nick,
				Content: fmt.Sprintf("%s ha salido del chat", nick),
			})
			fmt.Println("Has salido del chat.")
			break
		}

		if msgContent != "" {
			sendMessage(Message{
				Sender:  nick,
				Content: msgContent,
			})
		}
	}
}

func sendMessage(msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Error al codificar el mensaje:", err)
		return
	}

	_, err = http.Post(serverURL+"/send", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error al enviar el mensaje:", err)
	}
}

func receiveMessages(nick string) {
	for {
		resp, err := http.Get(serverURL + "/receive")
		if err != nil {
			fmt.Println("Error al conectarse al servidor:", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		for {
			var msg Message
			if err := decoder.Decode(&msg); err != nil {
				break
			}

			// Mostrar solo mensajes de otros usuarios
			if msg.Sender != nick {
				fmt.Printf("[%s]: %s\n", msg.Sender, msg.Content)
			}
		}
	}
}
