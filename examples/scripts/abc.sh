#!/bin/bash

echo "2, abc.sh: "

if [ $# -gt 0 ]; then
  echo "$@";
fi