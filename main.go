package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Action is not specified. Specify action: pub or sub")
		return
	}

	if len(os.Args) < 3 {
		fmt.Println("Subject is not specified")
		return
	}

	const (
		pubAction = "pub"
		subAction = "sub"
	)

	var (
		action  = os.Args[1]
		subject = os.Args[2]
	)

	switch action {
	case pubAction, subAction:
	default:
		fmt.Printf("Incorrect action %q\n", action)
		return
	}

	path := flag.String("config", "./config.json", "config file path")
	flag.Parse()

	file, err := os.Open(*path)
	if err != nil {
		fmt.Printf("can't open file (%s)\n", err)
		return
	}

	var cfg config

	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		fmt.Printf("can't decode file content (%s)\n", err)
		return
	}

	fmt.Println("Connecting Nats...")

	natsConn, err := nats.Connect(cfg.URL, nats.Timeout(20*time.Second))
	if err != nil {
		fmt.Printf("can't connect nats (%s)\n", err)
		return
	}

	stanConn, err := stan.Connect(cfg.ClusterID, strconv.Itoa(rand.Int()), stan.NatsConn(natsConn))
	if err != nil {
		fmt.Printf("can't connect stan (%s)\n", err)
		return
	}

	fmt.Println("Connected")

	switch action {
	case pubAction:
		if err := pub(subject, stanConn); err != nil {
			fmt.Printf("can't publish message (%s)\n", err)
			return
		}
	case subAction:
		if err := sub(subject, stanConn); err != nil {
			fmt.Printf("can't subscribe subject (%s)\n", err)
			return
		}
	}
}

func pub(subject string, stanConn stan.Conn) error {

	fmt.Println("Type or paste message here and then hit Return/Enter and Ctrl-D")

	payload, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("can't read from std input (%s)\n", err)
		return err
	}

	fmt.Println("Publishing...")

	if err := stanConn.Publish(subject, payload); err != nil {
		fmt.Printf("can't publish message (%s)\n", err)
		return err
	}

	fmt.Println("Published")

	return nil
}

func sub(subject string, stanConn stan.Conn) error {

	const group = "local"

	if _, err := stanConn.QueueSubscribe(subject, group, handle, stan.SetManualAckMode()); err != nil {
		fmt.Printf("can't subscribe stan (%s)\n", err)
		return err
	}

	<-make(chan struct{})

	return nil
}

func handle(msg *stan.Msg) {
	fmt.Printf("%s\n", msg.Data)
	_ = msg.Ack()
}

type config struct {
	URL       string `json:"url"`
	ClusterID string `json:"cluster_id"`
}
