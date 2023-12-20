#!/usr/bin/env bash

curl -H "Authorization: PDt4tr9mGFShyVPkBUZ7Mg" -X POST localhost:8080/api/tasks/new -d '{"title": "this is a new task yay", "priority": "1"}'

