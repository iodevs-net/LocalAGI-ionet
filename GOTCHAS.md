# GOTCHAS — IONET

Errores, bloqueos y lecciones aprendidas. Fuente de verdad para no caer dos veces.

## MCP Workspace

### Pool de configuración no se actualiza desde host
**Problema:** Editar `config/agents/ion.json` en el host NO propaga los cambios al contenedor. LocalAGI lee los agentes desde `/pool/agents/` (volumen persistente), no desde el bind mount.

**Fix:** Copiar manualmente al pool y reiniciar:
```bash
docker cp config/agents/ion.json localagi-ionet-ionet-1:/pool/agents/ion.json
docker compose restart ionet
```

### FastMCP devuelve 406/400 en curl directo
**Problema:** Al testear el MCP server con `curl -X POST ...`, FastMCP responde 406 o 400.

**Causa:** FastMCP espera cabeceras HTTP específicas (Accept, Session) que curl no envía por defecto.

**Fix:** No es un bug — LocalAGI envía las cabeceras correctas. Probar con `curl -X POST "http://localhost:8766/mcp" -H "Content-Type: application/json" -d '...'` o directamente desde ION.

### COPY *.py en Dockerfile no copia .json
**Problema:** El Dockerfile usa `COPY *.py .`, que NO copia archivos `.json`.

**Lección:** Para montar archivos de configuración/claves, usar volúmenes en `docker-compose.yaml` en vez de quemarlos en la imagen. Más seguro, más flexible.

## Google Drive OAuth

### Device Code Flow requiere sí o sí tipo "Desktop app"
**Problema:** `POST https://oauth2.googleapis.com/device/code` devuelve `401 invalid_client: Invalid client type`.

**Causa:** Aunque en consola diga "Escritorio", Google aplica políticas de organización o propagación que bloquean Device Code Flow.

**Estado:** No resuelto. Ver sección de Service Account.

### OOB (urn:ietf:wg:oauth:2.0:oob) deprecado
**Problema:** Google devuelve `400 invalid_request: Missing required parameter: redirect_uri`.

**Causa:** Google deprecó el flujo OOB (out-of-band) en 2023. Ya no permite redirect URIs `urn:ietf:wg:oauth:2.0:oob`.

**Fix (alternativo):** Usar `run_local_server()` con puerto mapeado en Docker, o usar redirect URI `http://localhost:PORT` registrada en la consola.

### Service account no puede crear archivos en My Drive
**Problema:** `Error: Service Accounts do not have storage quota.`

**Causa:** Las service accounts de Gmail gratis no tienen cuota de almacenamiento. Solo pueden leer archivos compartidos con ellas.

**Fix actual:** Service account en modo solo lectura para buscar/leer archivos en Drive. Escritura en workspace local (`/pool/ion-workspace/`).

**Fix futuro:** Autenticación OAuth real con `el.agent.ion@gmail.com` usando `run_local_server()` o rclone.

## Google Cloud Console

### Pantalla de consentimiento debe tener scopes
**Problema:** Device Code Flow falla aunque el cliente sea tipo Desktop.

**Causa potencial:** La pantalla de consentimiento OAuth necesita tener agregado el scope `.../auth/drive` incluso en modo Testing.

### rclone es más sencillo que OAuth manual
**Lección:** rclone tiene su propio cliente OAuth preconfigurado (millones de usuarios) que funciona con Device Code Flow sin configuración extra. Preferir rclone sobre implementación manual de OAuth.

## Docker

### docker exec -it falla en CI/automation
**Problema:** `docker exec -it` devuelve `cannot attach stdin to a TTY-enabled container because stdin is not a terminal`.

**Fix:** Usar `docker exec -i` (sin -t) o `docker exec` (sin flags) para comandos no interactivos.

### Port mapping del contenedor no es localhost del host
**Problema:** Un server escuchando en `127.0.0.1:PORT` dentro del contenedor NO es accesible desde `localhost:PORT` del host.

**Causa:** Docker mapea `0.0.0.0:PORT` del contenedor, no `127.0.0.1:PORT`.

**Fix:** Escuchar en `0.0.0.0:PORT` dentro del contenedor, o verificar bind address al hacer port forwarding.

### pgrep, ps, ss no existen en python:3.12-slim
**Problema:** La imagen slim no tiene herramientas de diagnóstico.

**Fix:** Usar Python para diagnósticos básicos:
```python
import os
os.listdir('/proc')  # listar PIDs
import socket
socket.socket().bind(('0.0.0.0', PORT))  # verificar puerto disponible
```

## Arquitectura

### Agente MCP vs modificar el código de LocalAGI
**Decisión:** En vez de forkear LocalAGI, se construyó un servidor MCP externo. Más mantenible, sin acoplamiento.

### Workspace local + Drive como aumentación
**Decisión:** ION escribe en workspace local (rápido, sin API calls) y lee de Drive vía service account. Si se necesita escritura en Drive, rclone sync es el camino.
