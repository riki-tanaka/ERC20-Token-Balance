package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type BalanceResponse struct {
	Name       string  `json:"name,omitempty"`
	Wallet     string  `json:"wallet,omitempty"`
	Symbol     string  `json:"symbol,omitempty"`
	Balance    float64 `json:"balance"`
	EthBalance float64 `json:"eth_balance,omitempty"`
	Decimals   uint8   `json:"decimals,omitempty"`
	Block      uint64  `json:"block,omitempty"`
}

type ErrorResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func getInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	contract := vars["contract"]
	wallet := vars["wallet"]

	log.Println("Fetching Wallet:", wallet, "at Contract:", contract)

	name, balance, token, decimals, ethAmount, block, err := GetAccount(contract, wallet)

	if err != nil {
		m := ErrorResponse{
			Error:   true,
			Message: "could not find contract address",
		}
		msg, _ := json.Marshal(m)
		w.Write(msg)
		return
	}

	new := BalanceResponse{
		Name:       name,
		Symbol:     token,
		Wallet:     wallet,
		Balance:    balance,
		EthBalance: ethAmount,
		Decimals:   decimals,
		Block:      block,
	}

	j, err := json.Marshal(new)

	if err == nil {
		w.Write(j)
	}
}

func getTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	contract := vars["contract"]
	wallet := vars["wallet"]

	log.Println("Fetching Wallet:", wallet, "at Contract:", contract)

	_, balance, _, _, _, _, err := GetAccount(contract, wallet)

	if err != nil {
		m := ErrorResponse{
			Error:   true,
			Message: "could not find contract address",
		}
		msg, _ := json.Marshal(m)
		w.Write(msg)
		return
	} else {

		new := BalanceResponse{
			Balance: balance,
		}
		j, _ := json.Marshal(new)

		w.Write(j)
	}

}

func StartServer() {

	r := mux.NewRouter()
	r.HandleFunc("/balance/{contract}/{wallet}", getInfoHandler).Methods("GET")
	r.HandleFunc("/token/{contract}/{wallet}", getTokenHandler).Methods("GET")

	log.Println("TokenBalance Server Running: http://" + UseIP + ":" + UsePort)

	http.Handle("/", r)
	http.ListenAndServe(UseIP+":"+UsePort, nil)

}
