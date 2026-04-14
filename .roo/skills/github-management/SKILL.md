---
name: github-management
description: Gestión completa de repositorios GitHub incluyendo issues, pull requests, comentarios, etiquetas y contenido de repositorios. Usa esta skill cuando necesites crear, editar, buscar o gestionar elementos en GitHub, incluyendo la automatización de flujos de trabajo de desarrollo como revisión de código, gestión de incidencias y coordinación de equipos. Incluye acciones para issue-opener, issue-closer, issue-comment, issue-edit, issue-labeler, issue-reader, issue-search, pr-reader, pr-comment, pr-reviewer, pr-creator, repository-get-content, repository-create-content, repository-readme, repository-list-files, repository-search-files.
---

# GitHub Management

## Cuándo usar esta skill

**Usar cuando:**
- Necesites crear, editar o cerrar issues en GitHub
- Requieras revisar pull requests y añadir comentarios o revisiones
- Busques issues o contenido específico en repositorios
- Necesites gestionar etiquetas (labels) en issues o PRs
- Requieras leer o crear contenido en repositorios GitHub
- Automatices flujos de trabajo de desarrollo (CI/CD manual)

**NO usar cuando:**
- Sea para tareas administrativas de GitHub (gestión de equipos, permisos)
- Requieras webhooks o integración con GitHub Actions
- La tarea sea específicamente de GitLab o Bitbucket

## Acciones disponibles

### Issues
| Acción | Descripción | Parámetros clave |
|--------|-------------|------------------|
| `issue-opener` | Crear nuevo issue | owner, repo, title, body, labels |
| `issue-closer` | Cerrar issue existente | owner, repo, issue_number |
| `issue-comment` | Añadir comentario a issue | owner, repo, issue_number, body |
| `issue-edit` | Editar issue | owner, repo, issue_number, title, body, state, labels |
| `issue-labeler` | Gestionar etiquetas | owner, repo, issue_number, labels |
| `issue-reader` | Leer detalles de issue | owner, repo, issue_number |
| `issue-search` | Buscar issues | query, owner, repo, state, labels |

### Pull Requests
| Acción | Descripción | Parámetros clave |
|--------|-------------|------------------|
| `pr-reader` | Leer detalles de PR | owner, repo, pr_number |
| `pr-comment` | Comentar en PR | owner, repo, pr_number, body |
| `pr-reviewer` | Realizar revisión de PR | owner, repo, pr_number, event, body, comments |
| `pr-creator` | Crear nuevo PR | owner, repo, title, body, head, base |

### Repositorio
| Acción | Descripción | Parámetros clave |
|--------|-------------|------------------|
| `repository-get-content` | Obtener contenido de archivo | owner, repo, path, ref |
| `repository-create-content` | Crear/actualizar archivo | owner, repo, path, content, message, sha |
| `repository-readme` | Obtener README del repo | owner, repo |
| `repository-list-files` | Listar archivos en repo | owner, repo, path, ref |
| `repository-search-files` | Buscar archivos en repo | owner, repo, query, file |

## Flujo de trabajo

### 1. Crear y gestionar issue

```
1. Verificar que tienes el owner, repo y credenciales necesarias
2. Preparar el contenido del issue (título, descripción, etiquetas)
3. Usar issue-opener para crear el issue
4. Opcional: Usar issue-labeler para añadir etiquetas adicionales
5. Confirmar la creación y proporcionar el link del issue
```

### 2. Revisar Pull Request

```
1. Obtener detalles del PR con pr-reader
2. Analizar los cambios (diffs, archivos modificados)
3. Preparar revisión con comentarios específicos
4. Usar pr-reviewer con event=COMMENT, APPROVE, REQUEST_CHANGES o CHANGES_REQUESTED
5. Añadir comentarios adicionales si es necesario
```

### 3. Buscar y rastrear issues

```
1. Construir query de búsqueda (puede usar sintaxis GitHub search)
2. Usar issue-search con filtros (state, labels, author, etc.)
3. Presentar resultados con enlaces y resumen
4. Ofrecer acciones de seguimiento (comentar, editar, cerrar)
```

## Ejemplos de uso

### Ejemplo 1: Crear issue con etiquetas
```json
{
  "action": "issue-opener",
  "owner": "usuario",
  "repo": "proyecto",
  "title": "Bug en autenticación con OAuth",
  "body": "## Descripción\nError al autenticar...\n\n## Pasos para reproducir\n1. Ir a /login\n2. Click en 'OAuth'\n\n## Entorno\n- Versión: 2.1.0\n- Navegador: Chrome 120",
  "labels": ["bug", "high-priority", "auth"]
}
```

### Ejemplo 2: Revisar PR
```json
{
  "action": "pr-reviewer",
  "owner": "usuario",
  "repo": "proyecto",
  "pr_number": 42,
  "event": "REQUEST_CHANGES",
  "body": "Buen trabajo en general. Solicito algunos cambios:",
  "comments": [
    {"path": "src/auth.go", "line": 15, "body": "Considerar validar este input primero"},
    {"path": "src/auth.go", "line": 23, "body": "Este error debería ser más específico"}
  ]
}
```

### Ejemplo 3: Buscar issues abiertos
```json
{
  "action": "issue-search",
  "query": "is:issue is:open label:bug",
  "owner": "usuario",
  "repo": "proyecto",
  "state": "open"
}
```

## Notas de implementación

- **Rate limits**: GitHub tiene límites de tasa (5000 requests/hora para autenticado)
- **Sintaxis markdown**: El body de issues y PRs soporta markdown completo
- **Autenticación**: Las acciones requieren token de GitHub configurado
- **Numeración**: Los issues y PRs tienen números únicos por repo, no globales
- **Estados**: Issues pueden estar "open" o "closed"; PRs tienen estados adicionales