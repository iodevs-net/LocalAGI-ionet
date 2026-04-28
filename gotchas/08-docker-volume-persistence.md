# Gotcha 08: Config No Persiste en Volumen Docker

## Sintoma
Despues de rebuildear la imagen, los cambios en `config/agents/` no
se reflejaban dentro del contenedor. El volumen persistente mantenia
la version anterior.

## Causa Raiz
`docker-compose.dev.yaml` monta `ionet_pool_dev:/pool` como volumen.
El volumen persiste entre rebuilds y su contenido prevalece sobre
el `COPY` del Dockerfile. El binding mount `.:/app` tampoco actualiza
`/pool/agents/` porque es un volumen separado.

## Solucion
Forzar actualizacion via `docker cp`:
```bash
docker cp config/agents/ion.json ionet-dev:/pool/agents/ion.json
```

O eliminar el volumen para que se regener desde cero:
```bash
docker volume rm ionet_pool_dev
```

## Leccion
Volumen Docker + COPY = el volumen gana. No asumir que rebuildear
actualiza datos persistentes. Estrategia: init container, entrypoint
script, o `docker cp` explicito.
