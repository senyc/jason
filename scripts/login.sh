#!/usr/bin/env bash

curl -i -X POST localhost:8080/api/user/login -d '{ "email": "testemail@gmail.com", "password": "thisismypassword"}'
