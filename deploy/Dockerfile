FROM kodmain/service:builder as builder
USER root:root
COPY / /tmp/builder
WORKDIR /tmp/builder
RUN make build

FROM scratch
ARG BINARY_NAME=NO_SERVICE
ENV BINARY_NAME=${BINARY_NAME}
COPY --from=builder /tmp/builder/deploy/build/optimized/${BINARY_NAME} /run/service
WORKDIR /run
ENTRYPOINT [ "/run/service" ]