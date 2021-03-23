
# How to Write Go Code


## Introduction

This document demonstrates the development of a simple Go package inside a module and introduces the go tool, the standard way to fetch, build, and install Go modules, packages, and commands.

> Note: This document assumes that you are using `Go 1.13 or later` and the `GO111MODULE` environment variable is not set. If you are looking for the older, pre-modules version of this document, it is archived [here - How to Write Go Code (with GOPATH)](https://golang.org/doc/gopath_code).


## Code organization

### Package

Go programs are organized into packages. A package is a collection of source files in the same directory that are compiled together. Functions, types, variables, and constants defined in one source file are visible to all other source files within the same package.

### Module

A repository contains one or more modules. A module is a collection of related Go packages that are released together. A Go repository typically contains only one module, located at the root of the repository. A file named `go.mod` there declares the module path: the import path prefix for all packages within the module. The module contains the packages in the directory containing its `go.mod` file as well as subdirectories of that directory, up to the next subdirectory containing another `go.mod` file (if any).

Note that you don't need to publish your code to a remote repository before you can build it. A module can be defined locally without belonging to a repository. However, it's a good habit to organize your code as if you will publish it someday.



# Refer

* [How to Write Go Code](https://golang.org/doc/code)
