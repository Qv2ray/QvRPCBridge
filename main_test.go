package main

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestDial(t *testing.T) {
	println(Dial("127.0.0.1:15490", 10000))
}

func TestGetStats(t *testing.T) {
	if errStr := Dial("127.0.0.1:15490", 10000); errStr != "" {
		log.Fatalf(errStr)
	}

	stats := GetStats("inbound>>>socks_IN>>>traffic>>>downlink", 1000)
	if stats == -1 {
		log.Fatalf("api call failed!")
	} else if stats == -666 {
		log.Fatalf("api call timed out!")
	} else {
		log.Printf("bandwidth usage: %d", stats)
	}
}

func TestGetStats_Evil(t *testing.T) {
	// notice we did not dial first
	assert.Equal(t, -999, GetStats("motherfucker", 1000))
}
