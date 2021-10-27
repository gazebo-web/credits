package main

import (
	"gitlab.com/ignitionrobotics/billing/credits/internal/server"
	"log"
	"os"
)

// main prepares the config and runs the credits HTTP server.
func main() {
	logger := log.New(os.Stdout, "[Credits API] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)

	// Prepare the config
	cfg, err := server.Setup(logger)
	if err != nil {
		logger.Fatalln("Failed to initialize server configuration:", err)
	}

	// Run the HTTP server with the given config
	if err = server.Run(cfg, logger); err != nil {
		logger.Fatalln("Failed to run HTTP server:", err)
	}

	logger.Println("Shutting HTTP server down...")
}
