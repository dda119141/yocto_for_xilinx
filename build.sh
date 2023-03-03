#!/bin/bash -e

#stop script when error occurs
set -o errexit 
set -o pipefail

readonly cur_dir=$(dirname $(readlink -e $0))
readonly re_c='^[0-9]+$'
readonly generated_dir="$cur_dir/generated"
readonly link_dir="$cur_dir/links"
readonly source_dir="${link_dir}/sources"
readonly build_dir="${link_dir}/builds"
readonly install_dir="${link_dir}/install"

usage() {
	echo "$0 is a wrapper for executing bitbake tasks."

	echo " "
	echo "Options:"
	echo "-c | --component <component> component to build  "
	echo "-s | --source_link create symlink source folder."
	echo "-b | --build_link create symlink build folder."
	echo "-i | --install_link create symlink install folder."
	echo "-h | --help                       show help"
}

if [ $# -lt 1 ]; then
	usage
	exit 1
fi

while [[ $# -gt 0 ]]; do
	key="$1"

	case ${key} in
		-c|--component)
			COMPONENT="$2"
			shift
			shift
			if [[ ${COMPONENT} =~ ${re_c} ]]; then
				echo "$COMPONENT"
				usage
				exit 1
			fi
			;;
		-s|--source_link)
			COMPONENT_S="$2"
			shift
			shift
			if [[ ${COMPONENT_S} =~ ${re_c} ]]; then
				usage
				exit 1
			fi
			;;
		-b|--build_link)
			COMPONENT_B="$2"
			shift
			shift
			if [[ ${COMPONENT_B} =~ ${re_c} ]]; then
				usage
				exit 1
			fi
			;;
		-i|--install_link)
			COMPONENT_I="$2"
			shift
			shift
			if [[ ${COMPONENT_I} =~ $re_c ]]; then
				usage
				exit 1
			fi
			;;
		-h|--help)
			usage
			exit 0
			;;
		*)
			usage
			exit 1
			;;
	esac
done

#######################################
# Cleanup files from the backup directory.
# Globals:
#   None
# Arguments:
#   DIRECTORY NAME
######################################
do_remove_not_relevant_char() {
	local _dir="$1"
	local l_prefix=$2

	#remove prefix and remove double quotes
	_dir_path=$(realpath $(echo ${_dir} | sed 's:"::g' | grep -oe '/.*' ))

	echo ${_dir_path}
}

#######################################
# Retrieve folder from Yocto recipe.
# Globals:
#   None
# Arguments:
#   Generated yocto folder creation.
######################################
do_retrieve_folder() {
	_pref="$3"
	source ${1}/setup_custom_project ${generated_dir}
	local _var="$(bitbake -e "$2" | grep $_pref)"
	local _dir_out=$(do_remove_not_relevant_char ${_var} ${_pref})

	echo ${_dir_out}
}

#######################################
# Build component.
# Globals:
#   None
# Arguments:
#   Generated yocto folder creation.
######################################
do_build() {
	source $1/setup_custom_project ${generated_dir}
	bitbake "$2"
}

#######################################
# Retrieve component source path.
# Globals:
#   None
# Arguments:
#   Generated yocto folder creation.
#   Component from which source path will be retrieved
######################################
do_get_sources() {
	local _component="$2"
	local out=$(do_retrieve_folder ${1} ${_component} "^S=")
	local _comp=$(echo ${_component/\//_})
	[ -d ${source_dir} ] && ln -s ${out} "${source_dir}/${_comp}"
}

#######################################
# Retrieve component install path.
# Globals:
#   None
# Arguments:
#   Generated yocto folder creation.
#   Component from which install path will be retrieved
######################################
do_get_install_folder() {
	local _component="$2"
	local out=$(do_retrieve_folder ${1} ${_component} "^D=")
	local _comp=$(echo ${_component/\//_})
	[ -d ${install_dir} ] && ln -s ${out} ${install_dir}/${_comp}
}

#######################################
# Retrieve component build path.
# Globals:
#   None
# Arguments:
#   Generated yocto folder creation.
#   Component from which install path will be retrieved
######################################
do_get_build_folder() {
	local _component="$2"
	local out=$(do_retrieve_folder ${1} ${_component} "^B=")
	local _comp=$(echo ${_component/\//_})
	[ -d ${build_dir} ] && ln -s ${out} ${build_dir}/${_comp}
}

#ToDo: function comment
do_prepare() {
   [ ! -d ${link_dir} ] && mkdir ${link_dir}
   [ ! -d ${build_dir} ] && mkdir ${build_dir}
   [ ! -d ${source_dir} ] && mkdir ${source_dir}
   [ ! -d ${install_dir} ] && mkdir ${install_dir}

   return 0
}

do_prepare

[ -z ${COMPONENT} ] || do_build $cur_dir ${COMPONENT};
[ -z ${COMPONENT_S} ] || do_get_sources $cur_dir ${COMPONENT_S};
[ -z ${COMPONENT_B} ] || do_get_build_folder $cur_dir ${COMPONENT_B};
[ -z ${COMPONENT_I} ] || do_get_install_folder $cur_dir ${COMPONENT_I};

