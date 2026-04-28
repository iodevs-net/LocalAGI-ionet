# Gotcha 06: Log de Exito SMTP Dentro del Bloque Error

## Sintoma
No se veia log de "Email sent successfully" en los logs. Imposible
confirmar si el SMTP habia enviado el correo o no.

## Causa Raiz
Bug de indentacion/logica en `sendMail()`: el `xlog.Info` de exito
estaba dentro del bloque `if err != nil {}`:

```go
if err != nil {
    xlog.Error(fmt.Sprintf("Email send err: %v", err))
    xlog.Info(fmt.Sprintf("Email sent successfully to %v", emails)) // <-- AQUI
}
```

Si el envio EXITOSO: no habia log. Si FALLABA: log de error Y falso exito.

## Solucion
Mover el log de exito al bloque `else`:

```go
if err != nil {
    xlog.Error(fmt.Sprintf("Email send err: %v", err))
} else {
    xlog.Info(fmt.Sprintf("Email sent successfully to %v", emails))
}
```

## Leccion
Siempre verificar que logs de exito/error esten en los bloques correctos.
El formateo puede enganar al ojo. Probar ambos caminos.
