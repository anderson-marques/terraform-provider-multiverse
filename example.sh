#!/bin/bash

echo bash > debug.txt

INPUT=$(cat)
MODE=$1
ID=$(echo -n $INPUT | jq .ID)
PAYLOAD=$(echo $INPUT | jq .Payload)
# remove wrapping quotes
if [ ! -z "$ID" ]
then
    ID=${ID:1:${#ID} - 2}
fi

echo "MODE: '$MODE'" >> debug.txt
echo "STDIN: '$INPUT'" >> debug.txt
echo "ID: '$ID'" >> debug.txt
echo "PAYLOAD: '$PAYLOAD'" >> debug.txt

read_resource() {
    echo "read_resource" >> debug.txt
    echo $(cat tmp/$ID.json) $(cat tmp/$ID.json | jq '{DeepObject: . | tojson}') | jq -s add
}

create_resource() {
    echo "create_resource" >> debug.txt
    ID=$(date +"%s")
    echo "ID: '$ID'" >> debug.txt
    JSON1=$(jq -n "{ID: \"$ID\", deep:{more:{here:123}}}")
    JSON2=$(jq -n "$PAYLOAD | fromjson")
    echo "$JSON1 $JSON2" | jq -s add > tmp/$ID.json
    read_resource
}

update_resource() {
    echo "update_resource" >> debug.txt
    echo "ID: '$ID'" >> debug.txt
    JSON1=$(jq -n "{ID: \"$ID\"}")
    JSON2=$(jq -n "$PAYLOAD | fromjson")
    echo "$JSON1 $JSON2" | jq -s add > tmp/$ID.json
    read_resource
}

delete_resource() {
    echo "update_resource" >> debug.txt
    echo "ID: '$ID'" >> debug.txt
    rm tmp/$ID.json
    echo '{}'
}

case $MODE in
    read)
        read_resource
        ;;
    create)
        create_resource
        ;;
    update)
        update_resource
        ;;
    delete)
        delete_resource
        ;;
    *)
        echo "Unknown mode '$MODE'" >> debug.txt
        ;;

esac


echo "finish" >> debug.txt