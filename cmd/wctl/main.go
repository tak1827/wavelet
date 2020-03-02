// Copyright (c) 2019 Perlin
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"time"
	// "github.com/davecgh/go-spew/spew"

	"github.com/perlin-network/wavelet"
	"github.com/perlin-network/wavelet/sys"
	"github.com/perlin-network/wavelet/wctl"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "wctl"
	app.Author = "Perlin"
	app.Email = "support@perlin.net"
	app.Version = sys.Version
	app.Usage = "a cli client to interact with the wavelet node"

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("Version:    %s\n", sys.Version)
		fmt.Printf("Go Version: %s\n", sys.GoVersion)
		fmt.Printf("Git Commit: %s\n", sys.GitCommit)
		fmt.Printf("OS/Arch:    %s\n", sys.OSArch)
		fmt.Printf("Built:      %s\n", c.App.Compiled.Format(time.ANSIC))
	}

	commonFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "api.host",
			Value: "localhost",
			Usage: "Host of the local HTTP API.",
		},
		cli.IntFlag{
			Name:  "api.port",
			Usage: "Port a local HTTP API.",
		},
		cli.StringFlag{
			Name:  "key",
			Usage: "Private key hex-encoded",
		},
		cli.StringFlag{
			Name:  "wallet",
			Usage: "path to file containing hex-encoded private key",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "poll_broadcaster",
			Usage: "continuously receive broadcaster updates",
			Flags: commonFlags,
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}

				evChan, err := client.PollLoggerSink(nil, wctl.RouteWSBroadcaster)
				if err != nil {
					return err
				}

				for ev := range evChan {
					output(ev)
				}
				return nil
			},
		},
		{
			Name:  "poll_consensus",
			Usage: "continuously receive consensus updates",
			Flags: commonFlags,
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}

				evChan, err := client.PollLoggerSink(nil, wctl.RouteWSConsensus)
				if err != nil {
					return err
				}

				for ev := range evChan {
					output(ev)
				}
				return nil
			},
		},
		{
			Name:  "poll_stake",
			Usage: "continuously receive stake updates",
			Flags: commonFlags,
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}

				evChan, err := client.PollLoggerSink(nil, wctl.RouteWSStake)
				if err != nil {
					return err
				}

				for ev := range evChan {
					output(ev)
				}
				return nil
			},
		},
		{
			Name:  "poll_accounts",
			Usage: "continuously receive account updates",
			Flags: append(commonFlags,
				[]cli.Flag{
					cli.StringFlag{
						Name:  "account_id",
						Usage: "account id to list (default: all)",
					},
				}...,
			),
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}

				var accountID *string
				if len(c.String("account_id")) > 0 {
					tmp := c.String("account_id")
					accountID = &tmp
				}

				evChan, err := client.PollAccounts(nil, accountID)
				if err != nil {
					return err
				}

				for ev := range evChan {
					output(ev)
				}
				return nil
			},
		},
		{
			Name:  "poll_contracts",
			Usage: "continuously receive contract updates",
			Flags: append(commonFlags,
				[]cli.Flag{
					cli.StringFlag{
						Name:  "contract_id",
						Usage: "contract id to list (default: all)",
					},
				}...,
			),
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}

				// get these optional variables
				var contractID *string
				if len(c.String("contract_id")) > 0 {
					tmp := c.String("contract_id")
					contractID = &tmp
				}

				evChan, err := client.PollContracts(nil, contractID)
				if err != nil {
					return err
				}

				for ev := range evChan {
					output(ev)
				}
				return nil
			},
		},
		{
			Name:  "poll_transactions",
			Usage: "continuously receive transaction updates",
			Flags: append(commonFlags,
				[]cli.Flag{
					cli.StringFlag{
						Name:  "tx_id",
						Usage: "transactions to list (default: all)",
					},
					cli.StringFlag{
						Name:  "sender_id",
						Usage: "sender id of transactions to list (default: all)",
					},
					cli.StringFlag{
						Name:  "creator_id",
						Usage: "creator id of transactions to list (default: all)",
					},
					cli.StringFlag{
						Name:  "tag",
						Usage: "tag of transactions to list (default: all)",
					},
				}...,
			),
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}

				// get these optional variables
				var (
					txID      *string
					senderID  *string
					creatorID *string
					tag       *byte
				)

				if len(c.String("tx_id")) > 0 {
					tmp := c.String("tx_id")
					txID = &tmp
				}
				if len(c.String("sender_id")) > 0 {
					tmp := c.String("sender_id")
					senderID = &tmp
				}
				if len(c.String("creator_id")) > 0 {
					tmp := c.String("creator_id")
					creatorID = &tmp
				}
				if len(c.String("tag")) > 0 {
					tmp := c.String("tag")
					t := byte(sys.TagLabels[tmp])
					tag = &t
				}

				evChan, err := client.PollTransactions(nil, txID, senderID, creatorID, tag)
				if err != nil {
					return err
				}

				for ev := range evChan {
					output(ev)
				}
				return nil
			},
		},
		{
			Name:  "ledger_status",
			Usage: "get the status of the ledger",
			Flags: append(commonFlags,
				[]cli.Flag{
					cli.StringFlag{
						Name:  "sender_id",
						Usage: "sender id of transactions to list (default: all)",
					},
					cli.StringFlag{
						Name:  "creator_id",
						Usage: "creator id of transactions to list (default: all)",
					},
					cli.IntFlag{
						Name:  "offset",
						Usage: "an offset of the number of transactions to list",
					},
					cli.IntFlag{
						Name:  "limit",
						Usage: "limit to max number of transactions to list",
					},
				}...,
			),
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}

				// get these optional variables
				var senderID *string
				var creatorID *string
				var offset *uint64
				var limit *uint64
				if len(c.String("sender_id")) > 0 {
					tmp := c.String("sender_id")
					senderID = &tmp
				}
				if len(c.String("creator_id")) > 0 {
					tmp := c.String("creator_id")
					creatorID = &tmp
				}
				if c.Uint("offset") > 0 {
					tmp := uint64(c.Uint("offset"))
					offset = &tmp
				}
				if c.Uint("limit") > 0 {
					tmp := uint64(c.Uint("limit"))
					limit = &tmp
				}

				res, err := client.GetLedgerStatus(senderID, creatorID, offset, limit)
				if err != nil {
					return err
				}

				buf, err := json.Marshal(res)
				if err != nil {
					fmt.Println(err)
				} else {
					output(buf)
				}

				return nil
			},
		},
		{
			Name:      "get_account",
			Usage:     "get an account",
			ArgsUsage: "<account ID>",
			Flags:     commonFlags,
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}
				acctID := c.Args().Get(0)

				res, err := client.GetAccount(acctID)
				if err != nil {
					return err
				}

				buf, err := json.Marshal(res)
				if err != nil {
					fmt.Println(err)
				} else {
					output(buf)
				}

				return nil
			},
		},
		{
			Name:      "get_contract_code",
			Usage:     "get the payload of a contract",
			ArgsUsage: "<contract ID>",
			Flags:     commonFlags,
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}
				contractID := c.Args().Get(0)

				res, err := client.GetContractCode(contractID)
				if err != nil {
					return err
				}

				fmt.Println(res)

				return nil
			},
		},
		{
			Name:      "get_contract_pages",
			Usage:     "get the page of a contract",
			ArgsUsage: "<contract ID>",
			Flags: append(commonFlags,
				[]cli.Flag{
					cli.StringFlag{
						Name:  "page_idx",
						Usage: "page offset of the contract",
					},
				}...,
			),
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}
				contractID := c.Args().Get(0)

				// get these optional variables
				var pageIdx *uint64
				if c.Uint("page_idx") > 0 {
					tmp := uint64(c.Uint("page_idx"))
					pageIdx = &tmp
				}

				res, err := client.GetContractPages(contractID, pageIdx)
				if err != nil {
					return err
				}

				fmt.Println(res)

				return nil
			},
		},
		{
			Name:      "send_transaction",
			Usage:     "send a transaction",
			ArgsUsage: "<tag> <json payload>",
			Flags: append(commonFlags,
				[]cli.Flag{
					cli.StringFlag{
						Name:  "payload",
						Usage: "the path to the payload file",
					},
				}...,
			),
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}

				tag, err := strconv.Atoi(c.Args().Get(0))
				if err != nil {
					return err
				}

				payload := bytes.NewBuffer(nil)

				if c.String("payload") != "" {
					payloadFile, err := ioutil.ReadFile(c.String("payload"))
					if err != nil {
						return err
					}

					parsedPayload, err := wavelet.ParseJSON(payloadFile, c.Args().Get(0)) // Parse payload file contents
					if err != nil {                                                       // Check for errors
						return err // Return found error
					}

					payload.Write(parsedPayload) // Write payload to buffer
				}

				res, err := client.SendTransaction(byte(tag), payload.Bytes())
				if err != nil {
					return err
				}

				buf, err := json.Marshal(res)
				if err != nil {
					fmt.Println(err)
				} else {
					output(buf)
				}

				return nil
			},
		},
		{
			Name:      "get_transaction",
			Usage:     "get a transaction",
			ArgsUsage: "<transaction ID>",
			Flags:     commonFlags,
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}
				txID := c.Args().Get(0)

				res, err := client.GetTransaction(txID)
				if err != nil {
					return err
				}

				buf, err := json.Marshal(res)
				if err != nil {
					fmt.Println(err)
				} else {
					output(buf)
				}

				return nil
			},
		},
		{
			Name:  "list_transactions",
			Usage: "list recent transactions",
			Flags: append(commonFlags,
				[]cli.Flag{
					cli.StringFlag{
						Name:  "sender_id",
						Usage: "sender id of transactions to list (default: all)",
					},
					cli.StringFlag{
						Name:  "creator_id",
						Usage: "creator id of transactions to list (default: all)",
					},
					cli.IntFlag{
						Name:  "offset",
						Usage: "an offset of the number of transactions to list",
					},
					cli.IntFlag{
						Name:  "limit",
						Usage: "limit to max number of transactions to list",
					},
				}...,
			),
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}

				// get these optional variables
				var senderID *string
				var creatorID *string
				var offset *uint64
				var limit *uint64
				if len(c.String("sender_id")) > 0 {
					tmp := c.String("sender_id")
					senderID = &tmp
				}
				if len(c.String("creator_id")) > 0 {
					tmp := c.String("creator_id")
					creatorID = &tmp
				}
				if c.Uint("offset") > 0 {
					tmp := uint64(c.Uint("offset"))
					offset = &tmp
				}
				if c.Uint("limit") > 0 {
					tmp := uint64(c.Uint("limit"))
					limit = &tmp
				}

				res, err := client.ListTransactions(senderID, creatorID, offset, limit)
				if err != nil {
					return err
				}

				buf, err := json.Marshal(res)
				if err != nil {
					fmt.Println(err)
				} else {
					output(buf)
				}

				return nil
			},
		},
		{
			Name:  "poll_metrics",
			Usage: "continuously receive metrics",
			Flags: commonFlags,
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}

				evChan, err := client.PollLoggerSink(nil, wctl.RouteWSMetrics)
				if err != nil {
					return err
				}

				for ev := range evChan {
					output(ev)
				}
				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Failed to parse configuration/command-line arguments: %+v\n", err)
	}
}

func setup(c *cli.Context) (*wctl.Client, error) {
	host := c.String("api.host")
	port := c.Uint("api.port")
	privateKeyFile := c.String("wallet")
	privateKey := c.String("key")

	if port == 0 {
		return nil, errors.New("port is missing")
	}

	var privateKeyBytes []byte
	var err error

	if len(privateKey) != 0 {
		privateKeyBytes = []byte(privateKey)
	} else if len(privateKeyFile) != 0 {
		privateKeyBytes, err = ioutil.ReadFile(privateKeyFile)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read private key %s", privateKeyFile)
		}
	}

	if len(privateKeyBytes) == 0 {
		return nil, errors.New("private key is missing")
	}

	rawPrivateKey, err := hex.DecodeString(string(privateKeyBytes))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to hex decode private key %s", privateKeyFile)
	}

	config := wctl.Config{
		APIHost:  host,
		APIPort:  uint16(port),
		UseHTTPS: false,
	}
	copy(config.PrivateKey[:], rawPrivateKey)

	client, err := wctl.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Write bytes to stdout; do JSON indent if possible.
func output(buf []byte) {
	var out bytes.Buffer

	if err := json.Indent(&out, buf, "", "\t"); err != nil {
		out.Write(buf)
	}

	fmt.Println(string(out.Bytes()))
}
