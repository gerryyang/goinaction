
# 1. Enable dependency tracking for your code 

```
go mod init example.com/hello
```

You will get `go.mod`

```
$ cat go.mod 
module example.com/hello

go 1.16
```

# 2. Create a file `hello.go`

``` go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

# 3. Run your code

```
$ go run .
Hello, World!
```

# 4. Call code in an external package

``` go
package main

import "fmt"

import "rsc.io/quote"

func main() {
    fmt.Println(quote.Go())
}
```

# 5. Add new module requirements and sums

```
go mod tidy
```

Go will add the quote module as a requirement, as well as a go.sum file for use in authenticating the module.

```
$ go help mod tidy
usage: go mod tidy [-e] [-v]

Tidy makes sure go.mod matches the source code in the module.
It adds any missing modules necessary to build the current module's
packages and dependencies, and it removes unused modules that
don't provide any relevant packages. It also adds any missing entries
to go.sum and removes any unnecessary ones.

The -v flag causes tidy to print information about removed modules
to standard error.

The -e flag causes tidy to attempt to proceed despite errors
encountered while loading packages.

See https://golang.org/ref/mod#go-mod-tidy for more about 'go mod tidy'.
```

# 6. Create a Go module

Go codes -> Packages -> Modules

```
mkdir greetings
cd greetings
go mod init example.com/greetings
```

Create a file `greetings.go`

``` go
package greetings

import "fmt"

// Hello returns a greeting for the named person.
func Hello(name string) string {
    // Return a greeting that embeds the name in a message.
    message := fmt.Sprintf("Hi, %v. Welcome!", name)
    return message
}
```

# 7. Call your module from another module

``` go
package main

import "fmt"

import "rsc.io/quote"

import "example.com/greetings"

func main() {
    fmt.Println("Hello, World!")
	fmt.Println(quote.Go());
	message := greetings.Hello("gerry")
	fmt.Println(message)
}
```

```
$ go mod edit -replace=example.com/greetings=./greetings
$ vi go.mod 
$ go mod tidy
go: found example.com/greetings in example.com/greetings v0.0.0-00010101000000-000000000000
$ cat go.mod 
module example.com/m

go 1.16

require (
        example.com/greetings v0.0.0-00010101000000-000000000000
        rsc.io/quote v1.5.2
)

replace example.com/greetings => ./greetings
$ go run .
Hello, World!
Don't communicate by sharing memory, share memory by communicating.
Hi, gerry. Welcome!
```

# 8. Return and handle an error

``` go
package main

import (
        "fmt"
        "log"
        "rsc.io/quote"
        "example.com/greetings"
)

func main() {

        // Set properties of the predefined Logger, including
        // the log entry prefix and a flag to disable printing
        // the time, source file, and line number.
        log.SetPrefix("greetings: ")
        log.SetFlags(0)

    fmt.Println("Hello, World!")

        fmt.Println(quote.Go());

        message, err := greetings.Hello("gerry")
        if err != nil {
                log.Fatal(err)
        }
        fmt.Println(message)

        message, err = greetings.Hello("")
        if err != nil {
                log.Fatal(err)
        }
}
```

# 9. Return a random greeting

``` go
package greetings

import (
        "fmt"
        "errors"
        "math/rand"
        "time"
)

// Hello returns a greeting for the named person.
func Hello(name string) (string, error) {
        if name == "" {
                return "", errors.New("empty name")
        }

    // Return a greeting that embeds the name in a message.
    message := fmt.Sprintf(randomFormat(), name)
    return message, nil
}

// Go executes init functions automatically at program startup, after global variables have been initialized. 
// Init sets initial values for variables used in the function
func init() {
        rand.Seed(time.Now().UnixNano())
}

// randomFormat returns one of a set of greeting messages. The returned
// message is selected at random.
func randomFormat() string {
    // A slice of message formats.
    formats := []string{
        "Hi, %v. Welcome!",
        "Great to see you, %v!",
        "Hail, %v! Well met!",
    }

    // Return a randomly selected message format by specifying
    // a random index for the slice of formats.
    return formats[rand.Intn(len(formats))]
}
```

# 10. Return greetings for multiple people

``` go
// Hellos returns a map that associates each of the named people
// with a greeting message.
func Hellos(names []string) (map[string]string, error) {
    // A map to associate names with messages.
    messages := make(map[string]string)
    // Loop through the received slice of names, calling
    // the Hello function to get a message for each name.
    for _, name := range names {
        message, err := Hello(name)
        if err != nil {
            return nil, err
        }
        // In the map, associate the retrieved message with
        // the name.
        messages[name] = message
    }
    return messages, nil
}
```

# 11. Add a test

Go's built-in support for unit testing makes it easier to test as you go. Specifically, using naming conventions, Go's testing package, and the go test command, you can quickly write and execute tests.

Ending a file's name with `_test.go` tells the go test command that this file contains test functions.

``` go
package greetings

import (
    "testing"
    "regexp"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHelloName(t *testing.T) {
    name := "Gladys"
    want := regexp.MustCompile(`\b`+name+`\b`)
    msg, err := Hello("Gladys")
    if !want.MatchString(msg) || err != nil {
        t.Fatalf(`Hello("Gladys") = %q, %v, want match for %#q, nil`, msg, err, want)
    }
}

// TestHelloEmpty calls greetings.Hello with an empty string,
// checking for an error.
func TestHelloEmpty(t *testing.T) {
    msg, err := Hello("")
    if msg != "" || err == nil {
        t.Fatalf(`Hello("") = %q, %v, want "", error`, msg, err)
    }
}
```

```
$ go test
PASS
ok      example.com/greetings   0.002s
$ go test -v
=== RUN   TestHelloName
--- PASS: TestHelloName (0.00s)
=== RUN   TestHelloEmpty
--- PASS: TestHelloEmpty (0.00s)
PASS
ok      example.com/greetings   0.002s
```

Break the greetings.Hello function to view a failing test.

``` go
// Hello returns a greeting for the named person.
func Hello(name string) (string, error) {
    // If no name was given, return an error with a message.
    if name == "" {
        return name, errors.New("empty name")
    }
    // Create a message using a random format.
    // message := fmt.Sprintf(randomFormat(), name)
    message := fmt.Sprint(randomFormat())
    return message, nil
}
```

```
$ go test
--- FAIL: TestHelloName (0.00s)
    greetings_test.go:15: Hello("Gladys") = "Hi, %!v(MISSING). Welcome!", <nil>, want match for `\bGladys\b`, nil
FAIL
exit status 1
FAIL    example.com/greetings   0.002s
```

# 11. Compile and install the application

While the `go run` command is a useful shortcut for compiling and running a program when you're making frequent changes, it doesn't generate a binary executable.

* The `go build` command compiles the packages, along with their dependencies, but it doesn't install the results.
* The `go install` command compiles and installs the packages.
* You can discover the install path by running the `go list` command.
* As an alternative, if you already have a directory like `$HOME/bin` in your shell path and you'd like to install your Go programs there, you can change the install target by setting the GOBIN variable using the `go env` command


```
$ go list -f '{{.Target}}'
/home/ubuntu/golang/workspace/bin/m
$ go env -w GOBIN=/home/ubuntu/bin
$ go list -f '{{.Target}}'
/home/ubuntu/bin/m
```

To unset a variable previously set by go env -w, use go env -u:

```
$ go env -u GOBIN
```

For added convenience, we'll add the install directory to our PATH to make running binaries easy:

```
$ export PATH=$PATH:$(dirname $(go list -f '{{.Target}}' .))
```


# 12. Managing dependencies

`go.mod` file reference

Each Go module is defined by a `go.mod` file that describes the module's properties, including its dependencies on other modules and on versions of Go.

* The current module's module path. This should be a location from which the module can be down loaded by Go tools, such as the module code's repository location. This serves as a unique identifier, when combined with the module's version number. It is also the prefix of the package path for all packages in the module. For more about how Go locates the module, see the Go Modules Reference.
* The minimum version of Go required by the current module.
* A list of minimum versions of other modules required by the current module.
* Instructions, optionally, to replace a required module with another module version or a local directory, or to exclude a specific version of a required module.

```
$ go mod init example.com/mymodule
```

```
$ cat go.mod 
module example.com/m

go 1.16

require (
        example.com/greetings v0.0.0-00010101000000-000000000000
        rsc.io/quote v1.5.2
)

replace example.com/greetings => ./greetings
```




# Refer

* https://golang.org/doc/tutorial/getting-started
* https://golang.org/ref/mod#go-mod-init
* https://pkg.go.dev/std
* https://pkg.go.dev
* https://pkg.go.dev/rsc.io/quote
* https://golang.org/ref/mod#authenticating
* https://golang.org/ref/mod#go-mod-edit
* https://blog.golang.org/maps
* https://golang.org/doc/effective_go#blank
* https://golang.org/doc/modules/gomod-ref

* https://golang.org/ref/mod
