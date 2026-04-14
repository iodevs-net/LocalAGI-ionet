---
name: system-operations
description: Operaciones de sistema incluyendo ejecución de comandos shell via SSH, diagnóstico de infraestructura y gestión de procesos. Usa esta skill cuando necesites ejecutar comandos en servidores remotos, diagnosticar problemas de sistema, realizar operaciones de mantenimiento, o automatizar tareas de infraestructura. Soporta autenticación por password y clave privada SSH. Incluye acción run_command.
---

# System Operations

## Cuándo usar esta skill

**Usar cuando:**
- Necesites ejecutar comandos en servidores remotos
- Requieras diagnosticar problemas de red o sistema
- Realices tareas de mantenimiento automatizado
- Monitorees servicios y procesos
- Gestiones configuraciones de servidores

**NO usar cuando:**
- El servidor no tenga acceso SSH disponible
- Requieras interfaz gráfica para la操作
- La tarea sea de muy bajo nivel (firmware, bootloader)
- Necesites ejecutar comandos con privilegios elevados frecuentemente

## Acciones disponibles

| Acción | Descripción | Parámetros clave |
|--------|-------------|------------------|
| `run_command` | Ejecutar comando en servidor | command, host, user |

## Configuración SSH

### Por password
```json
{
  "user": "admin",
  "password": "password123",
  "host": "192.168.1.100:22"
}
```

### Por clave privada
```json
{
  "user": "admin",
  "privateKey": "-----BEGIN RSA PRIVATE KEY-----\n...",
  "host": "192.168.1.100:22"
}
```

## Flujo de trabajo

### 1. Ejecutar comando simple

```
1. Verificar parámetros de conexión (host, user)
2. Preparar comando a ejecutar
3. Usar run_command
4. Analizar salida del comando
5. Reportar resultados
```

### 2. Diagnóstico de servidor

```
1. Verificar conectividad (ping)
2. Revisar uso de recursos (top, df, free)
3. Revisar logs del sistema
4. Verificar servicios activos
5. Documentar hallazgos
```

### 3. Mantenimiento automatizado

```
1. Ejecutar comandos de diagnóstico previos
2. Realizar acciones de mantenimiento
3. Verificar que los cambios fueron aplicados
4. Ejecutar comandos de verificación post-cambio
5. Reportar resultado final
```

## Ejemplos de uso

### Ejemplo 1: Ver estado de servicios
```json
{
  "action": "run_command",
  "command": "systemctl status nginx --no-pager",
  "host": "192.168.1.100:22",
  "user": "admin"
}
```

### Ejemplo 2: Ver uso de disco
```json
{
  "action": "run_command",
  "command": "df -h"
}
```

### Ejemplo 3: Reiniciar servicio
```json
{
  "action": "run_command",
  "command": "sudo systemctl restart docker"
}
```

### Ejemplo 4: Ver logs recientes
```json
{
  "action": "run_command",
  "command": "journalctl -n 50 --no-pager"
}
```

## Comandos útiles

### Diagnóstico de red
```bash
# Ping
ping -c 4 8.8.8.8

# Trace route
traceroute google.com

# DNS lookup
nslookup example.com

# Puertos abiertos
netstat -tulpn | grep LISTEN
```

### Monitoreo de recursos
```bash
# CPU y memoria
top -bn1

# Procesos
ps aux --sort=-%cpu | head

# Espacio en disco
df -h

# Uso de memoria
free -h
```

### Gestión de servicios
```bash
# Ver estado
systemctl status nginx

# Reiniciar
sudo systemctl restart nginx

# Ver logs
journalctl -u nginx -n 100
```

## Notas de implementación

- **SSH**: Usa golang.org/x/crypto/ssh para conexiones
- **Autenticación**: Soporta password y clave privada RSA
- **Host key**: Ignora verificación de host (InsecureIgnoreHostKey)
- **Output**: Devuelve stdout y stderr combinados
- **Errores**: Incluye mensajes de error si el comando falla

## Referencias

- [Comandos comunes](./references/common_commands.md) - Comandos frecuentemente usados
- [Diagnóstico](./references/diagnostics.md) - Guías de troubleshooting