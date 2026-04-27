# Gotcha 02: Email Autoreply Loop

## Sintoma
ION se enviaba autoreplies a si mismo infinitamente. Cada respuesta del agente
generaba un nuevo correo en INBOX que volvia a procesar.

## Causa Raiz
`processEmail()` no verificaba si el correo provenia de la propia cuenta IMAP.
Cuando el agente respondia a un correo desde ventas@ionet.cl, el SMTP enviaba
el reply, pero si el reply llegaba al INBOX (e.g., si el destinatario era
la misma cuenta IMAP), se procesaba como un nuevo mensaje entrante.

## Solucion
Filtrar correos donde `From` contenga el username IMAP:

```go
if strings.Contains(msg.Header.Get("From"), e.username) {
    xlog.Debug("Email from self, skipping to prevent loop")
    return
}
```

## Leccion
Siempre prevenir loops en sistemas que leen y escriben al mismo canal.
Check de identidad propia antes de procesar cualquier mensaje.
