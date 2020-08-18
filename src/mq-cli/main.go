package main

import (
	"bufio"
	"flag"
	"fmt"
	"mq-cli/queue"
	"os"
)

func main() {
	// define send or receive mode
	mode := flag.String("mode", "publish", "Mode: [publish, consume]")
	connectionString := flag.String("connection-string", "amqp://guest:guest@localhost:5672/", "The connection string used to interact with the queue.")
	channel := flag.String("channel", "default", "Channel name")
	flag.Parse()

	// open the queue
	var queue = new(queue.Queue)
	queue.OpenQueue(*connectionString, *channel)
	defer queue.CloseQueue()

	// pipe stdin/out to channel

	switch *mode {
	case "publish":
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			queue.PublishStringMessage(line)
		}
	case "consume":
		queue.ReadStringMessage()
	default:
		fmt.Println("[-] Error: unknown mode")
	}
}
