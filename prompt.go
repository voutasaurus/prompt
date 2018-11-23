package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func prompt(msg string) (string, error) {
	fmt.Fprint(os.Stderr, msg)
	return bufio.NewReader(os.Stdin).ReadString('\n')
}

func promptHidden(msg string) (string, error) {
	t := terminal.NewTerminal(struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stderr}, "")
	prev, err := terminal.MakeRaw(0)
	if err != nil {
		return "", err
	}
	defer terminal.Restore(0, prev)
	return t.ReadPassword(msg)
}
