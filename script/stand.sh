#!/bin/bash

ROOT=$(dirname "$(echo "$0" | grep -E "^/" -q && echo "$0" || echo "$PWD/${0#./}")")
. "${ROOT}/lib.sh" || exit 1

PROPS_FILE="${ROOT}/props.env"
LOCAL_OVERRIDE_PROPS_FILE="${ROOT}/.override-props.env"

export __ARG_MODE_DEBUG__= # режим отладки go приложения
OPER=
export __REPEAT__=
ERR=

help() {
  info "use: sh stand.sh  - development stand management"
  info "commands: [command [-option]]"
  info "    report  - run report (rest:${__API_REPORT_REST_PORT_EXPOSE__}, debug:${__API_REPORT_DEBUGGING_PORT_EXPOSE__}"
  info ""
  info "    build   - build image"
  info "    deploy  - run docker container"
  info "    down    - stop run docker container"
  info ""
  info "options:"
  info "    -debug        - golang application debug mode. Work for api_clients, api_admins etc"
}

func_check_and_create_override_props &&
  func_apply_whole_env_file "$PROPS_FILE" &&
  func_apply_env_file "$LOCAL_OVERRIDE_PROPS_FILE" || exit 1

# ------------------------------------------------
# --------------      arguments     --------------
# ------------------------------------------------
for p in "$@"; do
  case "$p" in
  "report") OPER="report" ;;
  "build") OPER="build" ;;
  "deploy") OPER="deploy" ;;
  "down") OPER="down" ;;
  "-r" | "-repeat") __REPEAT__=1 ;;
  "-debug") __ARG_MODE_DEBUG__=1 ;;
  *)
    ERR=1
    err "Unknown argument '${p}'"
    ;;
  esac
done

unset p

{ [ -n "$ERR" ] || [ -z "$OPER" ]; } && help && exit 0

# ------------------------------------------------
# -----------      services      -----------------
# ------------------------------------------------

case "$OPER" in
"report")
  func_create_network || exit 1

  CMD="dlv debug /debugging/cmd/*.go --headless --listen=:40000 --api-version=2 --accept-multiclient --output /tmp/__debug_bin"
  [ -z "$__ARG_MODE_DEBUG__" ] && CMD="$(func_run_cmd "go run /debugging/cmd/*.go")"

  docker run -it --rm \
    --name "report" \
    --hostname "report" \
    --network "$__NET__" \
    -p "${__API_REPORT_REST_PORT_EXPOSE__}:8080" \
    -p "${__API_REPORT_DEBUGGING_PORT_EXPOSE__}:40000" \
    -w "/debugging" \
    -v "${ROOT}/../:/debugging:ro" \
    --env-file "${ROOT}/../deployment/.local.env" \
    -e "HTML_DIR=/debugging/html" \
    \
    "$(func_get_work_image)" bash -c "$CMD"
  ;;

"build")
  docker build --no-cache -f "${ROOT}/../Dockerfile" -t "$__DOCKER_IMAGE__" "${ROOT}/.."
  ;;

"deploy")
  has=$(docker image ls --filter=reference="marat-report" -q) || return 1
  [ -n "$has" ] && info "already started" && exit 0

  docker run --detach --restart unless-stopped \
    --cpus=0.3 \
    --memory 300m \
    --memory-swap 300m \
    --name="marat-report" \
    -p "28080:8080" \
    -v "/etc/ssl/certs:/etc/ssl/certs:ro" \
    --env-file "${ROOT}/../deployment/.local.env" \
    "$__DOCKER_IMAGE__"

  docker logs -f "marat-report"
  ;;

"down")
  docker stop "marat-report"
  ;;

"curl")
  HAS=$(docker images --filter=reference="$__DOCKER_CURL_IMAGE__" -q) || exit 1
  if [ -z "$HAS" ]; then
    docker build -f "${ROOT}/docker-files/curl.dockerfile" -t "$__DOCKER_CURL_IMAGE__" "$ROOT" || exit 1
  fi

  docker run -it --rm --network "$__NET__" -w /app "$__DOCKER_CURL_IMAGE__" bash
  ;;
esac
