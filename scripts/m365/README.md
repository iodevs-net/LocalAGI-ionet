# IONET - Sincronización M365

Script para sincronizar documentos desde SharePoint/OneDrive de Microsoft 365 hacia el servidor IONET para indexación RAG.

## Requisitos Previos

### 1. Registrar App en Azure AD

1. Ir a [Azure Portal](https://portal.azure.com) > **Azure Active Directory** > **App registrations**
2. Click en **New registration**
3. Configurar:
   - **Name**: `IONET RAG Sync`
   - **Supported account types**: Accounts in this organizational directory only
   - **Redirect URI**: Web > `http://localhost`
4. Click **Register**
5. Copiar el **Application (client) ID** → `M365_CLIENT_ID`
6. Copiar el **Directory (tenant) ID** → `M365_TENANT_ID`

### 2. Generar Client Secret

1. En la app registrada, ir a **Certificates & secrets**
2. Click **New client secret**
3. Description: `Production`
4. Click **Add**
5. **IMPORTANTE**: Copiar el valor del secret inmediatamente (solo se muestra una vez)
6. Usar este valor como `M365_CLIENT_SECRET`

### 3. Configurar Permisos de API

1. Ir a **API permissions**
2. Click **Add a permission**
3. Seleccionar **Microsoft Graph**
4. Seleccionar **Application permissions**
5. Agregar:
   - `Sites.Read.All` (para SharePoint)
   - `Files.Read.All` (para OneDrive)
6. Click **Grant admin consent for IONET**

### 4. Compartir Documentos con la App

Para que la app pueda acceder a los documentos:

1. En SharePoint/OneDrive, ir al documento o biblioteca
2. Click en **Share** > **Manage access**
3. Agregar la app por nombre o email (ej: `ionet-rag-sync@ionet.onmicrosoft.com`)

## Configuración

### Variables de Entorno

```bash
# Obligatorias
M365_CLIENT_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
M365_CLIENT_SECRET=tu-secret-aqui
M365_TENANT_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx

# Opcionales (una de ellas)
M365_SITE_URL=https://ionet.sharepoint.com/sites/documentos  # SharePoint site
M365_DRIVE_ID=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx              # Drive ID específico

# Salida (opcionales)
RAG_DATA_DIR=/rag/raw  # Directorio donde se guardan los archivos
```

## Uso

### Desarrollo

```bash
# Desde el directorio del proyecto
cd scripts/m365

# Crear archivo .env con las credenciales
cp ../../.env .env

# Ejecutar
go run sync.go
```

### Producción

```bash
# Compilar
cd scripts/m365
go build -o sync-m365 sync.go

# Ejecutar manualmente
./sync-m365

# O con variables de entorno
M365_CLIENT_ID=xxx M365_CLIENT_SECRET=xxx M365_TENANT_ID=xxx ./sync-m365
```

### Cron (sincronización automática)

```bash
# Editar crontab
crontab -e

# Sincronizar cada 6 horas
0 */6 * * * /app/scripts/m365/sync-m365 >> /var/log/ionet-m365-sync.log 2>&1

# O diariamente a las 2 AM
0 2 * * * /app/scripts/m365/sync-m365 >> /var/log/ionet-m365-sync.log 2>&1
```

## Formatos Soportados

| Extensión | Tipo |
|-----------|------|
| `.pdf` | PDF |
| `.docx`, `.doc` | Word |
| `.xlsx`, `.xls` | Excel |
| `.pptx`, `.ppt` | PowerPoint |
| `.txt` | Texto plano |
| `.md` | Markdown |
| `.html`, `.htm` | HTML |

## Estructura de Salida

```
/rag/raw/
├── documentos/
│   ├── protocolo-instalacion.pdf
│   ├── guia-configuracion.docx
│   └── ...
├── procedimientos/
│   ├── ...
└── ...
```

## Troubleshooting

### Error: "Authentication failed"

- Verificar que `M365_CLIENT_ID`, `M365_CLIENT_SECRET` y `M365_TENANT_ID` son correctos
- Verificar que el secret no ha expirado
- Verificar que se dio **admin consent** a los permisos

### Error: "Access denied"

- La app no tiene permisos sobre los documentos
- Compartir explícitamente los documentos/bibliotecas con la app
- Verificar que el site/drive existe y la URL es correcta

### Error: "No such host"

- Verificar conexión a internet
- Verificar que `graph.microsoft.com` no está bloqueado por firewall

### Archivos no aparecen

- Verificar que las extensiones están soportadas
- Revisar logs en `/var/log/ionet-m365-sync.log`
- Verificar que el directorio de salida tiene permisos de escritura

## Docker Integration

Para correr en el contenedor de producción:

```bash
# Agregar al docker-compose.prod.yaml en la sección ionet:
#
# volumes:
#   - ./scripts/m365/sync.sh:/usr/local/bin/sync-m365:ro
# command: >
#   sh -c "sync-m365 && /usr/local/bin/localagi serve"
```

O crear un servicio systemd separado para el sync.

## Seguridad

- **NUNCA** commitear el archivo `.env` con credenciales reales
- Usar secretos de Docker/Kubernetes en producción
- Rotar el client secret periódicamente
- La app solo necesita permisos de lectura (`Read`)
