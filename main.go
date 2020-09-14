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
							     the program will attempt to shutdown gracefully.
							     This termination keyword can be changed on start (see --help)`
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
			Usage: "Max number of digits permitted for int input (max: 9)",
		},
		&cli.IntFlag{
			Name:  "interval, i",
			Value: 10,
			Usage: "Show statistics every * seconds",
		},
		&cli.IntFlag{
			Name:  "maxconn, c",
			Value: 5,
			Usage: "Max number of concurrent connections allowed",
		},
	}
	// Flag variables
	var port int
	var appender bool
	var logfile string
	var termination string
	var digits int
	var interval int
	var maxconn int
	// Parsing of flags
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
		maxconn = ctx.GlobalInt("maxconn")
		if maxconn < 0 {
			return errors.New("The number of max concurrent connections can't be negative")
		}
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("An error occurred while trying to parse options: %v\n", err)
		fmt.Println("Aborting...")
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
	tracker := NewNumberTracker()
	// Print final statistics on exit
	defer tracker.PrintStatistics()
	// Global context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// When shuttingdown
	exit := make(chan os.Signal)
	defer close(exit)
	signal.Notify(exit, os.Interrupt, os.Kill)
	go gracefulShutdown(exit, cancel, listener)

	// **** Handler ****
	// For periodic printing of statistics (interval is defined as flag at entrance)
	ticker := time.Tick(time.Second * time.Duration(interval))
	// Coordination channels
	intInput := make(chan int)
	defer close(intInput)
	processChan := tracker.ProcessNumber(ctx, intInput)
	// Rate limitting
	rateLimiter := make(chan struct{}, maxconn)
	defer close(rateLimiter)
	// Writing to logfile
	logger.StreamWrite(ctx, processChan)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Termination found, exiting...")
			return
		default:
			// Check-in to the rateLimiter (this will block if the queue is full)
			// Will be cleaned out on exit
			rateLimiter <- struct{}{}
			// Accepting connections
			conn, err := listener.Accept()
			if err != nil {
				fmt.Printf("The server stopped accepting connections (%v) \n", err)
				return
			}
			// Handling connection
			go func() {
				defer conn.Close()
				// Releasing connection's place in the queue
				defer func() { <-rateLimiter }()
				scanner := bufio.NewScanner(conn)
				// Reading each client's input
				for scanner.Scan() {
					select {
					// Checking context per connection
					case <-ctx.Done():
						fmt.Printf("Closing connection: %v\n", ctx.Err())
						closeAndFreeResources(conn, listener)
						return
					// Print statistics every 10 (or interval value) seconds
					case <-ticker:
						tracker.PrintStatistics()
					default:
						input := scanner.Text()
						if checker.CheckTermination(input) {
							// Cancelling global context, connection and server
							cancel()
							closeAndFreeResources(conn, listener)
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
						} else {
							// This will close connection on exit
							// (see deferred at the beginning of the goroutine)
							return
						}
					}
				}
			}()
		}
	}
}

func closeAndFreeResources(conn net.Conn, listener net.Listener) {
	conn.Close()
	listener.Close()
}

// Closes app resources for cleaner shutdown
func gracefulShutdown(exit <-chan os.Signal, cancel context.CancelFunc, listener net.Listener) {
	<-exit
	fmt.Println("Received kill/intrrupt signal...")
	cancel()
	listener.Close()
}
