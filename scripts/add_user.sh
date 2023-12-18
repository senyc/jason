#!/usr/bin/env bash

curl -i -X POST localhost:8080/newUser -d '{"firstName": "frank", "lastName": "Sinatra", "email": "testingcom"}'
