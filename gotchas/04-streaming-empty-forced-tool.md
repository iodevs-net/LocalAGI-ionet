# Gotcha 04: Forced Tool Retorna Vacio, Tratado como Exito

## Sintoma
`pickTool` con `forceTool` forzaba un tool call pero el streaming retornaba
0 tool calls. El codigo trataba esto como exito, produciendo respuestas
vACIAS o comportamiento inesperado.

## Causa Raiz
`tools.go` linea ~326: cuando `forceTool != ""` y `len(toolCalls) == 0`,
el codigo retornaba sin error. Deberia reintentar.

## Solucion
Retry cuando el streaming no produce tool calls para el tool forzado:

```go
if forceTool != "" && len(toolCalls) == 0 {
    return fmt.Errorf("streaming returned no tool calls for forced tool %q (content: %q)",
        forceTool, content)
}
```

## Leccion
Streaming puede fallar silenciosamente (tool call omitido). Validar
que el resultado del streaming contenga lo requerido antes de continuar.
