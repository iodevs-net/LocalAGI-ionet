---
name: knowledge-management
description: Gestión de memoria y conocimiento utilizando RAG (Retrieval Augmented Generation). Usa esta skill cuando necesites almacenar, buscar y recuperar información de la memoria del sistema, incluyendo notas, conceptos, procedimientos y datos históricos. El sistema utiliza Bleve como índice de búsqueda full-text para encontrar entradas relevantes. Incluye acciones para add_to_memory, list_memory, remove_from_memory, search_memory.
---

# Knowledge Management

## Cuándo usar esta skill

**Usar cuando:**
- Necesites guardar información importante para referencia futura
- Requieras buscar en el historial de conversaciones o datos almacenados
- Quieras mantener contexto entre sesiones
- Necesites almacenar procedimientos, notas o documentación
- Requieras búsqueda semántica en información previamente almacenada

**NO usar cuando:**
- La información sea extremadamente sensible (usar cifrado externo)
- Necesites almacenamiento estructurado (usar base de datos)
- La información sea efímera y no necesite persistencia
- Requieras sincronización entre múltiples instancias

## Acciones disponibles

| Acción | Descripción | Parámetros clave |
|--------|-------------|------------------|
| `add_to_memory` | Añadir nueva entrada a la memoria | name, content |
| `list_memory` | Listar todas las entradas en memoria | (sin parámetros) |
| `remove_from_memory` | Eliminar entrada por ID | id |
| `search_memory` | Buscar en memoria por query | query |

## Modelo de datos

### MemoryEntry
```json
{
  "id": "1699999999999999999",
  "name": "Protocolo de reinicio CPE",
  "content": "Para reiniciar un CPE en la red de IONET...",
  "created_at": "2024-01-15T10:30:00Z"
}
```

## Flujo de trabajo

### 1. Agregar información a memoria

```
1. Identificar información valiosa para guardar
2. Preparar name (título descriptivo) y content (contenido detallado)
3. Usar add_to_memory
4. Confirmar ID asignado y guardar para referencia futura
```

### 2. Buscar información en memoria

```
1. Formular query de búsqueda
2. Usar search_memory con el query
3. Revisar resultados devueltos (nombre, contenido, ID)
4. Usar información relevante para construir respuesta
```

### 3. Gestionar entradas de memoria

```
1. Listar todas las entradas con list_memory
2. Identificar entradas obsoletas o duplicadas
3. Eliminar con remove_from_memory usando el ID
4. Mantener memoria limpia y relevante
```

## Ejemplos de uso

### Ejemplo 1: Agregar protocolo
```json
{
  "action": "add_to_memory",
  "name": "Protocolo de reinicio de equipos Mikrotik",
  "content": "1. Acceder via Winbox al router\n2. Ir a System > Reboot\n3. Esperar 2 minutos\n4. Verificar conexión"
}
```

### Ejemplo 2: Buscar información
```json
{
  "action": "search_memory",
  "query": "reinicio CPE configuración"
}
```

### Ejemplo 3: Listar memoria
```json
{
  "action": "list_memory"
}
```

### Ejemplo 4: Eliminar entrada
```json
{
  "action": "remove_from_memory",
  "id": "1699999999999999999"
}
```

## Notas de implementación

- **Indexación**: Usa Bleve para búsqueda full-text
- **Búsqueda**: Busca tanto en name como en content
- **Límites**: ListMemory devuelve hasta 10000 entradas
- **IDs**: Generados automáticamente usando timestamp Unix en nanosegundos
- **Ordenamiento**: ListMemory ordena por created_at descendente

## Casos de uso para IONET

### Almacenar protocolos
```
add_to_memory con name="Protocolo X" y content=pasos detallados
```

### Buscar procedimientos
```
search_memory con query de problema o equipo específico
```

### Mantener documentación de cambios
```
add_to_memory después de cambios importantes en configuración
```

## Referencias

- [Guía de memoria](./references/memory_guide.md) - Mejores prácticas de organización
- [Patrones de búsqueda](./references/search_patterns.md) - Cómo formular queries efectivas