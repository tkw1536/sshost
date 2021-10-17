package main

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/tkw1536/sshost"
)

func main() {

	env, err := sshost.NewDefaultEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) <= 2 {
		log.Fatal(errors.New("need at least two args"))
	}

	client, closer, err := env.NewClient(nil, os.Args[1], context.Background())
	defer closer.Close()

	if err != nil {
		log.Fatal(err)
	}

	session, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	if err := session.Run(strings.Join(os.Args[2:], " ")); err != nil {
		log.Fatal(err)
	}

}
