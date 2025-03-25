FROM alpine/git AS downloader
ARG project_name="auth-id"
ARG project_repo="https://github.com/vtievsky/auth-id.git"
WORKDIR /tmp/${project_name}
RUN git clone --depth=1 ${project_repo} /tmp/${project_name}
#
FROM golang:alpine AS builder
ARG project_name="auth-id"
RUN apk add git openssl-dev musl-dev gcc
WORKDIR /go/${project_name}
COPY --from=downloader /tmp/${project_name} .
RUN go build -ldflags="-s -w" -o ./${project_name} ./cmd/${project_name}/main.go
#
FROM alpine
LABEL description="auth-id"
ARG project_name="auth-id"
COPY --from=builder /go/${project_name}/${project_name} /bin/${project_name}
#
WORKDIR /bin
ENTRYPOINT [ "/bin/auth-id" ]