ARG GO_BUILD_IMAGE

FROM $GO_BUILD_IMAGE AS builder

RUN apt-get update && \
    apt-get install -y \
  hwloc \
  jq \
  libhwloc-dev \
  mesa-opencl-icd \
  ocl-icd-opencl-dev \
  build-essential

WORKDIR /go/bin/app
COPY . /go/bin/app

RUN make amd64

RUN mkdir -vp /go/bin/service/config && \
    cp /go/bin/app/build/service-amd64 /go/bin/service/proxy && \
    cp /go/bin/app/config/service.toml /go/bin/service/config/service.toml && \
    cp -r /go/bin/app/static /go/bin/service/static

RUN mkdir -vp /go/bin/client/config/ && \
    cp /go/bin/app/build/client-amd64 /go/bin/client/proxy && \
    cp /go/bin/app/config/client.toml /go/bin/client/config/client.toml && \
    cp -r /go/bin/app/static /go/bin/client/static

FROM alpine as service

ENV TZ Asia/Shanghai

RUN apk add tzdata && apk add ca-certificates

COPY --from=builder /go/bin/service/ /build

EXPOSE 8080

WORKDIR /build

CMD ["./proxy"]

FROM alpine as client

ENV TZ Asia/Shanghai

RUN apk add tzdata && apk add ca-certificates

COPY --from=builder /go/bin/client/ /build

WORKDIR /build

CMD ["./proxy"]
