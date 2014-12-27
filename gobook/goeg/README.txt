Programming in Go by Mark Summerfield

ISBN: 0321774639

Copyright © 2011-12 Qtrac Ltd. 

All the programs, packages, and associated files in this archive are
licensed under the Apache License, Version 2.0 (the "License"); you may
not use these files except in compliance with the License. You can get a
copy of the License at: http://www.apache.org/licenses/LICENSE-2.0. (The
License is also included in this archive in file LICENSE-2.0.txt.)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

All the book's examples are designed to be educational, and many are
also designed to be useful. I hope that you find them helpful, and are
perhaps able to use some of them as starting points for your own
projects.

On Unix-like systems (e.g., Linux, FreeBSD, Mac OS X), once you have
installed Go, you can build all the examples in one go by executing:

    $ cd $HOME/goeg
    $ ./build.sh

On Windows it works similarly:

    C:\>cd goeg
    C:\goeg>build.bat

The build.sh (Unix) or build.bat (Windows) script sets GOPATH
temporarily just for the build and uses the go command (go build); both
assume that the go command (i.e., Go's bin directory) is in the PATH
which it will be if you installed a binary version.

If you want to build the examples individually and build your own Go
programs you will need to set GOPATH. This can be done temporarily by
running the accompanying gopath.sh (Unix) or gopath.bat (Windows) script
(after editing to change any paths to match your setup), or permanently
by adding the export lines from gopath.sh to your .bashrc file or on
Windows by creating a Go-specific console shortcut: see
gopath.sh or gopath.bat for more information.

Here is the list of programs and packages referred to in the book
grouped by chapter:

Chapter 1: An Overview in Five Examples
    hello
    bigdigits
    stack
    americanize
    polar2cartesian
    bigdigits_ans

Chapter 2: Identifiers, Booleans, and Numbers
    pi_by_digits
    statistics
    statistics_ans
    quadratic_ans1
    quadratic_ans2

Chapter 3: Strings
    m3u2pls
    playlist
    soundex

Chapter 4: Collection Types
    guess_separator
    wordfrequency
    chap4_ans

Chapter 5: Procedural Programming
    archive_file_list
    archive_file_list_ans
    statistics_nonstop
    statistics_nonstop2
    contains
    palindrome
    palindrome_ans
    memoize
    indent_sort
    common_prefix

Chapter 6: Object-Oriented Programming
    fuzzy
    fuzzy_immutable
    fuzzy_mutable
    fuzzy_value
    shaper1
    shaper2
    shaper3
    ordered_map
    qtrac.eu/omap
    font
    shaper_ans1
    shaper_ans2
    shaper_ans3

Chapter 7: Concurrent Programming
    filter
    cgrep1
    cgrep2
    cgrep3
    safemap
    apachereport1
    apachereport2
    apachereport3
    findduplicates
    safeslice
    apachereport4
    [apachereport5 added to examples after publication; see errata]
    imagetag1
    imagetag2
    sizeimages1
    sizeimages2

Chapter 8: File Handling
    invoicedata
    pack
    unpack
    unpack_ans
    utf16-to-utf8
    invoicedata_ans

Chapter 9: Packages
    qtrac.eu/omap
    cgrep3
    linkcheck
