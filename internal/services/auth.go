package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var client *http.Client // shared client for all Google APIs
func getProjectRoot() string {
	// exePath, err := os.Executable()
	// if err != nil {
	// 	log.Fatalf("Unable to get executable path: %v", err)
	// }

	// // Resolve root directory (assumes binary lives in cmd/server or similar)
	// return filepath.Dir(filepath.Dir(filepath.Dir(exePath))) // up 3 levels: cmd/server -> cmd -> root
	return "C:\\Users\\PC\\projects\\cutpanionKiosk"
}

func getCredentialsPath() string {

	return filepath.Join(getProjectRoot(), "credentials.json")
}

func initGoogleClient() {

	b, err := os.ReadFile(getCredentialsPath())
	if err != nil {
		log.Fatalf("Unable to read credentials file: %v", err)
	}

	// Include both scopes here
	config, err := google.ConfigFromJSON(b,
		sheets.SpreadsheetsReadonlyScope,
		calendar.CalendarReadonlyScope,
	)
	if err != nil {
		log.Fatalf("Unable to parse config from JSON: %v", err)
	}

	client = getClient(config)
}

func getClient(config *oauth2.Config) *http.Client {
	tokenFile := filepath.Join(getProjectRoot(), "token.json")
	token, err := tokenFromFile(tokenFile)
	if err != nil {
		token = getTokenFromWeb(config)
		saveToken(tokenFile, token)
	}
	return config.Client(context.Background(), token)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	return t, err
}

func saveToken(path string, token *oauth2.Token) {
	f, _ := os.Create(path)
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("Go to this URL in your browser:\n%v\n", authURL)

	var authCode string
	log.Print("Enter the authorization code: ")
	_, _ = fmt.Scan(&authCode)

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func GetSheetsService() *sheets.Service {
	if client == nil {
		initGoogleClient()
	}
	srv, _ := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	return srv
}

func GetCalendarService() *calendar.Service {
	if client == nil {
		initGoogleClient()
	}
	srv, _ := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	return srv
}
