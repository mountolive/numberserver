# New Relic Backend Test

The `number server` listed in this repository receives numbers and logs them into a log file 
(the path can be defined on start). A termination keyword can be set on start of the script. The server
will use that as a signal for a graceful shutdown attempt.

By default, the server will start in port `4000`, will write to the `./numbers.log` file,
will take numbers up to `999999999` and will recreate the log file per fresh restart.

It's written in [Go](https://golang.org/).

## Requirements

`go version go1.14.2` or above.

## Installation

First, clone this repository:

`git clone git@github.com:mountolive/newrelictest.git`

Build the project:

`go build`

This will create the script that would start the server, it should have execution permissions:

`./newrelictest`

## Usage

The script allows for some initial customization of the parameters which the server should use.
For a detailed list of the parameters, you can run `./newrelictest --help`. This will prompt the following:

```
NAME:
   Number logger - Writes numbers to defined log file.
               Numbers can have up to the max number
               of digits defined by the user: 9 by default.
               When the termination ("terminate") keyword is prompted,
               the program will attempt to shutdown gracefully

USAGE:
   newrelictest [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --port value, -p value         Port to be listened to (default: 4000)
   --append, -a                   Whether to append to existing log file or recreate on start
   --logfile value, -l value      Log file's path where the inputs would be written (default: "./numbers.log")
   --termination value, -t value  Terminate keyword, for shutting down the server (default: "terminate")
   --digits value, -d value       Max number of digits permitted for int input (default: 9)
   --interval value, -i value     Show statistics every * seconds (default: 10)
   --help, -h                     show help
```

## Testing

Test can be run with `go test` or, even better,  `go test --race` (this detects possible race conditions [here](https://golang.org/doc/articles/race_detector.html)). 

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

