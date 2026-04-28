#!/usr/bin/env python3
"""
MCP server for IONET Workspace.
Provides file operations and persona management for ION agent.

Path Safety: ALL file ops confined to WORKSPACE_ROOT (/pool/ion-workspace).
Traversal attempts rejected before any I/O.

Tools:
  workspace_listar, workspace_leer, workspace_escribir, workspace_editar,
  workspace_crear_carpeta, workspace_eliminar
  persona_cargar, persona_anotar, persona_limpiar
"""

import os
import re
import logging
import time
from datetime import datetime

from fastmcp import FastMCP
from drive import DriveManager

# ═══════════════════════════════════════════════════════════════════
# CONSTANTES
# ═══════════════════════════════════════════════════════════════════

WORKSPACE_ROOT = "/pool/ion-workspace"
PERSONAS_DIR = "personas"
MAX_NOTES = 30
MAX_FILE_SIZE_READ = 1_048_576  # 1 MB
HOST = "0.0.0.0"
PORT = 8766

mcp = FastMCP("IONET Workspace")

logging.basicConfig(level=logging.INFO, format="%(name)s %(message)s")
logger = logging.getLogger("mcp-ionet-workspace")


# ═══════════════════════════════════════════════════════════════════
# HELPERS SEGUROS
# ═══════════════════════════════════════════════════════════════════


def _safe_path(relative_path: str) -> str:
    """Resuelve path relativo dentro de WORKSPACE_ROOT. Rechaza traversal."""
    cleaned = relative_path.lstrip("/")
    full = os.path.normpath(os.path.join(WORKSPACE_ROOT, cleaned))
    if not full.startswith(WORKSPACE_ROOT):
        raise ValueError(f"Path traversal bloqueado: {relative_path}")
    return full


def _now() -> str:
    return datetime.now().strftime("%Y-%m-%d %H:%M")


def _ensure_parent(path: str):
    os.makedirs(os.path.dirname(path), exist_ok=True)


# ═══════════════════════════════════════════════════════════════════
# WORKSPACE TOOLS
# ═══════════════════════════════════════════════════════════════════


@mcp.tool()
def workspace_listar(path: str = "") -> dict:
    """Lista archivos y carpetas dentro del workspace. Path relativo al workspace root."""
    target = _safe_path(path)
    if not os.path.exists(target):
        return {"error": f"Path no encontrado: {path}"}
    if not os.path.isdir(target):
        return {"error": f"No es un directorio: {path}"}

    items = []
    for entry in sorted(os.listdir(target)):
        entry_path = os.path.join(target, entry)
        st = os.stat(entry_path)
        items.append({
            "nombre": entry,
            "tipo": "carpeta" if os.path.isdir(entry_path) else "archivo",
            "tamano": st.st_size,
            "modificado": datetime.fromtimestamp(st.st_mtime).isoformat(),
        })
    return {"items": items, "total": len(items)}


@mcp.tool()
def workspace_leer(path: str) -> dict:
    """Lee contenido de un archivo del workspace. Path relativo."""
    target = _safe_path(path)
    if not os.path.exists(target):
        return {"error": f"Archivo no encontrado: {path}"}
    if not os.path.isfile(target):
        return {"error": f"No es un archivo: {path}"}

    st = os.stat(target)
    if st.st_size > MAX_FILE_SIZE_READ:
        return {"error": f"Archivo demasiado grande ({st.st_size}b). Máx: {MAX_FILE_SIZE_READ}b"}

    with open(target, "r", encoding="utf-8") as f:
        content = f.read()

    return {
        "path": path,
        "contenido": content,
        "tamano": st.st_size,
        "modificado": datetime.fromtimestamp(st.st_mtime).isoformat(),
    }


@mcp.tool()
def workspace_escribir(path: str, contenido: str) -> dict:
    """Crea o sobrescribe un archivo. Crea carpetas intermedias si no existen."""
    target = _safe_path(path)
    _ensure_parent(target)
    with open(target, "w", encoding="utf-8") as f:
        f.write(contenido)
    logger.info(f"Escrito {path} ({len(contenido)}b)")
    return {"mensaje": f"Archivo guardado: {path}", "tamano": len(contenido)}


@mcp.tool()
def workspace_editar(path: str, buscar: str, reemplazar: str) -> dict:
    """Reemplaza texto en un archivo existente (exact match). Retorna conteo de reemplazos."""
    target = _safe_path(path)
    if not os.path.exists(target):
        return {"error": f"Archivo no encontrado: {path}"}

    with open(target, "r", encoding="utf-8") as f:
        content = f.read()

    count = content.count(buscar)
    if count == 0:
        return {"error": f"Texto no encontrado en {path}", "reemplazos": 0}

    content = content.replace(buscar, reemplazar)
    with open(target, "w", encoding="utf-8") as f:
        f.write(content)

    logger.info(f"Editado {path}: {count} reemplazo(s)")
    return {"mensaje": f"Reemplazos: {count}", "reemplazos": count}


@mcp.tool()
def workspace_crear_carpeta(path: str) -> dict:
    """Crea carpeta(s) dentro del workspace."""
    target = _safe_path(path)
    os.makedirs(target, exist_ok=True)
    logger.info(f"Carpeta creada: {path}")
    return {"mensaje": f"Carpeta creada: {path}"}


@mcp.tool()
def workspace_eliminar(path: str) -> dict:
    """Elimina archivo o carpeta vacía del workspace."""
    target = _safe_path(path)
    if target == WORKSPACE_ROOT:
        return {"error": "No se puede eliminar la raíz del workspace"}
    if not os.path.exists(target):
        return {"error": f"No encontrado: {path}"}

    if os.path.isdir(target):
        os.rmdir(target)
    else:
        os.remove(target)

    logger.info(f"Eliminado: {path}")
    return {"mensaje": f"Eliminado: {path}"}


# ═══════════════════════════════════════════════════════════════════
# PERSONA HELPERS
# ═══════════════════════════════════════════════════════════════════


def _persona_file(email: str) -> str:
    safe = email.replace("/", "_")
    return _safe_path(f"{PERSONAS_DIR}/{safe}.md")


def _parse_notes(content: str) -> list[str]:
    match = re.search(r"## Notas de ION\n(.*?)(?=\n## |\Z)", content, re.DOTALL)
    if not match:
        return []
    notes = []
    for line in match.group(1).strip().split("\n"):
        line = line.strip()
        if line.startswith("- "):
            notes.append(line[2:])
    return notes


def _replace_notes(content: str, notes: list[str]) -> str:
    section = "## Notas de ION\n" + "\n".join(f"- {n}" for n in notes) + "\n"
    if "## Notas de ION" in content:
        return re.sub(
            r"## Notas de ION\n.*?(?=\n## |\Z)",
            section,
            content,
            count=1,
            flags=re.DOTALL,
        )
    return content.rstrip() + "\n\n" + section


def _default_persona(email: str) -> str:
    return (
        f"# Persona\n\n"
        f"## Datos base\n"
        f"- email: {email}\n"
        f"\n"
        f"## Notas de ION\n"
        f"- {_now()}: Perfil creado\n"
    )


# ═══════════════════════════════════════════════════════════════════
# PERSONA TOOLS
# ═══════════════════════════════════════════════════════════════════


@mcp.tool()
def persona_cargar(email: str) -> dict:
    """Carga perfil de colaborador por email. Retorna contenido y metadatos."""
    target = _persona_file(email)
    if not os.path.exists(target):
        return {"error": f"Persona no encontrada: {email}", "email": email}

    with open(target, "r", encoding="utf-8") as f:
        content = f.read()

    st = os.stat(target)
    notes = _parse_notes(content)
    return {
        "email": email,
        "contenido": content,
        "notas": notes,
        "total_notas": len(notes),
        "modificado": datetime.fromtimestamp(st.st_mtime).isoformat(),
    }


@mcp.tool()
def persona_anotar(email: str, nota: str) -> dict:
    """Agrega nota al perfil del colaborador. Si excede MAX_NOTES, trunca la más antigua."""
    target = _persona_file(email)
    _ensure_parent(target)

    if os.path.exists(target):
        with open(target, "r", encoding="utf-8") as f:
            content = f.read()
    else:
        content = _default_persona(email)

    notes = _parse_notes(content)
    notes.append(f"{_now()}: {nota}")

    dropped = 0
    if len(notes) > MAX_NOTES:
        dropped = len(notes) - MAX_NOTES
        notes = notes[-MAX_NOTES:]

    content = _replace_notes(content, notes)
    with open(target, "w", encoding="utf-8") as f:
        f.write(content)

    logger.info(f"Persona {email}: {len(notes)} notas ({dropped} truncadas)")
    return {
        "mensaje": "Nota agregada" if not dropped else f"Nota agregada, {dropped} antigua(s) truncada(s)",
        "email": email,
        "total_notas": len(notes),
        "truncadas": dropped,
    }


@mcp.tool()
def persona_limpiar(email: str, resumen: str) -> dict:
    """Reemplaza TODAS las notas de ION por un resumen. Útil para compactar."""
    target = _persona_file(email)
    if not os.path.exists(target):
        return {"error": f"Persona no encontrada: {email}"}

    with open(target, "r", encoding="utf-8") as f:
        content = f.read()

    content = _replace_notes(content, [f"{_now()}: {resumen}"])
    with open(target, "w", encoding="utf-8") as f:
        f.write(content)

    logger.info(f"Persona limpiada: {email}")
    return {"mensaje": "Notas reemplazadas por resumen", "email": email}


# ═══════════════════════════════════════════════════════════════════
# DRIVE TOOLS
# ═══════════════════════════════════════════════════════════════════


def _drive() -> DriveManager:
    d = DriveManager()
    if not d.ready():
        raise ValueError("Google Drive no disponible. Verifica la clave de service account.")
    return d


@mcp.tool()
def drive_estado() -> dict:
    """Muestra estado de Google Drive (solo lectura por service account)."""
    d = DriveManager()
    s = d.get_auth_status()
    s["modo"] = "solo_lectura"
    s["nota"] = "Service account: solo lectura. Para escribir usa workspace_escribir."
    return s


@mcp.tool()
def drive_listar(carpeta_id: str = "root") -> list:
    """Lista archivos en una carpeta de Drive. Solo lectura."""
    return _drive().list_folder(carpeta_id)


@mcp.tool()
def drive_buscar(termino: str) -> list:
    """Busca archivos en Drive por nombre. Solo lectura."""
    return _drive().search(termino)


@mcp.tool()
def drive_leer(file_id: str) -> dict:
    """Lee contenido de un archivo de Drive (texto). Solo lectura."""
    return _drive().read_text(file_id)


@mcp.tool()
def drive_subir(nombre: str, contenido: str, carpeta_id: str = None) -> dict:
    """[SOLO LECTURA] La service account no puede escribir en Drive. Usa workspace_escribir."""
    return {"error": "Drive en modo solo lectura. Usa workspace_escribir para guardar archivos."}


@mcp.tool()
def drive_crear_carpeta(nombre: str, carpeta_id: str = None) -> dict:
    """[SOLO LECTURA] La service account no puede crear carpetas en Drive."""
    return {"error": "Drive en modo solo lectura. Usa workspace_crear_carpeta en su lugar."}


# ═══════════════════════════════════════════════════════════════════
# INIT
# ═══════════════════════════════════════════════════════════════════


@mcp.tool()
def ping() -> dict:
    """Verifica conectividad del workspace."""
    return {
        "status": "pong",
        "workspace_root": WORKSPACE_ROOT,
        "workspace_existe": os.path.isdir(WORKSPACE_ROOT),
        "time": time.time(),
    }


if __name__ == "__main__":
    # Crear directorios base si no existen
    os.makedirs(WORKSPACE_ROOT, exist_ok=True)
    os.makedirs(os.path.join(WORKSPACE_ROOT, PERSONAS_DIR), exist_ok=True)

    logger.info(f"Iniciando en puerto {PORT}, workspace: {WORKSPACE_ROOT}")
    mcp.run(transport="http", host=HOST, port=PORT)
