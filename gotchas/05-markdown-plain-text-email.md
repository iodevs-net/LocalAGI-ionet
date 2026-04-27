# Gotcha 05: Markdown en Correos Plain Text

## Sintoma
Los correos enviados por ION contenian sintaxis markdown visible:
`**Conclusión:**`, `- listas`, `## titulos`. El cliente de correo
mostraba el markup en lugar de renderizarlo.

## Causa Raiz
El agente responde en markdown (formato natural del LLM). El connector
solo convertia a HTML si el email original era HTML (`contentIsHTML`).
Para correos plain text, el reply se enviaba como text/plain con
markdown crudo.

## Solucion
Siempre convertir la respuesta del agente de markdown a HTML:

```go
// Convert agent markdown response to HTML for clean email rendering
p := parser.NewWithExtensions(parser.CommonExtensions | ...)
doc := p.Parse([]byte(replyContent))
opts := html.RendererOptions{Flags: html.CommonFlags | html.HrefTargetBlank}
renderer := html.NewRenderer(opts)
replyContent = string(markdown.Render(doc, renderer))
contentIsHTML = true
```

## Leccion
LLMs responden en markdown por defecto. Si el output va a un medio que
no lo renderiza (email, SMS, consola), convertir explicitamente.
