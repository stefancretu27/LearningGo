package go_didactical_apps


/*
* 1. Go toolchain:
* 	-> Is a set of tools (executables + the standard Go library) provided by the Go programming language to support the entire development lifecycle — 
* from writing and compiling code to testing, debugging, and managing dependencies.
*	-> Core components
*		1.1 go command - The central command-line tool that orchestrates most tasks.
*					   - Subcommands include:
*							+ go build: Compiles packages and dependencies.
*							+ go run: Compiles and runs Go programs.
*							+ go test: Runs unit tests.
*							+ go mod: Manages modules and dependencies.
*							+ ge get: Downloads or updates a module, making its packages available for import
*							+ go install: Installs binaries.
*							+ go clean: Removes build artifacts.
*							+ go fmt: Formats code according to Go standards.
*							+ go vet: Examines code for suspicious constructs.
* 		1.2 Compiler (gc) - Converts Go source code into machine code.
*					 	  - Optimized for speed and simplicity.
*						  - Supports SSA (Static Single Assignment) form for optimization.
*		1.3 Linker (ld) - Combines compiled packages into a single executable.
* 						- Handles symbol resolution and binary layout.
*		1.4 Assembler (asm) - Converts Go assembly code into machine code.
*							- Used for performance-critical parts of the standard library.
*		1.5 Debugger (delve) - Not part of the standard toolchain but widely used.
*							 - Helps with stepping through code, inspecting variables, and setting breakpoints.
*		1.6 Formatter (gofmt) - Automatically formats Go code to follow the language’s style guide.
*							  - Ensures consistency across teams and projects.
*		1.7 Linter (golint, staticcheck, etc.) - Checks for stylistic issues and potential bugs.
*											   - go vet is included in the toolchain and performs static analysis.
*		1.8 Package Manager (go mod) - Introduced in Go 1.11.
*									 - Handles versioning and dependency resolution using modules.
*/

/* 
* 2. Go package manager:
* 	-> Go toolchain contains a package management utilitary, built into the language.
*
* 	-> In detail, unlike python or rust, Go has no separate package manager. Instead, the Go toolchain itself 
* (the go command) includes commands which allow to manage the packages (external libs in Go):
*	- download modules - go get (go get github.com/some/module@v1.2.3)
*	- track dependencies - go.mod, go.sum files
*	- clean up unused packages - go mod tidy
*	- vendor dependencies - go mod vendor
* Example workflow:
*		- go mod init myapp         		# Start a module
*		- go get github.com/gin-gonic/gin   # Download a module in vendor directory
*		- go build                  		# Build the app
*/

/* 3. Go modules:
* 	-> A Go module is a folder containing all Go code files, eventually grouped as packages, with a go.mod file at its root. This file defines the module’s path and its dependencies,
* using "require" keyword and the list of dependencies, as github links to the repos
*	Example: 
*		require (
*			github.com/aws/aws-sdk-go-v2 v1.20.0
*			github.com/gorilla/mux v1.8.0
*		)
*
* 	-> Key files:
* 		- go.mod: Lists the module path and required dependencies.
*		- go.sum: Contains checksums for module versions to verify integrity.
*
* 	-> Key characteristics:
* 		- Defined by a go.mod file at the root of the project.
*		- Specifies the module path (usually a repo URL) and its dependencies (eg: module github.com/stefancretu/myapp)
*		- Can contain multiple packages (subfolders with .go files).
* Each newly created project must define a module.
*
* 	-> How go modules and package management work together?
* 	When you run:
* 		- go build: Go reads go.mod, downloads needed packages, and builds the app.
* 		- go mod tidy: Cleans up unused or missing dependencies.
*		- go mod vendor: Copies dependencies into a vendor/ folder for offline or controlled builds.
*
* 	-> In Go, vendor refers to a directory used to store local copies of dependencies (external packages) that your project relies on. 
* It's part of Go's module and dependency management system.
*/

/* 
* 4. Go packages
* 	-> In Go, a package is a collection of .go files kept in the same directory that share the same package name.
*
*	-> almost like a namespace in C++
*
*	-> package level compilation, not file-level compilation: each package is compiled separately, into intermmediate object file, representing a translation unit. 
*	Finally, all packages are linked together.
*
* 	-> If the app requires the usage of a certain package, the package manager will download the entire module of that package,
* but will compile and use only that package. So even though you only use one package, Go fetches the whole module because:
*		- Modules are the unit of versioning
* 		- Packages are the unit of code organization
*
* 	-> Key points:
* 		- Every .go file starts with a package declaration on the first line. (package package_name)
*		- A package can be imported by other packages.
*		- Packages help organize code into reusable components.
*		- a module cannot be imported, only packages, so a whole module's functionality is imported package by package
* 	
*	-> Visibility rules:
* 		- Exported identifiers have names (interfaces, structs, methods, functions etc) starting with capital letters, thus accessible in between packages of the same module,
*	but requires using import of the packages (path starting at module's level)
*		- Unexported identifiers (starting with a lowercase letter) are only accessible within the same package.
*/

func AboutGoPackageManagement() {

}