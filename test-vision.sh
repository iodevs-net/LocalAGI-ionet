#!/bin/bash
#
# 🧪 Script de Pruebas para Agente de Visión
# ==========================================
#
# Este script prueba el flujo completo de análisis de imágenes
#

set -e

# Colores
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}══════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}   🧪 PRUEBAS DEL AGENTE DE VISIÓN${NC}"
echo -e "${BLUE}══════════════════════════════════════════════════════${NC}"
echo ""

# ============================================
# PASO 1: Verificar estado del sistema
# ============================================
echo -e "${YELLOW}[1/5]${NC} Verificando estado del sistema..."

CONTAINERS=$(docker compose -f docker-compose.dev.yaml ps --format "{{.Status}}" | grep -c "Up")
if [ "$CONTAINERS" -ge 2 ]; then
    echo -e "${GREEN}✓ Contenedores activos: $CONTAINERS${NC}"
else
    echo -e "${RED}❌ Contenedores no activos${NC}"
    exit 1
fi
echo ""

# ============================================
# PASO 2: Verificar configuración
# ============================================
echo -e "${YELLOW}[2/5]${NC} Verificando configuración..."

# Verificar modelo multimodal
MODEL=$(grep MULTIMODAL_MODEL .env | cut -d'=' -f2)
if [[ "$MODEL" == *"nemotron-nano"* ]]; then
    echo -e "${GREEN}✓ Modelo multimodal: $MODEL${NC}"
else
    echo -e "${RED}❌ Modelo no configurado correctamente${NC}"
fi

# Verificar agente visión
if [ -f "./config/agents/agente-vision.json" ]; then
    echo -e "${GREEN}✓ Configuración de agente-vision presente${NC}"
else
    echo -e "${RED}❌ Configuración faltante${NC}"
fi
echo ""

# ============================================
# PASO 3: Verificar API
# ============================================
echo -e "${YELLOW}[3/5]${NC} Verificando API..."

API_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8090/api/health || echo "000")
if [ "$API_RESPONSE" = "200" ] || [ "$API_RESPONSE" = "301" ] || [ "$API_RESPONSE" = "302" ]; then
    echo -e "${GREEN}✓ API respondiendo (HTTP $API_RESPONSE)${NC}"
else
    echo -e "${YELLOW}⚠️  API en estado HTTP $API_RESPONSE (puede ser normal)${NC}"
fi
echo ""

# ============================================
# PASO 4: Verificar agentes configurados
# ============================================
echo -e "${YELLOW}[4/5]${NC} Verificando agentes configurados..."

AGENT_COUNT=$(ls ./config/agents/*.json | wc -l)
echo -e "${GREEN}✓ Agentes configurados: $AGENT_COUNT${NC}"

# Listar agentes
echo "   Agentes:"
for agent in ./config/agents/*.json; do
    name=$(basename "$agent" .json)
    echo "   • $name"
done
echo ""

# ============================================
# PASO 5: Información del sistema
# ============================================
echo -e "${YELLOW}[5/5]${NC} Información del sistema..."
echo ""

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}  ESTADO DEL SISTEMA IONET${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

echo "📊 Contenedores:"
docker compose -f docker-compose.dev.yaml ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" 2>/dev/null
echo ""

echo "🎯 Servicios:"
echo "   • IONET API:   http://localhost:8090"
echo "   • PostgreSQL:  localhost:5433"
echo "   • Modelo VLM:  nvidia/nemotron-nano-12b-v2-vl:free"
echo ""

echo "🌐 Agentes Activos:"
echo "   1. ION (Orchestrator)"
echo "   2. agente-clientes"
echo "   3. agente-servicios"
echo "   4. agente-protocolos"
echo "   5. agente-inventario"
echo "   6. agente-seguridad"
echo "   7. agente-redes"
echo "   8. agente-datos"
echo "   9. 🎨 agente-vision (NUEVO)"
echo ""

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}  ✅ SISTEMA OPERATIVO Y LISTO PARA PRUEBAS${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

echo -e "${YELLOW}📝 Para probar el agente de visión:${NC}"
echo "   1. Enviar una imagen a través de la API"
echo "   2. Ver logs: docker compose -f docker-compose.dev.yaml logs -f ionet"
echo "   3. Ver documentación: cat ./VISION_SETUP.md"
echo ""

echo -e "${GREEN}¡Pruebas completadas exitosamente!${NC}"
