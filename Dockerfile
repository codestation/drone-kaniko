FROM golang:1.14-alpine as builder

ARG DRONE_COMMIT_SHA=dev
ARG GOPROXY
ENV GOPROXY=${GOPROXY}
ENV CGO_ENABLED 0
WORKDIR /src

COPY go.mod go.sum /src/
RUN go mod download
COPY . /src/
RUN go build -o release/drone-kaniko \
   -ldflags "-w -s \
   -X main.version=${DRONE_COMMIT_SHA:0:8}" \
   -tags netgo \
   ./cmd/drone-kaniko

FROM gcr.io/kaniko-project/executor:debug-cdbd8af0578c56e2801b57461e9f417f9479d303

COPY --from=builder /src/release/drone-kaniko /drone-kaniko

ENTRYPOINT ["/drone-kaniko"]
