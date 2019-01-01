#!/bin/bash
set -e
sudo apt update

# tested on ubuntu 16.04
sudo apt -y install genisoimage qemu virtinst qemu-kvm libvirt-dev libguestfs-tools virt-manager
