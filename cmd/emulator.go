/*
Copyright © 2022 BoiseITGuru.find @Emerald City DAO

*/
package cmd

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/bjartek/overflow/overflow"
	"github.com/go-logfmt/logfmt"
	"github.com/gorilla/websocket"
	logrus "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//go:embed embed_files/flow.json
var flowConfig []byte

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
		checkFlowConfig()
		setupRoutes()
		fmt.Println("Now listening for connections to OurBetterPlayground on http://localhost:5050")
		log.Fatal(http.ListenAndServe(":5050", nil))
	},
}

var cliCmd *exec.Cmd
var tempFlowConfig string

// var o *overflow.Overflow

type jsonResponse struct {
	ResponseType string       `json:"responseType"`
	Data         logStructure `json:"data"`
}

type logStructure struct {
	Time  string `json:"time"`
	Level string `json:"level"`
	Msg   string `json:"msg"`
}

func checkFlowConfig() {
	if _, e := os.Stat("flow.json"); os.IsNotExist(e) {
		tempConfig, err := os.CreateTemp("", "flow-*.json")
		if err != nil {
			log.Fatal(err)
		}

		if _, err := tempConfig.Write(flowConfig); err != nil {
			log.Fatal(err)
		}

		tempConfig.Close()

		tempFlowConfig = tempConfig.Name()
		cliCmd = exec.Command("flow", "emulator", "--config-path", tempConfig.Name(), "--contracts")
	} else {
		cliCmd = exec.Command("flow", "emulator", "--contracts")
	}
}

func startEmulator(conn *websocket.Conn) {
	cmdReader, err := cliCmd.StdoutPipe()
	if err != nil {
		logrus.Error(err)
	}
	emulatorLogStream := logfmt.NewDecoder(cmdReader)
	go func() {
		for emulatorLogStream.ScanRecord() {
			var time string
			var level string
			var msg string
			var extra string

			for emulatorLogStream.ScanKeyval() {
				switch string(emulatorLogStream.Key()) {
				case "time":
					time = string(emulatorLogStream.Value())
				case "level":
					level = string(emulatorLogStream.Value())
				case "msg":
					msg = string(emulatorLogStream.Value())
				default:
					extra = string(emulatorLogStream.Key()) + " " + string(emulatorLogStream.Value())
				}
			}
			if emulatorLogStream.Err() != nil {
				panic(emulatorLogStream.Err())
			}

			emulatorLog := jsonResponse{
				ResponseType: "emulator-output",
				Data: logStructure{
					Time:  time,
					Level: level,
					Msg:   msg + ": " + extra,
				},
			}

			logrus.Info(msg + ": " + extra)
			conn.WriteJSON(emulatorLog)
		}
	}()

	if err := cliCmd.Start(); err != nil {
		logrus.Error(err)
	}

	// quick hack to make sure emulator is running first
	time.Sleep(3 * time.Second)

	var overflowConfig *overflow.OverflowBuilder

	if tempFlowConfig != "" {
		overflowConfig = overflow.NewOverflowEmulator().Config(tempFlowConfig)
	} else {
		overflowConfig = overflow.NewOverflowEmulator()
	}

	overflowConfig.Start()
}

func stopEmulator() {
	cliCmd.Process.Signal(syscall.SIGINT)
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

	formattedResponse := jsonResponse{
		ResponseType: "emulator-status",
		Data: logStructure{
			Time:  "",
			Level: "",
			Msg:   "started",
		},
	}

	if err = ws.WriteJSON(formattedResponse); err != nil {
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
	logrus.SetOutput(os.Stdout)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// emulatorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// emulatorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
