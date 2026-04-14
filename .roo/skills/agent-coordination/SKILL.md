---
name: agent-coordination
description: Orquestación y coordinación entre múltiples agentes IA. Usa esta skill cuando necesites delegar tareas a agentes especializados, coordinar flujos de trabajo complejos entre múltiples agentes, o construir sistemas multi-agente donde cada agente tiene responsabilidades específicas. El sistema permite llamadas síncronas entre agentes con paso de contexto y результаты. Incluye acción call_agent.
---

# Agent Coordination

## Cuándo usar esta skill

**Usar cuando:**
- Necesites delegar tareas específicas a agentes especializados
- Requieras coordinar múltiples agentes para tareas complejas
- Quieras construir flujos de trabajo de agentes encadenados
- Necesites agregar capacidades de otros agentes a tu contexto
- Coordines el agente ION (orquestador) con agentes especializados

**NO usar cuando:**
- La tarea pueda ser completada por un solo agente
- Requieras ejecución paralela verdadera (usar múltiples jobs)
- Necesites compartir estado mutable entre agentes
- La coordinación sea de tiempo real critica

## Arquitectura de agentes en IONET

```
┌─────────────────────────────────────────────────────────────┐
│                    ION (Orquestador)                        │
│           "Hola, necesito el protocolo para X"              │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│              Determinación de agente correcto               │
└─────────────────────────────────────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        ▼                   ▼                   ▼
┌──────────────┐   ┌──────────────┐   ┌──────────────────┐
│agente-       │   │agente-       │   │agente-           │
│protocolos    │   │redes         │   │base             │
│(RAG docs)    │   │(RAG docs)    │   │(RAG docs)       │
└──────────────┘   └──────────────┘   └──────────────────┘
```

## Agentes disponibles en IONET

| Agente | Descripción | Casos de uso |
|--------|-------------|--------------|
| `ion` | Orquestador principal | Routing inicial, coordinación |
| `agente-protocolos` | Protocolos y procedimientos | Consultas sobre流程 de trabajo |
| `agente-redes` | Networking y conectividad | Problemas de red, configuración |
| `agente-base` | Información general | Preguntas generales, ayuda |
| `agente-clientes` | Gestión de clientes | Soporte técnico clientes |
| `agente-inventario` | Inventario y equipos | Consulta de equipos |
| `agente-datos` | Análisis de datos | Reportes, métricas |
| `agente-seguridad` | Seguridad y permisos | Incidentes de seguridad |
| `agente-servicios` | Servicios y status | Estado de servicios |

## Acción disponible

| Acción | Descripción | Parámetros clave |
|--------|-------------|------------------|
| `call_agent` | Llamar a otro agente | agent_name, message |

## Flujo de trabajo

### 1. Delegar tarea a agente especializado

```
1. Identificar cuál agente es más apropiado para la tarea
2. Formular mensaje claro con contexto y pregunta específica
3. Usar call_agent con el agent_name y message
4. Recibir respuesta del agente
5. Synthesizar respuesta y proporcionar al usuario
```

### 2. Coordinación de múltiples agentes

```
1. Analizar tarea compleja y dividir en subtareas
2. Para cada subtarea, identificar agente apropiado
3. Ejecutar llamadas en secuencia o según dependencias
4. Recopilar resultados de cada agente
5. Combinar resultados en respuesta final
```

### 3. Escalamiento de tareas

```
1. Evaluar si la tarea requiere especialista
2. Si es así, usar call_agent para delegar
3. Si el agente no puede resolver, escalar a humano
4. Mantener contexto del flujo para reporte
```

## Ejemplos de uso

### Ejemplo 1: Consultar protocolo
```json
{
  "action": "call_agent",
  "agent_name": "agente-protocolos",
  "message": "Cuál es el procedimiento para reiniciar un CPE que no responde?"
}
```

### Ejemplo 2: Diagnosticar problema de red
```json
{
  "action": "call_agent",
  "agent_name": "agente-redes",
  "message": "Un cliente reporta que no tiene internet. IP del cliente: 192.168.100.45"
}
```

### Ejemplo 3: Obtener información de equipo
```json
{
  "action": "call_agent",
  "agent_name": "agente-inventario",
  "message": "Qué equipo está asociado al serial ION-2024-0543?"
}
```

## Notas de implementación

- **Llamadas síncronas**: El agente espera respuesta completa antes de continuar
- **Contexto compartido**: Las conversaciones son transmitidas al agente llamado
- **Whitelist/Blacklist**: Se pueden configurar listas de agentes permitidos/bloqueados
- **Metadata**: Resultados pueden incluir metadata de otros agentes
- **Errores**: Si el agente no existe, retorna error

## Configuración de llamadas

### Whitelist (solo permitir ciertos agentes)
```json
{
  "whitelist": "agente-protocolos,agente-redes"
}
```

### Blacklist (excluir ciertos agentes)
```json
{
  "blacklist": "agente-seguridad"
}
```

## Referencias

- [Guía de agentes](./references/agents_guide.md) - Descripción detallada de cada agente
- [Flujos comunes](./references/common_flows.md) - Patrones de coordinación frecuentes