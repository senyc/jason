#!/usr/bin/env bash

curl -i -X POST localhost:8080/api/user/new -d '{"firstName": "fred", "lastName": "Sinatra", "email": "newtest@gmail.com", "password": "thisismypassword"}'
