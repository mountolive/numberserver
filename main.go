package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli"
)

func main() {
	// Flag parsing
	app := cli.NewApp()
	app.Name = "Number logger"
	app.Usage = `Writes numbers to defined log file.
               Numbers can have up to the max number
							 of digits defined by the user; 9 by default.
							 When the terminatio keyword is typed,
							 the program will attempt to shutdown gracefully`
	app.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:  "port, p",
			Value: 4000,
			Usage: "Port to be listened to, default: 4000",
		},
		&cli.BoolFlag{
			Name:  "append, a",
			Usage: "Whether to append to existing log file or recreate o start",
		},
		&cli.StringFlag{
			Name:  "logfile, l",
			Value: "./numbers.log",
			Usage: "Log file's path where the inputs would be written",
		},
		&cli.StringFlag{
			Name:  "termination, t",
			Value: "terminate",
			Usage: "Terminate keyword, for shutting down the server",
		},
		&cli.IntFlag{
			Name:  "digits, d",
			Value: 9,
			Usage: "Max number of digits permitted for int input",
		},
	}
	var port int
	var appender bool
	var logfile string
	var termination string
	var digits int
	app.Action = func(ctx cli.Context) error {
		port = ctx.GlobalInt("port")
		if port < 0 || port > 65535 {
			return errors.New("Port can't be a negative number, nor greater than 65535")
		}
		appender = ctx.GlobalBool("append")
		logfile = ctx.GlobalString("logfile")
		termination = ctx.GlobalString("termination")
		digits = ctx.GlobalInt("digits")
		if digits < 0 || digits > 9 {
			return errors.New("Digits can't be a negative number, nor greater than 9")
		}
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("An error occurred while trying to parse options: %v\n", err)
		fmt.Print("Aborting...")
		return
	}
	// Actual server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("An error occurred when trying to create the connection: %v\n", err)
		fmt.Println("Aborting...")
		return
	}
	fmt.Println("Starting number server. Welcome!")
	// Creating Logger (contains statistics)
	logger := NewLogger(Filename(logfile), Appender(appender))
	// Creating Number Checker
	checker := NewDefaultNumberChecker()
	checker.SetTermination(termination)
	checker.SetNumLimit(digits)
	// Creating Number Tracker
	maxCapacity, err := strconv.Atoi(strings.Repeat("9", digits))
	if err != nil {
		fmt.Printf("An error occurred when trying to create number tracker's limit: %v\n", err)
		fmt.Println("Aborting...")
		return
	}
	tracker, err := NewNumberTracker(maxCapacity)
	if err != nil {
		fmt.Printf("An error occurred when trying to create number tracker: %v\n", err)
		fmt.Println("Aborting...")
		return
	}
	// Handler
}
