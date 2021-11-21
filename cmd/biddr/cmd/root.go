package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/MMarsolek/AuctionHouse/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"go.uber.org/zap"
)

const databaseFileName = "biddr.db"

var bunDB *bun.DB

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

func bootstrapDB(cmd *cobra.Command, args []string) error {
	rawDB, err := sql.Open(sqliteshim.ShimName, fmt.Sprintf("file:%s?cache=shared", databaseFileName))
	if err != nil {
		return errors.Wrap(err, "unable to open database")
	}

	bunDB = bun.NewDB(rawDB, sqlitedialect.New())
	return nil
}
