package main

import (
	"flag"
	"log"
	"sms-tcpserver/handlers"
	"sms-tcpserver/server"
)

func main() {
	// Parse command line arguments
	port := flag.String("port", "8080", "TCP server port")
	httpBaseURL := flag.String("http-url", "", "Base URL for HTTP requests (e.g., https://test.com/path)")
	flag.Parse()

	if *httpBaseURL == "" {
		log.Fatal("HTTP base URL is required. Use -http-url flag to specify it.")
	}

	// Initialize HTTP sender
	handlers.InitHTTPSender(*httpBaseURL)

	// Create and start the server
	srv := server.NewServer(":" + *port)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Run the server
	srv.Run()
}
