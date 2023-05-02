## go-cook &nbsp;&nbsp; [![Go Report Card](https://goreportcard.com/badge/git.sr.ht/~rottenfishbone/go-cook)](https://goreportcard.com/report/git.sr.ht/~rottenfishbone/go-cook)
This project aims to be a reimplementation, and eventually extension, of the original [CooklangCLI](https://github.com/cooklang/CookCLI), written in Go.  

The root package provides the parser and the requisite types under the package name `cook`. 

I will be maintaing a project structure that allows for module imports from `pkg.dev.go` to enable 
Cooklang utilities such as parsing, printing and creation, as a library.

### Motivation
The concept of cooklang has intrigued me and everyone I've described it to. However, I found the 
existing tooling fell short in some aspects (for my use case) and, additionally, was written in 
Swift; which is described as the creator as a cross-platform-misstep. As such, I figured a 
great launch point for learning Golang was to reimplement it in a way that would better suit my 
idea for what it can be (thinking lightweight Grocy). With that said, I hope to create a modular, 
easy to use system that is agnostic to the existing cooklang ecosystem.

-------

#### Current State
At present, the parser works and passes all canonical tests listed 
[here](https://github.com/cooklang/spec/tree/fa9bc51515b3317da434cb2b5a4a6ac12257e60b/tests). 

The CLI provides a read function to output stored recipes to terminal.

##### In progress
 - CLI
    - [x] `cook init`
    - [x] `cook read`
    - [ ] `cook server`
    - [ ] `cook shopping list`

##### Planned
 - Language Extensions (Images, Shopping List, Parse comments)
 - Shopping List Parser 

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
