# Gotcha 07: ion.json Corrupto por Re-serializacion Parcial

## Sintoma
ION no cargaba. Log: `agent file missing name, skipping file="ion.json"`.
El archivo existia pero el parseo fallaba.

## Causa Raiz
`json.dump()` de Python escribio un dict parcial (sin `name`, `description`,
`model`, `multimodal_model`) sobre el archivo original. Estos campos son
requeridos por el pool loader.

## Solucion
Reconstruir el JSON completo desde git history con Python, preservando
todos los campos requeridos.

## Leccion
Nunca sobrescribir archivos de config con serializacion parcial.
Validar JSON contra el schema esperado antes de escribir.
Hacer backup antes de modificar.
