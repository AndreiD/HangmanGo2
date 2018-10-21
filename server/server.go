package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/AndreiD/HangmanGo2/api"
	"github.com/AndreiD/HangmanGo2/server/configs"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// main start a gRPC server and waits for connection
func main() {

	// starts the server
	go func() {
		err := startGRPCServer(configs.Load(), "cert/server.crt", "cert/server.key")
		if err != nil {
			log.Fatalf("failed to start gRPC server: %s", err)
		}
	}()

	// forever
	select {}

}

// this will be called in a goroutine
func startGRPCServer(config *viper.Viper, certFile, keyFile string) error {

	//create a TCP listener
	address := config.GetString("hostname") + ":" + strconv.Itoa(config.GetInt("port"))

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// Create the TLS credentials
	creds, err := credentials.NewServerTLSFromFile("cert/server.crt", "cert/server.key")
	if err != nil {
		return fmt.Errorf("could not load TLS keys: %s", err)
	}

	// Create an array of gRPC options with the credentials & with login interceptor
	opts := []grpc.ServerOption{grpc.Creds(creds), grpc.UnaryInterceptor(UnaryInterceptor)}

	// create a gRPC server object
	grpcServer := grpc.NewServer(opts...)

	// attach the Ping service to the server
	api.RegisterHangmanServer(grpcServer, &hangman{})

	// graceful shutdown
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
		<-sigs
		// cleanup here
		grpcServer.GracefulStop()
	}()

	fmt.Println(" ")
	fmt.Println("██╗  ██╗ █████╗ ███╗   ██╗ ██████╗ ███╗   ███╗ █████╗ ███╗   ██╗")
	fmt.Println("██║  ██║██╔══██╗████╗  ██║██╔════╝ ████╗ ████║██╔══██╗████╗  ██║")
	fmt.Println("███████║███████║██╔██╗ ██║██║  ███╗██╔████╔██║███████║██╔██╗ ██║")
	fmt.Println("██╔══██║██╔══██║██║╚██╗██║██║   ██║██║╚██╔╝██║██╔══██║██║╚██╗██║")
	fmt.Println("██║  ██║██║  ██║██║ ╚████║╚██████╔╝██║ ╚═╝ ██║██║  ██║██║ ╚████║")
	fmt.Println("╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝ ╚═════╝ ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝")
	fmt.Println(" ")

	fmt.Printf(">>>>> Hangman Server Started on port %d\n", config.GetInt("port"))
	if (config.GetString("environment")) == "debug" {
		fmt.Println(">>>>> Server runs in DEBUG mode...")
	} else {
		fmt.Println(" >>>>> Server runs in production mode...")
	}
	fmt.Println(">>>>> Waiting for clients to connect. Hit CTRL+C to terminate")

	// start the server
	if err := grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve: %s", err)
	}

	return nil
}
