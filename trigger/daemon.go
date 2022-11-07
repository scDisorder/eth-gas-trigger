package trigger

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
	"log"
	"math/big"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type WatchOpts struct {
	Interval time.Duration `json:"interval"`
	// Gas is watching amount of gas value that triggers specified execution once gas price lower than this value
	Gas        *big.Int `json:"gas"`
	Cmd        string   `json:"cmd"`
	Repeatable bool     `json:"repeatable"`
}

func Run(ctx context.Context, opts WatchOpts) error {
	ticker := time.NewTicker(opts.Interval)

	ethProviderUrl := viper.GetString("eth.provider")
	client, err := ethclient.DialContext(ctx, strings.Replace(ethProviderUrl, "https", "wss", -1))
	if err != nil {
		log.Println(err)
		return err
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

tickerLoop:
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			log.Println("Context closed")
			break tickerLoop
		case _ = <-ticker.C:

			gasPrice, err := client.SuggestGasPrice(ctx)
			if err != nil {
				log.Printf("Unable to get gas price: %s", err)
			}

			triggers := gasPrice.Cmp(opts.Gas) < 0
			log.Printf("Current gas price: %d / Waiting for < %d", gasPrice, opts.Gas)
			// if gas price lower than limit
			if triggers {
				log.Printf("Gas price (%d) is lower than limit (%d)", gasPrice.Uint64(), opts.Gas.Uint64())
				cmd := exec.Command("/bin/sh", "-c", opts.Cmd)
				if out, err := cmd.CombinedOutput(); err != nil {

					log.Fatalf("Unable to execute command: %s", err)
				} else {
					fmt.Println(string(out))
				}
				log.Printf("Successfully executed: '%s'", opts.Cmd)

				if !opts.Repeatable {
					break tickerLoop
				}
			}
		case sig := <-sigchan:
			log.Printf("Signal received: %s", sig)
			break tickerLoop
		}
	}

	return nil
}
