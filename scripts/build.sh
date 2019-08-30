#!/bin/sh

export CIRCLE_BUILD_DATE=$(date -Ins --utc)

docker image build -t zenika/kotlin:$VERSION $BUILD_CONTEXT
