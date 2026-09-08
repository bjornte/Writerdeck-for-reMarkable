#!/bin/sh
# Called from the stock sidebar row (inside xochitl).
# Must return at once: a blocking wget that POST /api/lobby stops xochitl while
# this click is still running. On 3.28 that segfaults xochitl; OnFailure then
# runs remarkable-fail.sh and reboots the tablet.
setsid /bin/sh -c '
sleep 1
wget -qO- --post-data= --header=Content-Type:application/json http://127.0.0.1:8000/api/lobby >/dev/null 2>&1 \
  || wget -qO- --post-data= http://127.0.0.1:8000/api/launch >/dev/null 2>&1 \
  || true
' </dev/null >/dev/null 2>&1 &
exit 0
