#!/bin/bash

_RED_='\033[0;31m'
_GREEN_='\033[0;32m'
_NC_='\033[0m'

err() {
  echo -e "${_RED_}[ERR]${_NC_} $*" >&2
}

info() {
  echo -e "${_GREEN_}[INFO]${_NC_} $*"
}

func_apply_whole_env_file() {
  set -o allexport
  # shellcheck disable=SC1090
  . "$1" || exit 1
  set +o allexport
}

func_apply_env_file() {
  tmp_file=$(mktemp) || return 1

  for file in $@; do
    [ ! -f "$file" ] && err "file [$file] does not exist" && return 1
    grep -Ei '^[^=]+=[^[:space:]]+' "$file" >>"$tmp_file"
  done

  set -o allexport
  # shellcheck disable=SC1090
  . "$tmp_file" || exit 1
  set +o allexport
}

func_create_network() {
  net_id=$(docker network ls -f "name=${__NET__}" -q) || return 1
  [ -n "$net_id" ] && return 0
  docker network create --attachable "$__NET__" || return 1
  sleep 1
}

func_run_cmd() {
  [ -n "$__REPEAT__" ] && echo "while true; do $1; echo '++++++++ restart'; sleep 3; done" || echo "$1"
}

func_build_if_not_exist_dlv_image() {
  has=$(docker image ls --filter=reference="$__DOCKER_DLV_IMAGE__" -q) || return 1
  [ -n "$has" ] && return 0
  docker build -f "${ROOT}/docker-files/dlv.dockerfile" -t "$__DOCKER_DLV_IMAGE__" "$ROOT" || return 1
}

func_get_work_image() {
  [ -n "$__ARG_MODE_DEBUG__" ] && echo "$__DOCKER_DLV_IMAGE__" || echo "$__DOCKER_GOLANG_IMAGE__"
}

func_check_and_create_override_props() {
  [ -f "$LOCAL_OVERRIDE_PROPS_FILE" ] && return 0

  {
    echo "# переменные для переопределения"
    echo ""
    grep -Eio '^([^=]+=|^\s*#+.*)' "$PROPS_FILE" |
      grep -Ei '_REPLICAS__|_PORT_EXPOSE__|_REPO_DIR__|__OVERRIDE_|^\s*#+'
  } >"$LOCAL_OVERRIDE_PROPS_FILE"
}
