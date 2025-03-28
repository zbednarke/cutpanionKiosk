package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func getClient(config *oauth2.Config) *http.Client {
	tokenFile := "token.json"
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

	tok, _ := config.Exchange(context.Background(), authCode)
	return tok
}

func GetSheetsService() *sheets.Service {
	b, _ := os.ReadFile("credentials.json")
	config, _ := google.ConfigFromJSON(b, sheets.SpreadsheetsReadonlyScope)
	client := getClient(config)
	srv, _ := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	return srv
}

func GetCalendarService() *calendar.Service {
	b, _ := os.ReadFile("credentials.json")
	config, _ := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	client := getClient(config)
	srv, _ := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	return srv
}
