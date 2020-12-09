FROM envoyproxy/envoy-alpine:v1.15.3

RUN apk --no-cache add gettext coreutils

USER 10001

COPY config.yaml /envoy/envoy.yaml
COPY run-gateway.sh /envoy/run-gateway.sh

EXPOSE 8080 9901

RUN chmod -R 777 /envoy

CMD ["/envoy/run-gateway.sh"]