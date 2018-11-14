package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jasonlvhit/gocron"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"

	"github.com/hexoul/go-coinmarketcap/statistics"
	"github.com/hexoul/go-coinmarketcap/types"
)

var (
	gitID        string
	gitPW        string
	targetSymbol string
	targetAddr   string
	targetQuotes = "USD"
	targetSlugs  = "binance"
	accessKey    = map[string]string{}
	secretKey    = map[string]string{}
)

func init() {
	for _, val := range os.Args {
		arg := strings.Split(val, "=")
		if len(arg) < 2 {
			continue
		} else if arg[0] == "-targetSymbol" {
			targetSymbol = arg[1]
		} else if arg[0] == "-targetAddr" {
			targetAddr = arg[1]
		} else if arg[0] == "-targetQuotes" {
			targetQuotes = arg[1]
		} else if arg[0] == "-targetSlugs" {
			targetSlugs = arg[1]
		} else if strings.Contains(arg[0], "accesskey") {
			accessKey[strings.Split(arg[0], ":")[0][1:]] = arg[1]
		} else if strings.Contains(arg[0], "secretkey") {
			secretKey[strings.Split(arg[0], ":")[0][1:]] = arg[1]
		} else if arg[0] == "-gitID" {
			gitID = arg[1]
		} else if arg[0] == "-gitPW" {
			gitPW = arg[1]
		}
	}
}

// GitPushChanges commits log changes and pushs it
func GitPushChanges() error {
	if gitID == "" || gitPW == "" {
		return nil
	}

	// Open
	r, err := git.PlainOpen("./")
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	// Commit
	if _, err = w.Commit("Commit report.log changed", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "hexoul",
			Email: "crosien@gmail.com",
			When:  time.Now(),
		},
		All: true,
	}); err != nil {
		return err
	}

	// Push
	r.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: gitID,
			Password: gitPW,
		},
	})
	return nil
}

func main() {
	fmt.Println("Initialzing...")

	cryptoQuoteOptions := &types.Options{
		Symbol:  targetSymbol,
		Convert: targetQuotes,
	}

	var exchangeMarketPairsOptions []*types.Options
	for _, slug := range strings.Split(targetSlugs, ",") {
		exchangeMarketPairsOptions = append(exchangeMarketPairsOptions, &types.Options{
			Slug:    slug,
			Convert: targetQuotes,
		})
		// statistics.TaskGatherExchangeMarketPairs(exchangeMarketPairsOptions[i])
	}

	// statistics.TaskGatherCryptoQuote(cryptoQuoteOptions)
	// statistics.TaskGatherTokenMetric(targetSymbol, targetAddr)

	fmt.Printf("Done\nScheduling...\n")

	// CryptoQuote
	statistics.GatherCryptoQuote(cryptoQuoteOptions, gocron.Every(10).Minutes())
	statistics.GatherCryptoQuote(cryptoQuoteOptions, gocron.Every(1).Day().At("00:00"))

	// ExchangeMarketPairs
	for i := range exchangeMarketPairsOptions {
		statistics.GatherExchangeMarketPairs(exchangeMarketPairsOptions[i], targetSymbol, gocron.Every(10).Minutes())
		statistics.GatherExchangeMarketPairs(exchangeMarketPairsOptions[i], targetSymbol, gocron.Every(1).Day().At("00:00"))
	}

	// TokenMetric
	statistics.GatherTokenMetric(targetSymbol, targetAddr, gocron.Every(30).Minutes())
	statistics.GatherTokenMetric(targetSymbol, targetAddr, gocron.Every(1).Day().At("00:00"))

	// Git
	gocron.Every(1).Hour().Do(GitPushChanges)

	fmt.Printf("Done\nStart!!\n")
	<-gocron.Start()
}
