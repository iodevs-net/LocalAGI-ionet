# 🎯 AGENTE DE VISIÓN - Documentación Técnica

## 📋 Resumen

Se ha implementado el **Agente de Visión** (`agente-vision`) en IONET, especializado en análisis de imágenes y visión por computadora. Este agente actúa como puente entre el contenido visual y los demás agentes (que no tienen capacidad de visión).

## 🏗️ Arquitectura

### Agente: `agente-vision`

**Archivo de configuración:** `./config/agents/agente-vision.json`

**Características principales:**
- **Modelo principal:** Usa el modelo base de ION (`inclusionai/ling-2.6-1t:free`)
- **Modelo multimodal:** `nvidia/nemotron-nano-12b-v2-vl:free` ✅
- **Función:** Analizar imágenes y extraer información visual
- **Salida:** JSON estructurado para otros agentes

## 🔧 Funcionalidades

### 1. OCR (Reconocimiento de Texto)
- Extrae texto de imágenes en múltiples idiomas
- Lee documentos, facturas, contratos
- Interpreta pantallas y displays

### 2. Análisis Visual
- Identifica objetos, personas, escenas
- Detecta códigos QR y barras
- Analiza diagramas, gráficos, tablas
- Evalúa estado de equipos (daños, mantenimiento)

### 3. Traducción a No-Visuales
- Convierte información visual a texto estructurado
- Describe detalladamente el contenido
- Extrae datos numéricos y tabulares
- Traduce diagramas a descripciones paso a paso

## 📄 Formato de Salida JSON

```json
{
  "tipo_contenido": "[texto|documento|diagrama|código|error|equipo]",
  "idioma_detectado": "[es|en|otro]",
  "texto_extraido": "Texto OCR completo",
  "elementos_visuales": [
    {
      "tipo": "[objeto|persona|documento|código]",
      "descripcion": "...",
      "confianza": "[alta|media|baja]"
    }
  ],
  "datos_estructurados": {
    // Datos clave extraídos
  },
  "resumen_descriptivo": "Descripción detallada de la imagen",
  "contexto_tecnico": "Contexto para otros agentes"
}
```

## 🔄 Flujo de Trabajo

### 1. Recepción
- El agente ION recibe una consulta con imagen
- Identifica que hay contenido visual
- Deriva a `agente-vision`

### 2. Análisis
- `agente-vision` usa el modelo multimodal
- Procesa la imagen completa
- Extrae toda la información relevante

### 3. Estructuración
- Formatea los datos en JSON
- Determina el tipo de contenido
- Evalúa nivel de confianza

### 4. Derivación
- Identifica qué agente debe manejar el caso
- Pasa la información estructurada
- Incluye contexto técnico

### 5. Resolución
- El agente correspondiente actúa
- Usa la información de la imagen
- Proporciona la solución

## 🎯 Matriz de Derivación

| Contenido Visual | Agente Destino | Motivo |
|-----------------|----------------|--------|
| Documentos, facturas, contratos | agente-clientes | Validación legal/comercial |
| Errores, logs, pantallas | agente-servicios | Diagnóstico técnico |
| Planos, diagramas de red | agente-redes | Infraestructura |
| Códigos, configuraciones | agente-protocolos | Procedimientos |
| Inventario, etiquetas | agente-inventario | Activos |
| Alertas de seguridad | agente-seguridad | Riesgo |
| Backups, archivos | agente-datos | Recuperación |
| General/ambiguo | agente-base | Consulta general |

## 🛠️ Uso Práctico

### Ejemplo 1: Factura
```
Usuario envía: Foto de una factura

Proceso:
1. agente-vision analiza la imagen
2. Extrae: Monto, fecha, concepto, RFC
3. Deriva a: agente-clientes
4. Respuesta: Validación y registro
```

### Ejemplo 2: Error de Sistema
```
Usuario envía: Captura de pantalla de error

Proceso:
1. agente-vision lee el mensaje de error
2. Extrae: Código, descripción, componente
3. Deriva a: agente-servicios
4. Respuesta: Solución o escalamiento
```

### Ejemplo 3: Diagrama
```
Usuario envía: Foto de un diagrama de red

Proceso:
1. agente-vision identifica elementos
2. Extrae: IPs, conexiones, equipos
3. Deriva a: agente-redes
4. Respuesta: Análisis de configuración
```

## ⚙️ Configuración del Sistema

### Modelos Configurados

**En `.env`:**
```bash
MODEL_NAME=inclusionai/ling-2.6-1t:free
MULTIMODAL_MODEL=nvidia/nemotron-nano-12b-v2-vl:free
OPENAI_API_KEY=sk-or-v1-...
```

**Agentes Disponibles (8):**
1. ION (Orchestrator)
2. agente-clientes
3. agente-servicios
4. agente-protocolos
5. agente-inventario
6. agente-seguridad
7. agente-redes
8. agente-datos
9. **agente-vision** ✨ (Nuevo)

## 🔍 Reglas del Agente de Visión

### Obligatorias:
1. ✅ SIEMPRE usar el modelo multimodal para imágenes
2. ✅ SIEMPRE proporcionar formato JSON estructurado
3. ✅ EXTRAER todo el texto posible (OCR)
4. ✅ MANTENER consistencia con otros agentes
5. ✅ INDICAR confianza/incertidumbre

### Prohibidas:
1. ❌ Nunca alucinar información
2. ❌ Nunca asumir sin evidencia visual
3. ❌ Nunca omitir incertidumbre
4. ❌ Nunca saltar el formato JSON

## 📊 Ejemplos de Uso

### Caso 1: Factura para Validar

**Input:** Imagen de factura

**Proceso:**
```
1. ION recibe: "¿Esta factura es válida?" + [imagen]
2. ION deriva a: agente-vision
3. agente-vision analiza:
   - Tipo: documento
   - Texto: Extrae todos los datos
   - Datos: {monto: 5000, fecha: 2024-01-15, cliente: "ACME"}
4. Deriva a: agente-clientes
5. Respuesta: "Factura válida, registrada en sistema"
```

### Caso 2: Código de Error

**Input:** Captura de pantalla con error 500

**Proceso:**
```
1. ION recibe: "¿Qué significa este error?" + [captura]
2. ION deriva a: agente-vision
3. agente-vision analiza:
   - Tipo: error
   - Texto: "HTTP 500 - Internal Server Error"
   - Contexto: Apache/2.4.41, PHP
4. Deriva a: agente-servicios
5. Respuesta: "Error interno, revisar logs de PHP"
```

### Caso 3: Diagrama Técnico

**Input:** Foto de diagrama de red

**Proceso:**
```
1. ION recibe: "¿Cómo mejorar esta red?" + [diagrama]
2. ION deriva a: agente-vision
3. agente-vision analiza:
   - Tipo: diagrama
   - Elementos: Router, Switch, 3 PCs
   - Topología: Estrella
4. Deriva a: agente-redes
5. Respuesta: "Agregar redundancia con segundo switch"
```

## 🎨 Consideraciones Técnicas

### Rendimiento
- **Modelo multimodal:** Optimizado para eficiencia
- **Tamaño:** nano (12B parámetros)
- **Tiempo de respuesta:** 3-5 segundos promedio
- **Requisitos:** GPU/CPU moderada

### Calidad
- **Precisión OCR:** Alta para texto claro
- **Reconocimiento:** Bueno para elementos principales
- **Confianza:** Siempre indicada en la salida

### Limitaciones
- Modelo gratuito: rendimiento moderado
- Imágenes complejas: enfocarse en lo principal
- Texto pequeño: puede perder detalles
- Elementos muy densos: resumen general

## 🚀 Integración con IONET

### En ION (Orchestrator)

La matriz de derivación actualizada:

```
| Si el usuario menciona... | Derivar a... |
|---------------------------|--------------|
| ...                       | ...          |
| **imagen, foto, captura, visual, diagrama** | **agente-vision** |
| ...                       | ...          |
```

### En el Prompt de ION

```
## TÚ ERES LA PUERTA DE ENTRADA

Cada usuario que llega a IONET llega primero a ti. Tu trabajo es:
1. Entender qué necesita (incluyendo imágenes)
2. SI hay imagen: derivar a agente-vision
3. SI NO hay imagen: procesar normalmente
4. Encontrar la mejor respuesta posible
5. Hacer que se sienta bien atendido
6. Escalar cuando sea necesario
```

## 📈 Métricas Esperadas

| Métrica | Objetivo |
|---------|----------|
| Tiempo de análisis | < 5 segundos |
| Precisión OCR | > 90% |
| Derivación correcta | > 95% |
| Respuestas útiles | > 90% |

## 🔧 Mantenimiento

### Monitoreo
```bash
# Verificar logs del agente
docker compose -f docker-compose.dev.yaml logs agente-vision

# Ver uso de recursos
docker stats
```

### Actualizaciones
```bash
# Reiniciar si hay cambios
docker compose -f docker-compose.dev.yaml up -d --force-recreate
```

## 🎓 Conclusión

El **Agente de Visión** completa el ecosistema de IONET al agregar:
- ✅ Capacidad de procesar imágenes
- ✅ OCR para extraer texto
- ✅ Análisis visual automatizado
- ✅ Traducción a texto para otros agentes
- ✅ Integración fluida con ION

**Resultado:** IONET ahora puede "ver" y entender contenido visual, permitiendo a los usuarios enviar fotos, capturas, diagramas y documentos para su procesamiento automático.

---

## 📝 Historial de Cambios

- **v1.0** - Implementación inicial del agente de visión
  - Modelo multimodal: nvidia/nemotron-nano-12b-v2-vl:free
  - Integración completa con ION
  - 8 agentes totales en el sistema

