#!/bin/sh
# Run go test with coverage and output as junit report
workdir=.cover

generate_test_coverage_data() {
    rm -rf "$workdir"
    mkdir "$workdir"
    go test -cover -v | go-junit-report > "$workdir/go-results_tests.xml"
}

generate_test_coverage_data