#!/usr/bin/env bash
#==========================================================
# HEADER
#==========================================================
# SYNOPSIS
#    make_qemu_image_and_boot.sh [-hv] [--options]
#
# DESCRIPTION
#    This script is used to create and run an AMD64 linux filesystem image using QEMU that prints "hello world" after
#    bootup.
#
# OPTIONS
#    -h, --help           Print this help
#    -v, --verbose        Print script information
#    --version            Print script version
#    --
#
# EXAMPLES
#    ./make_qemu_image_and_boot.sh -h
#
#==========================================================
# IMPLEMENTATION
version="1.0"
#
#==========================================================
# END_OF_HEADER
#==========================================================

# Fail fast options
set -o errexit
set -o nounset
set -o pipefail

script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd -P)

# Show how the script should be used
usage() {
    cat <<USAGE_TEXT
Usage: $(basename "${BASH_SOURCE[0]}") [-h] [-v] [--version]

DESCRIPTION
    This script is used to to build and run an AMD64 Linux filesystem using QEMU. It prints "hello world" after
    startup.

    OPTIONS:
    -h, --help
        Print this help and exit.
    -v, --version
        Print script debug info.
    
USAGE_TEXT
    exit
}

# msg is only intended to be used for stderr calls
msg() {
    echo >&2 -e "${1-}"
}

# Show a message and exit the script with a specific code
die() {
    local -r msg="${1}"
    local -r code="${2:-90}"
    echo "${msg}" >&2
    exit "${code}"
}

# Parse the options provided with the script call
parse_user_options() {
    # default values of variables set from params
    # ci_flag=0

    while :; do
        case "${1-}" in
        -h | --help) usage ;;
        -v | --verbose) set -x ;;
        --version) echo "make_qemu_image_and_boot script version=${version}"; exit;;
        # Kill the script if an unknown option is given
        -?*) die "Unknown option: $1" ;;
        *) break ;;
        esac
        shift
    done

    return 0
}

parse_user_options "$@"

# Main script here

# Move to the script location
cd ${script_dir}


# End Main script