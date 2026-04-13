# IONET - Configuración de Agentes

Este directorio contiene la configuración de los agentes de IONET.

## Agentes

| Agente | Rol | Especialidad |
|--------|-----|--------------|
| **ION** | Orchestrator | Interlocutor principal, coordina a los demás agentes |
| **agente-protocolos** | Especializado | Protocolos, procedimientos, políticas |
| **agente-redes** | Especializado | Redes, comunicaciones, equipos |
| **agente-base** | Especializado | Información general, FAQs, contactos |

## Importar Agentes en LocalAGI

### Opción 1: Importar desde la UI

1. Acceder a `http://localhost:8080` (dev) o `https://agentes.ionet.cl` (prod)
2. Ir a **Agents** > **Create Agent**
3. En **Basic Info**, hacer click en **Import from JSON**
4. Seleccionar el archivo JSON del agente
5. Ajustar configuración si es necesario
6. Click **Save**

### Opción 2: Importar por API

```bash
# Importar ION
curl -X POST http://localhost:8080/api/agent/import \
  -H "Content-Type: application/json" \
  -d @ion.json

# Importar agentes especializados
curl -X POST http://localhost:8080/api/agent/import \
  -H "Content-Type: application/json" \
  -d @agente-protocolos.json
```

### Opción 3: Copiar manualmente al pool

```bash
# Los agentes se guardan en el pool de estado
# Por defecto: /pool/agents/

# Copiar archivos JSON
cp *.json /pool/agents/

# Reiniciar el servicio
docker compose restart ionet
```

## Orden de Importación

**IMPORTANTE**: Importar en este orden:

1. `agente-protocolos.json` (base)
2. `agente-redes.json` (base)
3. `agente-base.json` (base)
4. `ion.json` (orchestrator - requiere que existan los anteriores)

El agente ION tiene configuración de `call_agents` que referencia a los otros agentes.

## Configuración de Variables

Los archivos JSON usan variables de entorno que se sustituyen en runtime:

| Variable | Descripción | Ejemplo |
|----------|-------------|---------|
| `${MODEL_NAME}` | Modelo de texto principal | `MiniMax-M2.7` |
| `${MULTIMODAL_MODEL}` | Modelo multimodal | `MiniMax-M2.7` |

Para sobrescribir valores, editar el JSON o configurar en la UI después de importar.

## Permisos entre Agentes

El agente ION puede llamar a los agentes especializados. Esto se configura en la acción `call_agents`:

```json
"actions": [
  {
    "name": "call_agents",
    "config": "{}"
  }
]
```

Por defecto, `call_agents` permite llamar a todos los agentes del pool.

## Agregar Nuevos Agentes

1. Crear archivo JSON siguiendo el formato de `agente-protocolos.json`
2. Asegurar que el `name` sea único
3. Importar el agente
4. Si es necesario, actualizar ION para que pueda llamarlo

## RAG y Knowledge Base

Los agentes especializados (`agente-protocolos`, `agente-redes`, `agente-base`) tienen `enable_kb: true` para consultar la base de conocimiento RAG.

Para que funcione:

1. Asegurar que la base de datos RAG está configurada en `docker-compose.prod.yaml`
2. Los documentos de M365 sync se guardan en `/rag/raw/`
3. La indexación se hace automáticamente por el sistema RAG de LocalAGI

## Verificar Configuración

Después de importar, verificar en la UI:

1. **Agents** > Ver que todos aparecen
2. **Agent Settings** > ION > Ver que `call_agents` tiene acceso a los demás
3. **Knowledge** > Ver que las colecciones están creadas
4. Hacer una prueba de chat con ION
