#!/bin/bash

set -x

VERSION=0.0.1 #like 1.0.0
BROKER_BASE_URL=http://localhost
TAG=dev
BRANCH=stock-service-pact
PROVIDER_NAME=product
PACT_PATH=./app/${PROVIDER_NAME}/pacts/stockservice-${PROVIDER_NAME}service.json

pact-broker publish \
${PACT_PATH} \
--consumer-app-version=${VERSION} \
--broker-base-url=${BROKER_BASE_URL} \
--tag=${TAG} \
--branch=${BRANCH}
