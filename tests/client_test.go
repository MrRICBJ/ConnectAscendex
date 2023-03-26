package tests

import (
	"test/internal"
	"testing"
)

func TestConnection(t *testing.T) {
	// Create a new API client
	client := internal.NewAPIClient()

	// Attempt to establish a connection
	if err := client.Connection(); err != nil {
		t.Errorf("Failed to establish connection: %v", err)
	}

	// Close the connection
	client.Disconnect()
}

func TestSubscribeToChannel(t *testing.T) {
	// Create a new API client
	client := internal.NewAPIClient()

	// Attempt to establish a connection
	if err := client.Connection(); err != nil {
		t.Errorf("Failed to establish connection: %v", err)
	}

	// Subscribe to a channel
	symbol := "BTC/USDT"
	if err := client.SubscribeToChannel(symbol); err != nil {
		t.Errorf("Failed to subscribe to channel %s: %v", symbol, err)
	}

	// Close the connection
	client.Disconnect()
}
