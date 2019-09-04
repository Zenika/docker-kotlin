#!/bin/sh

for tag in $ADDITIONAL_TAGS; do
  echo Tagging $tag
  docker image tag $TAG $tag
done
