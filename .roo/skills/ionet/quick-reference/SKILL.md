---
name: ionet-quick-reference
description: Referencia rápida para técnicos de IONET. Esta skill proporciona una guía concisa de los comandos más usados, atajos mentales, y recordatorios importantes para el trabajo diario. Diseñada para consulta rápida durante operativos, sin necesidad de leer documentación extensa.
---

# IONET Quick Reference - Guía de Atajos

## Comandos frecuentes

### Gestión de equipos

| Acción | Comando/IP | Notas |
|--------|------------|-------|
| Acceso router | Winbox/Webfig | IP del cliente |
| Ping desde NOC | `ping IP_CLIENTE` | Verificar conectividad |
| Wake on LAN | Por MAC address | En caso de equipos apagados |

### Tickets

| Estado | Significado | Acción |
|--------|-------------|--------|
| Abiertos | Pendientes | Priorizar por SLA |
| En proceso | Asignados | Continuar trabajo |
| Resueltos | Solucionados | Verificar con cliente |
| Cerrados | Completados | Archivar |

### Servicios

| Servicio | Puerto | Verificación |
|----------|--------|--------------|
| PPPoE | 1723 | Verificar credenciales |
| DHCP | 67/68 | Check pool disponible |
| DNS | 53 | Test de resolución |
| NTP | 123 | Sincronización hora |

## Atajos mentales

### Problema: Cliente sin internet
```
1. Ping gateway - ¿responde?
2. Ping DNS - ¿responde?
3. Verificar PPPoE - ¿conectado?
4. Revisar señal - ¿niveles OK?
5. Resetear equipo si necesario
```

### Problema: Slow connection
```
1. Hacer speedtest
2. Verificar interferencia WiFi
3. Check QoS en router
4. Verificar límite de ancho de banda
5. Revisar equipos en la misma red
```

### Problema: Equipo no responde
```
1. Ping - ¿llega?
2. Snmpwalk - ¿responde SNMP?
3. Winbox - ¿acceso?
4. Reset físico - ¿último recurso?
5. Reemplazo - ¿requiere RMA?
```

## Números importantes

| Tipo | Número | Contacto |
|------|--------|----------|
| NOC principal | *ext* | 24/7 |
| Support nivel 2 | *ext* | Horario hábil |
| Proveedor upstream | *número* | Para fallas mayores |
| Emergency fiber | *número* | Roturas de fibra |

## Checklist de verificación

### Antes de escalar a nivel 2
- [ ] Ping a gateway OK
- [ ] DNS funciona
- [ ] PPPoE conectado
- [ ] Logs revisados
- [ ] Equipo reiniciado
- [ ] Cliente notificado

### Antes de crear ticket de vendor
- [ ] Equipo probado en diferentes condiciones
- [ ] Firmware actualizado
- [ ] Configuración verificada
- [ ] Logs recolectados
- [ ] SLA del cliente verificado

## Formatos de tickets estándar

### Ticket de equipo fallido
```
Cliente: [NOMBRE]
Equipo: [SERIAL/IP]
Síntoma: [DESCRIPCIÓN]
Intentos: [PASOS REALIZADOS]
Logs: [ATTACH]
```

### Ticket de performance
```
Cliente: [NOMBRE]
Ubicación: [SECTOR]
Velocidad medida: [DOWN/UP]
Velocidad contratada: [DOWN/UP]
Hora de medición: [HH:MM]
```

## Acciones rápidas

### Para usar en chat

```
- "protocolo reinicio" → Procedure de reinicio
- "status servicios" → Estado actual  
- "cliente [nombre]" → Info de cliente
- "equipo [serial]" → Info de equipo
```

## Notas de emergencia

### Línea de backup
Si falla OLT principal:
1. Verificar configuración alternativa
2. Activar puerto backup
3. Notificar a clientes afectados
4. Proceder con repair

### Fuentes de poder
Si hay corte eléctrico:
1. Verificar generador en sitio
2. Confirmar autonomía de UPS
3. Priorizar equipos críticos
4. Plan de comunicación a clientes

## Referencia rápida de LEDs

| LED | Estado | Significado |
|-----|--------|-------------|
| PWR | Verde fijo | Encendido OK |
| PWR | Rojo | Error de poder |
| LINK | Verde | Conexión establecida |
| LINK | Ámbar | Negociando |
| LINK | Apagado | Sin link |
| ACT | Parpadeo | Actividad de datos |