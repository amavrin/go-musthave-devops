#!/bin/bash

set -ue

ITER="${1:?}"
ITER_ARG="^TestIteration${ITER}$"

bash build.sh || { echo 'Build failed'; exit 1; }

./bin/devopstest \
	-test.v \
	-test.run="$ITER_ARG" \
	-agent-binary-path=cmd/agent/agent
