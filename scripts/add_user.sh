#!/usr/bin/env bash

curl -i -X POST localhost:8080/api/user/new -d '{"firstName": "frank", "lastName": "Sinatra", "email": "testemail@gmail.com", "password": "thisismypassword"}'
