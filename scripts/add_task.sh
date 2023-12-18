#!/usr/bin/env bash

curl -X POST localhost:8080/api/tasks/new/2 -d '{"title": "this is a new task yay", "priority": "2"}'
