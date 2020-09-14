# New Relic Backend Test

The `number server` listed in this repository receives numbers and logs them into a log file 
(the path can be defined on start). A termination keyword can be set on start of the script. The server
will use that as a signal for a graceful shutdown attempt.

By default, the server will start in port `4000`, will write to the `./numbers.log` file,
will take numbers up to `999999999` and will recreate the log file per fresh restart.

The server is limited to take up to 5 concurrent connections (although, this can be changed on start, also).

The server will prompt statistics to STDOUT every 10 seconds (by default, this interval can be [changed also](#usage)):

Example output (every 10 seconds):
```
Starting number server. Welcome!
Received 492203 unique numbers, 1969333 duplicates (Total processed: 2461536). Unique totals: 492203
Received 640896 unique numbers, 2566315 duplicates (Total processed: 3207211). Unique totals: 1133099
Received 636761 unique numbers, 2551738 duplicates (Total processed: 3188499). Unique totals: 1769860
Received 638486 unique numbers, 2560362 duplicates (Total processed: 3198848). Unique totals: 2408346
Received 619657 unique numbers, 2487099 duplicates (Total processed: 3106756). Unique totals: 3028003

```

It's written in [Go](https://golang.org/).

## Requirements

`go version go1.14.2` or above.

## Installation

First, clone this repository:

`git clone git@github.com:mountolive/newrelictest.git`

Build the project:

`go build`

This will create the script that would start the server, it should have execution permissions:

`./newrelictest`, or `./newrelictest.exe` for windows.

## Usage

The script allows for some initial customization of the parameters which the server should use.
For a detailed list of the parameters, you can run `./newrelictest --help` (`./newrelictest.exe --help` on windows). This will prompt the following:

```
NAME:
   Number logger - Writes numbers to defined log file.
                   Numbers can have up to the max number
                   of digits defined by the user: 9 by default.
                   When the termination ("terminate") keyword is prompted,
                   the program will attempt to shutdown gracefully.
                   This termination keyword can be changed on start (see --help)

USAGE:
   newrelictest [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --port value, -p value         Port to be listened to (default: 4000)
   --append, -a                   Whether to append to existing log file or recreate on start
   --logfile value, -l value      Log file's path where the inputs would be written (default: "./numbers.log")
   --termination value, -t value  Terminate keyword, for shutting down the server (default: "terminate")
   --digits value, -d value       Max number of digits permitted for int input (max: 9) (default: 9)
   --interval value, -i value     Show statistics every * seconds (default: 10)
   --maxconn value, -c value      Max number of concurrent connections allowed (default: 5)
   --help, -h
```

## Testing

Tests can be executed with `go test` or, even better,  `go test --race` (this detects possible race conditions, [check here](https://golang.org/doc/articles/race_detector.html)). 

In terms of actual execution, the following client could be of help for testing the results of the script:

```
package main

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:4000")
	defer conn.Close()
	if err != nil {
		fmt.Printf("%v \n", err)
		return
	}
	createRandomNum := func() string {
		num := rand.Intn(999999999)
		stringRep := fmt.Sprintf("%d", num)
		stringLength := len(stringRep)
		leadingZeroes := 9 - stringLength
		if leadingZeroes > 0 {
			stringRep = strings.Repeat("0", leadingZeroes) + stringRep
		}
		return stringRep + "\n"
	}
	timeout := time.After(time.Second * 60)
	for {
		select {
		case <-timeout:
			fmt.Println("Exiting")
			return
		default:
			conn.Write([]byte(createRandomNum()))
		}
	}
}
```

## Main assumptions

- Each input from a client ends in a carriage character (new-line)
- On shutdown, it's possible that the last few messages received won't be written to the log file.
    (When the global context is canceled, the first routine receiving the signal will close the file.)

