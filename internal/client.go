package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"test"
	"time"

	"github.com/gorilla/websocket"
)

type APIClient struct {
	Conn *websocket.Conn
}

func NewAPIClient() *APIClient {
	return &APIClient{}
}

func (a *APIClient) Connection() error {
	var err error

	// Define the websocket connection URL
	url := "wss://ascendex.com/1/api/pro/v1/stream"

	// Define the HTTP headers to be sent in the initial request
	headers := http.Header{}
	headers.Add("Origin", "https://ascendex.com")

	// Create the websocket dialer and set the headers
	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = 10 * time.Second

	// Dial the websocket endpoint
	a.Conn, _, err = dialer.Dial(url, headers)
	if err != nil {
		return err
	}

	// Set a read deadline on the websocket connection
	a.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	// Read the initial welcome message from the server
	_, welcome, err := a.Conn.ReadMessage()
	if err != nil {
		return err
	}
	fmt.Println("Received welcome message:", string(welcome))

	return nil
}

func (c *APIClient) Disconnect() {
	if c.Conn != nil {
		c.Conn.Close()
	}
}

func (с *APIClient) SubscribeToChannel(symbol string) error {
	err := с.Conn.WriteMessage(websocket.TextMessage, []byte(`{ "op": "sub", "id": "abcd1234", "ch": "bbo:BTC/USDT" }`))
	if err != nil {
		return err
	}

	return nil
}

func (c *APIClient) ReadMessagesFromChannel(ch chan<- project.BestOrderBook) {
	if c.Conn == nil {
		return
	}

	defer close(ch)

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Printf("error while reading message: %v", err)
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			fmt.Printf("error while unmarshalling message: %v", err)
			return
		}

		if data["m"] == "bbo" {
			var data struct {
				Symbol string `json:"symbol"`
				Data   struct {
					Ts  int64    `json:"ts"`
					Bid []string `json:"bid"`
					Ask []string `json:"ask"`
				} `json:"data"`
			}
			if err := json.Unmarshal(message, &data); err != nil {
				fmt.Printf("error while unmarshalling book: %v", err)
				return
			}
			bidPrice, err := strconv.ParseFloat(data.Data.Bid[0], 64)
			if err != nil {
				log.Println("Error parsing bid price:", err)
				continue
			}

			bidAmount, err := strconv.ParseFloat(data.Data.Bid[1], 64)
			if err != nil {
				log.Println("Error parsing bid amount:", err)
				continue
			}

			askPrice, err := strconv.ParseFloat(data.Data.Ask[0], 64)
			if err != nil {
				log.Println("Error parsing ask price:", err)
				continue
			}

			askAmount, err := strconv.ParseFloat(data.Data.Ask[1], 64)
			if err != nil {
				log.Println("Error parsing ask amount:", err)
				continue
			}

			bbo := project.BestOrderBook{
				Ask: project.Order{
					Price:  askPrice,
					Amount: askAmount,
				},
				Bid: project.Order{
					Price:  bidPrice,
					Amount: bidAmount,
				},
			}

			ch <- bbo
		}
	}
}

func (c *APIClient) WriteMessagesToChannel() {
	// No implementation needed for this method.
}
