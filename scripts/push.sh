#!/bin/sh

echo $DOCKER_PASSWORD | docker login --username $DOCKER_ACCOUNT --password-stdin

echo Pushing $TAG
docker image push $TAG

for tag in $ADDITIONAL_TAGS; do
  echo Pushing $tag
  docker image push $tag
done
