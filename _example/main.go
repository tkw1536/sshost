package main

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/tkw1536/sshost"
)

func main() {

	env, err := sshost.NewDefaultContext()
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) <= 2 {
		log.Fatal(errors.New("need at least two args"))
	}

	client, closer, err := env.NewClient(nil, os.Args[1])
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
