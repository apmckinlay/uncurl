/*
uncurl is a program intended to be used by Suneido
to do http and ftp.
It accepts requests on stdin and returns results on stdout.
Requests always result in a newline terminated response
possibly followed by additional data.
Errors are output as a newline terminated ERROR...

get http://...
*/
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

func main() {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		request := input.Text()
		execute(request)
	}
}

func execute(request string) {
	cmd := splitOnFirst(request, " ")
	trace("got:", cmd[0], cmd[1])
	switch cmd[0] {
	case "get":
		getCommand(cmd[1])
	case "q":
		os.Exit(0)
	default:
		fmt.Println("ERROR invalid command: " + cmd[0])
	}
}

func getCommand(url string) {
	get := splitOnFirst(url, "://")
	switch get[0] {
	case "http":
		httpGet(url)
	default:
		fmt.Println("ERROR unsupported scheme: " + get[0])
	}
}

func splitOnFirst(s, sep string) []string {
	split := strings.SplitN(s, sep, 2)
	if len(split) == 1 {
		split = append(split, "")
	}
	return split
}

func httpGet(request string) {
	trace("HTTP", request)
	response, err := http.Get(request)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return
	}
	defer response.Body.Close()
	b, _ := httputil.DumpResponse(response, false)
	write("HEADER", b)

	var buf [8192]byte
	reader := response.Body
	for {
		n, err := reader.Read(buf[0:])
		trace("\n" + strings.Repeat("=", 70))
		if n > 0 {
			write("DATA", buf[0:n])
		}
		if err == io.EOF {
			fmt.Println("END")
			return
		}
		if err != nil {
			fmt.Println("ERROR", err.Error())
			return
		}
	}
}

func write(s string, b []byte) {
	n := len(b)
	fmt.Println(s, n)
	nw, err := os.Stdout.Write(b)
	if nw != n || err != nil {
		log.Fatal("http: error writing data")
	}
}

func trace(args ...interface{}) {
	//fmt.Println(args...)
}
