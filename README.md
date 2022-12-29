## simone

New kind of tool.

## usage 

Simone is a package to build a Javascript engine with plugins in Go.
See `examples/http` for a minimal program. 

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