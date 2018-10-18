package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/spf13/viper"

	"github.com/AndreiD/HangmanGo2/api"
	"github.com/AndreiD/HangmanGo2/client/configs"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Config - let's store the config globally, for now
var Config *viper.Viper

var rootCmd = &cobra.Command{
	Use:   "Hangman Game",
	Short: "Guessing letters game",
	Long: `Hangman is a paper and pencil guessing game for two or more players.
One player thinks of a word, phrase or sentence and the other(s) tries to guess it by suggesting letters or numbers`,
}

// Execute (cobra wants it exported)
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var hangmanClient api.HangmanClient

func init() {

	// load the configs
	Config = configs.Load()

	if (Config.GetString("environment")) == "debug" {
		fmt.Println(">>>>> Client runs in DEBUG mode...")
	} else {
		fmt.Println(" >>>>> Client runs in production mode...")
	}

	// Create the client TLS credentials
	creds, err := credentials.NewClientTLSFromFile("cert/server.crt", "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}

	bindSocket := Config.GetString("hostname") + ":" + strconv.Itoa(Config.GetInt("port"))

	// Setup the login/pass
	auth := Authentication{
		Login:    Config.GetString("auth.username"),
		Password: Config.GetString("auth.password"),
	}

	// grpc dialing options
	dialOpts := []grpc.DialOption{}
	dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))
	dialOpts = append(dialOpts, grpc.WithPerRPCCredentials(&auth))

	ctx, cancel := AppContext()
	defer cancel()
	conn, err := grpc.DialContext(ctx, bindSocket, dialOpts...)

	if err != nil {
		log.Panicf("failed to dial %q: %s", bindSocket, err.Error())
	}
	hangmanClient = api.NewHangmanClient(conn)

	fmt.Printf(">>>>> Hangman Client Started on port %d\n", Config.GetInt("port"))
}

// AppContext returns the app context with a 2 second timeout
func AppContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 2*time.Second)
}

// Authentication holds the login/password
type Authentication struct {
	Login    string
	Password string
}

// GetRequestMetadata gets the current request metadata
func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"login":    a.Login,
		"password": a.Password,
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires transport security
func (a *Authentication) RequireTransportSecurity() bool {
	return true
}
