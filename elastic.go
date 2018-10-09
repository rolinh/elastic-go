//usr/bin/env go run $0 $@; exit
//
// Copyright 2015 The elastic.go authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Author: Robin Hahling <robin.hahling@gw-computing.net>

// elastic.go is a command line tool to query the Elasticsearch REST API.
package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/gilliek/go-xterm256/xterm256"
	"github.com/hokaccha/go-prettyjson"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "elastic"
	app.Usage = "A command line tool to query the Elasticsearch REST API"
	app.Version = "1.0.1"
	app.Author = "Robin Hahling"
	app.Email = "robin.hahling@gw-computing.net"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "baseurl",
			Value: "http://localhost:9200/",
			Usage: "Base API URL",
		},
		cli.BoolFlag{
			Name:  "trace",
			Usage: "Trace URLs called",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "cluster",
			ShortName: "c",
			Usage:     "Get cluster information ",
			Subcommands: []cli.Command{
				{
					Name:      "health",
					ShortName: "he",
					Usage:     "Get cluster health",
					Action: func(c *cli.Context) {
						out, err := getJSON(cmdCluster(c, "health"), c)
						if err != nil {
							fatal(err)
						}
						fmt.Println(out)
					},
				},
				{
					Name:      "state",
					ShortName: "s",
					Usage:     "Get cluster state (allows filter args)",
					Action: func(c *cli.Context) {
						out, err := getJSON(cmdCluster(c, "state"), c)
						if err != nil {
							fatal(err)
						}
						fmt.Println(out)
					},
				},
				{
					Name:      "stats",
					ShortName: "t",
					Usage:     "Get cluster stats",
					Action: func(c *cli.Context) {
						out, err := getJSON(cmdCluster(c, "stats"), c)
						if err != nil {
							fatal(err)
						}
						fmt.Println(out)
					},
				},
			},
		},
		{
			Name:      "index",
			ShortName: "i",
			Usage:     "Get index information",
			Subcommands: []cli.Command{
				{
					Name:      "docs-count",
					ShortName: "dc",
					Usage:     "Get index documents count",
					Action: func(c *cli.Context) {
						list, err := getRaw(cmdIndex(c, "list"), c)
						if err != nil {
							fatal(err)
						}
						for _, idx := range filteredDocsCountIndexes(list) {
							fmt.Println(idx)
						}
					},
				},
				{
					Name:      "list",
					ShortName: "l",
					Usage:     "List all indexes",
					Action: func(c *cli.Context) {
						list, err := getRaw(cmdIndex(c, "list"), c)
						if err != nil {
							fatal(err)
						}
						for _, idx := range filteredListIndexes(list) {
							fmt.Println(idx)
						}
					},
				},
				{
					Name:      "size",
					ShortName: "si",
					Usage:     "Get index size",
					Action: func(c *cli.Context) {
						list, err := getRaw(cmdIndex(c, "list"), c)
						if err != nil {
							fatal(err)
						}
						for _, idx := range filteredSizeIndexes(list) {
							fmt.Println(idx)
						}
					},
				},
				{
					Name:      "status",
					ShortName: "st",
					Usage:     "Get index status",
					Action: func(c *cli.Context) {
						list, err := getRaw(cmdIndex(c, "list"), c)
						if err != nil {
							fatal(err)
						}
						for _, idx := range filteredStatusIndexes(list) {
							fmt.Println(idx)
						}
					},
				},
				{
					Name:      "verbose",
					ShortName: "v",
					Usage:     "List indexes information with many stats",
					Action: func(c *cli.Context) {
						list, err := getRaw(cmdIndex(c, "list"), c)
						if err != nil {
							fatal(err)
						}
						fmt.Println(list)
					},
				},
			},
		},
		{
			Name:      "node",
			ShortName: "n",
			Usage:     "Get cluster nodes information",
			Subcommands: []cli.Command{
				{
					Name:      "list",
					ShortName: "l",
					Usage:     "List nodes information",
					Action: func(c *cli.Context) {
						out, err := getJSON(cmdNode(c, "list"), c)
						if err != nil {
							fatal(err)
						}
						fmt.Println(out)
					},
				},
				{
					Name:      "stats",
					ShortName: "s",
					Usage:     "List node stats (allows filter args)",
					Action: func(c *cli.Context) {
						out, err := getJSON(cmdNode(c, "stats"), c)
						if err != nil {
							fatal(err)
						}
						fmt.Println(out)
					},
				},
			},
		},
		{
			Name:      "query",
			ShortName: "q",
			Usage:     "Perform any ES API GET query",
			Action: func(c *cli.Context) {
				var out string
				var err error
				if strings.Contains(c.Args().First(), "_cat/") {
					out, err = getRaw(cmdQuery(c), c)
				} else {
					out, err = getJSON(cmdQuery(c), c)
				}
				if err != nil {
					fatal(err)
				}
				fmt.Println(out)
			},
		},
		{
			Name:      "stats",
			ShortName: "s",
			Usage:     "Get statistics",
			Subcommands: []cli.Command{
				{
					Name:      "size",
					ShortName: "s",
					Usage:     "Get index sizes",
					Action: func(c *cli.Context) {
						out, err := getJSON(cmdStats(c, "size"), c)
						if err != nil {
							fatal(err)
						}
						fmt.Println(out)
					},
				},
			},
		},
	}

	app.Run(os.Args)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func getJSON(route string, c *cli.Context) (string, error) {
	r, err := httpGet(route, isTraceEnabled(c))
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status code: %s", r.Status)
	}

	mediatype, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return "", err
	}
	if mediatype == "" {
		return "", errors.New("mediatype not set")
	}
	if mediatype != "application/json" {
		return "", fmt.Errorf("mediatype is '%s', 'application/json' expected", mediatype)
	}

	var b interface{}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		return "", err
	}
	out, err := prettyjson.Marshal(b)
	return string(out), err
}

func getRaw(route string, c *cli.Context) (string, error) {
	r, err := httpGet(route, isTraceEnabled(c))

	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status code: %s", r.Status)
	}

	out, err := ioutil.ReadAll(r.Body)
	return string(out), err
}

// processing functions
func filteredDocsCountIndexes(list string) []string {
	var out []string
	scanner := bufio.NewScanner(strings.NewReader(list))
	for scanner.Scan() {
		elmts := strings.Fields(scanner.Text())
		if len(elmts) < 6 {
			continue
		}
		out = append(out, fmt.Sprintf("%10s %s", colorizeStatus(elmts[5]), elmts[2]))
	}
	return out
}

func filteredListIndexes(list string) []string {
	var out []string
	scanner := bufio.NewScanner(strings.NewReader(list))
	for scanner.Scan() {
		elmts := strings.Fields(scanner.Text())
		if len(elmts) < 3 {
			continue
		}
		out = append(out, elmts[2])
	}
	return out
}

func filteredStatusIndexes(list string) []string {
	var out []string
	scanner := bufio.NewScanner(strings.NewReader(list))
	for scanner.Scan() {
		elmts := strings.Fields(scanner.Text())
		if len(elmts) < 3 {
			continue
		}
		out = append(out, fmt.Sprintf("%22s %s", colorizeStatus(elmts[0]), elmts[2]))
	}
	return out
}

func filteredSizeIndexes(list string) []string {
	var out []string
	scanner := bufio.NewScanner(strings.NewReader(list))
	for scanner.Scan() {
		elmts := strings.Fields(scanner.Text())
		if len(elmts) < 8 {
			continue
		}
		out = append(out, fmt.Sprintf("%10s %s", elmts[7], elmts[2]))
	}
	return out
}

func colorizeStatus(status string) string {
	var color xterm256.Color
	switch status {
	case "red":
		color = xterm256.Red
	case "green":
		color = xterm256.Green
	case "yellow":
		color = xterm256.Yellow
	default:
		return status
	}
	return xterm256.Sprint(color, status)
}

// command-line commands from now on
func cmdCluster(c *cli.Context, subCmd string) string {
	route := "_cluster/"
	url := c.GlobalString("baseurl")

	var arg string
	switch subCmd {
	case "health":
		arg = "health"
	case "state":
		arg = "state/" + strings.Join(c.Args(), ",")
	case "stats":
		arg = "stats/"
	default:
		arg = ""
	}
	return url + route + arg
}

func cmdIndex(c *cli.Context, subCmd string) string {
	var route string
	url := c.GlobalString("baseurl")
	switch subCmd {
	case "list":
		route = "_cat/indices?v"
	default:
		route = ""
	}
	return url + route
}

func cmdNode(c *cli.Context, subCmd string) string {
	var route string
	url := c.GlobalString("baseurl")
	switch subCmd {
	case "list":
		route = "_nodes/_all/host,ip"
	case "stats":
		route = "_nodes/_all/stats/" + strings.Join(c.Args(), ",")
	default:
		route = ""
	}
	return url + route
}

func cmdQuery(c *cli.Context) string {
	route := c.Args().First()
	url := c.GlobalString("baseurl")
	return url + route
}

func cmdStats(c *cli.Context, subCmd string) string {
	var route string
	url := c.GlobalString("baseurl")
	switch subCmd {
	case "size":
		route = "_stats/index,store"
	default:
		route = ""
	}
	return url + route
}

func httpGet(route string, trace bool) (*http.Response, error) {
	if trace {
		fmt.Fprintf(os.Stderr, "GET: %s", route)
	}
	r, err := http.Get(route)

	return r, err
}

func isTraceEnabled(c *cli.Context) bool {
	return c.GlobalBool("trace")
}
