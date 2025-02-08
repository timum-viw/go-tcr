#!/bin/bash

FILE="./supercoop"
COMMAND="echo 'File has been modified!'"
inotifywait -m -e modify "$FILE" | while read -r path action file; do
    echo "File $file has been modified."
    $COMMAND
done