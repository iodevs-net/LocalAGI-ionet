# =============================================================================
# IONET - Dockerfile Multi-stage
# Uso: 
#   DEV:    docker build --target dev -t ionet:dev .
#   PROD:   docker build --target prod -t ionet:prod .
# =============================================================================

# =============================================================================
# STAGE 1: Builder (compilación Go) - Compartido
# =============================================================================
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG LDFLAGS="-s -w"
RUN CGO_ENABLED=0 go build -ldflags="$LDFLAGS" -o localagi ./

# =============================================================================
# STAGE 2: UI Builder (compilación React) - Compartido
# =============================================================================
FROM oven/bun:1 AS ui-builder

WORKDIR /app

COPY webui/react-ui/package.json webui/react-ui/bun.lock* ./
RUN bun install --frozen-lockfile

COPY webui/react-ui/ ./
RUN bun run build

# =============================================================================
# STAGE 3: Producción
# =============================================================================
FROM ubuntu:24.04 AS prod

LABEL maintainer="IONET <info@ionet.cl>"
LABEL description="IONET - Producción"

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=America/Santiago

RUN apt-get update && apt-get install -y \
    ca-certificates curl tzdata openssh-server && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copiar binary y frontend
COPY --from=builder /build/localagi /usr/local/bin/
COPY --from=ui-builder /app/dist /app/webui/dist

# Copiar agentes pre-configurados (auto-importación)
RUN mkdir -p /pool/agents
COPY config/agents/*.json /pool/agents/

# Configurar SSH
RUN mkdir /var/run/sshd && \
    sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config && \
    sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config

# Entrypoint
COPY <<EOF /entrypoint.sh
#!/bin/bash
echo "[IONET] Iniciando en MODO PRODUCCIÓN"
echo "[IONET] Puerto SSH: 2222"
echo "[IONET] Puerto App: 8080"
/usr/sbin/sshd
exec /usr/local/bin/localagi serve
EOF
RUN chmod +x /entrypoint.sh

EXPOSE 8080 2222
ENTRYPOINT ["/entrypoint.sh"]

# =============================================================================
# STAGE 4: Desarrollo
# =============================================================================
FROM ubuntu:24.04 AS dev

LABEL maintainer="IONET <info@ionet.cl>"
LABEL description="IONET - Desarrollo local"

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=America/Santiago

RUN apt-get update && apt-get install -y \
    ca-certificates curl tzdata git vim htop bash && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copiar binary y frontend
COPY --from=builder /build/localagi /usr/local/bin/
COPY --from=ui-builder /app/dist /app/webui/dist

# Copiar agentes pre-configurados (auto-importación)
RUN mkdir -p /pool/agents
COPY config/agents/*.json /pool/agents/

# Crear estructura de directorios
RUN mkdir -p /app/core /app/pkg /app/services /app/webui /app/cmd /pool

EXPOSE 3000
ENTRYPOINT ["/usr/local/bin/localagi", "serve"]
