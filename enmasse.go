/*
 * enmasse.go
 *
 * Copyright (c) 2016 Pelle Beckman, http://beckman.io
 * 
 * Parts of this code provided by Google, licensed
 * under the Apache 2.0 License
 */
package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

type JSONDataFile []map[string]interface{}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("enmasse-credentials.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {

	templateFilePath := flag.String("template", "", "Template file")
	dataFilePath := flag.String("data", "", "JSON data file")

	flag.Parse()

	if *templateFilePath == "" {
		log.Fatalf("No template file provided!")
	} else if *dataFilePath == "" {
		log.Fatalf("No JSON data file provided!")
	}

	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	config, err := google.ConfigFromJSON(b, gmail.GmailComposeScope,
		gmail.GmailLabelsScope,
		gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	tmpl, err := template.ParseFiles(*templateFilePath)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	var emailData JSONDataFile
	dataFile, err := ioutil.ReadFile(*dataFilePath)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	err = json.Unmarshal(dataFile, &emailData)
	if err != nil {
		log.Fatalf("Unmarshal error: %v", err)
	}

	for _, element := range emailData {
		address := element["email"].(string)

		_, err = mail.ParseAddress(address)
		if err != nil {
			log.Printf("Bad recipient email address (%v) - skipping\n", address)
			continue
		}

		subject := element["subject"].(string)

		var doc bytes.Buffer

		err = tmpl.Execute(&doc, element)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		rendered := doc.String()

		var draft gmail.Draft
		var message gmail.Message

		messageStr := []byte(
			"To: " + address + "\r\n" +
				"Subject: " + subject + "\r\n\r\n" +
				rendered)

		message.Raw = base64.StdEncoding.EncodeToString(messageStr)
		draft.Message = &message
		_, err = srv.Users.Drafts.Create("me", &draft).Do()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}
}
