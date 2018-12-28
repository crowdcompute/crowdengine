#!/bin/bash
set -e
sudo apt update
sudo apt -y install genisoimage qemu virtinst qemu-kvm libvirt-dev libguestfs-tools virt-manager
