#!/bin/sh
export GOPATH=`pwd`
cd src
DIRS="americanise
    apachereport1
    apachereport2
    apachereport3
    apachereport4
    apachereport5
    archive_file_list
    archive_file_list_ans
    bigdigits
    bigdigits_ans
    cgrep1
    cgrep2
    cgrep3
    chap4_ans
    common_prefix
    contains
    filter
    findduplicates
    font
    fuzzy
    fuzzy_immutable
    fuzzy_mutable
    fuzzy_value
    guess_separator
    hello
    imagetag1
    imagetag2
    indent_sort
    invoicedata
    invoicedata_ans
    linkcheck
    m3u2pls
    memoize
    ordered_map
    pack
    palindrome
    palindrome_ans
    pi_by_digits
    playlist
    polar2cartesian
    quadratic_ans1
    quadratic_ans2
    safemap
    safeslice
    shaper1
    shaper2
    shaper3
    shaper_ans1
    shaper_ans2
    shaper_ans3
    sizeimages1
    sizeimages2
    soundex
    stacker
    statistics
    statistics_ans
    statistics_nonstop
    statistics_nonstop2
    wordfrequency
    unpack
    unpack_ans
    utf16-to-utf8"

for dir in $DIRS
do
    cd $dir
    echo building $dir...
    rm -f $dir
    go build
    cd ..
done
