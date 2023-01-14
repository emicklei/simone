# simone

[![GoDoc](https://godoc.org/github.com/emicklei/simone?status.svg)](https://pkg.go.dev/github.com/emicklei/simone?tab=doc)

New kind of tool.

## usage

Simone is a package to build a Javascript engine with plugins in Go.
By implementing exposed methods on plugins, functions are added to the runtime.
See `examples` for a minimal program.

## run flags

    -c  start a client REPL
    -h  show help
    -i  string
        run the script from filename as input
    -s  string
        run script from filename on startup
    -v  verbose logging

## optional plugins

| plugin | description
|-|-
| fs | access to the local filesystem

## builtin functions

    log(arg,arg,...);    // write log line
    include('lib.sim');  // evaluate script from file

## usage REPL

Commands start with the colon `:` prefix.

| command | description
|-|-
|:q| quit the tool
|:h| show help
|:d| toggle verbose logging
|:v| show list of global variables
|:p| show list of available plugin names

If you postfix a variable (or plugin) with a `?` then it will print all available functions.
If you postfix an entry with a `!` then it browse the result in a Web page.

## usage HTTP

    POST | GET http://localhost:9119/v1?action=<ACTION>&source=<ENTRY>&file=<FILE>&line=<LINE>

For POST the source is transported with the request body
For GET the source can be passed as a query parameter.