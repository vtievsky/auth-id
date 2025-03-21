#!/bin/bash

session_id='eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDI1ODQxNTUsImlhdCI6MTc0MjU4MzI1NSwic2Vzc2lvbl9pZCI6IjlhNTMzZWJhLTkwOWMtNDIwNy1iMzQ0LWI0YTdhMWYwNzYyYyIsImFjY2Vzc19vbmx5Ijp0cnVlfQ.nIRNOAebbzQX-AC7DO6Gogqw0YCmdlMate4yDNcx4Rw'

count=1
if [ -n "$1" ]
then
count=$1
fi

for (( i=1; i<=$count; i++ ))
do
curl --location '127.0.0.1:8080/v1/users' \
    --header "Authorization: bearer ${session_id}"

curl --location '127.0.0.1:8080/v1/roles' \
    --header "Authorization: bearer ${session_id}"

curl --location GET '127.0.0.1:8080/v1/roles/admin/users' \
    --header "Authorization: bearer ${session_id}"

curl --location GET '127.0.0.1:8080/v1/users/admin/privileges' \
    --header "Authorization: bearer ${session_id}"
done
