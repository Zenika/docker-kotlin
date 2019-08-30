#!/bin/sh

export CIRCLE_BUILD_DATE=$(date -Ins --utc)

docker image build -t docker-kotlin:$VERSION $BUILD_CONTEXT
