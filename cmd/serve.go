package cmd

import (
	"github.com/ehsundar/dopamine/internal/auth"
	"github.com/ehsundar/dopamine/internal/auth/token"
	"github.com/ehsundar/dopamine/internal/items"
	authMw "github.com/ehsundar/dopamine/pkg/middleware/auth"
	"github.com/ehsundar/dopamine/pkg/storage"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: serve,
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func serve(cmd *cobra.Command, args []string) {
	LoadConfig()

	s := storage.NewSqliteStorage(viper.GetString("dsn"))

	signingKey := token.LoadSigningKey()
	tokenManager := token.NewManager(signingKey)

	router := mux.NewRouter()

	items.RegisterHandlers(router, s)
	auth.RegisterHandlers(router, s, tokenManager)

	handlerFunc := http.Handler(router).ServeHTTP
	httpServer := authMw.NewAuthMiddleware(handlerFunc, tokenManager)

	serverAddr := "0.0.0.0:8080"
	log.Infof("serving http server at %s", serverAddr)
	err := http.ListenAndServe(serverAddr, httpServer)
	if err != nil {
		log.Error(err)
		panic(err)
	}
}
