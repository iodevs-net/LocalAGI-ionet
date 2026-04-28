#!/usr/bin/env python3
"""
Ultra Debug: Email Flow Analyzer for IONET/LocalAGI
====================================================
Simula el flujo completo de procesamiento de correo:
1. Conexión IMAP y fetch de email
2. Parsing MIME (multipart, text/plain, text/html)
3. Extracción de contenido (con detección de imágenes embebidas)
4. Medición de tokens y caracteres basura (base64, imágenes)
5. Construcción del prompt final para el LLM
6. (Opcional) Envío al LLM para verificar respuesta

Uso:
  python3 debug_email_flow.py                    # Último email no procesado
  python3 debug_email_flow.py --send-test        # Envía test + analiza
  python3 debug_email_flow.py --msg-id 42        # Email específico por seq
  python3 debug_email_flow.py --llm-test         # Incluye llamada real al LLM
"""

import imaplib
import email
import email.policy
import json
import os
import re
import sys
import base64
import quopri
import html
from datetime import datetime
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from email.mime.base import MIMEBase
from email.utils import parsedate_to_datetime

# ─── Config ───────────────────────────────────────────────────────────────────
CONFIG = {
    "imap_server": "imap.gmail.com",
    "imap_port": 993,
    "username": os.environ.get("ION_EMAIL_USER", "el.agente.ion@gmail.com"),
    "password": os.environ.get("ION_EMAIL_PASSWORD", "lydc yyja fnfr pooh"),
    "monitored_email": "ion@iodevs.net",
    "openai_api_key": os.environ.get("OPENAI_API_KEY", "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"),
    "openai_base_url": "https://openrouter.ai/api/v1",
    "model": os.environ.get("MODEL_NAME", "nvidia/nemotron-3-super-120b-a12b:free"),
}


# ─── ANSI Colors ──────────────────────────────────────────────────────────────
class C:
    BOLD = "\033[1m"
    DIM = "\033[2m"
    GREEN = "\033[92m"
    YELLOW = "\033[93m"
    RED = "\033[91m"
    CYAN = "\033[96m"
    MAGENTA = "\033[95m"
    BLUE = "\033[94m"
    RESET = "\033[0m"
    GRAY = "\033[90m"
    BGGRAY = "\033[100m"
    WHITE = "\033[97m"


def section(title):
    print(f"\n{C.BOLD}{C.CYAN}{'='*60}")
    print(f"  {title}")
    print(f"{'='*60}{C.RESET}\n")


def step(num, title, detail=""):
    print(f"  {C.BOLD}{C.GREEN}[PASO {num}]{C.RESET} {C.WHITE}{title}{C.RESET}")
    if detail:
        for line in detail.strip().split("\n"):
            print(f"    {C.GRAY}{line}{C.RESET}")
    print()


def warn(msg):
    print(f"  {C.BOLD}{C.YELLOW}⚠  {msg}{C.RESET}")


def ok(msg):
    print(f"  {C.BOLD}{C.GREEN}✓ {msg}{C.RESET}")


def fail(msg):
    print(f"  {C.BOLD}{C.RED}✗ {msg}{C.RESET}")


def data(label, value, max_chars=200):
    s = str(value)
    if len(s) > max_chars:
        s = s[:max_chars] + f"... ({len(value)} chars total)"
    print(f"    {C.GRAY}{label}:{C.RESET} {s}")
    print()


# ─── IMAP Fetch ───────────────────────────────────────────────────────────────

def fetch_email(msg_id=None):
    """Conecta vía IMAP, obtiene el último o un email específico."""
    section(f"IMAP: Conectando a {CONFIG['imap_server']}:{CONFIG['imap_port']}")

    conn = imaplib.IMAP4_SSL(CONFIG["imap_server"], CONFIG["imap_port"])
    conn.login(CONFIG["username"], CONFIG["password"])
    ok(f"Login OK: {CONFIG['username']}")

    conn.select("INBOX")
    result, data = conn.search(None, "ALL")
    ids = data[0].split()
    total = len(ids)
    ok(f"INBOX tiene {total} mensajes")

    if msg_id is None:
        msg_id = ids[-1]  # Último
    elif isinstance(msg_id, int):
        # Buscar por secuencia numérica
        msg_id = ids[msg_id - 1] if msg_id <= len(ids) else ids[-1]

    result, data = conn.fetch(msg_id, "(RFC822)")
    raw_email = data[0][1]
    conn.logout()

    ok(f"Fetch OK: mensaje #{msg_id.decode() if hasattr(msg_id, 'decode') else msg_id}")
    return email.message_from_bytes(raw_email, policy=email.policy.default)


# ─── MIME Analysis ────────────────────────────────────────────────────────────

def walk_mime(msg, depth=0, findings=None):
    """Recorre la estructura MIME y recolecta información."""
    if findings is None:
        findings = []

    ct = msg.get_content_type()
    cte = msg.get("Content-Transfer-Encoding", "")
    cd = msg.get("Content-Disposition", "")
    filename = msg.get_filename()

    # Detectar tipo de contenido embebido
    is_inline_image = (
        ct.startswith("image/")
        and "inline" in cd
    )
    is_attachment = (
        filename is not None
        and "attachment" in cd
    )
    is_base64 = "base64" in cte

    info = {
        "depth": depth,
        "content_type": ct,
        "encoding": cte,
        "disposition": cd,
        "filename": filename,
        "is_inline_image": is_inline_image,
        "is_attachment": is_attachment,
        "is_base64": is_base64,
        "size": len(str(msg)),
        "charset": msg.get_content_charset(),
        "boundary": msg.get_boundary(),
    }

    # Estimar payload size después de decode
    try:
        payload = msg.get_payload(decode=True)
        if payload:
            info["decoded_size"] = len(payload)
    except Exception:
        pass

    # Detectar contenido base64 inline en text/html (imágenes embebidas)
    if ct == "text/html":
        try:
            html_content = str(msg.get_payload(decode=True) or b"", "utf-8", errors="replace")
            # Buscar imágenes base64 embebidas
            b64_images = re.findall(
                r'<img[^>]+src=["\']data:image/[^;]+;base64,([^"\']+)["\']',
                html_content, re.IGNORECASE
            )
            info["embedded_b64_images"] = len(b64_images)
            info["embedded_b64_size"] = sum(len(b) for b in b64_images)
        except Exception:
            pass

    findings.append(info)

    if msg.is_multipart():
        for part in msg.walk():
            if part is msg:
                continue
            walk_mime(part, depth + 1, findings)

    return findings


def print_mime_tree(findings):
    """Print jerarquía MIME con indicadores visuales."""
    print(f"  {'  Depth':<8} {'Content-Type':<40} {'Size':<10} {'Flags':<20}")
    print(f"  {'─'*8} {'─'*40} {'─'*10} {'─'*20}")

    for f in findings:
        indent = "  " * f["depth"]
        icon = "📎" if f["is_attachment"] else "🖼" if f["is_inline_image"] else "📄" if f["content_type"] == "text/plain" else "🌐" if f["content_type"] == "text/html" else "  "
        flags = []
        if f["is_base64"]:
            flags.append(f"{C.YELLOW}B64{C.RESET}")
        if f["is_attachment"]:
            flags.append(f"{C.MAGENTA}ATTACH{C.RESET}")
        if f["is_inline_image"]:
            flags.append(f"{C.RED}IMG{C.RESET}")
        if f.get("embedded_b64_images", 0) > 0:
            flags.append(f"{C.RED}{f['embedded_b64_images']}EMBED-IMG({f['embedded_b64_size']}b){C.RESET}")

        size_str = f"{f['size'] // 1024}KB" if f['size'] > 1024 else f"{f['size']}B"
        decoded_str = ""
        if f.get("decoded_size"):
            decoded_str = f" → {f['decoded_size']//1024}KB decoded" if f['decoded_size'] > 1024 else f" → {f['decoded_size']}B dec"

        print(f"  {indent}{icon} {f['content_type'][:38]:38} {size_str:8}{decoded_str:20}")
        for flag in flags:
            print(f"  {'':8} {'':40} {'':10} {flag:20}")

        if f.get("filename"):
            print(f"  {'':8} {'Filename: ' + f['filename']:50}")


# ─── Content Extraction (como en email.go) ────────────────────────────────────

PREFIXES_HTML = ["<html", "<body", "<div", "<head"]

def extract_text_content(msg):
    """Replica la lógica de extractTextContent del email.go de LocalAGI."""
    parts = []

    if msg.is_multipart():
        for part in msg.walk():
            if part is msg:
                continue
            ct = part.get_content_type()
            if ct == "text/plain":
                payload = part.get_payload(decode=True)
                if payload:
                    try:
                        parts.append(("text/plain", payload.decode("utf-8", errors="replace")))
                    except Exception:
                        parts.append(("text/plain", str(payload)))
            elif ct == "text/html":
                payload = part.get_payload(decode=True)
                if payload:
                    try:
                        parts.append(("text/html", payload.decode("utf-8", errors="replace")))
                    except Exception:
                        parts.append(("text/html", str(payload)))
            elif ct.startswith("image/") or ct.startswith("application/"):
                # Skip images and attachments
                payload = part.get_payload(decode=True)
                if payload:
                    parts.append((ct, f"<ATTACHMENT/BINARY: {ct} {len(payload)} bytes>"))
    else:
        ct = msg.get_content_type()
        payload = msg.get_payload(decode=True)
        if payload:
            try:
                parts.append((ct, payload.decode("utf-8", errors="replace")))
            except Exception:
                parts.append((ct, str(payload)))

    # Primero buscar text/plain
    for ct, content in parts:
        if ct == "text/plain" and content.strip():
            return content, "text/plain", parts

    # Si no hay plain, buscar HTML
    for ct, content in parts:
        if ct == "text/html":
            return content, "text/html", parts

    # Fallback: primer contenido textual
    for ct, content in parts:
        if not ct.startswith("image/") and not ct.startswith("application/"):
            return content, ct, parts

    return "", "unknown", parts


def has_html_prefix(content):
    """Detecta si el contenido parece HTML."""
    stripped = content.strip().lower()
    for prefix in PREFIXES_HTML:
        if stripped.startswith(prefix.lower()):
            return True
    return False


def strip_base64_images(html_content):
    """Elimina imágenes base64 embebidas del HTML."""
    count_before = len(re.findall(r'data:image/[^;]+;base64,[^"\']+', html_content))
    cleaned = re.sub(
        r'<img[^>]+src=["\']data:image/[^;]+;base64,[^"\']+["\'][^>]*>',
        '<i>[IMAGEN BASE64 ELIMINADA]</i>',
        html_content,
        flags=re.IGNORECASE
    )
    # También limpiar URLs base64 directas
    cleaned = re.sub(
        r'data:image/[^;]+;base64,[a-zA-Z0-9+/=]+',
        '[BASE64_IMAGE_DATA_STRIPPED]',
        cleaned
    )
    count_after = len(re.findall(r'data:image/[^;]+;base64,[^"\']+', cleaned))
    return cleaned, count_before - count_after


def strip_html_tags(html_content):
    """Strip HTML tags (alternativa al markdown)."""
    text = re.sub(r'<style[^>]*>.*?</style>', '', html_content, flags=re.DOTALL)
    text = re.sub(r'<script[^>]*>.*?</script>', '', text, flags=re.DOTALL)
    text = re.sub(r'<[^>]+>', ' ', text)
    text = re.sub(r'\s+', ' ', text)
    return text.strip()


def size_in_tokens(text):
    """Estimación burda: ~4 chars por token."""
    return len(text) // 4


# ─── Allowed Sender Check ─────────────────────────────────────────────────────

ALLOWED_SENDERS = [
    "@ionet.cl",
    "@iodevs.net",
    "el.agente.ion@gmail.com",
    "ventas.ionet@gmail.com",
]

def check_allowed_sender(from_header):
    """Replica el filtro de allowedSenders del email.go."""
    for allowed in ALLOWED_SENDERS:
        if allowed in from_header:
            return True, allowed
    return False, None


def check_allowed_recipient(to_header, delivered_to):
    """Replica el filtro allowedTo del email.go."""
    target = CONFIG["monitored_email"]
    username = CONFIG["username"]
    return (
        target in to_header
        or username in to_header
        or target in delivered_to
        or username in delivered_to
    )


# ─── Main Analysis ────────────────────────────────────────────────────────────

def analyze_email(msg, run_llm=False):
    """Analiza un email completo y muestra todos los pasos."""
    section("📨 ANÁLISIS DE EMAIL")

    # ── Headers ──
    step(1, "HEADERS DEL EMAIL")
    for h in ["From", "To", "Cc", "Subject", "Date", "Message-ID",
              "Delivered-To", "Return-Path", "Content-Type", "MIME-Version"]:
        v = msg.get(h, "")
        if v:
            data(h, v)

    # ── Allowed checks ──
    step(2, "FILTROS DE SEGURIDAD (replicando email.go)")

    from_h = msg.get("From", "")
    to_h = msg.get("To", "")
    dt_h = msg.get("Delivered-To", "")
    is_allowed_sender, match = check_allowed_sender(from_h)
    if is_allowed_sender:
        ok(f"Sender autorizado: {from_h} (match: {match})")
    else:
        fail(f"Sender NO autorizado: {from_h}")
        warn("Email sería IGNORADO por filtro de remitentes")

    is_allowed_rcpt = check_allowed_recipient(to_h, dt_h)
    if is_allowed_rcpt:
        ok(f"Destinatario válido: To={to_h[:60]} / Delivered-To={dt_h[:60]}")
    else:
        fail(f"Destinatario NO válido para monitored_email={CONFIG['monitored_email']}")
        warn("Email sería IGNORADO por filtro de destinatarios")

    if not (is_allowed_sender and is_allowed_rcpt):
        fail("❌ FILTROS BLOQUEAN ESTE EMAIL - no se procesaría")
        return

    ok("✅ Filtros OK — email sería procesado")

    # ── MIME Structure ──
    step(3, "ESTRUCTURA MIME")
    findings = walk_mime(msg)
    print_mime_tree(findings)

    # ── Content extraction ──
    step(4, "EXTRACCIÓN DE CONTENIDO (extractTextContent)")

    content, content_type, all_parts = extract_text_content(msg)
    data("Tipo de contenido extraído", content_type)
    data("Longitud original (chars)", f"{len(content):,}")
    data("Tokens estimados original", f"{size_in_tokens(content):,}")

    # Detectar basura
    garbage_stats = {
        "base64_images": 0,
        "base64_chars": 0,
    }

    # Imágenes embebidas
    b64_images = re.findall(r'data:image/[^;]+;base64,([^"\']+)', content)
    if b64_images:
        total_b64 = sum(len(b) for b in b64_images)
        garbage_stats["base64_images"] = len(b64_images)
        garbage_stats["base64_chars"] = total_b64
        warn(f"Encontradas {len(b64_images)} imágenes base64 embebidas ({total_b64:,} chars)")
        warn(f"Tokens desperdiciados en imágenes: ~{total_b64 // 4:,}")
    else:
        ok("Sin imágenes base64 embebidas")

    # Mostrar partes individuales
    print()
    for i, (ct, c) in enumerate(all_parts):
        is_garbage = ct.startswith("image/") or ct.startswith("application/")
        icon = "🗑️" if is_garbage else "✅"
        label = f"   Parte {i+1}: {ct}"
        if is_garbage:
            warn(f"{icon} {label} ({len(c):,} chars, {c[:80]}...)")
        else:
            ok(f"{icon} {label} ({len(c):,} chars)")
            print(f"      {C.GRAY}Primeros 200 chars:{C.RESET}")
            print(f"      {c[:200]}")
            print()

    # ── Detección de HTML ──
    step(5, "DETECCIÓN HTML Y DEPURACIÓN")

    is_html = has_html_prefix(content)
    if is_html:
        warn(f"Contenido detectado como HTML (empieza con: {content.strip()[:30]}...)")
        # Stripear imágenes base64 embebidas
        cleaned, stripped_count = strip_base64_images(content)
        if stripped_count > 0:
            ok(f"Imágenes base64 eliminadas: {stripped_count}")
        else:
            ok("Sin imágenes base64 para eliminar")

        # Mostrar diff
        print(f"\n   {C.YELLOW}HTML sin tags:{C.RESET}")
        text_version = strip_html_tags(cleaned)
        print(f"   {text_version[:300]}")
        print(f"\n   {C.GRAY}HTML original ({len(content):,} chars) → Texto ({len(text_version):,} chars){C.RESET}")
        savings = size_in_tokens(content) - size_in_tokens(text_version)
        print(f"   {C.GREEN}Ahorro estimado: ~{savings:,} tokens{C.RESET}")
    else:
        ok("Contenido no HTML — texto plano")
        text_version = content
        cleaned = content

    # ── Construcción del prompt ──
    step(6, "CONSTRUCCIÓN DEL PROMPT (como en email.go)")

    date_str = msg.get("Date", datetime.now().isoformat())
    subject = msg.get("Subject", "(sin asunto)")

    # Replicar la construcción exacta de email.go:
    prompt = (
        f"This email thread was sent to you. You are {CONFIG['monitored_email']}:\n\n"
        f"From: {msg.get('From')}\n"
        f"Time: {date_str}\n"
        f"Subject: {subject}\n"
        f"=====\n"
        f"{text_version}"
    )

    data("Prompt length (chars)", f"{len(prompt):,}")
    data("Prompt estimated tokens", f"{size_in_tokens(prompt):,}")
    print(f"   {C.GRAY}Prompt preview (primeros 500 chars):{C.RESET}")
    print(f"   {C.DIM}{prompt[:500]}{C.RESET}")
    print()

    # ── Enviar al LLM (opcional) ──
    if run_llm:
        step(7, "LLM TEST: ENVIANDO A OPENROUTER")
        import http.client
        import json as j

        conn = http.client.HTTPSConnection("openrouter.ai")
        payload = j.dumps({
            "model": CONFIG["model"],
            "messages": [
                {
                    "role": "system",
                    "content": "Eres ION, el asistente de soporte TI de IONET. Responde de forma concisa y profesional."
                },
                {
                    "role": "user",
                    "content": prompt
                }
            ],
            "max_tokens": 500,
        })

        headers = {
            "Content-Type": "application/json",
            "Authorization": f"Bearer {CONFIG['openai_api_key']}",
            "HTTP-Referer": "https://iodevs.net",
            "X-Title": "LocalAGI",
        }

        print(f"   {C.GRAY}Enviando a OpenRouter (model={CONFIG['model']})...{C.RESET}")
        start = datetime.now()
        conn.request("POST", "/api/v1/chat/completions", payload, headers)
        resp = conn.getresponse()
        body = resp.read().decode()
        elapsed = (datetime.now() - start).total_seconds()

        if resp.status == 200:
            result = j.loads(body)
            choice = result["choices"][0]
            reply = choice["message"]["content"]
            finish = choice.get("finish_reason", "unknown")
            model_used = result.get("model", CONFIG["model"])

            ok(f"LLM respondió en {elapsed:.1f}s (model: {model_used}, finish: {finish})")
            print(f"\n   {C.GREEN}RESPUESTA:{C.RESET}")
            print(f"   {reply}")
            print()
            if choice.get("message", {}).get("reasoning"):
                reasoning = choice["message"]["reasoning"]
                print(f"\n   {C.GRAY}REASONING:{C.RESET}")
                print(f"   {reasoning}")
        else:
            fail(f"Error LLM ({resp.status}): {body[:500]}")
    else:
        step(7, "LLM TEST", "Saltado (usa --llm-test para activar)")
        print()

    # ── Resumen ──
    section("📊 RESUMEN")

    # Contar tokens basura
    total_garbage = 0
    for ct, c in all_parts:
        if ct.startswith("image/") or ct.startswith("application/"):
            total_garbage += len(c)

    print(f"  {C.BOLD}Tamaño total email:{C.RESET} {len(str(msg)):,} bytes")
    print(f"  {C.BOLD}Contenido texto útil:{C.RESET} {len(text_version):,} chars ~ {size_in_tokens(text_version):,} tokens")
    print(f"  {C.BOLD}Basura detectada:{C.RESET} {total_garbage:,} chars + {garbage_stats['base64_chars']:,} chars de imágenes embebidas")
    print(f"  {C.BOLD}Prompt final:{C.RESET} {len(prompt):,} chars ~ {size_in_tokens(prompt):,} tokens")

    if total_garbage > 0 or garbage_stats['base64_images'] > 0:
        print(f"\n  {C.RED}{C.BOLD}⚠  BASURA DETECTADA — OPTIMIZAR EXTRACCIÓN{MAGENTA}")
        print(f"  Recomendación:")
        print(f"  - Mejorar extractTextContent() para saltar attachments")
        print(f"  - Stripear imágenes base64 del HTML antes de convertir a markdown")
        print(f"  - Limitar tamaño máximo de contenido (ej: 100KB){C.RESET}")
    else:
        print(f"\n  {C.GREEN}{C.BOLD}✅ Contenido limpio — sin basura detectable{C.RESET}")

    print()
    return prompt


# ─── Send Test Email ──────────────────────────────────────────────────────────

def send_test_email():
    """Envía un email de test a sí mismo para probar el flujo completo."""
    section("📤 ENVIANDO EMAIL DE TEST")

    import smtplib

    msg = MIMEMultipart("alternative")
    msg["Subject"] = f"ULTRA DEBUG TEST - {datetime.now().strftime('%H:%M:%S')}"
    msg["From"] = CONFIG["username"]
    msg["To"] = CONFIG["username"]  # Self-send para evitar Cloudflare loop

    text_part = MIMEText(
        "Hola ION, este es un test automatizado del sistema de debug.\n"
        "Por favor confirma que recibes este mensaje correctamente.\n\n"
        "---\nMensaje generado por debug_email_flow.py"
    )
    msg.attach(text_part)

    html_part = MIMEText(
        "<html><body>"
        "<h2>Test ION Debug</h2>"
        "<p>Hola ION, este es un <b>test automatizado</b> del sistema de debug.</p>"
        "<p>Por favor confirma recepción.</p>"
        "<hr>"
        "<p><i>Mensaje generado por debug_email_flow.py</i></p>"
        "</body></html>",
        "html",
    )
    msg.attach(html_part)

    with smtplib.SMTP("smtp.gmail.com", 587) as s:
        s.starttls()
        s.login(CONFIG["username"], CONFIG["password"])
        s.send_message(msg)

    ok(f"Email enviado a {CONFIG['username']}")
    print()
    return True


# ─── CLI ──────────────────────────────────────────────────────────────────────

if __name__ == "__main__":
    send_test = "--send-test" in sys.argv
    llm_test = "--llm-test" in sys.argv
    msg_id = None

    for arg in sys.argv:
        if arg.startswith("--msg-id="):
            msg_id = int(arg.split("=")[1])

    print(f"{C.BOLD}{C.MAGENTA}")
    print(rf"""
  ╔══════════════════════════════════════════════════╗
  ║    IONET Email Flow Ultra Debugger v2            ║
  ║    Analiza cada paso del procesamiento de email  ║
  ╚══════════════════════════════════════════════════╝
    """)
    print(f"{C.RESET}")

    print(f"  Config:")
    print(f"    IMAP:     {CONFIG['imap_server']}:{CONFIG['imap_port']}")
    print(f"    Usuario:  {CONFIG['username']}")
    print(f"    Monitorea: {CONFIG['monitored_email']}")
    print(f"    Modelo:   {CONFIG['model']}")
    print(f"    LLM Test: {'SÍ' if llm_test else 'NO'}")
    print()

    if send_test:
        send_test_email()

    print(f"  {C.GRAY}Fetching email...{C.RESET}")
    msg = fetch_email(msg_id)
    analyze_email(msg, run_llm=llm_test)

    print(f"  {C.GRAY}Done.{C.RESET}\n")
