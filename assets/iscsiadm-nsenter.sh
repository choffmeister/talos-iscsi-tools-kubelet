#!/bin/sh

set -e

iscsid_pid=$(pidof iscsid)
nsenter --mount="/proc/${iscsid_pid}/ns/mnt" --net="/proc/${iscsid_pid}/ns/net" -- /usr/local/sbin/iscsiadm "$@"
