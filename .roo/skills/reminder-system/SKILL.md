---
name: reminder-system
description: Sistema de recordatorios y tareas programadas. Usa esta skill cuando necesites crear recordatorios, agendar tareas, o gestionar recordatorios periódicos. El sistema permite almacenar recordatorios con título, mensaje, fecha/hora y estado, con soporte para espacios de nombres (namespaces) para organizar diferentes tipos de recordatorios. Integración con scheduler para ejecución de tareas en el futuro.
---

# Reminder System

## Cuándo usar esta skill

**Usar cuando:**
- Necesites crear recordatorios para el usuario
- Requieras agendar acciones futuras
- Quieras hacer seguimiento de tareas pendientes
- Necesites notificaciones periódicas
- Gestiones recordatorios por categorías (namespace)

**NO usar cuando:**
- La tarea sea crítica y requiera garantía de ejecución (usar scheduler externo)
- Necesites persistencia a largo plazo (usar base de datos)
- Requieras integración con sistemas externos de calendario

## Modelo de datos

### Reminder
```json
{
  "id": "uuid-v4",
  "namespace": "personal",
  "title": "Reunión con equipo",
  "message": "Revisar avances del proyecto IONET",
  "datetime": "2024-01-15T14:00:00Z",
  "interval": null,
  "state": "pending",
  "created_at": "2024-01-10T10:00:00Z"
}
```

### Campos principales
- **id**: Identificador único UUID v4
- **namespace**: Espacio de nombres para organizar (ej: "personal", "trabajo", "ionet")
- **title**: Título corto del recordatorio
- **message**: Descripción detallada
- **datetime**: Cuándo ejecutar (formato ISO 8601)
- **interval**: Para recordatorios recurrentes (cron format o null)
- **state**: Estado actual (pending, done, cancelled)
- **created_at**: Timestamp de creación

## Flujo de trabajo

### 1. Crear recordatorio simple

```
1. Definir título claro y mensaje descriptivo
2. Especificar datetime para la notificación
3. Opcional: asignar namespace para organización
4. Usar new_reminder o acción equivalente
5. Confirmar ID asignado y datetime
```

### 2. Crear recordatorio recurrente

```
1. Definir título y mensaje
2. Especificar datetime inicial
3. Definir interval (ej: "0 9 * * *" para diario a las 9am)
4. Usar new_reminder con interval
5. El sistema recreará el recordatorio tras cada ejecución
```

### 3. Gestionar recordatorios

```
1. Listar recordatorios activos con list_reminders
2. Marcar como completado con done_reminder
3. Cancelar con cancel_reminder
4. Ver historial con list_reminders_history
```

## Acciones disponibles

| Acción | Descripción | Parámetros clave |
|--------|-------------|------------------|
| `new_reminder` | Crear nuevo recordatorio | namespace, title, message, datetime, interval |
| `list_reminders` | Listar recordatorios pendientes | namespace (opcional) |
| `done_reminder` | Marcar como completado | id |
| `cancel_reminder` | Cancelar recordatorio | id |
| `list_reminders_history` | Ver historial de recordatorios | namespace (opcional) |

## Ejemplos de uso

### Ejemplo 1: Recordatorio simple
```json
{
  "action": "new_reminder",
  "namespace": "ionet",
  "title": "Revisión semanal",
  "message": "Revisar estado de tickets abiertos en GitHub",
  "datetime": "2024-01-20T10:00:00Z"
}
```

### Ejemplo 2: Recordatorio diario
```json
{
  "action": "new_reminder",
  "namespace": "personal",
  "title": "Standup diario",
  "message": "Reunión de sincronización con el equipo",
  "datetime": "2024-01-15T09:00:00Z",
  "interval": "0 9 * * *"
}
```

### Ejemplo 3: Listar recordatorios
```json
{
  "action": "list_reminders",
  "namespace": "ionet"
}
```

### Ejemplo 4: Completar recordatorio
```json
{
  "action": "done_reminder",
  "id": "abc123-def456-ghi789"
}
```

## Namespaces recomendados para IONET

| Namespace | Uso |
|-----------|-----|
| `ionet` | Recordatorios del sistema IONET |
| `tickets` | Seguimiento de tickets de clientes |
| `mantenimiento` | Tareas de mantenimiento programadas |
| `personal` | Recordatorios personales del técnico |

## Notas de implementación

- **UUID**: Usa UUID v4 para identificadores únicos
- **Timezone**: Almacena en UTC, convierte a hora local para display
- **Interval**: Formato cron de 5 campos (minuto, hora, día, mes, día semana)
- **Estado**: pending → done/cancelled (transiciones válidas)
- **Persistencia**: Los recordatorios sobreviven reinicios del sistema

## Referencias

- [Guía de cron](./references/cron_guide.md) - Formato de intervalos
- [Patrones de uso](./references/usage_patterns.md) - Casos de uso comunes