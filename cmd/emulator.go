/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"syscall"

	// "github.com/bjartek/overflow/overflow"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

// emulatorCmd represents the emulator command
var emulatorCmd = &cobra.Command{
	Use:   "emulator",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		setupRoutes()
		fmt.Println("Now listening for connections to OurBetterPlayground on http://localhost:5050")
		log.Fatal(http.ListenAndServe(":5050", nil))
	},
}

var cliCmd = exec.Command("flow", "emulator", "--contracts")

func startEmulator(conn *websocket.Conn) {
	cmdReader, err := cliCmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
			if err := conn.WriteMessage(1, []byte(scanner.Text())); err != nil {
				log.Println(err)
				return
			}
		}
	}()

	if err := cliCmd.Start(); err != nil {
		log.Fatal(err)
	}
}

func stopEmulator() {
	cliCmd.Process.Signal(syscall.SIGTERM)
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Playground Connected")
	log.Println("Starting Emulator")
	startEmulator(ws)
	err = ws.WriteMessage(1, []byte("Emulator started successfully"))
	if err != nil {
		log.Println(err)
	}

	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	reader(ws)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				stopEmulator()
			}
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

func setupRoutes() {
	http.HandleFunc("/", wsEndpoint)
}

func init() {
	rootCmd.AddCommand(emulatorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// emulatorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// emulatorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
