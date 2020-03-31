package main

import (
  "bytes"
  "fmt"
  "os"
  "encoding/hex"
  "encoding/json"
  "github.com/perlin-network/wavelet/wctl"
  "log"
  // "github.com/davecgh/go-spew/spew"
)

var (
  host string = "127.0.0.1"
  port uint16 = 9000
  privateKey = "dd3a8952d78c56708d9598343dbfe5529c3659c1c073b763509b53c894dca0f5db5ed30943c90727ab5ac48b7052f5c1c18f4aa936acd3d24fdcf95ab5a892cf"
)


func main() {
  txID := os.Args[1]
  // Get host from cli args
  if len(os.Args) >= 3 {
    host = os.Args[2]
  }

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

  res, err := client.GetTransaction(txID)
  if err != nil {
    log.Fatal(err)
  }

  buf, err := json.Marshal(res)
  if err != nil {
    log.Fatal(err)
  } else {
    output(buf)
  }
}

// Write bytes to stdout; do JSON indent if possible.
func output(buf []byte) {
  var out bytes.Buffer

  if err := json.Indent(&out, buf, "", "\t"); err != nil {
    out.Write(buf)
  }

  fmt.Println(string(out.Bytes()))
}
