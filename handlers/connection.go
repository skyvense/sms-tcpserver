package handlers

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/url"
	"path"
	"sms-tcpserver/models"
	"sync"
	"time"
)

var (
	httpClient = &http.Client{
		Timeout: 120 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 允许不安全的HTTPS证书
			},
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	messageQueue = make(chan models.Message, 1000)
	httpWg       sync.WaitGroup
	httpBaseURL  string
)

// InitHTTPSender initializes the HTTP sender with the base URL
func InitHTTPSender(baseURL string) {
	httpBaseURL = baseURL
	// Start HTTP sender goroutine
	httpWg.Add(1)
	go httpSender()
}

// httpSender processes messages from the queue and sends HTTP requests
func httpSender() {
	defer httpWg.Done()

	for msg := range messageQueue {
		// Construct URL
		u, err := url.Parse(httpBaseURL)
		if err != nil {
			log.Printf("Error parsing base URL: %v", err)
			continue
		}
		u.Path = path.Join(u.Path, msg.Num, msg.Txt)

		// Send HTTP GET request
		resp, err := httpClient.Get(u.String())
		if err != nil {
			log.Printf("Error sending HTTP request: %v", err)
			continue
		}

		// Log result
		log.Printf("HTTP Response [url=%s status=%d]", u.String(), resp.StatusCode)
		resp.Body.Close()
	}
}

// HandleConnection processes incoming TCP connections
func HandleConnection(conn net.Conn) {
	defer conn.Close()

	// Set read deadline
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	// Create a buffer to store the incoming data
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Printf("Error reading from connection: %v", err)
		return
	}

	// Parse the JSON data
	var msg models.Message
	if err := json.Unmarshal(buffer[:n], &msg); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return
	}

	// Log the message
	log.Printf("Message [txt=%s num=%s cmd=%s tz=%d time=%02d:%02d:%02d date=%d-%02d-%02d seq=%d ref=%d max=%d]",
		msg.Txt, msg.Num, msg.Cmd,
		msg.Metas.Tz,
		msg.Metas.Hour, msg.Metas.Min, msg.Metas.Sec,
		2000+msg.Metas.Year, msg.Metas.Mon, msg.Metas.Day,
		msg.Metas.SeqNum, msg.Metas.RefNum, msg.Metas.MaxNum)

	// Put message in queue for HTTP request
	messageQueue <- msg
	log.Printf("Queued message for HTTP request [txt=%s num=%s]", msg.Txt, msg.Num)
}

// StopHTTPSender stops the HTTP sender goroutine
func StopHTTPSender() {
	close(messageQueue)
	httpWg.Wait()
}
