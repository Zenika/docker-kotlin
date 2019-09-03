#!/bin/sh

export CIRCLE_BUILD_DATE=$(date -Ins --utc)

docker image build --pull -t $TAG $BUILD_CONTEXT
