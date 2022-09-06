FROM golang:1.18-alpine AS builder
ARG SVC
ARG GOARCH
ARG GOARM

WORKDIR /go/src/github.com/0x6flab/jikoniApp/BackendApp
COPY . .
RUN apk update
RUN apk add make
RUN make $SVC
RUN mv build/jikoni-$SVC /exe

FROM scratch
# Certificates are needed so that mailing util can work.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /exe /
ENTRYPOINT ["/exe"]
