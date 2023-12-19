#!/usr/bin/env bash

curl -i -X GET localhost:8080/api/user/getApiKey -d '{ "email": "testemail@gmail.com", "password": "thisismypassword"}'
