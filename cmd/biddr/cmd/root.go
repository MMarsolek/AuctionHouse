package cmd

import (
	"context"
	"os"

	"github.com/MMarsolek/AuctionHouse/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:           "biddr",
	Short:         "Starts a server that can accept bids for auctions",
	SilenceErrors: true,
}

// Execute acts as the entry point into this program.
func Execute(defaultCommand string) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	remainingArgs := os.Args[1:]
	if len(remainingArgs) == 0 {
		args := append([]string{defaultCommand}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}

	ctx := log.WithLogger(context.Background(), logger.Sugar())
	if err = rootCmd.ExecuteContext(ctx); err != nil {
		rootCmd.Printf("%+v\n", err)
		os.Exit(1)
	}
}
