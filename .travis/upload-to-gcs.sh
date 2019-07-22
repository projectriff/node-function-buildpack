#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

gcloud auth activate-service-account --key-file <(echo ${GCLOUD_CLIENT_SECRET} | base64 --decode)

version=$(sed -n 's|version = \"\(.*\)\"|\1|p' buildpack.toml | head -n1)
git_sha=$(git rev-parse HEAD)
git_timestamp=$(TZ=UTC git show --quiet --date='format-local:%Y%m%d%H%M%S' --format="%cd")
git_branch=${TRAVIS_BRANCH}
slug=${version}-${git_timestamp}-${git_sha:0:16}

echo "Publishing buildpack"
gsutil cp -a public-read artifactory/io/projectriff/node/io.projectriff.node/${version}/io.projectriff.node-${slug}.tgz gs://projectriff/node-function-buildpack/
if [ "${git_branch}" = master ] ; then
    gsutil cp -a public-read artifactory/io/projectriff/node/io.projectriff.node/${version}/io.projectriff.node-${slug}.tgz gs://projectriff/node-function-buildpack/latest.tgz
fi

echo "Publishing version references"
gsutil -h 'Content-Type: text/plain' -h 'Cache-Control: private' cp -a public-read <(echo "${slug}") gs://projectriff/node-function-buildpack/versions/snapshots/${git_branch}
gsutil -h 'Content-Type: text/plain' -h 'Cache-Control: private' cp -a public-read <(echo "${slug}") gs://projectriff/node-function-buildpack/versions/snapshots/${version}
if [[ ${version} != *"-SNAPSHOT" ]] ; then
  gsutil -h 'Content-Type: text/plain' -h 'Cache-Control: private' cp -a public-read <(echo "${version}") gs://projectriff/node-function-buildpack/versions/releases/${git_branch}
fi
