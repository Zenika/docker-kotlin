#!/bin/sh

for tag in $ADDITIONAL_TAGS; do
  echo Tagging with $tag
  docker image tag $TAG $tag
done