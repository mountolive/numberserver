package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli"
)

func main() {
	// **** Flag parsing ****
	app := cli.NewApp()
	app.Name = "Number logger"
	app.Usage = `Writes numbers to defined log file.
               Numbers can have up to the max number
							 of digits defined by the user: 9 by default.
							 When the termination ("terminate") keyword is prompted,
							 the program will attempt to shutdown gracefully`
	app.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:  "port, p",
			Value: 4000,
			Usage: "Port to be listened to",
		},
		&cli.BoolFlag{
			Name:  "append, a",
			Usage: "Whether to append to existing log file or recreate on start",
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
		&cli.IntFlag{
			Name:  "interval, i",
			Value: 10,
			Usage: "Show statistics every * seconds",
		},
	}
	var port int
	var appender bool
	var logfile string
	var termination string
	var digits int
	var interval int
	app.Action = func(ctx *cli.Context) error {
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
		interval = ctx.GlobalInt("interval")
		if interval < 0 {
			return errors.New("Statistics' interval can't be negative")
		}
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("An error occurred while trying to parse options: %v\n", err)
		fmt.Print("Aborting...")
		return
	}
	// Using termination as a flag to terminate the script (if not set)
	// It won't be set, for example, if the user calls the --help subcommand
	if termination == "" {
		return
	}
	// **** Actual server ****
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
	// Global context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// When shuttingdown
	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt, os.Kill)
	go gracefulShutdown(exit, cancel, listener)
	// **** Handler ****
	// For periodic printing of statistics
	ticker := time.Tick(time.Second * time.Duration(interval))
	// Coordination channels
	intInput := make(chan int)
	processChan := tracker.ProcessNumber(ctx, intInput)
	// Writing to logfile
	logger.StreamWrite(ctx, processChan)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Termination found, exiting...")
			return
		default:
			// Handling connections
			conn, err := listener.Accept()
			if err != nil {
				fmt.Printf("The server stopped accepting connections (%v) \n", err)
				return
			}
			// Processing connection
			go func() {
				defer conn.Close()
				scanner := bufio.NewScanner(conn)
				// Reading each client's input
				for scanner.Scan() {
					select {
					// Print statistics every 10 seconds
					case <-ticker:
						tracker.PrintStatistics()
					default:
						input := scanner.Text()
						if checker.CheckTermination(input) {
							// Cancelling global context, connection and server
							cancel()
							conn.Close()
							listener.Close()
							return
						}
						if checker.ValidateInput(input) {
							value, err := strconv.Atoi(input)
							// Should be unreachable (given the ValidateInput)
							if err != nil {
								fmt.Printf("An error occurred while processing req: %s. Err: %v", input, err)
								return
							}
							intInput <- value
						}
					}
				}
			}()
		}
	}
}

func gracefulShutdown(exit <-chan os.Signal, cancel context.CancelFunc, listener net.Listener) {
	<-exit
	fmt.Println("Received kill/intrrupt signal...")
	cancel()
	listener.Close()
}
