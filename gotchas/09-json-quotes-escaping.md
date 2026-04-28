# Gotcha 09: JSON Quotes Escaping

## Sintoma
JSON invalido al editar system prompt. Error: `Expecting ',' delimiter`.

## Causa Raiz
Se inserto texto con comillas dobles sin escapar dentro de un string JSON:
```json
"system_prompt": "...Usa \"tú\" no \"vos\"..."
```
Al editar el JSON en el archivo, las comillas literales rompen el parseo.

## Solucion
Usar comillas angulares «» en lugar de dobles dentro de strings JSON:
```json
"system_prompt": "...Usa «tú» no «vos»..."
```

O escapar: `\"tú\"` no `"tú"`.

## Leccion
Toda comilla doble dentro de un string JSON debe escaparse (\"). Al editar
archivos JSON manualmente o con herramientas de texto, verificar que el
resultado sea parseable despues de cada cambio.
