# ION - Quick Start para Técnicos

⚡ **5 minutos para tener ION funcionando**

---

## REQUISITOS PREVIOS

- Docker instalado
- Git instalado
- Acceso a repositorio IONET

---

## PASO 1: Clonar y Configurar (2 minutos)

```bash
# Clonar repositorio
git clone https://github.com/ionet-cl/agentes-ionet.git
cd agentes-ionet

# Copiar configuración
cp .env.example .env
```

---

## PASO 2: Configurar 3 Variables Obligatorias (3 minutos)

Edita `.env` con tu editor favorito (nano, vim, VSCode):

```bash
nano .env
```

Configura SOLO estas 3 variables (el resto ya está pre-configurado):

```bash
# 1. API Key del LLM (OpenRouter, MiniMax, etc.)
OPENAI_API_KEY=sk-xxxxxx

# 2. Claves de acceso para ION (genera 2-3 claves)
LOCALAGI_API_KEYS=key1,key2,key3

# 3. Email de ION (para que técnicos escriban)
ION_EMAIL_USER=soporte@ionet.cl
ION_EMAIL_PASSWORD=tu-app-password
```

**¿Dónde obtenerlas?**

| Variable | ¿Dónde? |
|----------|---------|
| `OPENAI_API_KEY` | [OpenRouter](https://openrouter.ai) o [MiniMax](https://api.minimax.io) |
| `LOCALAGI_API_KEYS` | Genera strings aleatorios (ej: `abc123,xyz789`) |
| `ION_EMAIL_PASSWORD` | Office 365 > Security > App Password |

---

## PASO 3: Iniciar ION (1 minuto)

```bash
# Desplegar
./deploy.sh

# O con docker-compose directamente
docker compose -f docker-compose.prod.yaml up -d
```

---

## PASO 4: Verificar que Funciona

Accede a los siguientes canales:

### Email (principal)
```
Enviar email a: soporte@ionet.cl
Asunto: Prueba ION
Mensaje: Hola ION, ¿cuál es el protocolo para reiniciar un router?
```

### Web UI (opcional, para admin)
```
http://tu-servidor:8080
Usuario: ver en configuración
Contraseña: ver en configuración
```

### Teams (para notificaciones)
```
Configurar webhook en el canal de Teams:
1. Ir al canal
2. Connectors > Incoming Webhook > Configure
3. Copiar la URL generada
4. Pegar en .env como TEAMS_WEBHOOK_URL
```

---

## VERIFICACIÓN RÁPIDA

**¿ION responde por email?** ✅
- Envía un email a `soporte@ionet.cl`
- Deberías recibir respuesta en < 30 segundos

**¿ION envía notificaciones a Teams?** ✅
- Configura el webhook en Teams
- ION enviará alertas importantes al canal

**¿Los agentes especializados funcionan?** ✅
- Pregunta sobre: protocolos, redes, clientes, etc.
- ION derivará automáticamente al agente correcto

---

## AGENTES DISPONIBLES

| Agente | Para qué preguntar |
|--------|-------------------|
| **ION** | Punto de entrada, coordina todo |
| **agente-protocolos** | Procedimientos, políticas, manuales |
| **agente-redes** | Redes, servidores, IPs, VPN |
| **agente-clientes** | Clientes, contratos, usuarios nuevos |
| **agente-servicios** | Tickets, SLAs, solicitudes |
| **agente-inventario** | Equipos, licencias, software |
| **agente-seguridad** | Ciberseguridad, vulnerabilidades |
| **agente-datos** | Documentos, backups, versionado |

---

## COMUNICAR CON ION

### Desde Email (recomendado)
```
Para: soporte@ionet.cl
Asunto: [cualquier tema]
Mensaje: Tu consulta

Ejemplo:
"Asunto: Protocolo reset CPE
Mensaje: Cuál es el procedimiento exacto para reiniciar un CPE Huawei?"
```

### Desde Web UI (admin)
```
http://localhost:8080
Chatear directamente con ION
```

### Desde Teams (solo recibe notificaciones)
```
ION te enviará alertas automáticas al canal:
- Incidentes críticos
- Resúmenes de actividad
- Escalamientos a humanos
```

---

## EJEMPLOS DE CONSULTAS

**Protocolos:**
```
"Cuál es el procedimiento para reinstalar Office?"
"Necesito el protocolo de onboarding para nuevo cliente"
```

**Redes:**
```
"Cuál es la IP del servidor de respaldos?"
"Cómo configurar VPN para cliente X?"
```

**Clientes:**
```
"Qué contratos tiene el cliente Y?"
"Cómo agregar un nuevo usuario al cliente Z?"
```

**Seguridad:**
``"Tengo un alerta de ransomware, qué hago?"
"Cómo reportar un incidente de seguridad?"
```

---

## TROUBLESHOOTING

**ION no responde por email:**
```bash
# Ver logs
docker compose logs -f ionet

# Verificar configuración email
docker compose exec ionet env | grep ION_EMAIL
```

**ION no envía a Teams:**
```bash
# Verificar webhook URL configurada
cat .env | grep TEAMS_WEBHOOK

# Probar webhook manualmente
curl -X POST $TEAMS_WEBHOOK_URL \
  -H 'Content-Type: application/json' \
  -d '{"@type":"Message","text":"Test desde curl"}'
```

**Agents no funcionan:**
```bash
# Verificar que los agentes están importados
docker compose exec ionet ls -la /pool/agents/

# Reiniciar el servicio
docker compose restart ionet
```

---

## SOPORTE

- **Documentación completa**: [README.md](README.md)
- **Guía de agentes**: [config/agents/README.md](config/agents/README.md)
- **Problemas**: Abrir issue en GitHub

---

## PRÓXIMOS PASOS (OPCIONAL)

Una vez ION esté funcionando:

1. **Agregar más documentos a RAG**
   - Copiar archivos PDF/DOCX a `pool/rag/raw/`
   - O configurar M365 sync en `.env`

2. **Personalizar agentes**
   - Editar archivos en `config/agents/`
   - Ajustar prompts según necesidades

3. **Configurar más conectores**
   - Telegram, Slack, Discord
   - Ver documentación completa

---

## NOTAS DE SEGURIDAD

⚠️ **IMPORTANTE:**
- Nunca commitear el archivo `.env` con credenciales reales
- Usar contraseñas fuertes para `LOCALAGI_API_KEYS`
- En producción, usar Docker secrets o Kubernetes secrets
- Rotar contraseles periódicamente

---

## CHECKLIST DE DEPLOYMENT

Antes de considerar ION "listo para producción":

- [ ] `.env` configurado con todas las variables obligatorias
- [ ] Docker containers corriendo (`docker compose ps`)
- [ ] ION responde por email
- [ ] ION envía notificaciones a Teams
- [ ] Los 8 agentes están importados (verificar en UI)
- [ ] RAG funciona si hay documentos en `pool/rag/raw/`
- [ ] Logs sin errores (`docker compose logs ionet`)

---

**¿Listo?** 🚀

ION está listo para ayudar a los 10 técnicos con sus consultas diarias.

**Tiempo total estimado: 5-10 minutos**
