
# Using Go Modules

`Go 1.11` adds preliminary support for a new concept called `“modules”`, an alternative to `GOPATH` with integrated support for versioning and package distribution. Using modules, developers are no longer confined to working inside GOPATH, version dependency information is explicit yet lightweight, and builds are more reliable and reproducible.

`Modules` are how Go manages dependencies.

A module is a collection of packages that are released, versioned, and distributed together. Modules may be downloaded directly from version control repositories or from module proxy servers.

For a series of tutorials on modules, see https://golang.org/doc/tutorial/create-module.

For a detailed reference on modules, see https://golang.org/ref/mod.

By default, the go command may download modules from https://proxy.golang.org. It may authenticate modules using the checksum database at https://sum.golang.org. Both services are operated by the Go team at Google. The privacy policies for these services are available at https://proxy.golang.org/privacy and https://sum.golang.org/privacy, respectively.

The go command's download behavior may be configured using `GOPROXY`, `GOSUMDB`, `GOPRIVATE`, and other environment variables. See 'go help environment' and https://golang.org/ref/mod#private-module-privacy for more information.


A module is a collection of `Go packages` stored in a file tree with a `go.mod` file at its root. The `go.mod` file defines the module’s module path, which is also the import path used for the root directory, and its dependency requirements, which are the other modules needed for a successful build. Each dependency requirement is written as a module path and a specific semantic version.

As of `Go 1.11`, the go command enables the use of modules when the current directory or any parent directory has a `go.mod`, provided the directory is outside `$GOPATH/src`. (Inside `$GOPATH/src`, for compatibility, the go command still runs in the old `GOPATH` mode, even if a `go.mod` is found. See the [go command documentation](https://golang.org/cmd/go/#hdr-Preliminary_module_support) for details.) **Starting in `Go 1.13`, module mode will be the default for all development**.


## Creating a new module

Create a new, empty directory somewhere outside $GOPATH/src, cd into that directory, and then create a new source file, `hello.go`:

``` go
package hello

func Hello() string {
    return "Hello, world."
}
```

Let's write a test, too, in `hello_test.go`:

``` go
package hello

import "testing"

func TestHello(t *testing.T) {
    want := "Hello, world."
    if got := Hello(); got != want {
        t.Errorf("Hello() = %q, want %q", got, want)
    }
}
```

At this point, the directory contains a package, but not a module, because there is no `go.mod` file. 

Let's make the current directory the root of a module by using `go mod init` and then try `go test`:

```
$ go mod init github.com/gerryyang/goinaction/module/hello  
go: creating new go.mod: module github.com/gerryyang/goinaction/module/hello
go: to add module requirements and sums:
        go mod tidy
$ ls
go.mod  hello.go  hello_test.go
$ cat go.mod 
module github.com/gerryyang/goinaction/module/hello

go 1.16
$ go test
PASS
ok      github.com/gerryyang/goinaction/module/hello    0.006s
```

Congratulations! You’ve written and tested your first module.

The `go.mod` file only appears in the root of the module. Packages in subdirectories have import paths consisting of the module path plus the path to the subdirectory. For example, if we created a subdirectory `world`, we would not need to (nor want to) run `go mod init` there. The package would automatically be recognized as part of the `github.com/gerryyang/goinaction/module/hello` module, with import path `github.com/gerryyang/goinaction/module/hello/world`.

```
$ go test
PASS
ok      github.com/gerryyang/goinaction/module/hello/world      0.002s
```

## Adding a dependency

The primary motivation for Go modules was to improve the experience of using (that is, adding a dependency on) code written by other developers.

Let's update our `hello.go` to import `rsc.io/quote` and use it to implement Hello:

``` go
package hello

import "rsc.io/quote"

func Hello() string {
    return quote.Hello()
}

```


```
$ go test
hello.go:3:8: no required module provides package rsc.io/quote; to add it:
        go get rsc.io/quote
$ go get rsc.io/quote
go get: added rsc.io/quote v1.5.2
$ ls
go.mod  go.sum  hello.go  hello_test.go  world
$ cat go.mod 
module github.com/gerryyang/goinaction/module/hello

go 1.16

require rsc.io/quote v1.5.2 // indirect
ubuntu@VM-0-16-ubuntu:~/github/goinaction/module/hello$ cat go.sum 
golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c h1:qgOY6WgZOaTkIIMiVjBQcw93ERBE4m30iBm00nkL0i8=
golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
rsc.io/quote v1.5.2 h1:w5fcysjrx7yqtD/aO+QwRjYZOKnaM9Uh2b40tElTs3Y=
rsc.io/quote v1.5.2/go.mod h1:LzX7hefJvL54yjefDEDHNONDjII0t9xZLPXsUe+TKr0=
rsc.io/sampler v1.3.0 h1:7uVkIFmeBqHfdjD+gZwtXXI+RODJ2Wc4O7MPEh/QiW4=
rsc.io/sampler v1.3.0/go.mod h1:T1hPZKmBbMNahiBKFy5HrXp6adAjACjK9JXDnKaTXpA=
$ go test
PASS
ok      github.com/gerryyang/goinaction/module/hello    0.003s
```

The go command resolves imports by using the specific dependency module versions listed in `go.mod`. When it encounters an `import` of a package not provided by any module in `go.mod`, the go command automatically looks up the module containing that package and adds it to go.mod, using the latest version.

A second `go test` command will not repeat this work, since the `go.mod` is now up-to-date and the downloaded modules are cached locally (in `$GOPATH/pkg/mod`):

Note that while the go command makes adding a new dependency quick and easy, it is not without cost. Your module now literally depends on the new dependency in critical areas such as correctness, security, and proper licensing, just to name a few. For more considerations, see [Russ Cox's blog post, “Our Software Dependency Problem.”](https://research.swtch.com/deps)

As we saw above, adding one direct dependency often brings in other indirect dependencies too. The command `go list -m all lists` the current module and all its dependencies:

In the go list output, the current module, also known as the main module, is always the first line, followed by dependencies sorted by module path.

```
$ go list -m all
github.com/gerryyang/goinaction/module/hello
golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c
rsc.io/quote v1.5.2
rsc.io/sampler v1.3.0
```

The `golang.org/x/text` version `v0.0.0-20170915032832-14c0d48ead0c` is an example of a [pseudo-version](https://golang.org/cmd/go/#hdr-Pseudo_versions), which is the go command's version syntax for a specific untagged commit.

In addition to `go.mod`, the go command maintains a file named `go.sum` containing the expected [cryptographic hashes](https://golang.org/cmd/go/#hdr-Module_downloading_and_verification) of the content of specific module versions:

```
$ cat go.sum 
golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c h1:qgOY6WgZOaTkIIMiVjBQcw93ERBE4m30iBm00nkL0i8=
golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
rsc.io/quote v1.5.2 h1:w5fcysjrx7yqtD/aO+QwRjYZOKnaM9Uh2b40tElTs3Y=
rsc.io/quote v1.5.2/go.mod h1:LzX7hefJvL54yjefDEDHNONDjII0t9xZLPXsUe+TKr0=
rsc.io/sampler v1.3.0 h1:7uVkIFmeBqHfdjD+gZwtXXI+RODJ2Wc4O7MPEh/QiW4=
rsc.io/sampler v1.3.0/go.mod h1:T1hPZKmBbMNahiBKFy5HrXp6adAjACjK9JXDnKaTXpA=
```
The go command uses the `go.sum` file to ensure that future downloads of these modules retrieve the same bits as the first download, to ensure the modules your project depends on do not change unexpectedly, whether for malicious, accidental, or other reasons. **Both `go.mod` and `go.sum` should be checked into version control.**



## Upgrading dependencies

With Go modules, versions are referenced with semantic version tags. A semantic version has three parts: `major`, `minor`, and `patch`. For example, for `v0.1.2`, the major version is 0, the minor version is 1, and the patch version is 2. Let's walk through a couple minor version upgrades. In the next section, we’ll consider a major version upgrade.

From the output of go list -m all, we can see we're using an untagged version of golang.org/x/text. Let's upgrade to the latest tagged version and test that everything still works:

```
$ go list -m all
github.com/gerryyang/goinaction/module/hello
golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c
rsc.io/quote v1.5.2
rsc.io/sampler v1.3.0
```

```
$ go get golang.org/x/text
go: downloading golang.org/x/text v0.3.5
go get: upgraded golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c => v0.3.5
$ cat go.mod 
module github.com/gerryyang/goinaction/module/hello

go 1.16

require (
        golang.org/x/text v0.3.5 // indirect
        rsc.io/quote v1.5.2 // indirect
)
$ cat go.sum 
golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c h1:qgOY6WgZOaTkIIMiVjBQcw93ERBE4m30iBm00nkL0i8=
golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
golang.org/x/text v0.3.5 h1:i6eZZ+zk0SOf0xgBpEpPD18qWcJda6q1sxt3S0kzyUQ=
golang.org/x/text v0.3.5/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
rsc.io/quote v1.5.2 h1:w5fcysjrx7yqtD/aO+QwRjYZOKnaM9Uh2b40tElTs3Y=
rsc.io/quote v1.5.2/go.mod h1:LzX7hefJvL54yjefDEDHNONDjII0t9xZLPXsUe+TKr0=
rsc.io/sampler v1.3.0 h1:7uVkIFmeBqHfdjD+gZwtXXI+RODJ2Wc4O7MPEh/QiW4=
rsc.io/sampler v1.3.0/go.mod h1:T1hPZKmBbMNahiBKFy5HrXp6adAjACjK9JXDnKaTXpA=
$ go test
PASS
ok      github.com/gerryyang/goinaction/module/hello    0.003s
```
Woohoo! Everything passes. Let's take another look at `go list -m all` and the `go.mod` file:

```
$ go list -m all
github.com/gerryyang/goinaction/module/hello
golang.org/x/text v0.3.5
golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e
rsc.io/quote v1.5.2
rsc.io/sampler v1.3.0
```

The `golang.org/x/text` package has been upgraded to the latest tagged version. The `go.mod` file has been updated to specify v0.3.5 too. The `indirect` comment indicates a dependency is not used directly by this module, only indirectly by other module dependencies.

Now let's try upgrading the `rsc.io/sampler` minor version. Start the same way, by running go get and running tests:


```
$ go get rsc.io/sampler
go: downloading rsc.io/sampler v1.99.99
go get: upgraded rsc.io/sampler v1.3.0 => v1.99.99
$ go test
--- FAIL: TestHello (0.00s)
    hello_test.go:8: Hello() = "99 bottles of beer on the wall, 99 bottles of beer, ...", want "Hello, world."
FAIL
exit status 1
FAIL    github.com/gerryyang/goinaction/module/hello    0.007s
```

Uh, oh! The test failure shows that the latest version of `rsc.io/sampler` is **incompatible** with our usage. Let's list the available tagged versions of that module:


```
$ go list -m -versions rsc.io/sampler
rsc.io/sampler v1.0.0 v1.2.0 v1.2.1 v1.3.0 v1.3.1 v1.99.99
```

We had been using v1.3.0; v1.99.99 is clearly no good. Maybe we can try using v1.3.1 instead:


```
$ go get rsc.io/sampler@v1.3.1
go: downloading rsc.io/sampler v1.3.1
go get: downgraded rsc.io/sampler v1.99.99 => v1.3.1
$ go test
PASS
ok      github.com/gerryyang/goinaction/module/hello    0.009s
```

Note the explicit `@v1.3.1` in the `go get` argument. In general each argument passed to go get can take an explicit version; the default is `@latest`, which resolves to the latest version as defined earlier.

## Adding a dependency on a new major version

Let's add a new function to our package: `func Proverb` returns a Go concurrency proverb, by calling `quote.Concurrency`, which is provided by the module `rsc.io/quote/v3`. First we update `hello.go` to add the new function:

``` go
package hello

import (
    "rsc.io/quote"
    quoteV3 "rsc.io/quote/v3"
)

func Hello() string {
    return quote.Hello()
}

func Proverb() string {
    return quoteV3.Concurrency()
}
```

Then we add a test to `hello_test.go`:

``` go
func TestProverb(t *testing.T) {
    want := "Concurrency is not parallelism."
    if got := Proverb(); got != want {
        t.Errorf("Proverb() = %q, want %q", got, want)
    }
}
```

Then we can test our code:

```
$ go test
hello.go:5:2: no required module provides package rsc.io/quote/v3; to add it:
        go get rsc.io/quote/v3
$ go get rsc.io/quote/v3
go: downloading rsc.io/quote/v3 v3.1.0
go get: added rsc.io/quote/v3 v3.1.0
$ go test
PASS
ok      github.com/gerryyang/goinaction/module/hello    0.003s
```

Note that our module now depends on both `rsc.io/quote` and `rsc.io/quote/v3`:

```
$ cat go.mod 
module github.com/gerryyang/goinaction/module/hello

go 1.16

require (
        golang.org/x/text v0.3.5 // indirect
        rsc.io/quote v1.5.2 // indirect
        rsc.io/quote/v3 v3.1.0 // indirect
        rsc.io/sampler v1.3.1 // indirect
)
$ go list -m rsc.io/q...
rsc.io/quote v1.5.2
rsc.io/quote/v3 v3.1.0
```
Each different major version (v1, v2, and so on) of a Go module uses a different module path: starting at v2, the path must end in the major version. In the example, v3 of rsc.io/quote is no longer rsc.io/quote: instead, it is identified by the module path rsc.io/quote/v3. This convention is called [semantic import versioning](https://research.swtch.com/vgo-import), and it gives **incompatible packages (those with different major versions) different names. In contrast,`v1.6.0` of `rsc.io/quote` should be backwards-compatible with `v1.5.2`, so it reuses the name `rsc.io/quote`**. (In the previous section, rsc.io/sampler v1.99.99 should have been backwards-compatible with rsc.io/sampler v1.3.0, but bugs or incorrect client assumptions about module behavior can both happen.)

**The go command allows a build to include at most one version of any particular module path, meaning at most one of each major version: one rsc.io/quote, one rsc.io/quote/v2, one rsc.io/quote/v3, and so on.** This gives module authors a clear rule about possible duplication of a single module path: it is impossible for a program to build with both rsc.io/quote v1.5.2 and rsc.io/quote v1.6.0. At the same time, allowing different major versions of a module (because they have different paths) gives module consumers the ability to upgrade to a new major version incrementally. In this example, we wanted to use quote.Concurrency from rsc/quote/v3 v3.1.0 but are not yet ready to migrate our uses of rsc.io/quote v1.5.2. The ability to migrate incrementally is especially important in a large program or codebase.

## Upgrading a dependency to a new major version

Let's complete our conversion from using `rsc.io/quote` to using only `rsc.io/quote/v3`. Because of the major version change, we should expect that some APIs may have been removed, renamed, or otherwise changed in incompatible ways. Reading the docs, we can see that `Hello` has become `HelloV3`:

```
$ go doc rsc.io/quote/v3
package quote // import "rsc.io/quote/v3"

Package quote collects pithy sayings.

func Concurrency() string
func GlassV3() string
func GoV3() string
func HelloV3() string
func OptV3() string
```

We can update our use of `quote.Hello()` in hello.go to use `quoteV3.HelloV3()`:

``` go
package hello

import (
        //"rsc.io/quote"
        quoteV3 "rsc.io/quote/v3"
)

func Hello() string {
    //return "Hello, world."
        //return quote.Hello()
        return quoteV3.HelloV3()
}

func Proverb() string {
        return quoteV3.Concurrency()
}
```

Let's re-run the tests to make sure everything is working:

```
$ go test
PASS
ok      github.com/gerryyang/goinaction/module/hello    0.003s
```

## Removing unused dependencies

We've removed all our uses of rsc.io/quote, but it still shows up in go list -m all and in our go.mod file:

```
$ go list -m all
github.com/gerryyang/goinaction/module/hello
golang.org/x/text v0.3.5
golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e
rsc.io/quote v1.5.2
rsc.io/quote/v3 v3.1.0
rsc.io/sampler v1.3.1
$ cat go.mod
module github.com/gerryyang/goinaction/module/hello

go 1.16

require (
        golang.org/x/text v0.3.5 // indirect
        rsc.io/quote v1.5.2 // indirect
        rsc.io/quote/v3 v3.1.0 // indirect
        rsc.io/sampler v1.3.1 // indirect
)
```

Why? Because building a single package, like with go build or go test, can easily tell when something is missing and needs to be added, but not when something can safely be removed. Removing a dependency can only be done after checking all packages in a module, and all possible build tag combinations for those packages. An ordinary build command does not load this information, and so it cannot safely remove dependencies.

The `go mod tidy` command cleans up these unused dependencies:

```
$ go mod tidy
$ go list -m all
github.com/gerryyang/goinaction/module/hello
golang.org/x/text v0.3.5
golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e
rsc.io/quote/v3 v3.1.0
rsc.io/sampler v1.3.1
$ cat go.mod 
module github.com/gerryyang/goinaction/module/hello

go 1.16

require (
        golang.org/x/text v0.3.5 // indirect
        rsc.io/quote/v3 v3.1.0
        rsc.io/sampler v1.3.1 // indirect
)
$ go test
PASS
ok      github.com/gerryyang/goinaction/module/hello    0.003s
```

## Conclusion

Go modules are the future of dependency management in Go. Module functionality is now available in all supported Go versions (that is, in Go 1.11 and Go 1.12).

This post introduced these workflows using Go modules:

* `go mod init` creates a new module, initializing the `go.mod` file that describes it.
* `go build`, `go test`, and other package-building commands add new dependencies to `go.mod` as needed.
* `go list -m all` prints the current module’s dependencies.
* `go get` changes the required version of a dependency (or adds a new dependency).
* `go mod tidy` removes unused dependencies.

We encourage you to start using modules in your local development and to add `go.mod` and `go.sum` files to your projects. 


# Refer

* https://blog.golang.org/using-go-modules
* https://golang.org/doc/go1.11#modules
