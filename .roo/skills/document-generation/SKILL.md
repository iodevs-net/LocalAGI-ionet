---
name: document-generation
description: Generación de documentos incluyendo PDFs y imágenes. Usa esta skill cuando necesites crear documentos PDF formateados, generar imágenes con IA, o producir contenido visual. El sistema soporta markdown en PDFs (encabezados, listas, código, negrita) y generación de imágenes via API (DALL-E compatible). Los archivos generados pueden ser enviados automáticamente por conectores como Slack o Telegram. Incluye acciones para generate_pdf y generate_image.
---

# Document Generation

## Cuándo usar esta skill

**Usar cuando:**
- Necesites crear documentos PDF profesionales
- Requieras generar imágenes con descripciones textuales
- Quieras formatear contenido como documento exportable
- Necesites crear reportes, facturas, certificados
- Requieras contenido visual para presentaciones

**NO usar cuando:**
- El documento requiera formato muy específico (usar herramienta专门)
- Necesites edición de documentos existente
- Requieras gráficos复杂的 (usar herramientas especializadas)
- La generación de imágenes sea para contenido crítico (altos costos)

## Acciones disponibles

| Acción | Descripción | Parámetros clave |
|--------|-------------|------------------|
| `generate_pdf` | Generar documento PDF | title, content, filename |
| `generate_image` | Generar imagen con IA | prompt, size |

## Generación de PDF

### Características soportadas
- **Título**: Encabezado del documento (opcional)
- **Contenido**: Markdown completo con:
  - Encabezados (# ## ###)
  - Listas (ordenadas y no ordenadas)
  - Código (inline y bloques)
  - Negrita y cursiva
  - Links
- **Filename**: Nombre personalizado con extensión .pdf (opcional)

### Configuración
```json
{
  "outputDir": "/app/pdfs",
  "cleanOnStart": false
}
```

## Generación de imágenes

### Modelos disponibles
- `dall-e-3` (default)
- Modelos compatibles con API OpenAI-style

### Tamaños disponibles
- 256x256
- 512x512
- 1024x1024 (default)

### Configuración
```json
{
  "apiKey": "sk-...",
  "apiURL": "https://api.openai.com/v1",
  "model": "dall-e-3"
}
```

## Flujo de trabajo

### 1. Generar PDF simple

```
1. Preparar título y contenido en Markdown
2. Definir filename opcional
3. Usar generate_pdf
4. Confirmar path del archivo generado
5. Opcional: usar acción de envío para compartir
```

### 2. Generar reporte complejo

```
1. Recopilar datos a incluir en el reporte
2. Formatear como Markdown con encabezados y listas
3. Incluir tablas si es necesario
4. Usar generate_pdf con contenido completo
5. Distribuir via email o Telegram
```

### 3. Generar imagen

```
1. Formular prompt descriptivo para la imagen
2. Seleccionar tamaño apropiado
3. Usar generate_image
4. Recibir URL de la imagen generada
5. Usar en respuestas o enviar via conector
```

## Ejemplos de uso

### Ejemplo 1: PDF con título
```json
{
  "action": "generate_pdf",
  "title": "Reporte de Tickets - Enero 2024",
  "content": "# Resumen\n\n## Tickets abiertos: 45\n## Tickets cerrados: 38\n\n## Detalles\n\n1. Ticket #123 - CPE no responde\n2. Ticket #124 - Slow connection\n\n---\nGenerado automáticamente por IONET",
  "filename": "reporte-enero-2024.pdf"
}
```

### Ejemplo 2: PDF sin título
```json
{
  "action": "generate_pdf",
  "content": "## Protocolo de Reinicio de CPE\n\n### Pasos:\n1. Acceder a Routerboard\n2. Ir a System > Reboot\n3. Esperar 2 minutos\n4. Verificar conexión\n\n### Notas:\n- No presionar reset físico\n- Documentar hora del reinicio",
  "filename": "protocolo-cpe.md.pdf"
}
```

### Ejemplo 3: Generar imagen pequeña
```json
{
  "action": "generate_image",
  "prompt": "Diagrama de red mostrando routers conectados en estrella",
  "size": "256x256"
}
```

### Ejemplo 4: Generar imagen grande
```json
{
  "action": "generate_image",
  "prompt": "Logo para empresa de telecomunicaciones, estilo moderno azul y blanco",
  "size": "1024x1024"
}
```

## Integración con conectores

### Telegram
Los PDFs generados son enviados automáticamente como documentos.

### Slack
Los PDFs e imágenes son subidos al canal/thread correspondiente.

## Notas de implementación

- **PDF output**: Guardado en directorio configurado (default: /app/pdfs)
- **Seguridad filename**: Se limpia el filename para prevenir path traversal
- **Extensión**: Se asegura que el filename termine en .pdf
- **Imágenes**: URL devuelta puede ser temporal (DALL-E)
- **Límites**: Rate limits del proveedor de API aplican

## Casos de uso para IONET

### Reportes de tickets
```
generate_pdf con estadísticas de tickets del mes
```

### Protocolos documentados
```
generate_pdf con procedimientos paso a paso
```

### Diagramas de red
```
generate_image con descripción de topología
```

### Certificados
```
generate_pdf con formato formal para certificaciones
```

## Referencias

- [Guía de Markdown](./references/markdown_guide.md) - Sintaxis soportada en PDFs
- [Prompts efectivos](./references/prompt_tips.md) - Cómo formular prompts para imágenes