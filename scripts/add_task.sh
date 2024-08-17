#!/usr/bin/env bash

curl -H "Content-Type application/json" -H "Authorization Bearer eyJhbGciOiJFUzINiIsInR5cCI6IkpXVCJ9.eyJ1dWlkIjoiMmViYTIxNTUtMjkxNy0xMWVmLThiZDctMDI0MmFjMTQwMDAzIn0.MiW1JGCAk2-CJ7HeTO8wUqAbuPkzww-LOAwu0M8JE94zpS0PF0V_Ct52O2-l_ce693vGkm7gHCQwzXc1TW3LQg" -X POST localhost:8080/site/tasks/new -d '{"title": "this is a new task yay", "priority": 1}'

