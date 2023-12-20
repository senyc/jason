#!/usr/bin/env bash

curl -i -X POST localhost:8080/api/user/key/new -d '{ "email": "newtest@gmail.com"}'
