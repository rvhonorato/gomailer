package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
)

// --------------------------------------------------------------------------------
// Define the struct types
type user struct {
	Name  string
	Email string
}

type stringFlag struct {
	set   bool
	value string
}

// --------------------------------------------------------------------------------
// Functions
func check(e error) {
	// Check for errors
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}

func readUsers(filename string, delimiter string) []user {
	// Read the users from the file
	var ps []user

	path, err := filepath.Abs(filename)
	check(err)

	inFile, err := os.Open(path)
	check(err)
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	for _, l := range lines {
		// fmt.Println(l)
		s := strings.Split(l, delimiter)
		ps = append(ps, user{
			Email: s[0],
			Name:  strings.Join(s[1:], " "),
		})

	}

	return ps

}

func (sf *stringFlag) Set(x string) error {
	// Set the string flag
	sf.value = x
	sf.set = true
	return nil
}

func (sf *stringFlag) String() string {
	// Return the string flag
	return sf.value
}

func confirm() bool {
	scanner := bufio.NewScanner(os.Stdin)

	log.Println("Continue? (y/N)")
	var confirmation string
	if scanner.Scan() {
		confirmation = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return false
	}

	return confirmation == "y"
}

func main() {

	// Define the flags as custom stringFlag types
	var userfile stringFlag
	var templatefile stringFlag
	var subjectstring stringFlag

	flag.Var(&userfile, "userfile", "user file containing email and name as tab separated values")
	flag.Var(&templatefile, "template", "HTML template file of the email")
	flag.Var(&subjectstring, "subject", "Subject of the email")

	// Parse the command line flags
	flag.Parse()

	// Check if the userfile and templatefile flags are set
	if !userfile.set {
		log.Fatalln("-userfile not set")
	}

	if !templatefile.set {
		log.Fatalln("-template not set")
	}

	if !subjectstring.set {
		log.Fatalln("-subject not set")
	}

	// Check if the files exist
	if _, err := os.Stat(userfile.value); errors.Is(err, os.ErrNotExist) {
		log.Fatalln("user file", userfile.value, "does not exist")
	}

	if _, err := os.Stat(templatefile.value); errors.Is(err, os.ErrNotExist) {
		log.Fatalln("template file", templatefile.value, "does not exist")
	}

	// Check if the system variables are defined
	from := os.Getenv("GMAIL_USER")
	if from == "" {
		log.Fatalln("GMAIL_USER environment variable not set")
	}

	password := os.Getenv("GMAIL_PASSWORD")
	if password == "" {
		log.Fatal("GMAIL_PASSWORD environment variable not set")
	}

	// Start the mailer
	log.Println("Starting mailer")

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Read the users from the file
	log.Println("Looking for", userfile.value, "file...")
	ps := readUsers(userfile.value, "\t")
	if len(ps) == 0 {
		log.Fatalln("No users found")
	}

	// Report how many users found
	log.Printf("Found %d users", len(ps))

	// Display the users/emails and ask for confirmation
	for _, v := range ps {
		log.Println(">", "Name:", "<", v.Name, ">", "email:", "<", v.Email, ">")
	}

	confirmation := confirm()
	if !confirmation {
		log.Println("Exiting...")
		os.Exit(0)
	}

	// For each user, send the email
	log.Println("Sending emails...")
	for _, v := range ps {
		to := []string{v.Email}

		log.Println("> Sending to", "<", v.Email, ">")

		// Read the template
		t, _ := template.ParseFiles(templatefile.value)

		// Add the headers to the message
		mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
		var body bytes.Buffer
		body.Write([]byte(fmt.Sprintf("Subject: %s\n%s\n\n", subjectstring.value, mimeHeaders)))

		t.Execute(&body, struct {
			Name string
		}{
			Name: v.Name,
		})

		// Send email
		//  Sending "Bcc": messages is accomplished by including an email address in the to parameter
		//  but not including it in the msg headers.
		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
		if err != nil {
			log.Fatalln(err)
		}
	}
}
