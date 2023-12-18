#!/usr/bin/env bash

curl -X POST localhost:8080/newTask/2 -d '{"title": "this is a new task yay", "priority": "2"}'
