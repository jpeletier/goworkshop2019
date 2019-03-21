package main

import (
	"errors"
	"fmt"
	"goworkshop/gossip/cas/swarmcas"
	"goworkshop/gossip/kv"
	"goworkshop/gossip/kv/account"
	"goworkshop/gossip/objstore"
	"goworkshop/gossip/timeline"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/swarm/api/client"

	"github.com/urfave/cli"
)

func buildTimeline(c *cli.Context, account account.Account) timeline.Timeline {

	swarmClient := client.NewClient(c.GlobalString("swarmgateway"))

	kvservice := kv.New(&kv.Config{
		SwarmClient: swarmClient,
		Account:     account,
	})

	objstore := objstore.New(&objstore.Config{
		BackendStorage: swarmcas.New(swarmcas.Config{
			SwarmClient: swarmClient,
		}),
	})

	tm := timeline.New(&timeline.Config{
		KVService: kvservice,
		ObjStore:  objstore,
		Account:   account,
	})
	return tm
}

func getaddr(c *cli.Context) error {
	account := account.New(c.GlobalString("passphrase"))
	fmt.Printf("Your user address is %s\n", account.Addr().Hex())
	return nil
}

func post(c *cli.Context) error {
	if c.NArg() < 1 {
		return errors.New("Expected message")
	}
	account := account.New(c.GlobalString("passphrase"))

	args := c.Args()
	tm := buildTimeline(c, account)
	return tm.Post(args[0])
}

func view(c *cli.Context) error {
	var addr common.Address
	account := account.New(c.GlobalString("passphrase"))

	args := c.Args()
	if c.NArg() < 1 {
		addr = account.Addr()
	} else {
		addr = common.HexToAddress(args[0])
	}

	tm := buildTimeline(c, account)
	comments := tm.Dump(addr)
	for c := range comments {
		fmt.Printf("[%s] *** %s\n", time.Unix(c.Timestamp, 0), c.Text)
	}
	return nil
}

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "passphrase, p",
			Usage: "Passphrase to use for authentication",
		},
		cli.StringFlag{
			Name:  "swarmgateway, s",
			Usage: "URL of the Swarm Gateway to use",
			Value: "https://swarm.epiclabs.io",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:        "getaddr",
			Description: "Retrieves your user address",
			Action:      getaddr,
		},
		{
			Name:        "post",
			Description: "Allows you to post a message",
			UsageText:   "gossip post \"message you want to send\"",
			Action:      post,
		},
		{
			Name:        "view",
			Description: "Shows a user's timeline of posts. If a user address is not given, it shows your posts",
			UsageText:   "gossip view [user address]",
			Action:      view,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}