// =============================================================================
// IONET - Script de Sincronización M365
// Descarga documentos de SharePoint/OneDrive y los guarda para indexación RAG
// =============================================================================

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Configuración
type Config struct {
	ClientID     string
	ClientSecret string
	TenantID     string
	OutputDir    string
	SiteURL      string // SharePoint site URL (opcional)
	DriveID      string // OneDrive drive ID (opcional)
}

// Microsoft Graph API endpoints
const (
	graphAPIBase      = "https://graph.microsoft.com/v1.0"
	authEndpoint      = "https://login.microsoftonline.com/%s/oauth2/v2.0/token"
	tokenGrantType    = "client_credentials"
)

// Token response
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// Drive item (archivo o carpeta)
type DriveItem struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Size          int64      `json:"size"`
	ContentURL    string     `json:"@microsoft.graph.downloadUrl"`
	File          *FileInfo  `json:"file,omitempty"`
	Folder        *FolderInfo `json:"folder,omitempty"`
	LastModified  *time.Time `json:"lastModifiedDateTime"`
}

// FileInfo contiene metadatos de archivo
type FileInfo struct {
	MimeType string `json:"mimeType"`
}

// FolderInfo contiene metadatos de carpeta
type FolderInfo struct {
	ChildCount int `json:"childCount"`
}

func main() {
	log.SetPrefix("[IONET-M365-SYNC] ")

	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: No se encontró archivo .env, usando variables de entorno del sistema")
	}

	config := Config{
		ClientID:     os.Getenv("M365_CLIENT_ID"),
		ClientSecret: os.Getenv("M365_CLIENT_SECRET"),
		TenantID:     os.Getenv("M365_TENANT_ID"),
		OutputDir:    getEnvOrDefault("RAG_DATA_DIR", "/rag/raw"),
		SiteURL:      os.Getenv("M365_SITE_URL"),
		DriveID:      os.Getenv("M365_DRIVE_ID"),
	}

	// Validar configuración
	if config.ClientID == "" || config.ClientSecret == "" || config.TenantID == "" {
		log.Fatal("Error: M365_CLIENT_ID, M365_CLIENT_SECRET y M365_TENANT_ID son requeridos")
	}

	ctx := context.Background()

	// Obtener token de acceso
	accessToken, err := getAccessToken(ctx, config)
	if err != nil {
		log.Fatalf("Error autenticando con Microsoft Graph: %v", err)
	}

	log.Printf("Autenticación exitosa con Microsoft Graph API")

	// Sincronizar documentos
	if err := syncDocuments(ctx, accessToken, config); err != nil {
		log.Fatalf("Error sincronizando documentos: %v", err)
	}

	log.Println("Sincronización completada exitosamente")
}

// getAccessToken obtiene un token de acceso usando client credentials
func getAccessToken(ctx context.Context, config Config) (string, error) {
	authURL := fmt.Sprintf(authEndpoint, config.TenantID)

	// Preparar request body
	data := fmt.Sprintf(
		"grant_type=%s&client_id=%s&client_secret=%s&scope=https://graph.microsoft.com/.default",
		tokenGrantType, config.ClientID, config.ClientSecret,
	)

	req, err := http.NewRequestWithContext(ctx, "POST", authURL, strings.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("error creando request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error en request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("error autenticando: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("error decodificando respuesta: %w", err)
	}

	return tokenResp.AccessToken, nil
}

// syncDocuments sincroniza los documentos de SharePoint/OneDrive
func syncDocuments(ctx context.Context, accessToken string, config Config) error {
	// Crear directorio de salida si no existe
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return fmt.Errorf("error creando directorio de salida: %w", err)
	}

	// Determinar endpoint según configuración
	var endpoint string
	if config.DriveID != "" {
		// Usar drive específico
		endpoint = fmt.Sprintf("%s/drives/%s/root/children", graphAPIBase, config.DriveID)
	} else if config.SiteURL != "" {
		// Usar SharePoint site
		siteName := extractSiteName(config.SiteURL)
		endpoint = fmt.Sprintf("%s/sites/%s:/documents:/children", graphAPIBase, siteName)
	} else {
		// Usar OneDrive del usuario (primero obtener el drive del usuario)
		var err error
		endpoint, err = getUserDriveEndpoint(ctx, accessToken)
		if err != nil {
			return fmt.Errorf("error obteniendo endpoint de OneDrive: %w", err)
		}
	}

	log.Printf("Sincronizando desde: %s", endpoint)

	// Obtener lista de archivos
	items, err := listDriveItems(ctx, accessToken, endpoint)
	if err != nil {
		return fmt.Errorf("error listando archivos: %w", err)
	}

	log.Printf("Encontrados %d elementos", len(items))

	// Descargar archivos (excluyendo carpetas)
	downloaded := 0
	for _, item := range items {
		if item.Folder != nil {
			// Es una carpeta, procesar recursivamente
			folderEndpoint := fmt.Sprintf("%s/children", strings.Replace(endpoint, "/root/children", "", 1))
			folderEndpoint = fmt.Sprintf("%s/items/%s/children", strings.Split(endpoint, "/drives/")[0], item.ID)
			if err := downloadFolder(ctx, accessToken, folderEndpoint, config.OutputDir, item.Name); err != nil {
				log.Printf("Aviso: error procesando carpeta %s: %v", item.Name, err)
			}
			continue
		}

		if item.File == nil {
			continue // No es un archivo regular
		}

		// Verificar extensión (solo descargar documentos)
		ext := strings.ToLower(filepath.Ext(item.Name))
		supportedExts := map[string]bool{
			".pdf": true, ".docx": true, ".doc": true,
			".xlsx": true, ".xls": true, ".pptx": true, ".ppt": true,
			".txt": true, ".md": true, ".html": true, ".htm": true,
		}

		if !supportedExts[ext] {
			log.Printf("Omitiendo archivo no soportado: %s (%s)", item.Name, ext)
			continue
		}

		// Descargar archivo
		if err := downloadFile(ctx, accessToken, item, config.OutputDir); err != nil {
			log.Printf("Aviso: error descargando %s: %v", item.Name, err)
			continue
		}
		downloaded++
	}

	log.Printf("Descargados %d archivos", downloaded)
	return nil
}

// getUserDriveEndpoint obtiene el endpoint del OneDrive del usuario autenticado
func getUserDriveEndpoint(ctx context.Context, accessToken string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", graphAPIBase+"/me/drive", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	id, ok := result["id"].(string)
	if !ok {
		return "", fmt.Errorf("no se pudo obtener el drive ID")
	}

	return fmt.Sprintf("%s/drives/%s/root/children", graphAPIBase, id), nil
}

// listDriveItems lista los elementos en un path de drive
func listDriveItems(ctx context.Context, accessToken, endpoint string) ([]DriveItem, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error listing: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var result struct {
		Value []DriveItem `json:"value"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Value, nil
}

// downloadFolder descarga una carpeta recursivamente
func downloadFolder(ctx context.Context, accessToken, endpoint, baseDir, folderName string) error {
	items, err := listDriveItems(ctx, accessToken, endpoint)
	if err != nil {
		return err
	}

	folderPath := filepath.Join(baseDir, folderName)
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return err
	}

	for _, item := range items {
		if item.Folder != nil {
			// Carpeta anidada
			subFolderEndpoint := endpoint + "/" + item.ID + "/children"
			if err := downloadFolder(ctx, accessToken, subFolderEndpoint, folderPath, item.Name); err != nil {
				log.Printf("Aviso: error en subcarpeta %s: %v", item.Name, err)
			}
			continue
		}

		if item.File != nil {
			if err := downloadFile(ctx, accessToken, item, folderPath); err != nil {
				log.Printf("Aviso: error descargando %s: %v", item.Name, err)
			}
		}
	}

	return nil
}

// downloadFile descarga un archivo individual
func downloadFile(ctx context.Context, accessToken string, item DriveItem, outputDir string) error {
	// Crear estructura de directorios por fecha si hay lastModified
	outputPath := filepath.Join(outputDir, item.Name)

	// Verificar si ya existe y es más reciente
	if info, err := os.Stat(outputPath); err == nil {
		if item.LastModified != nil && info.ModTime().After(*item.LastModified) {
			log.Printf("Archivo ya actualizado: %s", item.Name)
			return nil
		}
	}

	log.Printf("Descargando: %s (%d bytes)", item.Name, item.Size)

	req, err := http.NewRequestWithContext(ctx, "GET", item.ContentURL, nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 300 * time.Second} // 5 min para archivos grandes
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error descargando: status=%d", resp.StatusCode)
	}

	// Crear archivo
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copiar contenido
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	log.Printf("Guardado: %s", outputPath)
	return nil
}

// extractSiteName extrae el nombre del site de una URL de SharePoint
func extractSiteName(siteURL string) string {
	// Ejemplo: https://ionet.sharepoint.com/sites/documentos
	parts := strings.Split(siteURL, "/sites/")
	if len(parts) > 1 {
		return strings.Split(parts[1], "/")[0]
	}
	return siteURL
}

// getEnvOrDefault obtiene variable de entorno o usa valor por defecto
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// downloadFileSimple versión simplificada que usa el API standard
func downloadFileSimple(ctx context.Context, accessToken, itemID, outputPath string) error {
	endpoint := fmt.Sprintf("%s/me/drive/items/%s/content", graphAPIBase, itemID)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 300 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	// Parsear Location header que contiene la URL de descarga real
	downloadURL := resp.Header.Get("Location")
	if downloadURL == "" {
		// El contenido viene directamente
		out, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer out.Close()
		_, err = io.Copy(out, resp.Body)
		return err
	}

	// Descargar desde la URL real
	downloadResp, err := http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer downloadResp.Body.Close()

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, downloadResp.Body)
	return err
}

// =============================================================================
// Funciones auxiliares para debugging y logging
// =============================================================================

// printJSON imprime un objeto como JSON formateado (para debugging)
func printJSON(label string, v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Printf("[DEBUG] %s:\n%s\n", label, string(b))
}

// downloadFileWithRetry descarga con reintentos
func downloadFileWithRetry(ctx context.Context, accessToken string, item DriveItem, outputDir string, maxRetries int) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if err := downloadFileSimple(ctx, accessToken, item.ID, filepath.Join(outputDir, item.Name)); err != nil {
			lastErr = err
			log.Printf("Reintento %d/%d para %s: %v", i+1, maxRetries, item.Name, err)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		return nil
	}
	return lastErr
}

// =============================================================================
// NOTAS DE USO
// =============================================================================
//
// 1. Registrar una app en Azure AD:
//    - Portal: https://portal.azure.com > Azure Active Directory > App registrations
//    - Nombre: IONET RAG Sync
//    - Tipo: Accounts in this organizational directory only
//
// 2. Permisos de API (Graph API):
//    - Sites.Read.All (SharePoint)
//    - Files.Read.All (OneDrive)
//    - User.Read (básico)
//
// 3. Client Secret:
//    - Certificates & secrets > New client secret
//    - Copiar el valor (solo se muestra una vez)
//
// 4. Variables de entorno requeridas:
//    M365_CLIENT_ID=<app-id>
//    M365_CLIENT_SECRET=<client-secret>
//    M365_TENANT_ID=<tenant-id>
//    M365_SITE_URL=https://ionet.sharepoint.com/sites/documentos  (opcional)
//
// 5. Ejecutar:
//
//    go run sync.go
//
//    o compilar y usar en producción:
//
//    go build -o sync-m365 sync.go
//    ./sync-m365
//
// 6. Para ejecutar en cron (Linux):
//    0 */6 * * * /app/scripts/m365/sync-m365 >> /var/log/ionet-m365-sync.log 2>&1
//
// =============================================================================
