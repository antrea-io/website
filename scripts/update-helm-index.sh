#!/usr/bin/env bash

# Copyright 2022 Antrea Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -eo pipefail

function echoerr {
    >&2 echo "$@"
}

_usage="Usage: $0 [--antrea-repo-url <URL>] --website-repo <DIR> --version <VERSION>
Update the Helm repo index file."

function print_usage {
    echoerr "$_usage"
}

function print_help {
    echoerr "Try '$0 --help' for more information."
}

ANTREA_REPO_URL="https://github.com/antrea-io/antrea"
WEBSITE_REPO=""
VERSION=""

while [[ $# -gt 0 ]]
do
key="$1"

case $key in
    --antrea-repo-url)
    ANTREA_REPO_URL="$2"
    shift 2
    ;;
    --website-repo)
    WEBSITE_REPO="$2"
    shift 2
    ;;
    --version)
    VERSION="$2"
    shift 2
    ;;
    -h|--help)
    print_usage
    exit 0
    ;;
    *)    # unknown option
    echoerr "Unknown option $1"
    exit 1
    ;;
esac
done

if [ "$WEBSITE_REPO" == "" ]; then
    echoerr "--website-repo is required"
    print_help
    exit 1
fi

if [ "$VERSION" == "" ]; then
    echoerr "--version is required"
    print_help
    exit 1
fi

THIS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

source $THIS_DIR/verify-helm.sh

if [ -z "$HELM" ]; then
    HELM="$(verify_helm)"
elif ! $HELM version > /dev/null 2>&1; then
    echoerr "$HELM does not appear to be a valid helm binary"
    print_help
    exit 1
fi

TMP_DIR=$(mktemp -d archives.XXXXXXXX)

RELEASE_ASSETS_URL="$ANTREA_REPO_URL/releases/download/$VERSION"
ARCHIVE_URL="$RELEASE_ASSETS_URL/antrea-chart.tgz"
INDEX_PATH="$WEBSITE_REPO/static/charts/index.yaml"

curl -sSfLo "$TMP_DIR/antrea-chart.tgz" "$ARCHIVE_URL"

$HELM repo index $TMP_DIR --merge $INDEX_PATH --url $RELEASE_ASSETS_URL

mv "$TMP_DIR/index.yaml" $INDEX_PATH

rm -rf $TMP_DIR
