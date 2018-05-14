package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mkideal/cli"
	"log"
	"math/big"
	"os"
)

var conn *ethclient.Client

var GethLocation string
var UsePort string
var UseIP string
var version string = "v0.0.1"

var decimals uint8

var help = cli.HelpCommand("display help information")

func main() {

	if err := cli.Root(root,
		cli.Tree(help),
		cli.Tree(child),
		cli.Tree(versionCli),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

type rootT struct {
	cli.Helper
}

var root = &cli.Command{
	Desc: "\n      #######################\n" +
		"           TokenBalance\n" +
		"      #######################\n\n" +
		"TokenBalance is an easy to use server that \n" +
		"give you your ERC20 token balance without \n" +
		"any troubles. Connects to your local geth \n" +
		"IPC and prints out a simple JSON response \n" +
		"for ethereum token balances.",
	// Argv is a factory function of argument object
	// ctx.Argv() is if Command.Argv == nil or Command.Argv() is nil
	Argv: func() interface{} { return new(rootT) },
	Fn: func(ctx *cli.Context) error {

		ctx.String("To start the tokenbalance server, use command:\ntokenbalance start --geth \"/root/ethereum/geth.ipc\" --port 8080 --ip 0.0.0.0\n * replace geth location with your own *\n")
		return nil
	},
}

// child command
type childT struct {
	cli.Helper
}

var child = &cli.Command{
	Name: "start",
	Desc: "run the tokenbalance http server",
	Argv: func() interface{} { return new(argT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		GethLocation = argv.Geth
		UsePort = argv.Port
		UseIP = argv.IP
		ConnectGeth()
		StartServer()

		return nil
	},
}

var versionCli = &cli.Command{
	Name: "version",
	Desc: "get the version of tokenbalance server",
	Argv: func() interface{} { return new(argT) },
	Fn: func(ctx *cli.Context) error {
		ctx.String(version + "\n")
		return nil
	},
}

type argT struct {
	cli.Helper
	Geth string `cli:"*g,geth" usage:"geth IPC location"`
	IP   string `cli:"ip" usage:"Bind to IP Address" dft:"0.0.0.0"`
	Port string `cli:"p,port" usage:"HTTP port for JSON" dft:"8080"`
}

func ConnectGeth() {
	var err error
	conn, err = ethclient.Dial(GethLocation)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v\n", err)
	} else {
		log.Printf("Connected to Geth at: %v\n", GethLocation)
	}
}

func GetAccount(contract string, wallet string) (string, string, string, uint8, string, int64, error) {
	var err error
	var symbol, tokenName string
	var address common.Address
	var getBlock *types.Block

	var ethAmount, balance *big.Int
	var tokenDecimals *big.Int

	token, err := NewTokenCaller(common.HexToAddress(contract), conn)
	if err != nil {
		log.Printf("Failed to instantiate a Token contract: %v\n", err)
		return "error", "0.0", "error", 0, "0.0", 0, err
	}

	getBlock, err = conn.BlockByNumber(context.TODO(), nil)
	if err != nil {
		log.Printf("Failed to get current block number: %v\n", err)
		getBlock = nil
	}

	maxBlock := getBlock.Number().Int64()

	address = common.HexToAddress(wallet)
	if err != nil {
		log.Println("Failed hex address: "+wallet, err)
		address = common.HexToAddress("0x")
	}

	ethAmount, err = conn.BalanceAt(context.TODO(), address, nil)
	if err != nil {
		log.Printf("Failed to get ethereum balance from address: %v | %v", address, err)
		ethAmount = big.NewInt(0)
	}

	balance, err = token.BalanceOf(nil, address)
	if err != nil {
		log.Println("Failed to get balance from contract: "+contract, err)
		balance = big.NewInt(0)
	}

	// the popular coin EOS doesn't have a symbol
	if common.HexToAddress(contract) == common.HexToAddress("0x86fa049857e0209aa7d9e616f7eb3b3b78ecfdb0") {
		symbol = "EOS"
	} else {
		symbol, err = token.Symbol(nil)
		if err != nil {
			log.Println("Failed to get symbol from contract: "+contract, err)
			symbol = "MISSING"
		}
	}
	tokenDecimals, err = token.Decimals(nil)
	if err != nil {
		log.Println("Failed to get decimals from contract: "+contract, err)
		tokenDecimals = big.NewInt(0)
	}
	tokenName, err = token.Name(nil)
	if err != nil {
		log.Println("Failed to retrieve token name from contract: "+contract, err)
		tokenName = "MISSING"
	}

	ethCorrected := BigIntDecimal(ethAmount, 18)
	tokenCorrected := BigIntDecimal(balance, int(tokenDecimals.Int64()))

	return tokenName, tokenCorrected, symbol, uint8(tokenDecimals.Uint64()), ethCorrected, maxBlock, err

}

func BigIntDecimal(balance *big.Int, decimals int) string {
	if balance.String() == "0" {
		return "0"
	}
	var newNum string
	for k, v := range balance.String() {
		if k == len(balance.String())-decimals {
			newNum += "."
		}
		newNum += string(v)
	}
	stringBytes := bytes.TrimRight([]byte(newNum), "0")
	newNum = string(stringBytes)
	if stringBytes[len(stringBytes)-1] == 46 {
		newNum += "0"
	}
	if stringBytes[0] == 46 {
		newNum = "0" + newNum
	}
	return newNum
}
