custom_project Yocto repo for the Yocto Base Build System
=============================================

This repository provides files to setup and build the Yocto base system for
supported custom custom_project products. The external source repositories are from the latest
yocto honister branch. Additionally, this repository can also be prepared for
installing debian packages from binaries/files being hosted in artifactory.


Getting Started
---------------

1. Clone the Xilinx repositories and switch to "honister" branch
    ```bash
    ./fetch_external_sources.sh

1. Yocto environment for Xilinx Zynq zcu102 can be sourced. This action enables bibatke commands to
   be executed. For instance if the custom_project minimal image has to be built, do the following:
    ```bash
    . setup_custom_project generated/
    bitbake custom_project-image-minimal
    ```

The custom_project-image-minimal contains a minimal root file system. In contrast to
core-image-minimal, the aptitude and systemd packages have been included.

1. If the linux kernel and device tree for Zynq zcu102 board has to be built, do:
    ```bash
    ./build.sh -c virtual/kernel
    ```

1. If the linux bootloader (u-boot) for Zynq zcu102 board has to be built, do:
    ```bash
    ./build.sh -c virtual/bootloader
    ```

1. If the linux bootloader (u-boot) for Zynq zcu102 board installation files
   have to be retrieved, do:
    ```bash
    ./build.sh -s virtual/bootloader
    ls -lht links/install/virtual_bootloader/
    ```

1. If the first stage bootloader for Zynq zcu102 board has to be built, do:
    ```bash
    ./build.sh -c fsbl-firmware
    ```

1. If the linux kernel source files for Zynq zcu102 board have to be 
   retrieved, do:
    ```bash
    ./build.sh -s virtual/kernel
    ls -lht links/sources/virtual_kernel/
    ```





