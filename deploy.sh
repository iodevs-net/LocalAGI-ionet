#!/bin/bash
# =============================================================================
# IONET Deployment Script
# =============================================================================
# Uso: ./deploy.sh [--gpu nvidia|amd|intel]
#
# Requiere: Docker, Docker Compose
# Configuración: editar .env antes de desplegar

set -e

GPU_TYPE="${1:-cpu}"

echo "=== IONET Deployment ==="
echo "GPU Type: $GPU_TYPE"

# Verificar .env
if [ ! -f .env ]; then
    echo "Error: .env no encontrado. Copia .env.example a .env y edítalo."
    exit 1
fi

# Verificar que OPENAI_API_KEY no sea el placeholder
source .env
if [ "$OPENAI_API_KEY" = "your-api-key-here" ]; then
    echo "Error: Configura tu OPENAI_API_KEY en .env"
    exit 1
fi

# Seleccionar docker-compose
case "$GPU_TYPE" in
    nvidia)
        COMPOSE_FILE="docker-compose.nvidia.yaml"
        ;;
    amd)
        COMPOSE_FILE="docker-compose.amd.yaml"
        ;;
    intel)
        COMPOSE_FILE="docker-compose.intel.yaml"
        ;;
    cpu)
        COMPOSE_FILE="docker-compose.yaml"
        ;;
    *)
        echo "GPU_type inválido: $GPU_TYPE"
        echo "Usar: cpu, nvidia, amd, intel"
        exit 1
        ;;
esac

echo "Usando: $COMPOSE_FILE"

# Build y start
docker compose -f "$COMPOSE_FILE" build
docker compose -f "$COMPOSE_FILE" up -d

echo ""
echo "=== Desplegado ==="
echo "Accede a: http://localhost:8080"
echo ""
echo "Logs: docker compose -f $COMPOSE_FILE logs -f"
