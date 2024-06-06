ARG TRAEFIK_VERSION=v3.0.0
ARG BASE_IMAGE=docker.io/traefik:${TRAEFIK_VERSION}
FROM ${BASE_IMAGE}

COPY . plugins-local/src/github.com/RiskIdent/traefik-remoteaddr-plugin/
ENV TRAEFIK_EXPERIMENTAL_LOCALPLUGINS_remoteaddr_MODULENAME="github.com/RiskIdent/traefik-remoteaddr-plugin"
