# syntax=docker/dockerfile:1

FROM gcr.io/distroless/base-debian10

WORKDIR ./

COPY battle_ship ./

USER nonroot:nonroot

ENTRYPOINT ["/battle_ship"]