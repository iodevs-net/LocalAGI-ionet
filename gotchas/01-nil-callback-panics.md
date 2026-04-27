# Gotcha 01: Nil Callback Panics

## Sintoma
SIGSEGV al procesar email. ION respondia pero el proceso crasheaba:
```
panic: runtime error: invalid memory address or nil pointer dereference
github.com/mudler/LocalAGI/core/state.(*AgentPool).startAgentWithConfig.func1()
    /build/core/state/pool.go:468
```

## Causa Raiz
`TeamsConnector.AgentReasoningCallback()` retornaba `nil`. El pool iteraba
connectors y llamaba `c.AgentReasoningCallback()(state)` sin nil-check.

Mismo patron para `AgentResultCallback()` que retornaba `nil` y crasheaba
en pool.go:514.

## Solucion
Dos capas de defensa:
1. **Origen**: `teams.go` — retornar funcion vacia no `nil`
2. **Defensa**: `pool.go` — nil-check antes de invocar callback

```go
// teams.go
func (t *TeamsConnector) AgentReasoningCallback() func(state types.ActionCurrentState) bool {
    return func(state types.ActionCurrentState) bool { return true }
}
func (t *TeamsConnector) AgentResultCallback() func(state types.ActionState) {
    return func(state types.ActionState) {}
}

// pool.go
for _, c := range connectors {
    if cb := c.AgentReasoningCallback(); cb != nil {
        if !cb(state) {
            return false
        }
    }
}
for _, c := range connectors {
    if cb := c.AgentResultCallback(); cb != nil {
        cb(state)
    }
}
```

## Leccion
Todo callback devuelto por interfaz debe tener nil-check en el llamante.
El contrato "puede retornar nil" no es obvio. Defensa en profundidad.
