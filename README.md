# elastic.go - A command-line tool to query the Elasticsearch REST API

`elastic.go` was born to turn commands that look like this:

```
curl -XGET 'http://localhost:9200/_nodes/_all/host,ip' | python -m json.tool
```

into this:

```
elastic node list
```

`elastic` fetches data from your Elasticsearch instance, formats the data
nicely when it is JSON compressed data, and adds a bit of colors to make it more
readable to the human eye. It aims at providing shortcuts for all default
Elasticsearch routes. For instance, you can get the cluster health with
`elastic cluster health` (or `elastic c he` for short).

Of course, you can still issue any `GET` requests with
`elastic query <YOUR REQUEST HERE>` (or `elastic q` for short), like
`elastic q twitter/tweet,user/_search?q=user:kimchy'`.
By design, only `GET` requests are allowed. I wanted to make it easy to query
Elasticseach indexes, not deleting them so use the good old `curl -XDELETE ...`
if this is what you want to achieve.
There is currently no support for `PUT` requests either since I have no use for
it. Pull requests are welcome however.

## Installation

Providing that [Go](https://golang.org) is installed and that `$GOPATH` is set,
simply use the following command:
```
go get -u github.com/Rolinh/elastic-go
```

Make sure that `$GOPATH/bin` is in your `$PATH`.

## Build a Docker image

A `Dockerfile` is provided. Simply run the following command to build the Docker
image:
```
docker build -t Rolinh/elastic-go .

## Usage

`elastic help` provides general help:
```
$ elastic help
NAME:
   elastic - A command line tool to query the Elasticsearch REST API

USAGE:
   elastic [global options] command [command options] [arguments...]

VERSION:
   1.0.0

AUTHOR(S):
   Robin Hahling <robin.hahling@gw-computing.net>

COMMANDS:
   cluster, c   Get cluster information
   index, i     Get index information
   node, n      Get cluster nodes information
   query, q     Perform any ES API GET query
   stats, s     Get statistics
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --baseurl "http://localhost:9200/"   Base API URL
   --help, -h                           show help
   --version, -v                        print the version
```

Help works for any subcommand as well. For instance:
```
$ elastic index help
NAME:
   elastic index - Get index information

USAGE:
   elastic index [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   docs-count, dc       Get index documents count
   list, l              List all indexes
   size, si             Get index size
   status, st           Get index status
   verbose, v           List indexes information with many stats
   help, h              Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h   show help

```
