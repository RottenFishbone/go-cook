### GoCook -- webapp
This directory holds the subproject that powers GoCook's `server` command.

#### Compilation
GoCook's webapp is built using svelte with vite.
In this directory, run the following in a terminal:

```
npm install
npm run build
```

This should place the compiled website into `./dist`, which will be embedded
into GoCook on compilation of the `cook` binary. As such, to make updates to the
binary you will be required to recompile.

Eventually I will include this as part of a build script.

#### License
This webapp shares the license of the larger `go-cook` repository.

