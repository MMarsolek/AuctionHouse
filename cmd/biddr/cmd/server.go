package cmd

import (
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts a fileserver and accepts bids for auctions",
	Long:  "Starts a fileserver and accepts bids for auctions",
	RunE:  startServer,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func startServer(cmd *cobra.Command, args []string) error {
	// fs := http.FileServer(http.Dir("./web"))
	// server := http.Server{}
	// mux := http.NewServeMux()
	// mux.HandleFunc("/whatever", func(rw http.ResponseWriter, r *http.Request) {
	// 	rw.WriteHeader(http.StatusTeapot)
	// })
	// mux.Handle("/web/", http.StripPrefix("/web", fs))
	// server.Handler = mux
	// server.Addr = ":8080"

	// if err := server.ListenAndServe(); err != nil {
	// 	fmt.Printf("Error shutting down: %v", err)
	// }
	cmd.Println("Hello world!")
	return nil
}
