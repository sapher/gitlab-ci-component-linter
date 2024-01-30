#!/bin/bash

rules_file=$1
readme_file=$2

# Check if jq is installed
if ! command -v jq &> /dev/null
then
	echo "jq could not be found!"
	exit 1
fi

# Check if rules file exists
if [ ! -f "$rules_file" ]; then
	echo "Rules file not found!"
	exit 1
fi

# Check if readme file exists
if [ ! -f "$readme_file" ]; then
	echo "Readme file not found!"
	exit 1
fi

# Read rules file
json_content=$(<"$rules_file")

# Build table
table="| Rule | Message | Severity |"
table+="\n|------|---------|----------|"
table+="\n$(echo "$json_content" | jq -r 'to_entries[] | "| \(.key) | \(.value.message) | `\(.value.severity)` |"')"

# Replace table in readme file
part1=$(sed '/<!-- BEGIN_HERE -->/q' "$readme_file")
part2=$(sed -n '/<!-- END_HERE -->/,$p' "$readme_file")

# Build new content
new_content="${part1}\n$table\n${part2}"

# Write new content to readme file
echo -e "$new_content" > "$readme_file"

# Add readme file to git index
git add "$readme_file"

exit 0
