FROM alpine:3.15.4

COPY ./backend.app ./backend.app

ARG BACKEND_PORT_ARG=9000
ENV BACKEND_PORT=${BACKEND_PORT_ARG}

ENTRYPOINT ["/backend.app"]


