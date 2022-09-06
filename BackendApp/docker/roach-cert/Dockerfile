# Stage 1 - generate certs
FROM cockroachdb/cockroach:latest AS generator

ENV REFRESHED_AT 2021_11_22

RUN mkdir -pv /tmp/certs/ca /tmp/certs/jikoni-roach-0 /tmp/certs/jikoni-roach-1 /tmp/safe /tmp/certs/jikoni-roach-client /tmp/certs/jikoni-roach-lb

RUN cockroach cert create-ca --certs-dir=/tmp/certs/ca --ca-key=/tmp/safe/ca.key \
    && cp -v /tmp/certs/ca/ca.crt /tmp/certs/jikoni-roach-client \
    && cp -v /tmp/certs/ca/ca.crt /tmp/certs/jikoni-roach-0 \
    && cp -v /tmp/certs/ca/ca.crt /tmp/certs/jikoni-roach-1 \
    && cp -v /tmp/certs/ca/ca.crt /tmp/certs/jikoni-roach-lb

RUN ./cockroach cert create-client root --certs-dir=/tmp/certs/jikoni-roach-client --ca-key=/tmp/safe/ca.key --also-generate-pkcs8-key
RUN ./cockroach cert create-client roach --certs-dir=/tmp/certs/jikoni-roach-client --ca-key=/tmp/safe/ca.key --also-generate-pkcs8-key

RUN cp -v /tmp/certs/jikoni-roach-client/client.* /tmp/certs/jikoni-roach-0 \
    && cp -v /tmp/certs/jikoni-roach-client/client.* /tmp/certs/jikoni-roach-1 \
    && cp -v /tmp/certs/jikoni-roach-client/client.* /tmp/certs/jikoni-roach-lb

RUN cockroach cert create-node jikoni-roach-0 jikoni-roach-lb --certs-dir=/tmp/certs/jikoni-roach-0 --ca-key=/tmp/safe/ca.key
RUN cockroach cert create-node jikoni-roach-1 jikoni-roach-lb --certs-dir=/tmp/certs/jikoni-roach-1 --ca-key=/tmp/safe/ca.key
RUN cockroach cert create-node jikoni-roach-lb --certs-dir=/tmp/certs/jikoni-roach-lb --ca-key=/tmp/safe/ca.key

# Stage 2 - share certs

FROM alpine:3

RUN mkdir -pv /certs/jikoni-roach-0 /certs/jikoni-roach-1 /certs/jikoni-roach-client /certs/jikoni-roach-lb

COPY --from=generator  /tmp/certs/jikoni-roach-0/* /certs/jikoni-roach-0/
COPY --from=generator  /tmp/certs/jikoni-roach-1/* /certs/jikoni-roach-1/
COPY --from=generator  /tmp/certs/jikoni-roach-client/* /certs/jikoni-roach-client/
COPY --from=generator /tmp/certs/jikoni-roach-lb/* /certs/jikoni-roach-lb/

CMD tail -f /dev/null
