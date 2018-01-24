#!/bin/bash

########################################################################################

main() {
  installCodeclimateTestReporter
  testWithCover
}

installCodeclimateTestReporter() {
  echo "Installing latest version of codeclimate-test-reporter..."

  npm install -g codeclimate-test-reporter

  if [[ $? -ne 0 ]] ; then
    echo "[ERROR] Can't install codeclimate-test-reporter"
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
    echo "Testing ../$subpkg..."

    go test -covermode=count -coverprofile=coverage.tmp "../$subpkg"

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

  echo -e "\nSending coverage data to Codebeat..."

  codeclimate-test-reporter < coverage.txt
}

########################################################################################

main "$@"