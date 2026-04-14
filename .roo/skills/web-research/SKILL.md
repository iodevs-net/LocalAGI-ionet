---
name: web-research
description: Investigación web incluyendo búsqueda en internet, scraping de páginas web, consulta a Wikipedia y navegación automatizada. Usa esta skill cuando necesites obtener información actualizada, analizar contenido de páginas web, buscar datos específicos en la web, o construir conocimiento basado en fuentes externas. Incluye acciones para search_internet, scrape, wikipedia.
---

# Web Research

## Cuándo usar esta skill

**Usar cuando:**
- Necesites buscar información actualizada en internet
- Requieras analizar contenido de páginas web específicas
- Busques definiciones o artículos en Wikipedia
- Necesites agregar enlaces y referencias a respuestas
- Construyas conocimiento basado en fuentes verificables

**NO usar cuando:**
- La información sea sensible o requiera acceso autenticado
- Necesites datos en tiempo real (precios, stocks)
- Requieras interacción con APIs propietarias
- El contenido esté protegido por CAPTCHA o similares

## Acciones disponibles

### Búsqueda
| Acción | Descripción | Parámetros clave |
|--------|-------------|------------------|
| `search_internet` | Buscar en DuckDuckGo | query |

### Scraping
| Acción | Descripción | Parámetros clave |
|--------|-------------|------------------|
| `scrape` | Obtener contenido completo de página web | url |

### Wikipedia
| Acción | Descripción | Parámetros clave |
|--------|-------------|------------------|
| `wikipedia` | Buscar artículos en Wikipedia | query |

## Flujo de trabajo

### 1. Búsqueda básica

```
1. Formular query de búsqueda clara y específica
2. Usar search_internet para obtener resultados
3. Analizar resultados y URLs devueltas
4. Si es necesario, usar scrape para obtener contenido detallado
5. Sintetizar información y proporcionar respuesta con referencias
```

### 2. Investigación profunda

```
1. Usar search_internet con query específica
2. Recopilar URLs relevantes de los resultados
3. Para cada URL importante, usar scrape
4. Extraer información clave de cada fuente
5. Comparar y validar información entre fuentes
6. Proporcionar respuesta citing fuentes
```

### 3. Consulta a Wikipedia

```
1. Formular query para Wikipedia
2. Usar wikipedia para obtener resumen
3. Si se necesita más detalle, usar scrape con enlace de Wikipedia
4. Proporcionar respuesta con contexto adicional
```

## Ejemplos de uso

### Ejemplo 1: Búsqueda web
```json
{
  "action": "search_internet",
  "query": "mejores prácticas CI/CD 2024"
}
```

### Ejemplo 2: Scraping de página
```json
{
  "action": "scrape",
  "url": "https://docs.docker.com/compose/"
}
```

### Ejemplo 3: Consulta Wikipedia
```json
{
  "action": "wikipedia",
  "query": "Kubernetes container orchestration"
}
```

## Notas de implementación

- **Resultados de búsqueda**: Devuelve URLs únicas procesadas
- **Scraping**: Obtiene todo el contenido de texto de la página
- **Wikipedia**: Devuelve resumen y enlaces relacionados
- **Rate limits**: DuckDuckGo tiene límites de uso
- **URLs de resultados**: Las URLs pueden contener跟踪 parámetros que son eliminados

## Referencias

- [Guía de búsqueda](./references/search_guide.md) - Mejores prácticas para queries
- [Técnicas de scraping](./references/scraping_tips.md) - Cómo obtener mejores resultados