#!/usr/bin/env bash

os=$(uname -s)
machine=$(uname -m)
api_url="https://api.github.com/repos/mritd/tmux-sysinfo/releases/latest"

case $os in
  "Linux")
    case $machine in
      "i386")
        sys_arch="linux-386"
        ;;
      "x86_64")
        if lscpu | grep -q "avx2"; then
          sys_arch="linux-amd64-v3"
        else
          sys_arch="linux-amd64"
        fi
        ;;
      "armv5"*)
        sys_arch="linux-armv5"
        ;;
      "armv6"*)
        sys_arch="linux-armv6"
        ;;
      "armv7"*)
        sys_arch="linux-armv7"
        ;;
      "aarch64")
        sys_arch="linux-arm64"
        ;;
      *)
        sys_arch="unknown"
        ;;
    esac
    ;;
  "Darwin")
    case $machine in
      "x86_64")
        sys_arch="darwin-amd64"
        ;;
      "arm64")
        sys_arch="darwin-arm64"
        ;;
      *)
        sys_arch="unknown"
        ;;
    esac
    ;;
  *)
    sys_arch="unknown"
    ;;
esac

# The directory where this plugin is located.
CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

if [ ! -f $CURRENT_DIR/tmux-sysinfo ] && ! $(builtin type -P "tmux-mem-cpu-load" &> /dev/null) ; then
    tmux run-shell "echo \"tmux-sysinfo not found. Attempting to download.\""

    pushd $CURRENT_DIR #Pushd to the directory where this plugin is located.

    curl -sSL $(curl -sSL $api_url | jq -r ".assets[].browser_download_url | select (. | test(\"${sys_arch}\"))") > tmux-sysinfo
    if [ "$?" != "0" ]; then
        tmux run-shell "echo \"tmux-sysinfo download failed!!!\""
    else
        chmod +x tmux-sysinfo
        tmux run-shell "echo \"tmux-sysinfo download success...\""
    fi
fi
