#!/bin/bash

########################################################################################

main() {
  testWithCover
}

testWithCover() {
  local has_errors

  # Remove coverage output if exist
  rm -f cover.tmp cover.out &> /dev/null

  # Fix coverage header
  echo "mode: count" > cover.out

  for subpkg in "spec" "check" ; do
    go test -covermode=count -coverprofile=cover.tmp "./$subpkg"

    if [[ $? -ne 0 ]] ; then
      [[ -z $has_errors ]] && has_errors=true
    fi

    if [[ -f cover.tmp ]] ; then
      egrep -v '^mode:' cover.tmp >> cover.out
      rm -f cover.tmp
    fi
  done

  if [[ $has_errors ]] ; then
    exit 1
  fi
}

########################################################################################

main "$@"
