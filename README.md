<p align="center">
  <img src="./webui/react-ui/public/logo_1.png" alt="IONET Logo" width="280"/>
</p>

<h1 align="center">IONET</h1>
<h3 align="center"><em>Tu IA. Tu Hardware. Tus Reglas.</em></h3>

<div align="center">

[![Go Report Card](https://goreportcard.com/badge/github.com/mudler/LocalAGI)](https://goreportcard.com/report/github.com/mudler/LocalAGI)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub stars](https://img.shields.io/github/stars/mudler/LocalAGI)](https://github.com/mudler/LocalAGI/stargazers)
[![GitHub issues](https://img.shields.io/github/issues/mudler/LocalAGI)](https://github.com/mudler/LocalAGI/issues)
[![Go Version](https://img.shields.io/badge/Go-1.26%2B-blue)](https://golang.org/)
[![React](https://img.shields.io/badge/React-18-green)](https://reactjs.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue)](https://www.docker.com/)

</div>

---

## 📖 Tabla de Contenidos

1. [Descripción General](#-descripción-general)
2. [Características Principales](#-características-principales)
3. [Arquitectura del Sistema](#-arquitectura-del-sistema)
4. [Quickstart](#-quickstart)
5. [Configuración](#-configuración)
6. [Estructura del Proyecto](#-estructura-del-proyecto)
7. [Sistema de Agentes](#-sistema-de-agentes)
8. [Acciones Disponibles](#-acciones-disponibles)
9. [Conectores](#-conectores)
10. [API REST](#-api-rest)
11. [Uso como Librería](#-uso-como-librería)
12. [Desarrollo](#-desarrollo)
13. [Hardware Soportado](#-hardware-soportado)
14. [Casos de Uso](#-casos-de-uso)
15. [Extendiendo IONET](#-extendiendo-ionet)
16. [FAQ](#-faq)
17. [Familia LocalAI](#-familia-localai)
18. [Capturas de Pantalla](#-capturas-de-pantalla)
19. [Licencia](#-licencia)

---

## 📝 Descripción General

**IONET** es una plataforma de agentes IA autohosteable y personalizable que se ejecuta **100% localmente** en el hardware del usuario. Diseñada para usuarios que valoran su privacidad, no requiere APIs en la nube ni suscripciones externas.

Con IONET puedes crear:

- 🤖 **Agentes IA personalizados** sin escribir código
- 🔗 **Conectores** para integrar con servicios existentes
- 📚 **Base de conocimiento RAG** para memoria a largo plazo
- ⚙️ **Automatizaciones** con tareas periódicas tipo cron
- 🎨 **Acciones personalizadas** escritas en Go (interpretadas, sin compilación)
- 🤝 **Equipos de agentes cooperativos** que trabajan juntos

### ¿Por qué IONET?

| Característica | IONET | Soluciones en la nube |
|----------------|-------|----------------------|
| Privacidad | ✅ 100% local | ❌ Datos en servidores externos |
| Costo | ✅ Gratuito (hardware propio) | ❌ Suscripciones mensuales |
| Control | ✅ Total sobre tus datos | ❌ Dependencia del proveedor |
| Flexible | ✅ Código abierto y extensible | ❌ Limitado por el servicio |

---

## ✨ Características Principales

### 🛡️ Privacidad Total
- Todos los datos se procesan localmente
- Sin telemetría ni conexiones externas
- Tus API keys solo se usan para modelos locales

### 🎛️ Agentes Sin Código
- Configuración visual mediante interfaz web intuitiva
- Sistema de plantillas para casos de uso comunes
- Exportación e importación de configuraciones

### 🤖 Teamwork de Agentes
- Creación de equipos cooperativos de agentes
- Comunicación entre agentes para tareas complejas
- Especialización de agentes por función

### 📡 Conectores Integrados
- **Mensajería**: Discord, Slack, Telegram, IRC, Matrix, Email
- **Desarrollo**: GitHub Issues, GitHub PRs
- **Redes Sociales**: Twitter/X
- Configuración simple mediante JSON

### 🧠 Memoria Inteligente
- **Memoria de corto plazo**: Contexto conversacional
- **Memoria de largo plazo**: Base de conocimiento RAG con búsqueda semántica
- **Resumen automático**: Síntesis de conversaciones anteriores

### 🔄 Tareas Periódicas
- Sintaxis cron para programar tareas
- Recordatorios únicos y recurrentes
- Ejecución en background sin intervención

### 🖼️ Soporte Multimodal
- Visión por computadora
- Generación de imágenes
- Generación de audio
- Análisis de documentos PDF

### 🔧 Acciones Personalizadas
- Scripts Go interpretados (sin compilación)
- Carga automática de acciones desde directorio
- Biblioteca extensible de funciones

### 🛠️ MCP (Model Context Protocol)
- Soporte para servidores MCP locales (STDIO)
- Soporte para servidores MCP remotos (HTTP)
- Integración con ecosistema MCP existente

### 📚 Skills (Habilidades)
- Gestión visual de habilidades reutilizables
- Sincronización con repositorios Git
- Inyección automática en agentes

---

## 🏗️ Arquitectura del Sistema

```
┌─────────────────────────────────────────────────────────────────┐
│                         IONET                                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────┐    ┌─────────────────────────────────┐   │
│  │   Web UI (React)│    │      API Layer (Fiber v2)       │   │
│  │   localhost:8080│◄──►│   REST API + SSE + WebSocket    │   │
│  └─────────────────┘    └─────────────────────────────────┘   │
│                                    │                           │
│                                    ▼                           │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    Agent Core                          │   │
│  │  ┌─────────┐  ┌──────────┐  ┌───────────┐  ┌────────┐ │   │
│  │  │ Agent   │  │Scheduler │  │ Connector │  │ Action │ │   │
│  │  │ Engine  │  │(Cron)    │  │ Manager   │  │ Loader │ │   │
│  │  └─────────┘  └──────────┘  └───────────┘  └────────┘ │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                    │                           │
│                                    ▼                           │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                  State & Memory                         │   │
│  │  ┌─────────┐  ┌──────────┐  ┌───────────┐              │   │
│  │  │ JSON    │  │RAG (Bleve│  │PostgreSQL │              │   │
│  │  │ State   │  │/Chroma)  │  │(Opcional) │              │   │
│  │  └─────────┘  └──────────┘  └───────────┘              │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                    │                           │
│                                    ▼                           │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │              LLM Provider (OpenAI-compatible)           │   │
│  │  ┌──────────┐  ┌────────────┐  ┌──────────────────┐  │   │
│  │  │ LocalAI  │  │ MiniMax    │  │ Otros proveedores │  │   │
│  │  └──────────┘  └────────────┘  └──────────────────┘  │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Stack Tecnológico

| Componente | Tecnología | Versión |
|------------|------------|---------|
| Backend | Go | 1.26+ |
| Framework HTTP | Fiber | v2 |
| Frontend | React | 18 |
| Bundler | Vite | 7.x |
| Package Manager | Bun | 1.2+ |
| CLI | Cobra | v1.10 |
| Base de datos (RAG) | PostgreSQL / Bleve | - |
| Estado | JSON files | - |
| Containers | Docker | - |

---

## 🚀 Quickstart

### Prerrequisitos

- Docker y Docker Compose instalados
- Git
- (Opcional) GPU NVIDIA/AMD/Intel con drivers apropiados

### 1. Clonar el Repositorio

```bash
git clone https://github.com/mudler/LocalAGI.git
cd LocalAGI
```

### 2. Configurar Variables de Entorno

```bash
# Copiar archivo de ejemplo
cp .env.example .env

# Editar con tu editor preferido
nano .env
```

### 3. Iniciar con Docker

```bash
# Configuración CPU (por defecto)
docker compose up

# Configuración NVIDIA GPU
docker compose -f docker-compose.nvidia.yaml up

# Configuración AMD GPU
docker compose -f docker-compose.amd.yaml up

# Configuración Intel GPU
docker compose -f docker-compose.intel.yaml up
```

### 4. Acceder a la Interfaz Web

Abre tu navegador en: **[http://localhost:8080](http://localhost:8080)**

### Videos Tutoriales

| Video | Descripción | Enlace |
|-------|-------------|--------|
| Crear un Agente Básico | Tutorial paso a paso | [YouTube](https://youtu.be/HtVwIxW3ePg) |
| Observabilidad | Monitoreo de agentes en tiempo real | [YouTube](https://youtu.be/v82rswGJt_M) |
| Filtros y Triggers | Automatización avanzada | [YouTube](https://youtu.be/d_we-AYksSw) |
| RAG y Matrix | Base de conocimiento | [YouTube](https://youtu.be/2Xvx78i5oBs) |

---

## ⚙️ Configuración

### Variables de Entorno

| Variable | Descripción | Valor por defecto |
|----------|-------------|-------------------|
| `OPENAI_API_KEY` | Clave API del proveedor LLM | - |
| `OPENAI_BASE_URL` | URL base del API compatible con OpenAI | `https://api.minimax.io/v1` |
| `MODEL_NAME` | Nombre del modelo de texto principal | `MiniMax-M2.7` |
| `MULTIMODAL_MODEL` | Modelo para capacidades multimodales | `MiniMax-M2.7` |
| `LOCALAGI_LLM_API_URL` | URL del servidor LLM | - |
| `LOCALAGI_LLM_API_KEY` | Clave API para autenticación | - |
| `LOCALAGI_TIMEOUT` | Timeout para requests | `5m` |
| `LOCALAGI_STATE_DIR` | Directorio para estado de agentes | `/pool` |
| `LOCALAGI_BASE_URL` | URL base de la aplicación | `http://localhost:3000` |
| `LOCALAGI_SSHBOX_URL` | URL de SSHBox (user:pass@host:port) | - |
| `LOCALAGI_ENABLE_CONVERSATIONS_LOGGING` | Habilitar logging de conversaciones | `false` |
| `LOCALAGI_API_KEYS` | Lista de API keys separadas por coma | - |
| `LOCALAGI_CUSTOM_ACTIONS_DIR` | Directorio con acciones Go personalizadas | - |

### Variables para Base de Conocimiento (Opcional)

| Variable | Descripción | Valor por defecto |
|----------|-------------|-------------------|
| `VECTOR_ENGINE` | Motor de vectores (`chromem`, `postgres`) | `chromem` |
| `EMBEDDING_MODEL` | Modelo para embeddings | - |
| `DATABASE_URL` | URL de PostgreSQL (para VECTOR_ENGINE=postgres) | - |
| `COLLECTION_DB_PATH` | Path de la base de datos de colecciones | - |
| `FILE_ASSETS` | Directorio para archivos de assets | - |

### Ejemplo de `.env` Completo

```bash
# =============================================================================
# IONET Configuration
# =============================================================================

# API Configuration - MiniMax (configuración por defecto)
OPENAI_API_KEY=tu-minimax-api-key
OPENAI_BASE_URL=https://api.minimax.io/v1
MODEL_NAME=MiniMax-M2.7
MULTIMODAL_MODEL=MiniMax-M2.7

# Alternative: LocalAI (para modelos locales)
# OPENAI_BASE_URL=http://localai:8080
# OPENAI_API_KEY=localai
# MODEL_NAME=gemma-3-4b-it-qat

# PostgreSQL Knowledge Base (opcional)
# VECTOR_ENGINE=postgres
# DATABASE_URL=postgresql://user:pass@postgres:5432/db?sslmode=disable
# EMBEDDING_MODEL=granite-embedding-107m-multilingual

# MCP Servers (opcional)
# LOCALAGI_MCP_SERVERS=[{"url":"https://mcp.server.com","token":"token"}]

# SSHBox (opcional)
# LOCALAGI_SSHBOX_URL=root:root@sshbox:22

# General Settings
LOCALAGI_TIMEOUT=5m
LOCALAGI_ENABLE_CONVERSATIONS_LOGGING=false
```

---

## 📁 Estructura del Proyecto

```
LocalAGI/
├── cmd/                          # Punto de entrada CLI
│   ├── root.go                   # Comando raíz (Cobra)
│   ├── serve.go                  # Comando serve (servidor web)
│   ├── agent.go                  # Comandos de gestión de agentes
│   ├── agent_run.go             # Ejecución de agente individual
│   └── env.go                   # Gestión de variables de entorno
│
├── core/                         # Núcleo del sistema
│   ├── agent/                    # Motor de agentes
│   ├── conversations/            # Seguimiento de conversaciones
│   ├── scheduler/                # Programador de tareas (cron)
│   │   ├── interfaces.go
│   │   ├── json_store.go
│   │   ├── scheduler.go
│   │   └── task.go
│   ├── state/                    # Estado de agentes
│   │   ├── config.go
│   │   ├── internal.go
│   │   ├── pool.go
│   │   └── compaction.go
│   ├── sse/                      # Server-Sent Events
│   ├── action/                   # Sistema de acciones
│   │   ├── custom.go            # Acciones personalizadas Go
│   │   ├── reminder.go          # Recordatorios
│   │   ├── state.go             # Acciones de estado
│   │   └── newconversation.go
│   └── types/                    # Tipos y estructuras
│       ├── actions.go
│       ├── conversation.go
│       ├── filters.go
│       ├── job.go
│       ├── prompts.go
│       ├── result.go
│       └── state.go
│
├── services/                     # Servicios y acciones
│   ├── actions.go               # Registro de acciones
│   ├── common.go                # Utilidades comunes
│   ├── connectors.go            # Gestor de conectores
│   ├── filters.go               # Filtros de mensajes
│   ├── prompts.go               # Prompts dinámicos
│   ├── skills/                  # Sistema de skills
│   │   ├── service.go
│   │   └── prompt.go
│   ├── actions/                 # Acciones predefinidas
│   │   ├── browse.go           # Navegación web
│   │   ├── search.go          # Búsqueda web
│   │   ├── scrape.go          # Extracción de contenido
│   │   ├── wikipedia.go       # Consulta a Wikipedia
│   │   ├── githubissue*.go    # Gestión de issues GitHub
│   │   ├── githubpr*.go       # Gestión de PRs GitHub
│   │   ├── githubrepository*.go # Gestión de repositorios
│   │   ├── memory.go          # Operaciones de memoria RAG
│   │   ├── sendmail.go        # Envío de emails
│   │   ├── sendtelegrammessage.go # Mensajes Telegram
│   │   ├── twitter_post.go    # Publicar en Twitter
│   │   ├── webhook.go         # Llamadas webhook
│   │   ├── shell.go           # Comandos shell
│   │   ├── counter.go         # Contador
│   │   ├── pikvm.go           # Control PiKVM
│   │   ├── genimage.go        # Generación de imágenes
│   │   ├── gensong.go         # Generación de audio
│   │   ├── genpdf.go          # Generación de PDFs
│   │   └── callagents.go      # Llamadas entre agentes
│   ├── connectors/             # Conectores de servicios
│   │   ├── telegram.go        # Integración Telegram
│   │   ├── discord.go         # Integración Discord
│   │   ├── slack.go          # Integración Slack
│   │   ├── githubissue.go    # GitHub Issues
│   │   ├── githubpr.go       # GitHub PRs
│   │   ├── twitter.go        # Twitter/X
│   │   ├── irc.go            # IRC
│   │   ├── matrix.go         # Matrix
│   │   ├── email.go          # Email (SMTP/IMAP)
│   │   └── common/           # Utilidades compartidas
│   └── filters/               # Filtros de contenido
│       ├── classifier.go
│       └── regex.go
│
├── webui/                       # Interfaz web
│   ├── app.go                  # Aplicación Fiber
│   ├── options.go             # Opciones de configuración
│   ├── elements.go            # Elementos UI
│   ├── collections_*.go       # Backend de colecciones RAG
│   ├── collections/           # Módulo de colecciones
│   │   ├── types.go
│   │   ├── state.go
│   │   ├── rag_provider.go
│   │   └── inprocess.go
│   ├── public/                # Archivos estáticos
│   │   ├── logo_1.png
│   │   └── css/
│   └── react-ui/              # Frontend React
│       ├── package.json
│       ├── vite.config.js
│       ├── index.html
│       ├── src/
│       │   ├── main.jsx
│       │   ├── App.jsx
│       │   ├── router.jsx
│       │   ├── components/
│       │   │   ├── AgentForm.jsx
│       │   │   ├── ActionForm.jsx
│       │   │   ├── ConnectorForm.jsx
│       │   │   ├── ConfigForm.jsx
│       │   │   ├── FilterForm.jsx
│       │   │   ├── Sidebar.jsx
│       │   │   ├── ThemeToggle.jsx
│       │   │   └── agent-form-sections/
│       │   ├── pages/
│       │   │   ├── Home.jsx
│       │   │   ├── Chat.jsx
│       │   │   ├── AgentsList.jsx
│       │   │   ├── AgentSettings.jsx
│       │   │   ├── AgentStatus.jsx
│       │   │   ├── CreateAgent.jsx
│       │   │   ├── ImportAgent.jsx
│       │   │   ├── GroupCreate.jsx
│       │   │   ├── ActionsPlayground.jsx
│       │   │   ├── Knowledge.jsx
│       │   │   ├── Skills.jsx
│       │   │   └── SkillEdit.jsx
│       │   ├── hooks/
│       │   │   ├── useAgent.js
│       │   │   ├── useChat.js
│       │   │   └── useSSE.js
│       │   └── contexts/
│       │       └── ThemeContext.jsx
│       └── public/
│
├── tests/                       # Tests E2E
│   └── e2e/
│       ├── e2e_suite_test.go
│       └── e2e_test.go
│
├── example/                     # Ejemplos
│   └── custom_actions/
│       └── hello.go           # Ejemplo de acción personalizada
│
├── docker-compose.yaml         # Config CPU
├── docker-compose.nvidia.yaml  # Config NVIDIA GPU
├── docker-compose.amd.yaml     # Config AMD GPU
├── Dockerfile.webui            # Imagen Docker principal
├── Dockerfile.sshbox           # Imagen Docker para SSH
├── Dockerfile.realtimesst      # Imagen para SSE en tiempo real
├── Makefile                    # Targets de construcción
├── go.mod                      # Dependencias Go
├── go.sum                      # Checksums Go
└── README.md                   # Este archivo
```

---

## 🤖 Sistema de Agentes

### Conceptos Fundamentales

#### 1. Agente
Un agente es una instancia IA autónoma con:
- **Configuración**: Modelo, prompts, conectores, acciones
- **Estado**: Memoria de corto plazo, historial, progreso
- **Ciclo de vida**: Crear → Iniciar → Pausar → Detener → Eliminar

#### 2. Características del Agente

```json
{
  "name": "mi-agente",
  "model": "MiniMax-M2.7",
  "multimodal_model": "MiniMax-M2.7",
  "system_prompt": "Eres un asistente útil especializado en...",
  "identity_guidance": "Descripción de la personalidad del agente",
  "enable_kb": true,
  "enable_reasoning": true,
  "enable_planning": true,
  "long_term_memory": true,
  "summary_long_term_memory": false,
  "periodic_runs": "0 */6 * * *",
  "permanent_goal": "Ayudar con tareas de desarrollo de software",
  "can_stop_itself": false,
  "initiate_conversations": true,
  "kb_results": 5,
  "hud": true
}
```

#### 3. Tipos de Memoria

| Tipo | Descripción | Persistencia |
|------|-------------|--------------|
| **Corto plazo** | Contexto actual de la conversación | En memoria |
| **Largo plazo** | Base de conocimiento RAG | Vector DB |
| **Resumen** | Síntesis de conversaciones previas | Archivos JSON |
| **Estado interno** | `NowDoing`, `DoingNext`, `Goal`, `Memories` | JSON |

### Configuración de Agente (Referencia Completa)

```bash
# Obtener esquema de configuración
curl -X GET "http://localhost:8080/api/meta/agent/config"
```

#### Campos Disponibles

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `name` | string | Identificador único del agente |
| `model` | string | Modelo LLM principal |
| `multimodal_model` | string | Modelo para visión |
| `system_prompt` | string | Prompt del sistema |
| `identity_guidance` | string | Guía de personalidad |
| `enable_kb` | boolean | Habilitar base de conocimiento |
| `enable_reasoning` | boolean | Habilitar razonamiento |
| `enable_planning` | boolean | Habilitar planificación |
| `long_term_memory` | boolean | Memoria de largo plazo |
| `summary_long_term_memory` | boolean | Resumen de memoria |
| `periodic_runs` | string | Expresión cron |
| `permanent_goal` | string | Objetivo permanente |
| `can_stop_itself` | boolean | Puede detenerse solo |
| `initiate_conversations` | boolean | Inicia conversaciones |
| `kb_results` | int | Resultados de búsqueda KB |
| `hud` | boolean | Mostrar HUD de estado |
| `standalone_job` | boolean | Tarea independiente |
| `random_identity` | boolean | Identidad aleatoria |
| `actions` | array | Acciones habilitadas |
| `connectors` | object | Configuración de conectores |
| `filters` | object | Filtros de mensajes |
| `mcp_servers` | array | Servidores MCP |
| `skills` | array | Habilidades habilitadas |

---

## ⚡ Acciones Disponibles

### Acciones de GitHub

| Acción | Descripción | Campos |
|--------|-------------|--------|
| `github-issue-opener` | Crear issue | `owner`, `repo`, `title`, `body` |
| `github-issue-closer` | Cerrar issue | `owner`, `repo`, `issue_number` |
| `github-issue-comment` | Comentar en issue | `owner`, `repo`, `issue_number`, `comment` |
| `github-issue-edit` | Editar issue | `owner`, `repo`, `issue_number`, `title`, `body` |
| `github-issue-labeler` | Añadir labels | `owner`, `repo`, `issue_number`, `labels` |
| `github-issue-reader` | Leer issue | `owner`, `repo`, `issue_number` |
| `github-issue-search` | Buscar issues | `query`, `owner`, `repo` |
| `github-pr-opener` | Crear PR | `owner`, `repo`, `title`, `body`, `head`, `base` |
| `github-pr-closer` | Cerrar PR | `owner`, `repo`, `pr_number` |
| `github-pr-comment` | Comentar en PR | `owner`, `repo`, `pr_number`, `comment` |
| `github-pr-reader` | Leer PR | `owner`, `repo`, `pr_number` |
| `github-pr-reviewer` | Revisar PR | `owner`, `repo`, `pr_number`, `event`, `comment` |
| `github-repository-get-content` | Obtener archivo | `owner`, `repo`, `path`, `ref` |
| `github-repository-get-all-content` | Listar contenido | `owner`, `repo`, `path` |
| `github-repository-list-files` | Listar archivos | `owner`, `repo`, `path` |
| `github-repository-create-content` | Crear archivo | `owner`, `repo`, `path`, `content`, `message` |
| `github-repository-update-content` | Actualizar archivo | `owner`, `repo`, `path`, `content`, `message` |
| `github-repository-readme` | Leer README | `owner`, `repo` |
| `github-repository-search-files` | Buscar archivos | `owner`, `repo`, `query` |

### Acciones Web

| Acción | Descripción | Campos |
|--------|-------------|--------|
| `web-search` | Búsqueda web | `query` |
| `web-scrape` | Extraer contenido | `url` |
| `web-browse` | Navegar URL | `url` |
| `web-wikipedia` | Buscar en Wikipedia | `query` |

### Acciones Multimedia

| Acción | Descripción | Campos |
|--------|-------------|--------|
| `generate-image` | Generar imagen | `prompt`, `model` |
| `generate-song` | Generar audio | `prompt`, `duration` |
| `generate-pdf` | Generar PDF | `content`, `title` |

### Acciones de Comunicación

| Acción | Descripción | Campos |
|--------|-------------|--------|
| `send-mail` | Enviar email | `to`, `subject`, `body` |
| `send-telegram-message` | Mensaje Telegram | `chat_id`, `text` |
| `twitter-post` | Publicar tweet | `text` |
| `webhook-call` | Llamada webhook | `url`, `method`, `headers`, `body` |

### Acciones de Sistema

| Acción | Descripción | Campos |
|--------|-------------|--------|
| `shell-command` | Ejecutar comando | `command`, `timeout` |
| `counter` | Incrementar contador | `name`, `amount` |
| `pikvm-command` | Comando PiKVM | `command` |

### Acciones de Memoria

| Acción | Descripción | Campos |
|--------|-------------|--------|
| `memory-add` | Añadir a memoria | `text`, `collection` |
| `memory-search` | Buscar en memoria | `query`, `collection`, `limit` |
| `memory-remove` | Eliminar de memoria | `id`, `collection` |

### Acciones de Recordatorios

| Acción | Descripción | Campos |
|--------|-------------|--------|
| `set-recurring-reminder` | Recordatorio recurrente | `message`, `cron_expr` |
| `set-onetime-reminder` | Recordatorio único | `message`, `delay` |

### Acciones de Agentes

| Acción | Descripción | Campos |
|--------|-------------|--------|
| `call-agent` | Llamar a otro agente | `agent_name`, `message` |
| `new-conversation` | Nueva conversación | (ninguno) |
| `no-reply` | No responder | (ninguno) |

### Acciones de Estado

| Acción | Descripción | Campos |
|--------|-------------|--------|
| `set-doing` | Establecer estado | `task` |
| `set-goal` | Establecer objetivo | `goal` |
| `add-memory` | Añadir a estado | `memory` |
| `add-to-history` | Añadir al historial | `item` |

---

## 🔗 Conectores

### Configuración de Conectores

Cada conector se configura mediante un objeto JSON en la definición del agente.

#### Telegram

```json
{
  "token": "tu-bot-token-de-telegram",
  "group_mode": true,
  "mention_only": true,
  "admins": "usuario1,usuario2",
  "channel_id": "-1001234567890"
}
```

> **Nota**: Para que funcione en grupos, debes desactivar el "Privacy Mode" en @BotFather.

#### Discord

```json
{
  "token": "Bot TU_TOKEN_DE_DISCORD",
  "defaultChannel": "ID_DEL_CANAL"
}
```

> **Importante**: Habilita "Message Content Intent" en la configuración del bot.

#### Slack

```json
{
  "botToken": "xoxb-tu-bot-token",
  "appToken": "xapp-tu-app-token"
}
```

> Usa el manifest `slack.yaml` incluido para crear la app fácilmente.

#### GitHub Issues

```json
{
  "token": "TU_GITHUB_PAT_TOKEN",
  "repository": "nombre-del-repo",
  "owner": "propietario",
  "botUserName": "nombre-del-bot"
}
```

#### GitHub PRs

```json
{
  "token": "TU_GITHUB_PAT_TOKEN",
  "repository": "nombre-del-repo",
  "owner": "propietario"
}
```

#### Twitter/X

```json
{
  "apiKey": "tu-api-key",
  "apiSecret": "tu-api-secret",
  "accessToken": "tu-access-token",
  "accessSecret": "tu-access-secret"
}
```

#### IRC

```json
{
  "server": "irc.example.com",
  "port": "6667",
  "nickname": "IONETBot",
  "channel": "#canal",
  "alwaysReply": false
}
```

#### Matrix

```json
{
  "homeserver": "https://matrix.org",
  "user_id": "@usuario:matrix.org",
  "password": "tu-password",
  "room_id": "!habitacion:matrix.org"
}
```

#### Email (SMTP/IMAP)

```json
{
  "smtpServer": "smtp.gmail.com:587",
  "imapServer": "imap.gmail.com:993",
  "smtpInsecure": false,
  "imapInsecure": false,
  "username": "tu@gmail.com",
  "email": "tu@gmail.com",
  "password": "tu-password-o-app-password",
  "name": "IONET Agent"
}
```

---

## 🌐 API REST

### Endpoints de Gestión de Agentes

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `GET` | `/api/agents` | Listar todos los agentes |
| `POST` | `/api/agent/create` | Crear un nuevo agente |
| `GET` | `/api/agent/:name` | Obtener detalles del agente |
| `DELETE` | `/api/agent/:name` | Eliminar un agente |
| `PUT` | `/api/agent/:name/pause` | Pausar un agente |
| `PUT` | `/api/agent/:name/start` | Reanudar un agente |
| `GET` | `/api/agent/:name/status` | Ver historial de estado |
| `GET` | `/api/agent/:name/config` | Obtener configuración |
| `PUT` | `/api/agent/:name/config` | Actualizar configuración |

### Endpoints de Chat

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `POST` | `/api/chat/:name` | Enviar mensaje |
| `POST` | `/api/notify/:name` | Enviar notificación |
| `GET` | `/api/sse/:name` | Stream SSE en tiempo real |
| `POST` | `/v1/responses` | API compatible con OpenAI |

### Endpoints de Acciones y Grupos

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `GET` | `/api/actions` | Listar acciones disponibles |
| `POST` | `/api/action/:name/run` | Ejecutar una acción |
| `POST` | `/api/agent/group/generateProfiles` | Generar perfiles de grupo |
| `POST` | `/api/agent/group/create` | Crear grupo de agentes |

### Endpoints de Configuración

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `GET` | `/api/meta/agent/config` | Esquema de configuración |
| `GET` | `/settings/export/:name` | Exportar agente |
| `POST` | `/settings/import` | Importar agente |

### Ejemplos con cURL

#### Listar Agentes

```bash
curl -X GET "http://localhost:8080/api/agents"
```

#### Crear Agente

```bash
curl -X POST "http://localhost:8080/api/agent/create" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "mi-asistente",
    "model": "MiniMax-M2.7",
    "system_prompt": "Eres un asistente útil especializado en tecnología.",
    "enable_kb": true,
    "enable_reasoning": true,
    "actions": ["web-search", "web-scrape"]
  }'
```

#### Enviar Mensaje

```bash
curl -X POST "http://localhost:8080/api/chat/mi-asistente" \
  -H "Content-Type: application/json" \
  -d '{"message": "¿Cuál es el clima hoy?"}'
```

#### Stream SSE

```bash
curl -N -X GET "http://localhost:8080/api/sse/mi-asistente"
```

#### Ejecutar Acción

```bash
curl -X POST "http://localhost:8080/api/action/web-search/run" \
  -H "Content-Type: application/json" \
  -d '{"query": "últimas noticias de IA"}'
```

#### Exportar Agente

```bash
curl -X GET "http://localhost:8080/settings/export/mi-asistente" \
  --output mi-asistente.json
```

#### Importar Agente

```bash
curl -X POST "http://localhost:8080/settings/import" \
  -F "file=@mi-asistente.json"
```

---

## 📚 Uso como Librería

### Instalación

```bash
go get github.com/mudler/LocalAGI@latest
```

### Ejemplo: Agente Simple

```go
package main

import (
    "context"
    "log"
    
    "github.com/mudler/LocalAGI/core/agent"
    "github.com/mudler/LocalAGI/core/types"
)

func main() {
    // Crear nuevo agente
    a, err := agent.New(
        agent.WithModel("MiniMax-M2.7"),
        agent.WithLLMAPIURL("https://api.minimax.io/v1"),
        agent.WithLLMAPIKey("tu-api-key"),
        agent.WithSystemPrompt("Eres un asistente útil."),
        agent.WithCharacter(agent.Character{
            Name:        "asistente",
            Description: "Un asistente IA helpful",
        }),
        agent.WithActions([]types.Action{
            // Añadir acciones personalizadas aquí
        }),
        agent.WithStateFile("./state/asistente.state.json"),
        agent.WithCharacterFile("./state/asistente.character.json"),
        agent.WithTimeout("10m"),
        agent.EnableKnowledgeBase(),
        agent.EnableReasoning(),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Iniciar el agente
    go func() {
        if err := a.Run(context.Background()); err != nil {
            log.Printf("Agente detenido: %v", err)
        }
    }()

    // Detener el agente cuando sea necesario
    a.Stop()
}
```

### Ejemplo: Pool de Agentes

```go
package main

import (
    "context"
    "log"
    
    "github.com/mudler/LocalAGI/core/state"
    "github.com/mudler/LocalAGI/core/types"
)

func main() {
    pool, err := state.NewAgentPool(
        "default-model",
        "default-multimodal-model",
        "transcription-model",
        "en",
        "tts-model",
        "https://api.minimax.io/v1",
        "tu-api-key",
        "./state",
        func(config *types.AgentConfig) func(ctx context.Context, pool *state.AgentPool) []types.Action {
            return func(ctx context.Context, pool *state.AgentPool) []types.Action {
                return []types.Action{}
            }
        },
        func(config *types.AgentConfig) []types.Connector {
            return []types.Connector{}
        },
        func(config *types.AgentConfig) []types.DynamicPrompt {
            return []types.DynamicPrompt{}
        },
        func(config *types.AgentConfig) types.JobFilters {
            return types.JobFilters{}
        },
        "10m",
        true,
        nil,
    )
    if err != nil {
        log.Fatal(err)
    }

    agentConfig := &types.AgentConfig{
        Name:                 "mi-agente",
        Model:                "MiniMax-M2.7",
        SystemPrompt:         "Eres un asistente útil.",
        EnableKnowledgeBase:  true,
        EnableReasoning:      true,
    }

    err = pool.CreateAgent("mi-agente", agentConfig)
    if err != nil {
        log.Fatal(err)
    }

    err = pool.StartAll()
    if err != nil {
        log.Fatal(err)
    }

    status := pool.GetStatusHistory("mi-agente")
    log.Printf("Estado del agente: %v", status)

    pool.Stop("mi-agente")
    err = pool.Remove("mi-agente")
}
```

---

## 💻 Desarrollo

### Requisitos

- Go 1.26+
- Bun 1.2+
- Git
- Docker (para contenedores)
- Node.js 18+ (para desarrollo frontend)

### Construcción desde Código Fuente

```bash
# Clonar repositorio
git clone https://github.com/mudler/LocalAGI.git
cd LocalAGI

# Construir frontend
cd webui/react-ui
bun install
bun run build
cd ../..

# Construir backend
go build -o ionet main.go

# Ejecutar
./ionet
```

### Desarrollo con Hot-Reload

**Terminal 1 - Frontend:**
```bash
cd webui/react-ui
bun install
bun run build
bun run dev
```

**Terminal 2 - Backend:**
```bash
mkdir -p pool

export OPENAI_BASE_URL=https://api.minimax.io/v1
export OPENAI_API_KEY=tu-api-key
export MODEL_NAME=MiniMax-M2.7
export MULTIMODAL_MODEL=MiniMax-M2.7
export LOCALAGI_LLM_API_URL=$OPENAI_BASE_URL
export LOCALAGI_LLM_API_KEY=$OPENAI_API_KEY
export LOCALAGI_STATE_DIR=./pool
export LOCALAGI_TIMEOUT=5m
export LOCALAGI_ENABLE_CONVERSATIONS_LOGGING=false
export LOCALAGI_SSHBOX_URL=root:root@sshbox:22

go run main.go
```

---

## 🔧 Extendiendo IONET

### Acciones Personalizadas en Go

Las acciones personalizadas permiten extender la funcionalidad de IONET con código Go interpretado (sin necesidad de compilación).

#### Estructura de una Acción

Una acción personalizada requiere tres funciones:

1. **`Run(config map[string]interface{}) (string, map[string]interface{}, error)`** - Función principal de ejecución
2. **`Definition() map[string][]string`** - Define los parámetros de la acción
3. **`RequiredFields() []string`** - Lista de campos requeridos

#### Ejemplo: Acción de Clima

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "io"
)

type WeatherParams struct {
    City    string `json:"city"`
    Country string `json:"country"`
}

type WeatherResponse struct {
    Main struct {
        Temp     float64 `json:"temp"`
        Humidity int     `json:"humidity"`
    } `json:"main"`
    Weather []struct {
        Description string `json:"description"`
    } `json:"weather"`
}

func Run(config map[string]interface{}) (string, map[string]interface{}, error) {
    p := WeatherParams{}
    b, err := json.Marshal(config)
    if err != nil {
        return "", map[string]interface{}{}, err
    }
    if err := json.Unmarshal(b, &p); err != nil {
        return "", map[string]interface{}{}, err
    }

    url := fmt.Sprintf(
        "https://api.openweathermap.org/data/2.5/weather?q=%s,%s&appid=TU_API_KEY&units=metric",
        p.City, p.Country,
    )
    resp, err := http.Get(url)
    if err != nil {
        return "", map[string]interface{}{}, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", map[string]interface{}{}, err
    }

    var weather WeatherResponse
    if err := json.Unmarshal(body, &weather); err != nil {
        return "", map[string]interface{}{}, err
    }

    result := fmt.Sprintf(
        "Clima en %s, %s: %.1f°C, %s, Humedad: %d%%",
        p.City, p.Country, weather.Main.Temp, weather.Weather[0].Description, weather.Main.Humidity,
    )

    return result, map[string]interface{}{}, nil
}

func Definition() map[string][]string {
    return map[string][]string{
        "city":    {"string", "Ciudad para consultar el clima"},
        "country": {"string", "Código de país (ej: US, ES, MX)"},
    }
}

func RequiredFields() []string {
    return []string{"city", "country"}
}
```

#### Carga Automática de Acciones

Para cargar automáticamente acciones personalizadas:

```bash
mkdir -p /path/to/custom-actions
export LOCALAGI_CUSTOM_ACTIONS_DIR=/path/to/custom-actions
```

### MCP (Model Context Protocol)

MCP permite conectar IONET con servidores externos para acceder a herramientas adicionales.

#### Servidores MCP Locales (STDIO)

```json
{
  "mcpServers": {
    "github": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "-e", "GITHUB_PERSONAL_ACCESS_TOKEN", "ghcr.io/github/github-mcp-server"],
      "env": {"GITHUB_PERSONAL_ACCESS_TOKEN": "tu-token"}
    }
  }
}
```

#### Servidores MCP Remotos (HTTP)

```json
{
  "mcpServers": [
    {"url": "https://mcp.server.com", "token": "tu-token"}
  ]
}
```

### Skills (Habilidades)

Las skills son conjuntos reutilizables de instrucciones y recursos que los agentes pueden utilizar.

1. Abrir la sección **Skills** en el menú lateral
2. Los skills se almacenan en `STATE_DIR/skills`
3. Crear, editar, buscar, importar y exportar skills
4. Sincronizar con repositorios Git

Para habilitar skills por agente, ir a **Advanced Settings** y activar **Enable Skills**.

---

## 🖥️ Hardware Soportado

### Configuraciones Disponibles

| Configuración | Descripción | Uso Recomendado |
|---------------|-------------|-----------------|
| **CPU** | Solo CPU (sin GPU) | Desarrollo, testing |
| **NVIDIA GPU** | CUDA | Alto rendimiento, modelos grandes |
| **AMD GPU** | ROCm | Alto rendimiento, modelos grandes |
| **Intel GPU** | SYCL | Media performance, laptops |

### Modelos Recomendados

| Modelo | Tamaño | Uso |
|--------|--------|-----|
| `gemma-3-4b-it-qat` | 4B | CPU, inicio rápido |
| `gemma-3-12b-it` | 12B | GPU media |
| `gemma-3-27b-it` | 27B | GPU alta gama |
| `qwen_qwq-32b` | 32B | Mejor coordinación de agentes |

---

## 💡 Casos de Uso

### 1. Automatización de GitHub
- Gestionar issues automáticamente
- Revisar PRs y dejar comentarios
- Mantener documentación actualizada

### 2. Bot de Discord/Telegram
- Atención al cliente 24/7
- Moderación de servidores
- Información y consultas

### 3. Asistente Personal
- Organización de tareas
- Búsqueda de información
- Gestión de calendario

### 4. Base de Conocimiento RAG
- Documentación empresarial
- Biblioteca de investigación
- FAQs automatizadas

### 5. Equipos de Agentes Cooperativos
- Especialización por tarea
- Comunicación entre agentes
- Resolución de problemas complejos

---

## ❓ FAQ

### ¿IONET requiere GPU?
**No es obligatorio**, pero se recomienda para mejor rendimiento. IONET puede funcionar en CPU con modelos pequeños.

### ¿Puedo usar mis propios modelos?
Sí. IONET es compatible con cualquier modelo que exponga un API compatible con OpenAI.

### ¿Mis datos salen de mi hardware?
**No.** Todos los datos se procesan localmente.

### ¿Cómo reporto bugs?
Abre un issue en el [repositorio de GitHub](https://github.com/mudler/LocalAGI/issues).

### ¿IONET es compatible con MCP?
Sí. Soporta servidores MCP locales (STDIO) y remotos (HTTP).

---

## 🌟 Familia LocalAI

IONET es parte de un ecosistema de herramientas de IA local:

| Proyecto | Descripción |
|----------|-------------|
| **[LocalAI](https://github.com/mudler/LocalAI)** | Alternativa Open Source a OpenAI. API REST compatible con OpenAI. |
| **[LocalRecall](https://github.com/mudler/LocalRecall)** | API RESTful y sistema de gestión de base de conocimiento. |
| **[skillserver](https://github.com/mudler/skillserver)** | Servidor para gestión de skills reutilizables. |

---

## 📸 Capturas de Pantalla

### Dashboard Principal
![Dashboard](https://github.com/user-attachments/assets/a40194f9-af3a-461f-8b39-5f4612fbf221)

### Configuración de Agente
![Agent Settings](https://github.com/user-attachments/assets/fb3c3e2a-cd53-4ca8-97aa-c5da51ff1f83)

### Crear Grupo de Agentes
![Create Group](https://github.com/user-attachments/assets/102189a2-0fba-4a1e-b0cb-f99268ef8062)

### Observabilidad del Agente
![Agent Observability](https://github.com/user-attachments/assets/f7359048-9d28-4cf1-9151-1f5556ce9235)

---

## 🤝 Conectores Soportados

<p align="center">
  <img src="https://github.com/user-attachments/assets/4171072f-e4bf-4485-982b-55d55086f8fc" alt="Telegram" width="60"/>
  <img src="https://github.com/user-attachments/assets/9235da84-0187-4f26-8482-32dcc55702ef" alt="Discord" width="220"/>
  <img src="https://github.com/user-attachments/assets/a88c3d88-a387-4fb5-b513-22bdd5da7413" alt="Slack" width="220"/>
  <img src="https://github.com/user-attachments/assets/d249cdf5-ab34-4ab1-afdf-b99e2db182d2" alt="IRC" width="220"/>
  <img src="https://github.com/user-attachments/assets/52c852b0-4b50-4926-9fa0-aa50613ac622" alt="GitHub" width="220"/>
</p>

---

## 📄 Licencia

MIT License - Ver archivo [LICENSE](LICENSE) para más detalles.

---

<p align="center">
  <strong>PROCESAMIENTO LOCAL. PENSAMIENTO GLOBAL.</strong><br>
  Hecho con ❤️ por <a href="https://github.com/mudler">mudler</a>
</p>