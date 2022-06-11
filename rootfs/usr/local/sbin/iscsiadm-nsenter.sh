#!/bin/sh

set -e

iscsid_pid=$(pidof -s iscsid)
if [ -z "${iscsid_pid}" ]; then
  echo "Unable to find process id of iscsid"
  exit 1
fi
nsenter --mount="/proc/${iscsid_pid}/ns/mnt" --net="/proc/${iscsid_pid}/ns/net" -- /usr/local/bin/iscsiadm "$@"
