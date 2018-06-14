#!/bin/bash

########################################################################################

main() {
  installCodeclimateTestReporter
  testWithCover
}

installCodeclimateTestReporter() {
  echo "Installing latest version of codeclimate-test-reporter…"

  npm install -g codeclimate-test-reporter

  if [[ $? -ne 0 ]] ; then
    echo "[ERROR] Can't install codeclimate-test-reporter"
    exit 1
  fi

  echo "Installing latest version of goveralls…"

  go get -v github.com/mattn/goveralls

  if [[ $? -ne 0 ]] ; then
    echo "[ERROR] Can't install goveralls"
    exit 1
  fi

  echo ""
}

testWithCover() {
  local has_errors

  # Remove coverage output if exist
  rm -f coverage.tmp coverage.txt &> /dev/null

  # Fix coverage header
  echo "mode: count" > coverage.txt

  for subpkg in "spec" "check" ; do
    go test -covermode=count -coverprofile=coverage.tmp "./$subpkg"

    if [[ $? -ne 0 ]] ; then
      [[ -z $has_errors ]] && has_errors=true
    fi

    if [[ -f coverage.tmp ]] ; then
      egrep -v '^mode:' coverage.tmp >> coverage.txt
      rm -f coverage.tmp
    fi
  done

  if [[ $has_errors ]] ; then
    exit 1
  fi

  echo -e "\nSending coverage data to Coveralls…"

  goveralls -service travis-ci -repotoken $COVERALLS_TOKEN -coverprofile coverage.txt

  echo -e "\nSending coverage data to Codebeat…"

  codeclimate-test-reporter < coverage.txt
}

########################################################################################

main "$@"
