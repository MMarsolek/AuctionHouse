package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/MMarsolek/AuctionHouse/log"
	"github.com/MMarsolek/AuctionHouse/model"
	"github.com/MMarsolek/AuctionHouse/server"
	"github.com/MMarsolek/AuctionHouse/server/controller/auth"
	"github.com/MMarsolek/AuctionHouse/storage"
	"github.com/MMarsolek/AuctionHouse/storage/relational"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	serverParamPort             = "port"
	serverParamAdminUser        = "admin-username"
	serverParamAdminDisplayName = "admin-display-name"
	serverParamAdminPassword    = "admin-password"
)

var serverCmd = &cobra.Command{
	Use:               "server",
	Short:             "Starts a fileserver and accepts bids for auctions",
	Long:              "Starts a fileserver and accepts bids for auctions",
	PersistentPreRunE: bootstrapDB,
	RunE:              startServer,
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntP(serverParamPort, "p", 8080, "The port to use.")
	serverCmd.Flags().StringP(serverParamAdminUser, "u", "admin", "The default admin username")
	serverCmd.Flags().StringP(serverParamAdminDisplayName, "d", "admin", "The default admin display name")
	serverCmd.Flags().StringP(serverParamAdminPassword, "P", "admin", "The default admin password")
}

func startServer(cmd *cobra.Command, args []string) error {
	err := relational.CreateSchema(cmd.Context(), bunDB)
	if err != nil {
		return errors.Wrap(err, "unable to create database")
	}

	port, err := cmd.Flags().GetInt(serverParamPort)
	if err != nil {
		return errors.Wrap(err, "unable to get port")
	}

	userClient := relational.NewUserClient(bunDB)
	err = tryCreateDefaultAdmin(cmd, userClient)
	if err != nil {
		return errors.Wrap(err, "unable to create default admin")
	}

	address := fmt.Sprintf(":%d", port)
	ahServer := server.NewAuctionHouseServer(
		cmd.Context(),
		userClient,
		relational.NewAuctionItemClient(bunDB),
		relational.NewAuctionBidClient(bunDB),
		address,
	)

	incomingSignals := make(chan os.Signal, 1)
	signal.Notify(incomingSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-incomingSignals
		log.Info(cmd.Context(), "Caught signal", "signal", sig.String())
		if err := ahServer.Shutdown(cmd.Context()); err != nil {
			log.Error(cmd.Context(), "Error when trying to shut down server: %v", err)
		}
	}()

	log.Info(cmd.Context(), "Web server started", "address", address)
	if err := ahServer.ListenAndServe(); err != nil {
		log.Info(cmd.Context(), "Shutting down the server", "err", err)
	}

	return nil
}

func tryCreateDefaultAdmin(cmd *cobra.Command, userClient storage.UserClient) error {
	defaultAdminUser, err := cmd.Flags().GetString(serverParamAdminUser)
	if err != nil {
		return errors.Wrap(err, "unable to get default admin username")
	}
	_, err = userClient.Get(cmd.Context(), defaultAdminUser)
	if err != nil {
		if !errors.Is(err, storage.ErrEntityNotFound) {
			return errors.Wrap(err, "unable to connect to database storage")
		}

		defaultAdminDisplayName, err := cmd.Flags().GetString(serverParamAdminDisplayName)
		if err != nil {
			return errors.Wrap(err, "unable to get default admin display name")
		}

		defaultAdminPassword, err := cmd.Flags().GetString(serverParamAdminPassword)
		if err != nil {
			return errors.Wrap(err, "unable to get default admin password")
		}

		encodedPassword, err := auth.GenerateEncodedPassword(defaultAdminPassword)
		if err != nil {
			return errors.Wrap(err, "unable to hash password")
		}

		err = userClient.Create(cmd.Context(), &model.User{
			Username:       defaultAdminUser,
			DisplayName:    defaultAdminDisplayName,
			HashedPassword: encodedPassword,
			Permission:     model.PermissionLevelAdmin,
		})
	}

	return nil
}
