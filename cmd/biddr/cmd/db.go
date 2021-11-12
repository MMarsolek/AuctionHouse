package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/MMarsolek/AuctionHouse/storage/relational"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

const databaseFileName = "biddr.db"

var bunDB *bun.DB

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database related sub commands",
	Long:  "Database related sub commands",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		rawDB, err := sql.Open(sqliteshim.ShimName, fmt.Sprintf("file:%s?cache=shared", databaseFileName))
		if err != nil {
			return errors.Wrap(err, "unable to open database")
		}

		bunDB = bun.NewDB(rawDB, sqlitedialect.New())
		return nil
	},
}

var dbInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes the database",
	Long:  "Initializes the database",
	RunE:  initDatabase,
}

func init() {
	dbCmd.AddCommand(dbInitCmd)
	rootCmd.AddCommand(dbCmd)
}

func initDatabase(cmd *cobra.Command, args []string) error {
	if _, err := os.Stat(databaseFileName); err == nil {
		var input string
		fmt.Print("Database already exists. Drop and recreate tables? [y/N]: ")
		fmt.Scanln(&input)
		input = strings.ToLower(input)
		if !strings.HasPrefix(input, "y") {
			return nil
		}

		err = os.Remove(databaseFileName)
		if err != nil {
			return errors.Wrapf(err, "could not delete file %s", databaseFileName)
		}
	}
	err := relational.CreateSchema(cmd.Context(), bunDB)
	if err != nil {
		return errors.Wrap(err, "unable to create schema")
	}
	return nil
}
