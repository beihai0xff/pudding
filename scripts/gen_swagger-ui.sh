#!/bin/bash

set -e

[[ -z "$SWAGGER_UI_VERSION" ]] && echo "missing \$SWAGGER_UI_VERSION" && exit 1
[[ -z "$APP" ]] && echo "missing \$APP" && exit 1


SWAGGER_UI_GIT="https://github.com/swagger-api/swagger-ui.git"
CACHE_DIR="./.cache/swagger-ui/$SWAGGER_UI_VERSION"
GEN_DIR="./third_party/swagger-ui"

rm -f ./third_party/swagger-ui/*.swagger.json # ignore nonexistent files, never prompt
cp -r ./api/http-spec/pudding/${APP}/v1/* $GEN_DIR/

escape_str() {
  echo "$1" | sed -e 's/[]\/$*.^[]/\\&/g'
}

# do caching if there's no cache yet
if [[ ! -d "$CACHE_DIR" ]]; then
  mkdir -p "$CACHE_DIR"
  tmp="$(mktemp -d)"
  git clone --depth 1 --branch "$SWAGGER_UI_VERSION" "$SWAGGER_UI_GIT" "$tmp"
  cp -r "$tmp/dist/"* "$CACHE_DIR"
  cp -r "$tmp/LICENSE" "$CACHE_DIR"
  rm -rf "$tmp"
fi

# populate swagger.json
tmp="    urls: ["
for i in $(find "$GEN_DIR" -name "*.swagger.json"); do
  escaped_gen_dir="$(escape_str "$GEN_DIR/")"
  path="${i//$escaped_gen_dir/}"
  tmp="$tmp{\"url\":\"$path\",\"name\":\"$path\"},"
done
# delete last characters from $tmp
tmp="${tmp//.$//}"
tmp="${tmp}],"

# recreate swagger-ui, delete all file except swagger.json
find "$GEN_DIR" -type f -not -name "*.swagger.json" -delete
mkdir -p "$GEN_DIR"
cp -r "$CACHE_DIR/"* "$GEN_DIR"

# replace the default URL
line="$(grep -n "url" "$GEN_DIR/swagger-initializer.js" | cut -f1 -d:)"
escaped_tmp="$(escape_str "$tmp")"
sed -i'' -e "$line s/^.*$/$escaped_tmp/" "$GEN_DIR/swagger-initializer.js"
rm -f "$GEN_DIR/swagger-initializer.js-e"