// Package auth implements the Gitlab authentication with tokens
package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Tokener interface {
	Get(string) (string, error)
	Set(string, string) error
	Remove(string) error
}

type Token struct {
	store string
	token map[string]string
}

func New(filename string) (Tokener, error) {
	t := &Token{
		store: filename,
		token: make(map[string]string),
	}

	switch err := t.read(); err {
	case nil, io.EOF:
		return t, nil
	default:
		return nil, err
	}
}

func (t *Token) read() error {
	f, err := os.OpenFile(t.store, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(&t.token)
}

func (t *Token) write() error {
	f, err := os.OpenFile(t.store, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	e := json.NewEncoder(f)
	e.SetIndent("", "  ")

	return e.Encode(t.token)
}

// Get token for hostname
func (t *Token) Get(hostname string) (string, error) {
	v, ok := t.token[hostname]
	if !ok {
		return "", fmt.Errorf("token for %s not available", hostname)
	}

	return v, nil
}

// Set token for hostname
func (t *Token) Set(hostname, token string) error {
	t.token[hostname] = token

	return t.write()
}

// Remove token for hostname
func (t *Token) Remove(hostname string) error {
	delete(t.token, hostname)

	return t.write()
}
