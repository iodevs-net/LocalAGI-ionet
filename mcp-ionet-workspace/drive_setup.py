#!/usr/bin/env python3
"""One-time Google Drive OAuth setup for IONET Workspace.
Uses Device Code Flow - no browser in container needed.

Usage:
  docker exec -it mcp-ionet-workspace python drive_setup.py

Then visit the URL shown, authorize with el.agent.ion@gmail.com.
Token auto-saves to /pool/ion-workspace/.drive-token.json
"""

import os
import sys
import json
import time
import requests

TOKEN_PATH = "/pool/ion-workspace/.drive-token.json"
SCOPES = "https://www.googleapis.com/auth/drive"


def main():
    client_id = os.environ.get("GOOGLE_DRIVE_CLIENT_ID")
    client_secret = os.environ.get("GOOGLE_DRIVE_CLIENT_SECRET")

    if not client_id or not client_secret:
        print("ERROR: Faltan variables GOOGLE_DRIVE_CLIENT_ID y GOOGLE_DRIVE_CLIENT_SECRET")
        print("Agrega las credenciales en docker-compose.yaml bajo mcp-ionet-workspace.environment")
        sys.exit(1)

    # Step 1: Get device code
    r = requests.post("https://oauth2.googleapis.com/device/code", data={
        "client_id": client_id,
        "scope": SCOPES,
    })
    d = r.json()

    if "error" in d:
        print(f"ERROR: {d.get('error_description', d['error'])}")
        sys.exit(1)

    print("\n" + "=" * 60)
    print("  AUTENTICACIÓN GOOGLE DRIVE PARA ION")
    print("=" * 60)
    print(f"\n  1. Visita:  {d['verification_uri']}")
    print(f"  2. Código:  {d['user_code']}")
    print(f"  3. Cuenta:  el.agent.ion@gmail.com")
    print(f"\n  Tienes {d['expires_in'] // 60} minutos antes que expire.")
    print("\n  Esperando autorización...\n")

    # Step 2: Poll until user authorizes
    start = time.time()
    while time.time() - start < d["expires_in"]:
        r = requests.post("https://oauth2.googleapis.com/token", data={
            "client_id": client_id,
            "client_secret": client_secret,
            "device_code": d["device_code"],
            "grant_type": "urn:ietf:params:oauth:grant-type:device_code",
        })
        t = r.json()

        if "refresh_token" in t:
            os.makedirs(os.path.dirname(TOKEN_PATH) or ".", exist_ok=True)
            with open(TOKEN_PATH, "w") as f:
                json.dump(t, f)
            print(f"  ✓ Token guardado en {TOKEN_PATH}")
            print(f"  ✓ ION ya puede usar Google Drive")
            return

        if t.get("error") == "access_denied":
            print("  ✗ Autorización denegada por el usuario.")
            sys.exit(1)

        if t.get("error") == "expired_token":
            print("  ✗ Código expirado. Ejecuta nuevamente.")
            sys.exit(1)

        time.sleep(5)

    print("  ✗ Tiempo de espera agotado.")
    sys.exit(1)


if __name__ == "__main__":
    main()
