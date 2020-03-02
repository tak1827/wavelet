package main

import (
	"log"
	"bytes"
	"encoding/hex"
	"github.com/perlin-network/wavelet/wctl"
	// "github.com/davecgh/go-spew/spew"
)

var (
	host string = "127.0.0.1"
	port uint16 = 9000
	privateKey string = "dd3a8952d78c56708d9598343dbfe5529c3659c1c073b763509b53c894dca0f5db5ed30943c90727ab5ac48b7052f5c1c18f4aa936acd3d24fdcf95ab5a892cf"

	tag int = 0

	numWorkers int = 8
)


func main() {

	// Create wctl client
	rawPrivateKey, err := hex.DecodeString(privateKey)
	if err != nil {
		log.Fatal(err)
	}
	config := wctl.Config{
		APIHost:  host,
		APIPort:  port,
		UseHTTPS: false,
	}
	copy(config.PrivateKey[:], rawPrivateKey)
	// spew.Dump(config)
	client, err := wctl.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// Loop sending tx
	for {
		chErr := make(chan error, numWorkers)

		for i := 0; i < numWorkers; i++ {
			go sendTx(client, chErr)
		}


		for i := 0; i < numWorkers; i++ {
			if e := <-chErr; e != nil {
				log.Fatal(e.Error())
			}
		}

	}
}

func sendTx(
	client *wctl.Client,
	chErr chan error) {
	// Send Non type tx
	payload := bytes.NewBuffer(nil)
	_, err := client.SendTransaction(byte(tag), payload.Bytes())
	chErr <- err
}
