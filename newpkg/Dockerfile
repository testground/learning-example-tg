# BUILD_BASE_IMAGE is the base image to use for the build. It contains a rolling
# accumulation of Go build/package caches.
ARG BUILD_BASE_IMAGE=golang:1.17.9-buster

# This Dockerfile performs a multi-stage build and RUNTIME_IMAGE is the image
# onto which to copy the resulting binary.
#
# Picking a different runtime base image from the build image allows us to
# slim down the deployable considerably.
#
# The user can override the runtime image by passing in the appropriate builder
# configuration option.
ARG RUNTIME_IMAGE=busybox:1.31.1-glibc

#:::
#::: BUILD CONTAINER
#:::
FROM ${BUILD_BASE_IMAGE} AS builder

# PLAN_DIR is the location containing the plan source inside the container.
ENV PLAN_DIR /plan

# SDK_DIR is the location containing the (optional) sdk source inside the container.
ENV SDK_DIR /sdk

# Delete any prior artifacts, if this is a cached image.
RUN rm -rf ${PLAN_DIR} ${SDK_DIR} /testground_dep_list

# TESTPLAN_EXEC_PKG is the executable package of the testplan to build.
# The image will build that package only.
ARG TESTPLAN_EXEC_PKG="."

# GO_PROXY is the go proxy that will be used, or direct by default.
ARG GO_PROXY="direct"

# BUILD_TAGS is either nothing, or when expanded, it expands to "-tags <comma-separated build tags>"
ARG BUILD_TAGS

# TESTPLAN_EXEC_PKG is the executable package within this test plan we want to build. 
ENV TESTPLAN_EXEC_PKG ${TESTPLAN_EXEC_PKG}

# We explicitly set GOCACHE under the /go directory for more tidiness.
ENV GOCACHE /go/cache


# Copy only go.mod files and download deps, in order to leverage Docker caching.
COPY /plan/go.mod ${PLAN_DIR}/go.mod


# Download deps.
RUN echo "Using go proxy: ${GO_PROXY}" \
    && cd ${PLAN_DIR} \
    && go env -w GOPROXY="${GO_PROXY}" \
    && go mod download


# Now copy the rest of the source and run the build.
COPY . /


RUN cd ${PLAN_DIR} \
    && go env -w GOPROXY="${GO_PROXY}" \
    && CGO_ENABLED=${CgoEnabled} GOOS=linux GOARCH=amd64 go build -o ${PLAN_DIR}/testplan.bin ${BUILD_TAGS} ${TESTPLAN_EXEC_PKG}

# Store module dependencies
RUN cd ${PLAN_DIR} \
  && go list -m all > /testground_dep_list

#:::
#::: (OPTIONAL) RUNTIME CONTAINER
#:::

## The 'AS runtime' token is used to parse Docker stdout to extract the build image ID to cache.
FROM ${RUNTIME_IMAGE} AS runtime

# PLAN_DIR is the location containing the plan source inside the build container.
ENV PLAN_DIR /plan


COPY --from=builder /testground_dep_list /
COPY --from=builder ${PLAN_DIR}/testplan.bin /testplan

ARG RABBIT_URL="rabbitmq"

ENV ENV_TYPE=docker
ENV PORT=8081
ENV RABBIT_BROKER_URL=amqp://guest:guest@${RABBIT_URL}:5672/
ENV REDIS_URL=redis:6379

EXPOSE 6060
ENTRYPOINT [ "/testplan"]