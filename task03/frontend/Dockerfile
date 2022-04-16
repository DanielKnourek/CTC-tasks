FROM alpine:3.15.4

COPY ./frontend_test.app ./frontend_test.app
COPY ./static ./static

ARG FRONTEND_PORT_ARG=9001
ENV FRONTEND_PORT=${FRONTEND_PORT_ARG}

ENTRYPOINT ["/frontend_test.app"]


