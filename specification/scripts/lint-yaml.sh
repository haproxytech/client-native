#!/bin/sh

# - Ensure single `tags` are used(openapi-generator doesn't work well with multiple tags)
# - Ensure Pascal case is used for `name` field

set -e pipefail

haproxy_spec="./specification/build/haproxy_spec.yaml"
found="false"
countTags=0
lineNo=0

# checkCamelCaseInNameField checks if usage of camel case name is supported
# and if not, it exit with error code 1
checkCamelCaseInNameField()
{
  local name="$(echo $@ | grep '^name: ')"
  if [ -n "$name" ]; then
    local upperCased="$(echo "$name" | cut -d ":" -f 2 | grep -e '[[:upper:]]')"

    if [ -n "$upperCased" ]; then
      # Add camel case names we want to support
      local allowedCamelCaseNameFields=" HAProxy Support X-Runtime-Actions forceDelete"

      local allowed="$(echo "$allowedCamelCaseNameFields" | grep "$upperCased")"
      if [ -z "$allowed" ] ; then
        echo "Camel case \"$trimLine\" used at line no: $lineNo. Use Pascal case instead."
        exit 1
      fi
    fi
  fi
}

while IFS= read -r line
do
  lineNo=$((lineNo+1))
  # skip root tags
  if [ "$line" = "tags:" ]; then
    continue
  fi

  trimLine="${line#"${line%%[![:space:]]*}"}"

  # ensure camel case is not used in name fields
  checkCamelCaseInNameField $trimLine

  if [ "$found" = "true" ] && [ -n "$trimLine" ] && [ -z "$(echo "$trimLine" | cut -d "-" -f 1)" ]; then
    countTags=$((countTags+1))
  else
    found="false"
    countTags=0
  fi

  if [ "$countTags" -gt 1 ]; then
    echo "Multiple tags are not supported. Additional tag: $line at line no: $lineNo"
    exit 1
  fi

  # handle only tags from path
  if [ "$trimLine" = "tags:" ]; then
    found="true"
  fi
done < "$haproxy_spec"

echo "Linting YAML PASSED"
exit 0
