package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const MONO_API_URL = "https://api.monobank.ua/personal"

type MonoClient struct {
	ClientId    string `json:"clientId"`
	Name        string `json:"name"`
	Permissions string `json:"permissions"`
	WebHookUrl  string `json:"webHookUrl,omitempty"`
	Accounts    []struct {
		Id           string `json:"id"`
		SendId       string `json:"sendId"`
		Balance      int64  `json:"balance"`
		CreditLimit  int64  `json:"creditLimit"`
		Type         string `json:"type"`
		CurrencyCode int32  `json:"currencyCode"`
		Cashback     string `json:"cashbackType"`
	} `json:"accounts"`
	Jars []struct {
		Id           string `json:"id"`
		SendId       string `json:"sendId"`
		Title        string `json:"title"`
		Description  string `json:"description"`
		CurrencyCode int32  `json:"currencyCode"`
		Balance      int64  `json:"balance"`
		Goal         int64  `json:"goal"`
	} `json:"jars"`
}

type Statement struct {
	Id              string `json:"id"`
	Time            int64  `json:"time"`
	Description     string `json:"description"`
	Mcc             int32  `json:"mcc"`
	OriginalMcc     int32  `json:"originalMcc"`
	Hold            bool   `json:"bool"`
	Amount          int64  `json:"amount"`
	OperationAmount int64  `json:"operationAmount"`
	CurrencyCode    int32  `json:"currencyCode"`
	CommissionRate  int64  `json:"commissionRate"`
	CashbackAmount  int64  `json:"cashbackAmount"`
	Balance         int64  `json:"balance"`
	Comment         string `json:"comment,omitempty"`
	ReceiptId       string `json:"receiptId,omitempty"`
	InvoiceId       string `json:"invoiceId,omitempty"`
	CounterEdrpou   string `json:"counterEdrpou,omitempty"`
	CounterIban     string `json:"counterIban,omitempty"`
}

func initRequest(method string, urlElem ...string) (*http.Request, error) {
	api_url, err := url.JoinPath(MONO_API_URL, urlElem...)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, api_url, nil)
	if err != nil {
		return req, err
	}
	req.Header.Add("X-Token", os.Getenv("MONO_API_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	return req, err
}

func makeRequest(method string, objectToUnmarshal any, urlElem ...string) error {
	req, err := initRequest("GET", urlElem...)
	if err != nil {
		return err
	}
	hc := &http.Client{}
	r, err := hc.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	var body []byte
	if body, err = io.ReadAll(r.Body); err != nil {
		return err
	}

	if err = json.Unmarshal(body, objectToUnmarshal); err == nil {
		return err
	}
	return nil
}

func getClientInfo() (*MonoClient, error) {
	client := MonoClient{}
	if err := makeRequest("GET", &client, "client-info"); err != nil {
		return nil, err
	}
	return &client, nil
}

func getStatement(account, from, to string) ([]Statement, error) {
	var stat []Statement
	if err := makeRequest("GET", &stat, "statement", account, from, to); err != nil {
		return nil, err
	}
	return stat, nil
}

func main() {
	var client *MonoClient
	client, err := getClientInfo()
	if err != nil {
		log.Fatalf("Cannot get client info: %s\n", err)
	}
	log.Printf("Got client data: %v\n", *client)

	var statements []Statement
	from := time.Date(2023, time.March, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, time.March, 10, 0, 0, 0, 0, time.UTC)
	log.Printf("from: %v to: %v", from, to)

	statements, err = getStatement(os.Getenv("MONO_ACCOUNT_ID"), strconv.FormatInt(from.Unix(), 10), strconv.FormatInt(to.Unix(), 10))
	if err != nil {
		log.Fatalf("Cannot get statement")
	}
	log.Printf("Got statement: %v\n", statements)
}
