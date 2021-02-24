FROM envoyproxy/envoy-alpine:v1.17.0

RUN apk --no-cache add gettext coreutils

ENV ENVOY_UID=0
USER 10001

COPY config.yaml /envoy/envoy.yaml
COPY run-gateway.sh /envoy/run-gateway.sh

EXPOSE 8080 9901

USER root
RUN chmod -R 777 /envoy
USER 10001

CMD ["/envoy/run-gateway.sh"]