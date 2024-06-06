# When using this Dockerfile, make sure to also add static Traefik config for
# loading the local plugin. Such as via CLI:
#   --experimental.localplugins.remoteaddr.modulename=github.com/RiskIdent/traefik-remoteaddr-plugin

ARG TRAEFIK_VERSION=v3.0.0
ARG BASE_IMAGE=docker.io/traefik:${TRAEFIK_VERSION}
FROM ${BASE_IMAGE}

COPY . plugins-local/src/github.com/RiskIdent/traefik-remoteaddr-plugin/
