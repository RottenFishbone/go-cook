## go-cook &nbsp;&nbsp; [![Go Report Card](https://goreportcard.com/badge/git.sr.ht/~rottenfishbone/go-cook)](https://goreportcard.com/report/git.sr.ht/~rottenfishbone/go-cook)
This project aims to be a reimplementation, and eventually extension, of the original 
[CooklangCLI](https://github.com/cooklang/CookCLI), written in Go.  

The root package provides the parser and the requisite types under the package name `cook`. 
Thus, despite the scope of the project, importing the root will provide a 
lightweight cook parsing library.

-------

#### Current State
At present, the parser is almost stable. It has all canonical features, however 
it will need to be extended in the future to include photos and possibly comments.
Canonical tests [here](https://github.com/cooklang/spec/tree/main/tests) are generated
using `internal/cmd/test_gen` to pull the latest tests and generate test code.

The CLI provides basic recipe reading as well as access to a webserver and an API server.

(Implemented) commands are as follows
```
  help        Help about any command
  init        Creates the default config file.
  read        Parses a recipe file and pretty prints it to stdout
  server      Hosts a local webserver to view/manage recipes.
```

Most focus so far has been on making a usable web interface. While non-final, it
does provide a usable and pleasant interface already.

This functionality will be extended to the CLI once necessary components are finished.

### Compilation 
To use the parsing library, you can simply run `go build`, 
or import the package in a module. 

To use the command line utility *with a server* it is required that you build the 
[Svelte](https://svelte.dev/) webapp *prior* to compiling the binary. 
This requires `npm` is installed.

The easiest way is to simply run `make` in the project root which will build the 
webapp and then compile the binary (embedding the webapp into the binary).

Tests can be regenerated/updated using `make canonical`

If you don't want to use make, or don't have access, you can view the Makefile 
to examine the build process for various tasks (build, test gen, formatting etc.).

To build the CLI binary without the webapp, simply run `go build` in the `cmd/cook/`
directory.

----------

### Author
Jayden Dumouchel -- jdumouch@ualberta.ca | rottenfishbone@pm.me

### License
This project is licensed under the MIT License, see `LICENSE` file for details.
