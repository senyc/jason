#!/usr/bin/env bash

curl -H "Authorization: 74b416ed91cb7a819ba17cd49b8bfdddf4e2869a00099d67917657ef524b2e39" -X POST localhost:8080/api/tasks/new -d '{"title": "this is a new task yay", "priority": "1"}'

