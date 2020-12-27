#!/bin/bash

function removeEmptyDirectories {
  directory=$1

  for dir in $(cd $directory && find . -type d)
  do
    if [ dir == "." ]
    then
      continue
    fi

    path="$directory/$(echo $dir | cut -c3-)"

    ## Check THIS directory for "no files" after all subfolders are deleted
    printf "%-50s" "Checking $path"
    if [ -z "$(find ${path} -type f)" ]
    then
      rmdir "${path}"
      printf "Deleted\n"
    else
      printf "Keeping\n"
    fi
  done
}

for target in $(cd modules && (find . -type f -name *.go | grep -v "test"))
do
  printf "%-50s" "Mocking ${target}"
  mockgen -source="modules/${target}" -destination "mocks/${target}"
  printf "%-25s" "Generated: Yes"
  printf "Status: "
  if [ ! -z "$(go build "mocks/${target}" 2> >(grep -i 'not used'))" ]
  then
    rm "mocks/${target}"
    echo "Deleted"
  else
    echo "Verified"
  fi
done

printf "\n\n"
printf "Checking mocks for empty folders\n"
removeEmptyDirectories "./mocks"
