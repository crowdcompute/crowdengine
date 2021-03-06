#!/bin/bash
set -e

# Variables
VMNAME=$1
ARTIFACT=$2
IFACE=$(route | grep '^default' | grep -o '[^ ]*$')
HOSTIP=$(ifconfig ${IFACE} | sed -En 's/127.0.0.1//;s/.*inet (addr:)?(([0-9]*\.){3}[0-9]*).*/\2/p')

# colors and helpers
bold() { echo -e "\e[1m$@\e[0m" ; }
red() { echo -e "\e[31m$@\e[0m" ; }
green() { echo -e "\e[32m$@\e[0m" ; }
yellow() { echo -e "\e[33m$@\e[0m" ; }
die() { red "ERR: $@" >&2 ; exit 2 ; }
ok() { green "${@:-OK}" ; }

set_defaults() {
    AUTOSTART=false                 # Automatically start VM at boot time
    CPUS=1                          # Number of virtual CPUs
    FEATURE=host                    # Use host cpu features to the guest
    MEMORY=1024                     # Amount of RAM in MB
    DISK_SIZE=""                    # Disk Size in GB
    DNSDOMAIN=example.local         # DNS domain
    GRAPHICS=none                  # Graphics type
    IMAGEDIR=${HOME}/virt/images    # Directory to store images
    VMDIR=${HOME}/virt/vms          # Directory to store virtual machines
    BRIDGE=virbr0                   # Hypervisor bridge
    QCOW="ubuntu-16.04-server-cloudimg-amd64-disk1.img"
    NETWORK_PARAMS="bridge=${BRIDGE},model=virtio"
    OS_VARIANT="ubuntu16.04"
    
    USER_DATA=user-data
    META_DATA=meta-data
    # IMAGE_URL=https://cloud-images.ubuntu.com/releases/16.04/release
    RESIZE_DISK=false               # Resize disk (boolean)
    PUBKEY=""                       # SSH public key
    DISTRO=centos7                  # Distribution
    MACADDRESS=""                   # MAC Address
    PORT=-1                         # Console port
    TIMEZONE=Europe/Athens          # Timezone
    ADDITIONAL_USER=${USER}         # User
    # Reset OPTIND
    OPTIND=1
}

clean_vms() {
    (virsh destroy ${VMNAME} > /dev/null 2>&1 || true )
    (virsh undefine ${VMNAME} > /dev/null 2>&1 || true )
    (rm -rf ${VMDIR}/${VMNAME} > /dev/null 2>&1 || true )
}

provision_vm() {
    check_vmname_set
    
    [ -d "${VMDIR}/${VMNAME}" ] && rm -rf ${VMDIR}/${VMNAME}
    mkdir -p ${VMDIR}/${VMNAME}
    
    green "[OK] VM directory created"
    
    pushd ${VMDIR}/${VMNAME}
    touch ${VMNAME}.log
    
    prepare_cloudinit_iso
    
    IMAGE=${IMAGEDIR}/${QCOW}
    if [ ! -f ${IMAGEDIR}/${QCOW} ]
    then
        die "Cloud image not found. Please download it"
    fi
    
    green "[OK] Base image found"
    
    # Check if domain already exists
    domain_exists "${VMNAME}"
    if [ "${DOMAIN_EXISTS}" -eq 1 ]; then
        die "${VMNAME} already exists."
    fi
    
    
    # storpool_exists "${VMNAME}"
    # if [ "${STORPOOL_EXISTS}" -eq 1 ]; then
    #     die "Storage pool ${VMNAME} already exists"
    # fi
    
    # make the dir if it doesnt exits
    # mkdir -p ${VMDIR}
    # Start clean
    
    
    # copy image to the destination directory
    DISK=${VMNAME}.qcow2
    cp $IMAGE "${VMDIR}/${VMNAME}/${DISK}" && ok

    sudo qemu-img resize "${VMDIR}/${VMNAME}/${DISK}" 10G && ok

    import_vm
    
    # copy artifacts to image using the -a arg
    sudo virt-copy-in -a "${VMDIR}/${VMNAME}/${DISK}" ${ARTIFACT} /home/ubuntu/

    sleep 1

    (virsh start ${VMNAME}  &>> ${VMNAME}.log && ok )
    
    sleep 1

    # Eject cdrom
    virsh change-media ${VMNAME} hda --eject --config &>> ${VMNAME}.log

    if [ -f "/var/lib/libvirt/dnsmasq/${BRIDGE}.status" ]
    then
        yellow "Waiting for domain to get an IP address"
        MAC=$(virsh dumpxml ${VMNAME} | awk -F\' '/mac address/ {print $2}')
        while true
        do
            IP=$(grep -B1 $MAC /var/lib/libvirt/dnsmasq/$BRIDGE.status | head \
            -n 1 | awk '{print $2}' | sed -e s/\"//g -e s/,//)
            if [ "$IP" = "" ]
            then
                sleep 1
            else
                ok
                break
            fi
        done
    else
        yellow "Bridge looks like a layer 2 bridge, get the domain's IP address from your DHCP server"
        IP="<IP address>"
    fi
    

    green "SSH to ${VMNAME}: 'ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no ubuntu@${IP}'"
    
    # Remove the unnecessary cloud init files
    green "Cleaning up cloud-init files"
    rm $USER_DATA $META_DATA $CI_ISO && ok

    popd

}

import_vm() {
    (virt-install --import \
        --name ${VMNAME} \
        --memory ${MEMORY} \
        --vcpus ${CPUS} \
        --cpu ${FEATURE} \
        --disk "${VMDIR}/${VMNAME}/${DISK}",format=qcow2,bus=virtio \
        --disk ${CI_ISO},device=cdrom \
        --network ${NETWORK_PARAMS} \
        --os-type=linux \
        --os-variant=${OS_VARIANT} \
        --graphics none \
        --noreboot \
        --noautoconsole || true )
}

check_vmname_set() {
    [ -n "${VMNAME}" ] || die "VMNAME not set."
}

domain_exists() {
    virsh dominfo "${1}" > /dev/null 2>&1 \
    && DOMAIN_EXISTS=1 \
    || DOMAIN_EXISTS=0
}

storpool_exists() {
    virsh pool-info "${1}" > /dev/null 2>&1 \
    && STORPOOL_EXISTS=1 \
    || STORPOOL_EXISTS=0
}

prepare_cloudinit_iso() {
    CI_ISO=${VMNAME}-cidata.iso
    # cloud-init config: set hostname, remove cloud-init package,
    # and add ssh-key
    cat > $USER_DATA << _EOF_
Content-Type: multipart/mixed; boundary="==BOUNDARY=="
MIME-Version: 1.0
--==BOUNDARY==
Content-Type: text/cloud-config; charset="us-ascii"

#cloud-config

# Hostname management
preserve_hostname: False
hostname: ${VMNAME}
fqdn: ${VMNAME}.${DNSDOMAIN}

password: 1234
chpasswd: { expire: False }
ssh_pwauth: True

# Users
#users:
#    - default
#    - name: ${ADDITIONAL_USER}
#      groups: ['${SUDOGROUP}']
#      shell: /bin/bash
#      sudo: ALL=(ALL) NOPASSWD:ALL
#      ssh-authorized-keys:
#        - ${KEY}

# Configure where output will go
output:
  all: ">> /var/log/cloud-init.log"

# configure interaction with ssh server
#ssh_genkeytypes: ['ed25519', 'rsa']

# Install my public ssh key to the first user-defined user configured
# in cloud.cfg in the template (which is centos for CentOS cloud images)
#ssh_authorized_keys:
#  - ${KEY}

timezone: ${TIMEZONE}

# Remove cloud-init when finished with it
runcmd:
  - mkdir -p /home/ubuntu/uploads
  - sudo chmod 777 /home/ubuntu/uploads
  - rm -r /home/ubuntu/go1.11.4.linux-amd64.tar.gz
  - sudo usermod -a -G docker ubuntu
  - sudo touch /etc/cloud/cloud-init.disabled
  - sudo systemctl stop networking && systemctl start networking
  - sudo systemctl disable cloud-init.service
  - export HOSTIP=${HOSTIP}
  - echo "export HOSTIP=${HOSTIP}" >> /home/ubuntu/.profile
  - cd /home/ubuntu/ && ./gocc --rpc --http --httpport 8085 --httpaddr 0.0.0.0
_EOF_
    
    if [ ! -z "${SCRIPTNAME+x}" ]
    then
        SCRIPT=$(< $SCRIPTNAME)
        cat >> $USER_DATA << _EOF_

--==BOUNDARY==
Content-Type: text/x-shellscript; charset="us-ascii"
${SCRIPT}

--==BOUNDARY==--
_EOF_
    else
       cat >> $USER_DATA << _EOF_

--==BOUNDARY==--
_EOF_
    fi
    
    { echo "instance-id: ${VMNAME}"; echo "local-hostname: ${VMNAME}"; } > $META_DATA
    
    # Create ISO with cloud-init config
    if command -v genisoimage &>/dev/null
    then
        genisoimage -output $CI_ISO \
        -volid cidata \
        -joliet -r $USER_DATA $META_DATA &>> ${VMNAME}.log \
        && ok \
        || die "Could not generate ISO."
    else
        mkisofs -o $CI_ISO -V cidata -J -r $USER_DATA $META_DATA &>> ${VMNAME}.log \
        && ok \
        || die "Could not generate ISO."
    fi
    
}

# Main
kvm_group="$(ls -l /dev/kvm | awk '{ print $4 }')"
if groups $username | grep &>/dev/null "\b$kvm_group\b"; then
    green "[OK] Permissions are correct"
    clean_vms
    set_defaults
    provision_vm
    
else
    red "[ERROR] You don't have access to /dev/kvm. Add '${USER}' to the '${kvm_group}' group: "
    yellow "sudo add ${USER} to ${kvm_group} group"
fi
