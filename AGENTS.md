# Repository Guidelines - IONET

**Contexto**: IONET es una empresa de soporte informático chilena con 10 técnicos. ION es el agente principal (orchestrator) que atiende a los trabajadores y deriva a agentes especializados.

---

## Contexto de Negocio

```
┌─────────────────────────────────────────────────────────────────┐
│                    10 TÉCNICOS IONET                            │
│         (usuarios que consultan a ION)                         │
└─────────────────────────┬───────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────┐
│  ION (Orchestrator)                                             │
│  - Clasifica consultas                                          │
│  - Deriva al agente correcto                                    │
│  - Sintetiza respuestas                                         │
└─────────────────────────┬───────────────────────────────────────┘
                          │
        ┌─────────────────┼─────────────────┐
        ▼                 ▼                 ▼
   ┌─────────┐      ┌──────────┐      ┌──────────┐
   │Agentes  │      │Agentes   │      │Agentes   │
   │Especial.│      │Técnicos  │      │Datos      │
   └─────────┘      └──────────┘      └──────────┘
```

## Agentes Definidos

| # | Agente | Dominio | Quand Use |
|---|--------|---------|-----------|
| 1 | **ION** | Orchestrator | Punto de entrada, clasificación, síntesis |
| 2 | **agente-clientes** | Clientes, onboarding, inducciones | Nombres, contratos, usuarios nuevos |
| 3 | **agente-servicios** | Service desk, tickets, SLAs | Solicitudes, seguimiento, catálogo |
| 4 | **agente-protocolos** | Procedimientos, políticas | Procesos, cómo hacer |
| 5 | **agente-inventario** | HW, SW, licencias, activos | Equipos, software, licencias |
| 6 | **agente-seguridad** | Ciberseguridad, incidentes | Seguridad, vulnerabilidades, compliance |
| 7 | **agente-redes** | Redes, servidores, infra | Conectividad, IPs, configuración |
| 8 | **agente-datos** | Documentos, backups, retención | Archivos, versionado, backup |

### Personalidades Comunes

Todos los agentes comparten:
- **Nivel**: ULTRA SENIOR (20+ años experiencia)
- **Tono**: Profesional, conciso, accionable
- **Formato**: Siempre结论 primera, luego análisis
- **Derivan**: Cuando la consulta requiere otro dominio o escalamiento humano

---

## Arquitectura & Data Flow

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Connectors │────▶│    ION     │────▶│  Agentes    │
│ (Discord,   │     │ (Orches-   │     │ Especializ. │
│  Telegram, │     │  trator)   │     │             │
│  Slack...)  │     └─────────────┘     └─────────────┘
└─────────────┘            │                    │
                          ▼                    ▼
                   ┌─────────────┐     ┌─────────────┐
                   │    State    │     │     RAG      │
                   │  (Memory)   │     │  (KB docs)   │
                   └─────────────┘     └─────────────┘
```

### Rutas de Derivación (por palabra clave)

| Si el usuario menciona... | Derivar a... |
|---------------------------|--------------|
| cliente, contrato, onboarding, usuario nuevo | agente-clientes |
| ticket, SLA, solicitud, seguimiento, service desk | agente-servicios |
| procedimiento, proceso, política, cómo hacer | agente-protocolos |
| equipo, licencia, software, hardware, inventario | agente-inventario |
| seguridad, hack, virus, breach, vulnerability | agente-seguridad |
| red, IP, servidor, conexión, VPN, internet | agente-redes |
| documento, backup, archivo, versión, retención | agente-datos |
| contacto, organigrama, FAQ, info general | agente-base |

---

## Desarrollo de Agentes

### Estructura de Configuración

Los agentes se definen en `config/agents/*.json`:

```json
{
  "name": "agente-nombre",
  "description": "Descripción breve",
  "model": "${MODEL_NAME}",
  "system_prompt": "## IDENTIDAD: ...",
  "enable_kb": true,
  "kb_results": 8,
  "kb_auto_search": true,
  "kb_as_tools": true,
  "long_term_memory": false,
  "enable_reasoning": true,
  "actions": [{ "name": "call_agents", "config": "{}" }]
}
```

### Campos Obligatorios

| Campo | Descripción |
|-------|-------------|
| `name` | Identificador único (kebab-case) |
| `description` | Descripción breve para ION |
| `system_prompt` | Personalidad, contexto, formato de respuestas |
| `model` | Usar `${MODEL_NAME}` (variable de entorno) |
| `enable_kb` | Habilitar búsqueda en base de conocimiento |
| `kb_results` | Cantidad de resultados a retornar |

### Sistema Prompt - Plantilla

```markdown
## IDENTIDAD: [nombre] - [título]

Eres [nombre], [descripción de expertise].

## DOMINIO DE EXPERTISE

### 1. [Área 1]
### 2. [Área 2]
### 3. [Área 3]

## CONSULTAS QUE RESPONDES

- "[tipo de pregunta 1]"
- "[tipo de pregunta 2]"

## FORMATO DE RESPUESTA

[Plantilla específica para el dominio]

## INTEGRACIÓN RAG

### Fuentes Indexadas
- [Lista de documentos]

### Búsqueda por
- [Criterios]

## REGLAS OPERATIVAS

### Tiempo de Respuesta
- Simple: < X segundos
- Complejo: < Y segundos

### Escalamiento
- [Cuándo derivar a humano]

### Deriva
- [Cuándo derivar a otro agente]
```

---

## Desarrollo de Acciones

### Inter-Agente: `call_agents`

La acción principal para que ION derive a otros agentes:

```json
{
  "name": "call_agents",
  "config": "{}"
}
```

Permite a ION invocar a cualquier agente especializado.

### Acciones Disponibles

| Acción | Descripción |
|--------|-------------|
| `search` | Búsqueda web |
| `scrape` | Web scraping |
| `memory` | Memoria persistente |
| `browse` | Navegación web |
| `call_agents` | Invocar otros agentes |
| `sendmail` | Enviar email |
| `shell` | Ejecutar comandos |
| `webhook` | Webhook HTTP |
| `wikipedia` | Consulta Wikipedia |
| `genimage` | Generar imagen |
| `genpdf` | Generar PDF |
| `githubissue*` | Gestión de issues |
| `githubpr*` | Gestión de PRs |

### Agregar Nueva Acción

1. Crear en `services/actions/nueva_accion.go`
2. Implementar interfaz `Action`:
   ```go
   type Action interface {
       Name() string
       Execute(ctx context.Context, opts ...ActionOption) (string, error)
   }
   ```
3. Registrar en `services/actions.go`

---

## Directorios Clave

| Directorio | Propósito |
|------------|-----------|
| `config/agents/` | Definiciones JSON de agentes |
| `core/agent/` | Motor de ejecución de agentes |
| `core/state/` | Estado, pool, memoria |
| `services/connectors/` | Discord, Telegram, Slack, etc. |
| `services/actions/` | Acciones ejecutables |
| `services/prompts/` | Plantillas de prompts |
| `pool/rag/` | Documentos para RAG |

---

## Comandos de Desarrollo

```bash
# Build
make build

# Run local
make run

# Tests
make tests

# Docker dev
docker compose -f docker-compose.dev.yaml up

# Docker prod
docker compose -f docker-compose.prod.yaml up -d
```

### Variables de Entorno

```bash
# LLM (MiniMax para IONET)
OPENAI_BASE_URL=https://api.minimax.io/v1
OPENAI_API_KEY=tu-api-key
MODEL_NAME=MiniMax-M2.7

# Seguridad
LOCALAGI_API_KEYS=key1,key2,key3

# RAG
VECTOR_ENGINE=chromem  # dev
EMBEDDING_MODEL=sentence-transformers/all-MiniLM-L6-v2
```

---

## Testing

Framework: **Ginkgo v2 + Gomega**

```go
import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

func TestAgent(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Agent test suite")
}
```

Requiere Docker services:
```bash
docker compose up -d --build
LOCALAGI_MODEL="..." LOCALAI_API_URL="..." go test ./...
```

---

## Escalamiento y Derivation Rules

### Escalamiento a Humano (ION)

Deriva cuando:
- P1/P2 Incidentes de seguridad
- Decisiones contractuales
- Incidentes con cliente VIP
- Situaciones legales/compliance
- Cambios de configuración de producción
- Más de 3 agentes requieren coordinación

### Escalamiento a Humano (agente-seguridad)

**SIEMPRE** escala:
- Ransomware activo
- Brecha de datos confirmada
- Ataque en progreso
- Compromiso de credenciales privilegiadas
- Requerimiento legal/regulatorio

---

## Consideraciones para Desarrolladores

1. **Todo debe estar en español chileno** - Los 10 técnicos esperan respuestas en español natural de Chile

2. **Concisión es clave** - "Cero verbosidad, cada palabra justifica su existencia"

3. **Conclusión primero** - Nunca会让用户等待，用户永远不应该等待答案的开头

4. **RAG es vital** - Los agentes especializados deben buscar en documentación interna primero

5. **Derivar干净的** - Si no es tu dominio, deriva inmediatamente al agente correcto

6. **Formato consistente** - Cada agente tiene plantillas de respuesta específicas, respetarlas

7. **Métricas claras** - Tiempo de respuesta, precisión de clasificación, escalamientos apropiados