DESCRIPTION = "custom_project image definition for board vairants"
LICENSE = "BSD"

require custom_project-image-common.inc

inherit extrausers

# Configure default users/groups
# Default rules (assumed no debug-tweaks image feature):
# * disabled root login (set by system default)
# * Add a user 'custom_project' with no password
#   - Set to immediately expire
#   - Add to the sudoers file
IMAGE_CLASSES += "extrausers"
EXTRA_USERS_PARAMS ?= "\
    useradd -p '' custom_project;passwd-expire project; \
"
EXTRA_USERS_SUDOERS ?= "custom_project ALL=(ALL:ALL) ALL;"

