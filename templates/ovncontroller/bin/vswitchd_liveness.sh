#!/bin/bash

set -e

error_exit() {
    echo "$1" >&2
    exit 1
}

# Check if ovs_vswitchd is running
check_ovs_vswitchd_pid() {
    if ! pidof -q ovs-vswitchd; then
        error_exit "ERROR - ovs_vswitchd is not running"
    fi
}

# Function to check the running status of ovs_vswitchd
check_ovs_vswitchd_status() {
    if ! /usr/bin/ovs-appctl bond/show 2>&1; then
        error_exit "ERROR - Failed to get output from ovs_vswitchd, ovs-appctl exit status: $?"
    fi
}


check_ovs_vswitchd_pid
check_ovs_vswitchd_status