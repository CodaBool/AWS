#!/bin/bash

arr=()
for i in $(find . -type f \( -path "*node_modules*" -o -path "./.git" \) -prune -o -name ".env" -print)
do
  # echo "+ $i"
  arr+=($i)
done

if [ -z "$arr" ]; then
  # echo "No .env files found"
  exit
fi

comma_arr=()
for i in "${arr[@]}"
do
  # echo "sync on $i"
  # TODO: test if "=" can be in the value side
  res=$(awk '{$1=$1}1' FS='\n' OFS=',' RS= $i)
  # echo "result = $res"
  comma_arr+=("F=$i,$res")
done

final_string=$(printf ",%s" "${comma_arr[@]}")
final_string=${final_string:1}

jq -n --arg env "$final_string" '{"env":$env}'