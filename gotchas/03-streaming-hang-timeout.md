# Gotcha 03: Streaming Hang por Falta de Timeout HTTP

## Sintoma
El agente se colgaba 120s en cada llamada streaming a OpenRouter.
`decisionWithStreaming` → `localai_stream` quedaba bloqueado hasta que
OpenRouter eventualmente respondia o el contexto expiraba.

## Causa Raiz
`localai_client.go` usaba `http.DefaultClient` que NO tiene timeout.
El `bufio.Scanner` del SSE reader no tenia proteccion contra respuestas
lentas del servidor. OpenRouter free tier es impredecible.

## Solucion
Timeout explicito de 60s en el HTTP client:

```go
client: &http.Client{Timeout: 60 * time.Second}
```

## Leccion
`http.DefaultClient` en produccion es un bug. Siempre configurar timeout.
60s es balance entre paciencia y contencion de dano.
