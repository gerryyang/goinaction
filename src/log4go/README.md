
log4go
[https://code.google.com/archive/p/log4go/](https://code.google.com/archive/p/log4go/)

---

Logging package similar to log4j for the Go programming language

Mission:

The goal of log4go is to be a robust, configurable, powerful logging package to empower Go developers to debug their programs more effectively on the fly and diagnose problems in the field without hampering their effectiveness during development or hampering the performance of their applications.

Overview:

This package is a replacement logging package which will be both a drop-in replacement for and a significant extension of the built-in logging functionality in Go.

Features: * File logging with rotation (size, linecount, daily) and custom output formats * Console logging * Network logging via JSON and TCP/UDP * XML Logger * Closure logging for defered parameter expansion * Automatic log filtering based on log levels on a per-output basis * XML configuration available for no-compile changes to logging * Wrapper functions and global loggers for easy configuration and rapid deployment * Drop-in compatibility with code using the standard log package

Installation: goinstall log4go.googlecode.com/hg

Getting Started:

Please consult the Wiki GettingStarted page for how to get started using log4go, including more detailed installation and updating instructions


