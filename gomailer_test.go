package main

import (
	"io"
	"log"
	"os"
	"testing"
)

// Test the functions

func TestReadUsers(t *testing.T) {
	// Test the readUsers function
	ps := readUsers("recipients.txt", "\t")

	if len(ps) == 0 {
		t.Fatalf("No users found")
	}

}

func TestConfirm(t *testing.T) {
	// Test the confirm function
	log.SetOutput(io.Discard)
	content := []byte("y")

	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }() // Restore original Stdin

	os.Stdin = tmpfile
	confirmation := confirm()
	if !confirmation {
		t.Errorf("confirm returned false, expected true")
	}

	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

}
