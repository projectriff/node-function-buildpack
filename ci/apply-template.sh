#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

template=${1}
comment=${2:-#}

if [ -f $template ] ; then
  echo "${comment} DO NOT EDIT - this file is the output of the '$template' template "
  while IFS= read -r line
  do
    expressions=$(echo $line | grep -oE '\{\{[^}]+\}\}') || true
    while read -r expression; do
      output=$(eval $(echo $expression | sed -e 's/^{{//g' | sed -e 's/}}$//g'))
      line=${line//$expression/$output}
    done <<< "$expressions"
    echo "$line"
  done < "$template"
fi
