#!/bin/bash

tag="latest"

if [ -z "$DOCKER_USERNAME" ]; then 
        echo >&2 "Docker Username Not Set"; exit 1;
fi

if [ -z "$DOCKER_PASSWORD" ]; then 
        echo >&2 "Docker password Not Set"; exit 1;
fi

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

if [ ! -z "$TRAVIS" ]; then 

        if [ "$TRAVIS_PULL_REQUEST" = "false" ]; then

                if [ -n "$TRAVIS_TAG" ]; then
                        tag="$TRAVIS_TAG"
                fi

                echo "PUSHING TO DOCKERHUB"
                echo "--------------------"
                echo "TAG: $tag"

                docker tag faker:latest uknth/faker:$tag

                docker push uknth/faker:$tag
        fi
fi



