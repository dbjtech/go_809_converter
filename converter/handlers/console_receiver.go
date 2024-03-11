package handlers

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

var consoleChanel = "809converter_console_chanel"
var specialCommands = []string{"display_all"}

type ConsoleReceiver struct {
	rds    *redis.Client
	pubSub *redis.PubSub
}

func (c *ConsoleReceiver) StartListen() {
	if c.rds == nil {
		log.Println("Redis has gone, can not listen console's" +
			" command")
		return
	}
	c.pubSub = c.rds.PSubscribe(context.Background(), consoleChanel)

}
func (c *ConsoleReceiver) ListenChange() {
	log.Println("start console listen")
	ctx := context.Background()
	for {
		resp, err := c.pubSub.ReceiveMessage(ctx)
		if err != nil {

		}
		fmt.Println(resp.Channel, resp.Payload)
	}
}
func (c *ConsoleReceiver) Execute() {

}
func (c *ConsoleReceiver) DisplayAll() {

}
