# =============================================================================
# IONET - Dockerfile Unificado
# Uso: docker build --build-arg MODE=dev|prod -t ionet .
# =============================================================================

# Modo de construcción: dev o prod (default: prod)
ARG MODE=prod

# =============================================================================
# STAGE 1: Builder (compilación Go)
# =============================================================================
FROM golang:1.26-alpine AS builder

# Instalar git para go mod
RUN apk add --no-cache git

WORKDIR /build

# Copiar go.mod y go.sum primero (optimiza cache de Docker)
COPY go.mod go.sum ./
RUN go mod download

# Copiar todo el código fuente
COPY . .

# =============================================================================
# STAGE 2: UI Builder (compilación React)
# =============================================================================
FROM oven/bun:1 AS ui-builder

WORKDIR /app

# Copiar package files
COPY webui/react-ui/package.json webui/react-ui/bun.lock* ./
RUN bun install --frozen-lockfile

# Copiar fuente React
COPY webui/react-ui/ ./

# Build de producción del frontend
RUN bun run build

# =============================================================================
# STAGE 3: Binary final (sin frontend - se usa en todos los modos)
# =============================================================================
FROM golang:1.26-alpine AS binary-builder

WORKDIR /build

# Copiar módulos de go
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Arguments para linker
ARG LDFLAGS="-s -w"

# Build del binary
RUN CGO_ENABLED=0 go build -ldflags="$LDFLAGS" -o localagi ./

# =============================================================================
# STAGE 4: Prod Image (optimizada para producción)
# =============================================================================
FROM ubuntu:24.04 AS prod

LABEL maintainer="IONET <info@ionet.cl>"
LABEL description="IONET - Agente IA interno para ionet.cl"

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=America/Santiago

# Instalar runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    curl \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copiar binary compilado
COPY --from=binary-builder /build/localagi /usr/local/bin/

# Copiar frontend buildado (desde ui-builder)
COPY --from=ui-builder /app/dist /app/webui/dist

# Crear usuario para SSH (mejor práctica de seguridad)
RUN useradd -m -s /bin/bash ionet && \
    echo "ionet:ionet" | chpasswd

# Instalar y configurar SSH
RUN apt-get update && apt-get install -y openssh-server && \
    rm -rf /var/lib/apt/lists/*

# Configurar SSH
RUN mkdir /var/run/sshd && \
    sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config && \
    sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config && \
    sed -i 's/#Banner none/Banner \/etc\/ssh\/banner/' /etc/ssh/sshd_config && \
    echo "========================================" >> /etc/ssh/banner && \
    echo "  IONET - Acceso Restringido" >> /etc/ssh/banner && \
    echo "  Solo personal autorizado" >> /etc/ssh/banner && \
    echo "========================================" >> /etc/ssh/banner

# Scripts de inicio
COPY <<EOF /entrypoint.sh
#!/bin/bash
set -e

echo "[IONET] Iniciando en modo PRODUCCION"
echo "[IONET] Puerto SSH: 2222"
echo "[IONET] Puerto App: 8080"

# Iniciar SSH en background
/usr/sbin/sshd

# Iniciar aplicación
exec /usr/local/bin/localagi serve
EOF

RUN chmod +x /entrypoint.sh

EXPOSE 8080 2222

ENTRYPOINT ["/entrypoint.sh"]

# =============================================================================
# STAGE 5: Dev Image (para desarrollo local)
# =============================================================================
FROM alpine:3.19 AS dev

LABEL maintainer="IONET <info@ionet.cl>"
LABEL description="IONET - Desarrollo local"

ENV TZ=America/Santiago

# Instalar runtime mínimos + herramientas de desarrollo
RUN apk add --no-cache \
    ca-certificates \
    curl \
    tzdata \
    git \
    vim \
    htop \
    bash

WORKDIR /app

# Copiar binary y frontend (pre-buildados desde stages anteriores)
COPY --from=binary-builder /build/localagi /usr/local/bin/
COPY --from=ui-builder /app/dist /app/webui/dist

# Crear estructura de directorios para volumes
RUN mkdir -p /app/core /app/pkg /app/services /app/webui /app/cmd /pool

# Copiar código fuente para hot-reload (montado como volumen)
# Los volúmenes en docker-compose sobreescribirán estos

ENTRYPOINT ["/usr/local/bin/localagi", "serve"]
