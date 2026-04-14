---
name: ionet-orchestration
description: Orquestación principal del sistema IONET. Usa esta skill cuando necesites entender cómo el agente ION coordina los agentes especializados, cómo se estructura el flujo de información, o cómo solicitar ayuda de manera efectiva. El orquestador es el punto de entrada para todos los usuarios y determina qué agente especializado debe manejar cada consulta. Esta skill documenta la arquitectura de coordinación multi-agente de IONET.
---

# IONET Orchestration - Guía del Orquestador ION

## Cuándo usar esta skill

**Usar cuando:**
- Quieras entender cómo funciona la coordinación de agentes
- Necesites saber cuándo usar cada agente especializado
- Requieras optimizar queries para mejor respuesta
- Depures problemas de routing entre agentes
- Configures nuevos agentes o modifiques flujos

**NO usar cuando:**
- Ya sepas exactamente qué agente necesitas (usar directamente)
- La consulta sea solo para un agente específico
- Necesites debugging de bajo nivel (ver logs)

## Arquitectura del sistema

```
┌────────────────────────────────────────────────────────────────┐
│                         USUARIO (técnico)                      │
│              "Tengo un problema con el cliente X"              │
└────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌────────────────────────────────────────────────────────────────┐
│              ION (Orquestador / Interlocutor)                  │
│  • Recibe consulta inicial                                     │
│  • Determina tipo de problema                                   │
│  • Identifica agente apropiado                                  │
│  • Coordina respuesta final                                     │
└────────────────────────────────────────────────────────────────┘
                                │
         ┌──────────────────────┼──────────────────────┐
         │                      │                      │
         ▼                      ▼                      ▼
┌─────────────────┐   ┌─────────────────┐   ┌─────────────────┐
│ agente-protocolos│   │  agente-redes    │   │  agente-base    │
│ (procedimientos)│   │  (networking)    │   │  (general info) │
└─────────────────┘   └─────────────────┘   └─────────────────┘
         │                      │                      │
         └──────────────────────┼──────────────────────┘
                                │
                                ▼
┌────────────────────────────────────────────────────────────────┐
│              Respuesta integrada al usuario                    │
└────────────────────────────────────────────────────────────────┘
```

## Agentes especializados disponibles

| Agente | Responsabilidad | Cuándo usarlo |
|--------|-----------------|---------------|
| `ion` | Orquestación | Consultas generales, routing |
| `agente-protocolos` | Procedimientos | Cómo hacer X, protocolos |
| `agente-redes` | Networking | Problemas de conexión, configs |
| `agente-clientes` | Gestión clientes | Info de cliente, tickets |
| `agente-inventario` | Equipos | Serial, modelo, ubicación |
| `agente-datos` | Análisis | Métricas, reportes |
| `agente-seguridad` | Incidentes | Problemas de security |
| `agente-servicios` | Status servicios | Estado de servicios |
| `agente-base` | Info general | Preguntas generales |

## Flujo de enrutamiento

### Paso 1: Recepción
```
El técnico envía mensaje al orquestador ION
ION analiza el mensaje y extrae:
- Intención del usuario
- Entidad mencionada (cliente, equipo, servicio)
- Tipo de solicitud (consulta, acción, reporte)
```

### Paso 2: Clasificación
```
ION clasifica la consulta en una categoría:
- PROTOCOLO: ¿Cómo hacer algo? → agente-protocolos
- RED: ¿Problema de red? → agente-redes
- CLIENTE: ¿Info de cliente? → agente-clientes
- EQUIPO: ¿Info de equipo? → agente-inventario
- GENERAL: ¿Pregunta abierta? → agente-base
```

### Paso 3: Delegación
```
ION llama al agente apropiado usando call_agent
Pasa contexto relevante del mensaje original
Espera respuesta del agente
```

### Paso 4: Integración
```
ION recibe respuesta del agente
Integra en respuesta coherente
Agrega contexto adicional si es necesario
Devuelve al usuario
```

## Cómo hacer consultas efectivas

### Bueno ✓
```
"Cuál es el protocolo para reiniciar un CPE que no responde?"
"Qué equipo tiene el serial ION-2024-0543?"
"Estado de los servicios en la region norte"
```

### Mejorable ✗
```
"Help" (muy vago)
"Problema" (sin contexto)
"El cliente del ticket 123" (requiere más info)
```

## Configuración del orquestador

### Variables de entorno relevantes
```bash
IONET_ORCHESTRATOR_MODE=strict|flexible
IONET_DEFAULT_AGENT=agente-base
IONET_ROUTING_LOG=enabled
```

## Métricas de monitoreo

El orquestador puede reportar:
- Tiempo de respuesta promedio
- Distribución de consultas por agente
- Consultas no resueltas
- Intentos de escalamiento

## Casos de escalamiento

### Escalamiento a agente especializado
```
1. Usuario hace consulta general
2. ION detecta que requiere especialista
3. Llama a agente apropiado
4. Integra respuesta
```

### Escalamiento a humano
```
1. Consulta no puede ser resuelta por agentes
2. ION notifica al supervisor
3. Técnico humano recibe ticket
4. Respuesta integrate después
```

## Referencias

- [Flujos de trabajo](./references/workflows.md) - Patrones de coordinación
- [Configuración de agentes](./references/agent_config.md) - Cómo agregar/modificar agentes