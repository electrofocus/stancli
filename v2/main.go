package main

import (
	"bytes"
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

	flg := parseFlags()

	natsConn, err := nats.Connect(flg.url, nats.Timeout(20*time.Second))
	if err != nil {
		fmt.Printf("Can't connect Nats (%s).\n", err)
		return
	}

	stanConn, err := stan.Connect(flg.clusterID, strconv.Itoa(rand.Int()), stan.NatsConn(natsConn))
	if err != nil {
		fmt.Printf("Can't connect STAN (%s).\n", err)
		return
	}

	fmt.Printf("Connected to STAN %q cluster via %q URL.\n", flg.clusterID, flg.url)

	switch flg.action {
	case pubAction:
		if err := pub(flg.subject, stanConn); err != nil {
			fmt.Printf("Can't publish message (%s).\n", err)
			return
		}
	case subAction:
		if err := sub(flg.subject, stanConn); err != nil {
			fmt.Printf("Can't subscribe subject (%s).\n", err)
			return
		}
	}
}

func pub(subject string, stanConn stan.Conn) error {

	fmt.Println("Type or paste message here and then hit Return/Enter and Ctrl-D")

	payload, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("Can't read from std input (%s).\n", err)
		return err
	}

	fmt.Println("Publishing...")

	if err := stanConn.Publish(subject, payload); err != nil {
		fmt.Printf("Can't publish message (%s).\n", err)
		return err
	}

	fmt.Println("Published")

	return nil
}

func sub(subject string, stanConn stan.Conn) error {

	const group = "local"

	if _, err := stanConn.QueueSubscribe(subject, group, handle, stan.SetManualAckMode()); err != nil {
		fmt.Printf("Can't subscribe stan (%s).\n", err)
		return err
	}

	<-make(chan struct{})

	return nil
}

func handle(msg *stan.Msg) {
	defer func() {
		_ = msg.Ack()
	}()

	var buf bytes.Buffer
	if err := json.Indent(&buf, msg.Data, "", "    "); err != nil {
		fmt.Printf("%s.\n", msg.Data)
		return
	}

	fmt.Println(buf.String())
}

const (
	pubAction = "pub"
	subAction = "sub"
)

type flags struct {
	url       string
	clusterID string
	action    string
	subject   string
}

func parseFlags() flags {
	var (
		url       = flag.String("url", "nats://0.0.0.0:4222", "Nats URL")
		clusterID = flag.String("cluster-id", "test-cluster", "Nats cluster ID")
		pub       = flag.Bool("pub", false, "Publish action")
		sub       = flag.Bool("sub", false, "Subscribe action")
		subject   = flag.String("subject", "", "Nats subject")
	)

	flag.Parse()

	var action string
	if *pub {
		action = pubAction
	} else if *sub {
		action = subAction
	} else {
		fmt.Println(`Specify "pub" or "sub" action using
--pub or --sub flag respectively.
For expample:
    ./stancli --pub --subject=abc`)
		os.Exit(1)
	}

	if *subject == "" {
		fmt.Println(`Specify subject (or topic) using --subject flag.
For expample:
    ./stancli --pub --subject=abc`)
		os.Exit(1)
	}

	return flags{
		url:       *url,
		clusterID: *clusterID,
		action:    action,
		subject:   *subject,
	}
}
