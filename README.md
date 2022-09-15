# GoMAILER

[![License: 0BSD](https://img.shields.io/badge/license-0BSD-informational)](https://opensource.org/licenses/0BSD)
[![CI](https://github.com/rvhonorato/gomailer/actions/workflows/unittest.yml/badge.svg)](https://github.com/rvhonorato/gomailer/actions/workflows/unittest.yml)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/53caceaf22504b2c8020ca62a6d45367)](https://www.codacy.com/gh/rvhonorato/gomailer/dashboard?utm_source=github.com&utm_medium=referral&utm_content=rvhonorato/gomailer&utm_campaign=Badge_Grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/53caceaf22504b2c8020ca62a6d45367)](https://www.codacy.com/gh/rvhonorato/gomailer/dashboard?utm_source=github.com&utm_medium=referral&utm_content=rvhonorato/gomailer&utm_campaign=Badge_Coverage)

This is a simple CLI to send batch e-mails using a html-template and a list of recipients.

## Installation

1.  [Install go](https://go.dev/doc/install)

2.  Clone this repository
    ```text
    git clone https://github.com/rvhonorato/gomailer.git
    cd gomailer
    ```

3.  Run or build
    ```bash
    go run .
    # or
    go build
    ./gomailer
    ```

## Usage

1.  Create a template file (e.g. `template.html`)

     - The template must have a `{{.Name}}` placeholder for the recipient's name

2.  Create a recipient list file (e.g. `recipients.txt`)

     - In each line, the e-mail is captured as with a regular expression, all the rest is considered to be the recipients name

3.  Define `GMAIL_USER` and `GMAIL_PASSWORD` system variables, by default it uses gmail as the server; if you want to use another server, you need to change it in the source code.

     - `GMAIL_PASSWORD` is the [app password](https://support.google.com/accounts/answer/185833?hl=en), not the account password

4.  Execute;

    ```text
    $ ./gomailer -h
    Usage of ./gomailer:
    -s value
          Subject of the email
    -t value
          HTML template file of the email
    -u value
          user file containing email and name as tab separated values

    $ ./gomailer -t template.html -u recipients.txt -s "Hello!"
    ```

## Testing

```text
go test -v
```
