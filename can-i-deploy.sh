#!/bin/bash

set -x

VERSION=0.0.1
PACTICIPANT=BasketService
BROKER_BASE_URL=http://localhost

pact-broker can-i-deploy \
--pacticipant=${PACTICIPANT} \
--version=${VERSION} \
--broker-base-url=${BROKER_BASE_URL}