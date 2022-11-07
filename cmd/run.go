package cmd

import (
	"context"
	"github.com/ethereum/go-ethereum/params"
	"github.com/scDisorder/eth-gas-trigger/trigger"
	"github.com/spf13/cobra"
	"math/big"
	"time"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run gas price watcher and execute callback",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		interval, _ := cmd.Flags().GetDuration("interval")
		priceInGwei, _ := cmd.Flags().GetInt64("gwei")
		command, _ := cmd.Flags().GetString("cmd")
		repeat, _ := cmd.Flags().GetBool("repeatable")

		opts := trigger.WatchOpts{
			Interval:   interval,
			Gas:        new(big.Int).SetUint64(uint64(priceInGwei) * params.GWei),
			Cmd:        command,
			Repeatable: repeat,
		}

		return trigger.Run(ctx, opts)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	runCmd.Flags().DurationP("interval", "i", 15*time.Second, "Interval for price check call")
	runCmd.Flags().BoolP("repeatable", "r", false, "Repeatable execution")
	runCmd.Flags().Int64("gwei", 0, "Gas price in GWei")
	runCmd.Flags().StringP("cmd", "c", "", "Command to execute")
	_ = runCmd.MarkFlagRequired("gwei")
	_ = runCmd.MarkFlagRequired("cmd")
}
