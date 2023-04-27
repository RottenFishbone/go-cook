## go-cook &nbsp;&nbsp; [![Go Report Card](https://goreportcard.com/badge/git.sr.ht/~rottenfishbone/go-cook)](https://goreportcard.com/report/git.sr.ht/~rottenfishbone/go-cook)
This project aims to be a reimplementation, and eventually extension, of the original [CooklangCLI](https://github.com/cooklang/CookCLI), written in Go.  

The root package provides the parser and the requisite types under the package name `cooklang`. 

I will be maintaing a project structure that allows for module imports from `pkg.dev.go` to enable Cooklang utilities such as parsing, printing and creation, as a library.

##### Current State
At present, the parser works and passes all canonical tests listed [here](https://github.com/cooklang/spec/tree/fa9bc51515b3317da434cb2b5a4a6ac12257e60b/tests).

-------
##### in progress
 - CLI
    - [x] `cook init`
    - [x] `cook read`
    - [ ] `cook shopping list`
    - [ ] `cook server`

##### todo
 - Language Extensions (Images, Shopping List, Parse comments)
 - Shopping List Parser 

--------

### Author
Jayden Dumouchel -- jdumouch@ualberta.ca | rottenfishbone@pm.me
