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
    focal_flag=0

    while :; do
        case "${1-}" in
        -h | --help) usage ;;
        -v | --verbose) set -x ;;
        --version) echo "make_qemu_image_and_boot script version=${version}"; exit;;
        --focal) focal_flag=1 ;;
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

debootstrap_dir=debootstrap
root_filesystem=debootstrap.ext4.qcow2

# Install needed packages and cleanup
# debootstrap for debootstrap
# qemu-system-x86 for qemu
# linux-image-generic for supermin bug
sudo apt-get update && apt-get install -y --no-install-recommends \
  debootstrap \
  qemu-system-x86 \
  linux-image-generic \
  && rm -rf /var/lib/apt/lists/*

# Check if we already made debootstrap dir
if [ ! -d "$debootstrap_dir" ]; then
  # Create debootstrap directory.
  # - linux-image-generic: downloads the kernel image under /boot
  # If user wants focal ubuntu then use that, else jammy
  if [ ${focal_flag} = 0]; then
    sudo debootstrap \
        --include linux-image-generic \
        focal \
        "$debootstrap_dir" \
        http://archive.ubuntu.com/ubuntu \
    ;
  else
    sudo debootstrap \
        --include linux-image-generic \
        jammy \
        "$debootstrap_dir" \
        http://archive.ubuntu.com/ubuntu \
    ;
fi

linux_image="$(printf "${debootstrap_dir}/boot/vmlinuz-"*)"

if [! -f "$root_filesystem" ]; then
  # Create init script to print hello world and spin to stop user prompt
  sudo touch "${debootstrap_dir}/init"
  cat <<'EOF' | sudo tee "${debootstrap_dir}/init"
#!/bin/bash
echo "hello world"
while :; do :; done
EOF
  sudo chmod +x "${debootstrap_dir}/init"

  # Avoid supermin bug
  sudo dpkg-statoverride --force-statoverride-add --update --add root root 0644 "$linux_image"

  # Generate image file from debootstrap directory with 1G size
  sudo virt-make-fs \
    --format qcow2 \
    --size +1G \
    --type ext4 \
    "$debootstrap_dir" \
    "$root_filesystem" \
  ;
  sudo chmod 666 "$root_filesystem"
fi

qemu-system-x86_64 \
  -append 'console=ttyS0 root=/dev/sda init=/init' \
  -drive "file=${root_filesystem},format=qcow2" \
  -serial mon:stdio \
  -m 2G \
  -kernel "${linux_image}" \
;

# End Main script