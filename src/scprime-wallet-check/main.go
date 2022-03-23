package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-go/types"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	Server string
	Port int
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "scprime-wallet-check",
			Short:    "Check if ScPrime Wallet is unlocked.",
			Keyspace: "sensu.io/plugins/sensu-scprime-check/config",
		},
	}

	options = []*sensu.PluginConfigOption{
		&sensu.PluginConfigOption{
			Path:      "port",
			Env:       "CHECK_PORT",
			Argument:  "port",
			Shorthand: "p",
			Default:   4280,
			Usage:     "ScPrime server port to check",
			Value:     &plugin.Port,
		},
		&sensu.PluginConfigOption{
			Path:      "server",
			Env:       "CHECK_SERVER",
			Argument:  "server",
			Shorthand: "s",
			Default:   "localhost",
			Usage:     "ScPrime server to check",
			Value:     &plugin.Server,
		},
	}
)

func main() {
	useStdin := false
	fi, err := os.Stdin.Stat()
	if err != nil {
		fmt.Printf("Error check stdin: %v\n", err)
		panic(err)
	}
	//Check the Mode bitmask for Named Pipe to indicate stdin is connected
	if fi.Mode()&os.ModeNamedPipe != 0 {
		log.Println("using stdin")
		useStdin = true
	}

	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, useStdin)
	check.Execute()
}

func checkArgs(event *types.Event) (int, error) {
	
	return sensu.CheckStateOK, nil
}

func executeCheck(event *types.Event) (int, error) {
	// Check state of Wallet
	res, err := httpScPrime("/wallet")
	if err != nil {
		fmt.Printf("%s CRITICAL: failed to run check, error: %v\n", plugin.PluginConfig.Name, err)
		return sensu.CheckStateCritical, nil
	}
	walletState, err := checkWallet(res)
	if err != nil {
		fmt.Printf("%s CRITICAL: failed to run check, error: %v\n", plugin.PluginConfig.Name, err)
		return sensu.CheckStateCritical, nil
	}
	
	if walletState {
		fmt.Printf("%s OK: Wallet unlocked\n", plugin.PluginConfig.Name)
		return sensu.CheckStateOK, nil
	} else {
		fmt.Printf("%s WARNING: Wallet not unlocked\n", plugin.PluginConfig.Name)
		return sensu.CheckStateWarning, nil
	}
}

func httpScPrime (endpoint string) ([]byte, error) {
	baseURL := "http://"+plugin.Server+":"+strconv.Itoa(plugin.Port)

	scprimeClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}
	req, err := http.NewRequest(http.MethodGet, baseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "ScPrime-Agent")

	res, getErr := scprimeClient.Do(req)
	if getErr != nil {
		return nil, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	return body, nil
}

func checkWallet (body []byte) (bool, error) {
	type wallet struct {
		Unlocked bool `json:"unlocked"`
	}

	walletState := wallet{}
	jsonErr := json.Unmarshal(body, &walletState)
	if jsonErr != nil {
		return false, jsonErr
	}

	return walletState.Unlocked, nil
}