package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func GetService(ctx context.Context) (*gmail.Service, *http.Client, error) {
	client, err := GetClient(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not retrieve token: %v", err)
	}
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to retrieve Gmail client: %v", err)
	}
	return srv, client, err
}

// Retrieve a token, saves the token, then returns the generated client.
func GetClient(ctx context.Context) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	clientInfo, err := os.ReadFile("./secrets/gmail_oauth.json")
	if err != nil {
		return nil, fmt.Errorf("Unable to read client file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(clientInfo, gmail.GmailReadonlyScope)
	tokFile := "./secrets/gmail_token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return nil, fmt.Errorf("Error while retrieving token from web: %v", err)
		}
		saveToken(tokFile, tok)
	}
	return config.Client(ctx, tok), nil
}

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	// Listen on a random available port
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, err
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	redirectURL := fmt.Sprintf("http://localhost:%d/callback", port)
	config.RedirectURL = redirectURL

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Opening Google OAuth link in your browser:\n%v\n", authURL)
	OpenBrowser(authURL)

	codeChan := make(chan string)

	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/callback" {
				code := r.URL.Query().Get("code")
				codeChan <- code
				fmt.Fprintf(w, "Authorization successful! You can close this window now.")
			}
		}),
	}

	// Start the server in a goroutine
	go func() {
		if err := server.Serve(listener); err != http.ErrServerClosed {
			log.Printf("HTTP server Serve: %v", err)
		}
	}()

	// Wait for the auth code
	authCode := <-codeChan

	if err = server.Shutdown(context.Background()); err != nil {
		log.Printf("Error server sutdown %v", err)
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		return nil, err
	}
	return tok, nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)

	if tok.Expiry.Second() < time.Now().Second() {
		fmt.Println("Token expired")
		return nil, fmt.Errorf("Token expired")
	}

	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
