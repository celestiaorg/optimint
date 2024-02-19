## rollkit

A modular framework for rollups, with an ABCI-compatible client interface.

### Synopsis


Rollkit is a modular framework for rollups, with an ABCI-compatible client interface.
The rollkit-cli uses the environment variable "RKHOME" to point to a file path where the node keys, config, and data will be stored. 
If a path is not specified for RKHOME, the rollkit command will create a folder "~/.rollkit" where it will store said data.


### Options

```
  -h, --help               help for rollkit
      --home string        directory for config and data (default "HOME/.rollkit")
      --log_level string   set the log level; default is info. other options include debug, info, error, none (default "info")
      --trace              print out full stack trace on errors
```

### SEE ALSO

* [rollkit completion](rollkit_completion.md)	 - Generate the autocompletion script for the specified shell
* [rollkit docs-gen](rollkit_docs-gen.md)	 - Generate documentation for rollkit CLI
* [rollkit start](rollkit_start.md)	 - Run the rollkit node
* [rollkit version](rollkit_version.md)	 - Show version info
