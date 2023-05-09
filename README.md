## go-cook &nbsp;&nbsp; [![Go Report Card](https://goreportcard.com/badge/git.sr.ht/~rottenfishbone/go-cook)](https://goreportcard.com/report/git.sr.ht/~rottenfishbone/go-cook)
This project aims to be a reimplementation, and eventually extension, of the original [CooklangCLI](https://github.com/cooklang/CookCLI), written in Go.  

The root package provides the parser and the requisite types under the package name `cook`. 
Thus, despite the scope of the project, importing the root will provide a lightweight cook parsing library.

-------

#### Current State
At present, the parser works and passes all canonical tests listed 
[here](https://github.com/cooklang/spec/tree/fa9bc51515b3317da434cb2b5a4a6ac12257e60b/tests). 

The CLI provides basic recipe reading as well as access to a webserver and an API server.

(Implemented) commands are as follows
```
  help        Help about any command
  init        Creates the default config file.
  read        Parses a recipe file and pretty prints it to stdout
  server      Hosts a local webserver to view/manage recipes.
```

Most focus so far has been on making a usable web interface. While unfinished, it presently
can perform any action required to manage or view a recipe directory.

This functionality will be extended to the CLI once necessary components are finished.

### Compilation 
To use the libraries, you can simply run `go build` as needed. 

To use the command line utility *with a server* it is required that you build the 
[Svelte](https://svelte.dev/) webapp *prior* to compiling the binary. This requires `npm` is installed.

The easiest way is to simply run `make` in the project root which will build the webapp and then 
compile the binary.

Alternatively (for unix-based OS):
```
cd internal/web/
npm install
npm run build
cd ../../
cd cmd/cook/
go build
```
Will produce a usable binary in `cmd/cook`.

----------

### Author
Jayden Dumouchel -- jdumouch@ualberta.ca | rottenfishbone@pm.me

### License
This project is licensed under the MIT License, see `LICENSE` file for details.
