# IONET - Agente IA Interno para ionet.cl

Fork de [LocalAGI](https://github.com/mudler/LocalAGI) personalizado para uso interno de IONET.

## Propósito

Plataforma de agentes IA autohosteable que permite a los técnicos de IONET:
- Consultar protocolos y procedimientos sin leer manuales extensos
- Obtener ayuda personalizada de agentes especializados
- Acceder a información confidencial sin enviar datos a servicios externos
- Resolver dudas técnicas de manera autónoma

## Arquitectura

```
┌─────────────────────────────────────────────────────────────┐
│                    USUARIO (técnico)                        │
│                       ↑    ↓                                │
│                       │    │                                │
└───────────────────────┼────┼────────────────────────────────┘
                        │    │
                        ↓    │
┌─────────────────────────────────────────────────────────────┐
│              ION (Interlocutor/Orchestrator)                │
│  "Hola, necesito saber el protocolo para reinicio de CPE"  │
│                       ↓    ↑                                │
│     ┌────────────────┼────┴────────────────┐              │
│     ↓                ↓                     ↓              │
│ ┌──────────┐   ┌──────────┐    ┌──────────────────┐      │
│ │ Protocolos│   │  Redes   │    │   Base General   │      │
│ │  (RAG)   │   │  (RAG)   │    │     (RAG)       │      │
│ └──────────┘   └──────────┘    └──────────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

## Quick Start

### Desarrollo Local

```bash
# 1. Copiar variables de entorno
cp .env.example .env
# Editar .env con tus credenciales

# 2. Iniciar servicios
docker compose -f docker-compose.dev.yaml up

# 3. Acceder a la UI
open http://localhost:8080
```

### Producción (Hetzner)

```bash
# 1. Subir al servidor
git clone <repo> /opt/ionet
cd /opt/ionet

# 2. Configurar variables
cp .env.example .env
nano .env  # Llenar credenciales reales

# 3. Iniciar producción
docker compose -f docker-compose.prod.yaml up -d

# 4. Acceder via Cloudflare
# https://agentes.ionet.cl
```

## Estructura del Proyecto

```
├── Dockerfile              # Unificado: dev + prod
├── docker-compose.dev.yaml  # Desarrollo local
├── docker-compose.prod.yaml # Producción Hetzner
├── config/
│   └── agents/              # Configuración de agentes
│       ├── ion.json                    # Orchestrator
│       ├── agente-protocolos.json      # Protocolos
│       ├── agente-redes.json           # Redes
│       └── agente-base.json            # Info general
├── scripts/
│   └── m365/               # Sync SharePoint/OneDrive
│       ├── sync.go
│       └── README.md
└── .env.example            # Template de configuración
```

## Configuración

### Variables Requeridas

```bash
# LLM Provider (requerido)
OPENAI_BASE_URL=https://api.minimax.io/v1
OPENAI_API_KEY=tu-api-key
MODEL_NAME=MiniMax-M2.7
MULTIMODAL_MODEL=MiniMax-M2.7

# Seguridad
LOCALAGI_API_KEYS=key1,key2,key3

# Producción
POSTGRES_PASSWORD=password-seguro
```

### M365 Sync (opcional)

Para sincronizar documentos desde SharePoint/OneDrive:

```bash
M365_CLIENT_ID=<azure-app-id>
M365_CLIENT_SECRET=<azure-secret>
M365_TENANT_ID=<azure-tenant-id>
M365_SITE_URL=https://ionet.sharepoint.com/sites/documentos
```

Ver [scripts/m365/README.md](scripts/m365/README.md) para configuración detallada.

## Importar Agentes

1. Acceder a la UI de administración
2. Ir a **Agents** > **Create Agent**
3. Importar desde JSON:
   - `config/agents/agente-protocolos.json`
   - `config/agents/agente-redes.json`
   - `config/agents/agente-base.json`
   - `config/agents/ion.json` (importar al final)

Ver [config/agents/README.md](config/agents/README.md) para detalles.

## SSH al Contenedor (Producción)

```bash
# Desde el servidor Hetzner
ssh root@hetzner-ip

# Conectar al contenedor directamente
docker exec -it ionet bash

# O vía SSH directo al contenedor (puerto 2222)
ssh ionet@agentes.ionet.cl -p 2222
# Password: ionet (cambiar en producción!)
```

## Docker Modes

| Modo | Target | Uso | Puerto SSH |
|------|--------|-----|-----------|
| `dev` | `dev` | Desarrollo local | No |
| `prod` | `prod` | Producción Hetzner | 2222 |

```bash
# Desarrollo
docker build --build-arg MODE=dev -t ionet:dev .
docker compose -f docker-compose.dev.yaml up

# Producción
docker build --build-arg MODE=prod -t ionet:prod .
docker compose -f docker-compose.prod.yaml up -d
```

## Troubleshooting

### "Agent not found" al llamar otros agentes

Verificar que todos los agentes están importados y activos.

### RAG no devuelve resultados

1. Verificar que los documentos están en `/rag/raw/`
2. Verificar que el motor RAG está configurado (`VECTOR_ENGINE`)
3. Reiniciar el servicio para re-indexar

### Error de autenticación M365

1. Verificar que la app de Azure tiene permisos `Sites.Read.All`
2. Verificar que se dio **Admin consent**
3. Comprobar que los documentos están compartidos con la app

## Licencia

MIT License - Ver [LICENSE](../LICENSE)

## Créditos

- Basado en [LocalAGI](https://github.com/mudler/LocalAGI) por mudler
- Fork personalizado para IONET Chile
