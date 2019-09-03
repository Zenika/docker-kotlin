#!/bin/sh

docker image build \
  --pull \
  -t $TAG \
  --build-arg SOURCE=$SOURCE \
  --build-arg CIRCLE_SHA1=$CIRCLE_SHA1 \
  --build-arg TAG=$TAG \
  --build-arg COMPILER_URL=$COMPILER_URL \
  --build-arg CIRCLE_BUILD_DATE=$(date -Ins --utc) \
  $BUILD_CONTEXT

docker image inspect $TAG
