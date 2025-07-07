FROM node:22.17-alpine3.21 AS node_build

FROM golang:1.23-alpine3.21 AS build
WORKDIR /server
COPY . .
COPY --from=node_build /usr/lib /usr/lib
COPY --from=node_build /usr/local/lib /usr/local/lib
COPY --from=node_build /usr/local/include /usr/local/include
COPY --from=node_build /usr/local/bin /usr/local/bin
RUN apk update && \
    apk add just && \
    just setup-dev-env && \
    just build-prod

FROM scratch
WORKDIR /server
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /server/server /server/server
EXPOSE 8080
CMD [ "./app" ]