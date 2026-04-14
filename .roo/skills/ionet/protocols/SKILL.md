---
name: ionet-protocols
description: Acceso a protocolos y procedimientos internos de IONET. Usa esta skill cuando necesites consultar procedimientos de trabajo, configuraciones estándar, o guías de operación. El sistema utiliza RAG (Retrieval Augmented Generation) para encontrar protocolos relevantes basados en la consulta del usuario. Esta skill es específica del proyecto IONET y contiene conocimiento interno sobre procesos de trabajo de la empresa.
---

# IONET Protocols - Agente Especializado

## Cuándo usar esta skill

**Usar cuando:**
- Necesites consultar procedimientos de trabajo de IONET
- Requieras información sobre configuraciones estándar de equipos
- Busques guías de operación para equipos específicos
- Necesites resolver dudas sobre流程 de trabajo interno
- Quieras documentar nuevos procedimientos

**NO usar cuando:**
- La consulta sea sobre temas externos a IONET
- Requieras información de redes general (usar agente-redes)
- Necesites datos de clientes específicos (usar agente-clientes)

## Agente especializado

| Campo | Valor |
|-------|-------|
| Nombre | `agente-protocolos` |
| Tipo | Agente RAG especializado |
| Base de conocimiento | Protocolos y procedimientos internos |

## Capacidades

### Consultar protocolos existentes
```
1. Formular pregunta específica sobre el protocolo buscado
2. El sistema busca en el índice RAG
3. Devuelve el protocolo más relevante
4. Si hay múltiples resultados, presenta opciones
```

### Mantener actualizados los protocolos
```
1. Cuando se identifique un nuevo procedimiento
2. Documentar usando add_to_memory
3. Incluir namespace "protocolos" para organización
4. Verificar que sea accesible mediante búsqueda
```

## Protocolos comunes en IONET

### Gestión de CPE (Customer Premises Equipment)

#### Reinicio de CPE
```
1. Acceder via Winbox/Webfig al router del cliente
2. Ir a System > Reboot
3. Confirmar reinicio
4. Esperar 2-3 minutos
5. Verificar conexión desde el NOC
```

#### Cambio de firmware
```
1. Descargar firmware compatible de MikroTik
2. Subir firmware al router via Files
3. Ir a System > Packages
4. Upload nuevo paquete
5. Reiniciar router
6. Verificar funcionamiento
```

### Configuración de servicios

#### Activar servicio de internet
```
1. Verificar línea en OLT (si es fibra)
2. Crear PPPoE credentials en RADIUS
3. Configurar router del cliente
4. Verificar conexión y velocidad
5. Documentar en sistema de tickets
```

#### Cambio de plan
```
1. Verificar elegibilidad del cliente
2. Actualizar quotas en RADIUS
3. Ajustar velocidad en router
4. Confirmar con test de velocidad
5. Actualizar contrato
```

## Flujo de trabajo

### 1. Consultar procedimiento

```
1. El orquestador (ION) recibe la consulta
2. Identifica que es sobre protocolos
3. Delega a agente-protocolos
4. agente-protocolos busca en RAG
5. Devuelve respuesta con protocolo o guía
```

### 2. Resolver duda técnica

```
1. Técnico hace pregunta al orquestador
2. ION determina que requiere protocolo específico
3. Consulta a agente-protocolos
4. Devuelve procedimiento paso a paso
5. Técnico sigue los pasos
```

## Ejemplos de uso

### Ejemplo: Consultar protocolo de reinicio
```
Entrada: "Cómo reinicio un equipo Mikrotik que no responde?"

El agente responde con el procedimiento de reinicio estándar
```

### Ejemplo: Configurar QoS
```
Entrada: "Cuál es el procedimiento para configurar QoS en un router MikroTik?"

El agente busca y devuelve configuración estándar de QoS
```

## Métricas de uso

El agente puede devolver información sobre:
- Frecuencia de uso de cada protocolo
- Número de consultas por categoría
- Protocolos más buscados

## Referencias

- [Catálogo de equipos](./references/equipos.md) - Lista de equipos y sus configuraciones
- [Procedimientos estándar](./references/procedimientos.md) - Documentación detallada