#!/usr/bin/env bash

set -e

jar=/tmp/DynamoDBLocal.jar

# download the dynamo jar if necessary
if [ ! -e "$jar" ]
then
    if [ `uname` = "Darwin" ] ; then
	# this will prompt for java to be installed if necessary
	java -version
    else
    	sudo apt-get update && sudo apt-get install -y default-jre
    fi
    echo "Downloading dynamodb-local..."
    curl -L -k --url https://s3-us-west-2.amazonaws.com/dynamodb-local/dynamodb_local_latest.tar.gz -o /tmp/dynamodb_local_latest.tar.gz
    tar -zxvf /tmp/dynamodb_local_latest.tar.gz -C /tmp/
fi

exec java -jar "$jar" -sharedDb -inMemory -port 8002
