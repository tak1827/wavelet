package main

import (
  "log"
  "os"
  "strconv"
  "encoding/hex"
  "github.com/perlin-network/wavelet/wctl"
  "github.com/perlin-network/wavelet/conf"

  "crypto/sha512"
	"encoding/base64"
  // "github.com/davecgh/go-spew/spew"
)

func main() {
	argHost := os.Args[1]
	argPort := os.Args[2]
  argPeer := os.Args[3]
  // spew.Dump(argPeer)

  // Parse port
  port, err := strconv.Atoi(argPort)
  if err != nil {
  	log.Fatal("Failed parse port")
  }

  // Config client secret
	h, err := hex.DecodeString("")
	if err != nil {
		log.Fatal(err.Error())
	}
	sha := sha512.Sum512_224(h)
	secret := base64.StdEncoding.EncodeToString(sha[:])
  conf.Update(
		conf.WithSecret(secret),
	)

	// Create wctl client
  config := wctl.Config{
    APIHost:  argHost,
    APIPort:  uint16(port),
    UseHTTPS: false,
  }
  config.APISecret = conf.GetSecret()
  client, err := wctl.NewClient(config)
  if err != nil {
    log.Fatal(err)
  }

  // Connect to peer
  m, err := client.Connect(argPeer)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Print(m.Message)
}
