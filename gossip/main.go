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

func buildTimeline(swarmClient *client.Client, account account.Account) timeline.Timeline {

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
	swarmClient := client.NewClient(c.GlobalString("swarmgateway"))

	args := c.Args()
	tm := buildTimeline(swarmClient, account)
	return tm.Post(args[0])
}

func view(c *cli.Context) error {
	var addr common.Address
	account := account.New(c.GlobalString("passphrase"))
	swarmClient := client.NewClient(c.GlobalString("swarmgateway"))

	args := c.Args()
	if c.NArg() < 1 {
		addr = account.Addr()
	} else {
		addr = common.HexToAddress(args[0])
	}

	tm := buildTimeline(swarmClient, account)
	comments := tm.Dump(addr)
	for c := range comments {
		fmt.Printf("[%s] *** %s\n", time.Unix(c.Timestamp, 0), c.Text)
	}
	return nil
}

func setnick(c *cli.Context) error {
	if c.NArg() < 1 {
		return errors.New("Expected nickname")
	}
	account := account.New(c.GlobalString("passphrase"))
	swarmClient := client.NewClient(c.GlobalString("swarmgateway"))

	args := c.Args()
	kvservice := kv.New(&kv.Config{
		SwarmClient: swarmClient,
		Account:     account,
	})
	return kvservice.Put("nickname", []byte(args[0]))
}

func getnick(c *cli.Context) error {
	account := account.New(c.GlobalString("passphrase"))
	swarmClient := client.NewClient(c.GlobalString("swarmgateway"))

	args := c.Args()
	var addr common.Address
	if c.NArg() < 1 {
		addr = account.Addr()
	} else {
		addr = common.HexToAddress(args[0])
	}

	kvservice := kv.New(&kv.Config{
		SwarmClient: swarmClient,
		Account:     account,
	})
	nicknameBytes, err := kvservice.Get(addr, "nickname")
	if err != nil {
		return err
	}
	fmt.Printf("Nickname for %s is %s\n", addr.Hex(), string(nicknameBytes))
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
		{
			Name:        "setnick",
			Description: "Sets your nickname",
			UsageText:   "setnick <nickname>",
			Action:      setnick,
		},
		{
			Name:        "getnick",
			Description: "Prints a user's nickname. If a user address is not given, it shows your nickname",
			UsageText:   "getnick [user address]",
			Action:      getnick,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
