uri() {
  sed 's|gs://|https://storage.googleapis.com/|' "${ROOT}"/dependency/url
}

sha256() {
  shasum -a 256 "${ROOT}"/dependency/node-function-invoker-*.tgz | cut -f 1 -d ' '
}
