FROM golang:1.24-bookworm AS base

RUN apt-get update; \
# RUN apt-get install openssl-dev musl-dev gcc
# RUN apt-get install gcc
    apt-get install \
    libssl-dev

WORKDIR /app

RUN mkdir /out

COPY go.mod .
COPY go.sum .
RUN go mod download

FROM base AS builder
ENV APP_NAME="auth-id"
ARG GIT_BRANCH="GIT_BRANCH"
ARG BUILD_DATE="BUILD_DATE"
ARG BUILD_NUMBER="BUILD_NUMBER"
ARG GIT_SHORT_HASH="GIT_SHORT_HASH"
ARG GIT_TAG="GIT_TAG"
ARG VERSION_PACKAGE_PATH="0.0.1"

WORKDIR /app
ADD . /app

# COPY build/app/migrations /out/migrations

# RUN apk add openssl-dev musl-dev gcc

RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    go build \
    -ldflags \
    "-X ${VERSION_PACKAGE_PATH}.gitBranch=${GIT_BRANCH} \
    -X ${VERSION_PACKAGE_PATH}.buildDate=${BUILD_DATE} \
    -X ${VERSION_PACKAGE_PATH}.buildNumber=${BUILD_NUMBER} \
    -X ${VERSION_PACKAGE_PATH}.shortCommitHash=${GIT_SHORT_HASH} \
    -X ${VERSION_PACKAGE_PATH}.gitTag=${GIT_TAG}"\
    -o /out/${APP_NAME} cmd/${APP_NAME}/*.go

# FROM */golang-migrate:bookworm as migrator

# FROM */golang-runtime:bookworm AS release
FROM golang:1.24 AS release

ENV APP_NAME="auth-id"

# COPY --from=migrator /bin/migrate /bin/migrate
# COPY --from=build /out/migrations /migrations

WORKDIR /app

COPY --from=builder /out/${APP_NAME} /app/

CMD ["/app/auth-id"]
