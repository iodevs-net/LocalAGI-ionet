---
name: messaging-integration
description: Integración con plataformas de mensajería incluyendo Telegram, Email (SMTP/IMAP), Slack y Discord. Usa esta skill cuando necesites enviar mensajes, gestionar conversaciones multicanal, o coordinar notificaciones entre diferentes plataformas de comunicación. Soporta mensajes de texto, imágenes, documentos PDF, audio y respuestas automáticas.
---

# Messaging Integration

## Cuándo usar esta skill

**Usar cuando:**
- Necesites enviar mensajes a través de Telegram, Email, Slack o Discord
- Requieras gestionar conversaciones en múltiples plataformas
- Automatices notificaciones y alertas
- Configures respuestas automáticas basadas en disparadores
- Coordines comunicación entre equipos

**NO usar cuando:**
- La mensajería sea en tiempo real con requisitos de baja latencia estrictos
- Requieras funcionalidades de chat en vivo con agentes humanos
- Necesites integración con APIs propietarias específicas

## Acciones disponibles

### Email
| Acción | Descripción | Parámetros clave |
|--------|-------------|------------------|
| `send_email` | Enviar email via SMTP | to, subject, message |

### Telegram
| Acción | Descripción | Parámetros clave |
|--------|-------------|------------------|
| `send_telegram_message` | Enviar mensaje a Telegram | chat_id, message |

## Conectores disponibles

### Telegram
- **Tipo**: Receptor y emisor bidireccional
- **Características**: 
  - Soporta imágenes (fotos en mensajes)
  - Soporta mensajes de voz (transcripción automática)
  - Respuestas de audio (TTS)
  - Documentos PDF
  - Respuesta a mensajes de grupo (con mención)
- **Configuración**: Token de bot, admins (opcional), channel_id

### Slack
- **Tipo**: Receptor y emisor bidireccional
- **Características**:
  - Mensajes en canales y threads
  - Soporta imágenes, PDFs, canciones
  - Muestra URLs de búsqueda
  - Reemplaza menciones de usuario con nombres
- **Configuración**: App Token, Bot Token, Channel ID

### Discord
- **Tipo**: Receptor y emisor bidireccional
- **Características**:
  - Mensajes en canales y threads
  - Responde solo cuando se menciona o en canal default
- **Configuración**: Token del bot, Default Channel

### Email
- **Tipo**: Receptor (IMAP) y emisor (SMTP)
- **Características**:
  - Conversión HTML a Markdown para procesamiento
  - Soporta respuestas encadenadas (Reply All)
  - Incluye conversaciones previas en respuestas
  - Soporta contenido HTML en respuestas
- **Configuración**: SMTP Server, IMAP Server, credenciales

## Flujo de trabajo

### 1. Enviar mensaje a Telegram

```
1. Verificar que tienes el chat_id del destinatario
2. Preparar el mensaje ( Markdown soportado)
3. Usar send_telegram_message
4. Confirmar envío y proporcionar ID del mensaje
```

### 2. Enviar email con respuesta automática

```
1. Configurar el connector de email con credenciales SMTP/IMAP
2. El sistema recibe emails automáticamente
3. Procesa el contenido y genera respuesta
4. Envía respuesta al remitente original
```

### 3. Coordinar mensajes entre plataformas

```
1. Identificar la plataforma más apropiada para el destinatario
2. Usar la acción correspondiente (send_email, send_telegram_message)
3. Considerar formato y limitaciones de cada plataforma
4. Registrar la conversación en el tracker correspondiente
```

## Ejemplos de uso

### Ejemplo 1: Enviar mensaje por Telegram
```json
{
  "action": "send_telegram_message",
  "chat_id": 123456789,
  "message": "📢 *Notificación importante*\n\nEl servidor está experimentando problemas.\n\nPor favor, revisa los logs."
}
```

### Ejemplo 2: Enviar email
```json
{
  "action": "send_email",
  "to": "equipo@ejemplo.com",
  "subject": "Reporte Semanal",
  "message": "Adjunto el reporte semanal de métricas.\n\nSaludos"
}
```

## Notas de implementación

- **Límites de mensaje**: Telegram tiene límite de 3000 caracteres por mensaje
- **Formato Markdown**: Telegram usa MarkdownV2, Slack usa mrkdwn
- **Rate limits**: Cada plataforma tiene sus propios límites de tasa
- **Archivos**: Slack puede subir archivos locales (PDFs, imágenes)
- **Hilos**: Discord crea threads automáticamente para respuestas

## Referencias

- [Guía de conectores](./references/connectors.md) - Configuración detallada de cada conector
- [Tabla de compatibilidad](./references/platform_features.md) - Comparativa de características entre plataformas