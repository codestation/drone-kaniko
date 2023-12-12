FROM golang:1.21 as builder

ARG CI_COMMIT_TAG
ARG GOPROXY
ENV GOPROXY=${GOPROXY}

WORKDIR /src

COPY go.mod go.sum /src/
RUN go mod download
COPY . /src/

RUN  set -ex; \
	CGO_ENABLED=0 go build -o release/drone-kaniko \
    -trimpath \
    -tags netgo \
    -ldflags "-w -s \
    -X go.megpoid.dev/drone-kaniko/cmd/main.Tag=${CI_COMMIT_TAG}" \
    ./cmd/drone-kaniko

FROM mplatform/manifest-tool:2.1.4 as manifest-tool

# use the debug image since it comes with /kaniko/warmer
FROM gcr.io/kaniko-project/executor:v1.18.0-debug
LABEL maintainer="Codestation <codestation@megpoid.dev>"

COPY --from=builder /src/release/drone-kaniko /kaniko/
COPY --from=manifest-tool /manifest-tool /kaniko/manifest-tool

ENTRYPOINT ["/kaniko/drone-kaniko"]
