# Gotchas — ION Agent

Problemas encontrados durante desarrollo y sus soluciones documentadas.
Cada archivo cubre un incidente con: sintoma, causa raiz, solucion, leccion.

## Indice

- [01-nil-callback-panics.md](01-nil-callback-panics.md) — SIGSEGV en pool.go por callbacks nil
- [02-email-autoreply-loop.md](02-email-autoreply-loop.md) — Loop infinito de autoreply
- [03-streaming-hang-timeout.md](03-streaming-hang-timeout.md) — Streaming colgado 120s por falta de timeout HTTP
- [04-streaming-empty-forced-tool.md](04-streaming-empty-forced-tool.md) — forced tool retorna vacio y se trata como exito
- [05-markdown-plain-text-email.md](05-markdown-plain-text-email.md) — Correos con sintaxis markdown visible
- [06-logging-bug-sendmail.md](06-logging-bug-sendmail.md) — Log de exito SMTP dentro del bloque error
- [07-ion-json-corrupted.md](07-ion-json-corrupted.md) — ion.json invalido por re-serializacion parcial
- [08-docker-volume-persistence.md](08-docker-volume-persistence.md) — Config no persiste en volumen Docker
- [09-json-quotes-escaping.md](09-json-quotes-escaping.md) — JSON invalido por comillas dobles sin escapar
