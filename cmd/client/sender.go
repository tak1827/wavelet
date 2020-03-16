package main

import (
	"os"
	"log"
	"bytes"
	"encoding/hex"
	"strconv"
	"github.com/perlin-network/wavelet/wctl"
	"github.com/davecgh/go-spew/spew"
	"sync"
	"time"
)

var (
	host string = "127.0.0.1"
	port int = 9000
	// privateKey string = "dd3a8952d78c56708d9598343dbfe5529c3659c1c073b763509b53c894dca0f5db5ed30943c90727ab5ac48b7052f5c1c18f4aa936acd3d24fdcf95ab5a892cf"

	tag int = 0

	numWorkers int = 10
)


func main() {
	argPrivateKey := os.Args[1]

  // Get host from cli args
  if len(os.Args) >= 3 {
    host = os.Args[2]
  }

  if len(os.Args) >= 4 {
  	port, _ = strconv.Atoi(os.Args[3])
    // if err != nil {
    // 	log.Fatal(err)
    // }
  }

	// Create wctl client
	rawPrivateKey, err := hex.DecodeString(argPrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	config := wctl.Config{
		APIHost:  host,
		APIPort:  uint16(port),
		UseHTTPS: false,
	}
	copy(config.PrivateKey[:], rawPrivateKey)
	spew.Dump(config)
	client, err := wctl.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// Loop sending tx
	for {
		var wg sync.WaitGroup
		wg.Add(numWorkers)
		chErr := make(chan error, numWorkers)

		for i := 0; i < numWorkers; i++ {
			go sendTx(client, &wg, chErr)
		}

		wg.Wait()

		// var counter int
		// for err = range chErr {
		// 	if err != nil {
		// 		log.Fatal(err.Error())
		// 	}

		// 	// Close chan if the last
		// 	counter++
		// 	if counter == numWorkers {
		// 		close(chErr)
		// 	}
		// }

		for i := 0; i < numWorkers; i++ {
			if e := <-chErr; e != nil {
				log.Print(e.Error())
			}
		}
	}
}

func sendTx(
	client *wctl.Client,
	wg *sync.WaitGroup,
	chErr chan error) {
time.Sleep(time.Millisecond * 100)
	defer wg.Done()
	// Send Non type tx
	payload := bytes.NewBuffer(nil)
	_, err := client.SendTransaction(byte(tag), payload.Bytes())
	chErr <- err
}
