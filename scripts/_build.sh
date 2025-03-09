using colors

_bin_dir=${_root_dir}/bin

run_build_usage() {
  echo "  where $(color -bold op) is:"
  echo "    $(color -lt_green \<empty\>) : build the app"
  echo "    $(color -lt_green test)    : build the app for test"
}

run_build() {
  local _op=$1
  shift
  case ${_op} in
    -h|--help|help)
      run_handler_usage build
      exit 1
      ;;
    ""|postal)
      _build_app $*
      ;;
    generate)
      _build_generate $*
      ;;
    test)
      echo "Build for test ..."
      go build -tags testutils ./... $*
      ;;
    *)
      error "Unrecognized $(color -bold op): '${_op}'"
      exit 1
      ;;
  esac
}

_build_app() {
  echo "Build ..."
  mkdir -p ${_bin_dir}
  #cd ${_root_dir}/${_name}
  _version=$(_run bump -r ${_root_dir} -d -b alpha -s)
  go build -ldflags="-X postal/cmd.appVersion=${_version}" -o ${_bin_dir}/postal $* ./cmd/postal/main.go
}

_build_generate() {
  go generate ./...
}
