package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Trade struct {
	Symbol    string `json:"symbol"`
	TradeID   string `json:"tradeId"`
	Price     string `json:"price"`
	Size      string `json:"size"`
	Side      string `json:"side"`
	Timestamp int64  `json:"timestamp"`
}

type Message struct {
	Topic  string  `json:"topic"`
	Symbol string  `json:"symbol"`
	Data   []Trade `json:"data"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Permitir conexiones de cualquier origen (para desarrollo)
	},
}

var (
	rnd       = rand.New(rand.NewSource(time.Now().UnixNano()))
	symbols   = []string{"BTC_USDT", "ETH_USDT", "XRP_USDT", "ADA_USDT", "DOGE_USDT"}
	lastPrice = 50000.0
)

func generateMockTrade(tipestamp int64) Trade {
	priceChange := rnd.Float64()*20 - 10
	lastPrice += priceChange
	return Trade{
		Symbol:    symbols[rnd.Intn(len(symbols))],
		TradeID:   strconv.Itoa(rnd.Intn(10000)),
		Price:     strconv.FormatFloat(lastPrice, 'f', 2, 64),
		Side:      []string{"BUY", "SELL"}[rnd.Intn(2)],
		Timestamp: tipestamp,
	}
}

func tradeStreamHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	currentTime := time.Now().UnixMilli()

	for {
		trades := []Trade{}

		for i := 0; i < rnd.Intn(5)+1; i++ {
			trades = append(trades, generateMockTrade(currentTime))
		}

		message := Message{
			Topic:  "TRADE",
			Symbol: "MULTI",
			Data:   trades,
		}

		msgBytes, _ := json.Marshal(message)
		if err := conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
			break
		}

		currentTime += 6000
		time.Sleep(1 * time.Second)
	}
}

func main() {
	router := gin.Default()

	router.GET("/ws/trades", tradeStreamHandler)

	router.Run(":8080")
}
