"""Google Drive wrapper for IONET Workspace MCP.
Uses service account auth. Key file mounted in /app/drive-key.json"""

import os
import json
import logging
import io

from google.oauth2.service_account import Credentials
from googleapiclient.discovery import build
from googleapiclient.http import MediaIoBaseUpload
from googleapiclient.errors import HttpError

logger = logging.getLogger("mcp-ionet-workspace.drive")

KEY_PATH = "/app/drive-key.json"
SCOPES = ["https://www.googleapis.com/auth/drive"]
MAX_FILE_SIZE = 5 * 1024 * 1024  # 5 MB


class DriveManager:
    """Thin wrapper around Google Drive v3 API. Uses service account."""

    def __init__(self):
        self.service = None
        self._init_service()

    def _init_service(self):
        if not os.path.exists(KEY_PATH):
            logger.info("Service account key not found at %s", KEY_PATH)
            return
        try:
            creds = Credentials.from_service_account_file(KEY_PATH, scopes=SCOPES)
            self.service = build("drive", "v3", credentials=creds, cache_discovery=False)
            logger.info("Drive service initialized via service account")
        except Exception as e:
            logger.warning("Drive init failed: %s", e)

    def ready(self) -> bool:
        return self.service is not None

    def list_files(self, query: str = "trashed=false", page_size: int = 20) -> list:
        results = self.service.files().list(
            q=query,
            pageSize=page_size,
            fields="files(id,name,mimeType,size,modifiedTime,webViewLink)",
        ).execute()
        return [
            {
                "id": f["id"],
                "nombre": f["name"],
                "tipo": f["mimeType"],
                "tamano": int(f.get("size", 0)),
                "modificado": f.get("modifiedTime", ""),
                "url": f.get("webViewLink", ""),
            }
            for f in results.get("files", [])
        ]

    def list_folder(self, folder_id: str = "root") -> list:
        return self.list_files(query=f"'{folder_id}' in parents and trashed=false")

    def search(self, term: str) -> list:
        safe = term.replace("'", "\\'")
        return self.list_files(query=f"name contains '{safe}' and trashed=false")

    def read_text(self, file_id: str) -> dict:
        try:
            f = self.service.files().get(fileId=file_id, fields="id,name,mimeType,size").execute()
            mime = f["mimeType"]
            if mime == "application/vnd.google-apps.document":
                content = self.service.files().export_media(fileId=file_id, mimeType="text/plain").execute()
            else:
                content = self.service.files().get_media(fileId=file_id).execute()
            text = content.decode("utf-8") if isinstance(content, bytes) else content
            if len(text) > MAX_FILE_SIZE:
                return {"error": "Archivo demasiado grande para leer como texto"}
            return {"id": file_id, "nombre": f["name"], "contenido": text, "tamano": len(text)}
        except HttpError as e:
            return {"error": f"Error al leer: {e.reason}"}

    def upload_text(self, name: str, content: str, parent_id: str = None) -> dict:
        file_metadata = {"name": name, "mimeType": "text/plain"}
        if parent_id:
            file_metadata["parents"] = [parent_id]

        media = MediaIoBaseUpload(
            io.BytesIO(content.encode("utf-8")),
            mimetype="text/plain",
            resumable=True,
        )
        try:
            f = self.service.files().create(
                body=file_metadata, media_body=media,
                fields="id,name,webViewLink"
            ).execute()
            return {"id": f["id"], "nombre": f["name"], "url": f.get("webViewLink", "")}
        except HttpError as e:
            return {"error": f"Error al subir: {e.reason}"}

    def create_folder(self, name: str, parent_id: str = None) -> dict:
        file_metadata = {
            "name": name,
            "mimeType": "application/vnd.google-apps.folder",
        }
        if parent_id:
            file_metadata["parents"] = [parent_id]
        try:
            f = self.service.files().create(
                body=file_metadata, fields="id,name,webViewLink"
            ).execute()
            return {"id": f["id"], "nombre": f["name"], "url": f.get("webViewLink", "")}
        except HttpError as e:
            return {"error": f"Error al crear carpeta: {e.reason}"}

    def get_auth_status(self) -> dict:
        if not self.service:
            return {"autenticado": False, "email": None, "nombre": None}
        try:
            about = self.service.about().get(fields="user").execute()
            user = about.get("user", {})
            return {
                "autenticado": True,
                "email": user.get("emailAddress", ""),
                "nombre": user.get("displayName", ""),
            }
        except HttpError:
            return {"autenticado": False, "email": None, "nombre": None}
