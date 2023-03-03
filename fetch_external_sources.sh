#!/bin/bash -e

readonly cur_dir=$(dirname $(readlink -e $0))
readonly xilinx_git="https://github.com/Xilinx/"
readonly array_repos=("meta-xilinx" "meta-openembedded" "poky" \
	"meta-xilinx-tools" "meta-xilinx-tsn")
readonly external_dir="${cur_dir}/external"
readonly xilinx_branch="rel-v2022.2"

readonly host_tools=(\
"gawk" "wget" "git" "diffstat" "unzip" "texinfo" "gcc" "build-essential" \
"chrpath" "socat" "cpio" "python3" "python3-pip" "python3-pexpect" "xz-utils" \
"debianutils" "iputils-ping" "python3-git" "python3-jinja2" "libegl1-mesa" \
"libsdl1.2-dev" "pylint3" "xterm" "python3-subunit" "mesa-common-dev" "zstd" "liblz4-tool" \
)

usage() {
	echo "$0 is a script for fetching xilinx yocto sources."

	echo " "
	echo "Options:"
	echo "-h | --help                       show help"
}

do_install_host_packages() {
	for tool in "${host_tools[@]}"; do
		local installed=$(dpkg -l | grep ${tool})
		if [[ -z ${installed} ]]; then
			echo "${tool} to be installed on host"
			sudo apt install ${tool}
		fi
	done
}


execute_if_dir_exits () {
	local variable=$1
	local command=$2

	if [[ -d ${variable} ]]; then
		echo ${command}
	fi
}

do_check_if_repo_exists() {
	local _dir="${external_dir}"/"${1}"

	local check_ext_dir=$(execute_if_dir_exits ${_dir} $(echo "exist"))
	if [[ -z ${check_ext_dir} ]]; then
		echo ""; return 1
	fi

	pushd ${_dir} 2>&1 > /dev/null
	local _repo=$(execute_if_dir_exits ${check_ext_dir} $(git remote -v | awk '{print $2}'))
	if [[ -n ${_repo} ]]; then
		echo "repo cloned"
	else
		echo ""
	fi
	popd
}


do_fetch_repos() {
	for rep in "${array_repos[@]}"; do
		local repc=${xilinx_git}/"${rep}"
		local repo_exists=$(do_check_if_repo_exists "${rep}")

		if [[ -z ${repo_exists} ]]; then
			pushd ${external_dir}
			git clone "${repc}.git" -b ${xilinx_branch}
		fi
	done
}


[ ! -d ${external_dir} ] && mkdir ${external_dir};

do_install_host_packages
do_fetch_repos


