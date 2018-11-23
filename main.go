package main

import (
	"errors"
	"log"
	"os"
	"os/user"
)

var usage = `usage: 
	prompt (set|get) servicename`

func main() {
	log.SetFlags(0)
	if len(os.Args) < 2 {
		log.Fatal(usage)
	}

	cmd := os.Args[1]
	service := os.Args[2]

	switch cmd {
	case "set":
		if err := setCreds(service); err != nil {
			log.Fatal(err)
		}
	case "get":
		u, p, err := getCreds(service)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(u, p)
	default:
		log.Println(usage)
	}
}

func getCreds(service string) (username, password string, err error) {
	u, err := user.Current()
	if err != nil {
		return "", "", err
	}
	f := &kvfile{Path: u.HomeDir + "/." + service}
	username, ok := f.Read("USERNAME")
	if !ok {
		return "", "", errors.New("USERNAME not set")
	}
	password, ok = f.Read("PASSWORD")
	if !ok {
		return "", "", errors.New("USERNAME not set")
	}
	return username, password, nil
}

func setCreds(service string) error {
	name, err := prompt("username: ")
	if err != nil {
		return err
	}
	pass, err := promptHidden("password: ")
	if err != nil {
		return err
	}

	u, err := user.Current()
	if err != nil {
		return err
	}

	f := &kvfile{Path: u.HomeDir + "/." + service}
	if err := f.Write("USERNAME", name); err != nil {
		return err
	}

	if err := f.Write("PASSWORD", pass); err != nil {
		return err
	}

	return nil
}
