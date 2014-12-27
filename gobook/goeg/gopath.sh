#!/bin/sh
# Programming in Go by Mark Summerfield ISBN: 0321774639

# Execute this shell script in a console before using the Go tools.
# Or, if you want to be able to use the tools at anytime, copy the
# uncommented export lines into your .bashrc file.

# Uncomment the following two exports if you installed and built from
# source instead of installing a binary package, and change the path if
# necessary.
#export GOROOT=/usr/local/go
#export PATH=$PATH:$GOROOT/bin

# Tell your shell where to find the Go book's examples ($HOME/goeg) and
# your own code ($HOME/app/go); change the paths if necessary. If you
# don't need the book's examples anymore just change it to have a single
# path. Important: wherever you put your own Go programs must have a src
# directory, e.g., $HOME/app/go/src, with your programs and packages
# inside it, e.g., $HOME/app/go/src/myapp.
export GOPATH=$HOME/app/go:$HOME/goeg
