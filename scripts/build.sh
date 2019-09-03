#!/bin/sh

export CIRCLE_BUILD_DATE=$(date -Ins --utc)

docker image build -t $TAG $BUILD_CONTEXT
