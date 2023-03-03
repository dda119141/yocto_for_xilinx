inherit core-image

IMAGE_AUTOLOGIN = "1"

AUTOLOGIN = "${@ 'autologin' if d.getVar('IMAGE_AUTOLOGIN') == '1' else '' }"

inherit ${AUTOLOGIN}

COMMON_FEATURES = " \
    package-management \
    ssh-server-dropbear \
    "

IMAGE_FEATURES += "${COMMON_FEATURES}"

COMMON_INSTALL = " \
    mtd-utils \
    bridge-utils \
    pciutils \
    kernel-modules \
    "

IMAGE_INSTALL = " \
    packagegroup-core-boot \
    ${COMMON_INSTALL} \
    ${CORE_IMAGE_EXTRA_INSTALL} \
    "

CORE_IMAGE_EXTRA_INSTALL = " "
IMAGE_LINGUAS = " "


