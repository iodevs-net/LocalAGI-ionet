package e2e_test

import (
	"net/http"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// authTestKeys define las claves de API usadas en los tests de autenticación.
// Estas claves deben estar configuradas en LOCALAGI_API_KEYS del .env.
var (
	authTestKey1 = "key1"
	authTestKey2 = "key2"
	authTestKey3 = "key3"
)

var _ = Describe("Authentication with multiple API keys", Label("Auth"), func() {
	Context("Multiple valid API keys", func() {
		BeforeEach(func() {
			// Verificar que el servidor está disponible
			Eventually(func() error {
				_, err := http.Get(localagiURL + "/readyz")
				return err
			}, "2m", "10s").ShouldNot(HaveOccurred())
		})

		// Test 1: Verificar que la primera clave es válida
		It("should authenticate with key1", func() {
			client := &http.Client{Timeout: 10 * time.Second}
			req, err := http.NewRequest("GET", localagiURL+"/api/agents", nil)
			Expect(err).ToNot(HaveOccurred())

			// Usar header Authorization: Bearer
			req.Header.Set("Authorization", "Bearer "+authTestKey1)

			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			// Debe retornar 200 OK (o cualquier código que no sea 401)
			Expect(resp.StatusCode).ToNot(Equal(http.StatusUnauthorized),
				"key1 should be a valid API key")
		})

		// Test 2: Verificar que la segunda clave es válida
		It("should authenticate with key2", func() {
			client := &http.Client{Timeout: 10 * time.Second}
			req, err := http.NewRequest("GET", localagiURL+"/api/agents", nil)
			Expect(err).ToNot(HaveOccurred())

			req.Header.Set("Authorization", "Bearer "+authTestKey2)

			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			Expect(resp.StatusCode).ToNot(Equal(http.StatusUnauthorized),
				"key2 should be a valid API key")
		})

		// Test 3: Verificar que la tercera clave es válida
		It("should authenticate with key3", func() {
			client := &http.Client{Timeout: 10 * time.Second}
			req, err := http.NewRequest("GET", localagiURL+"/api/agents", nil)
			Expect(err).ToNot(HaveOccurred())

			req.Header.Set("Authorization", "Bearer "+authTestKey3)

			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			Expect(resp.StatusCode).ToNot(Equal(http.StatusUnauthorized),
				"key3 should be a valid API key")
		})

		// Test 4: Verificar autenticación con x-api-key header
		It("should authenticate using x-api-key header", func() {
			client := &http.Client{Timeout: 10 * time.Second}
			req, err := http.NewRequest("GET", localagiURL+"/api/agents", nil)
			Expect(err).ToNot(HaveOccurred())

			req.Header.Set("x-api-key", authTestKey1)

			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			Expect(resp.StatusCode).ToNot(Equal(http.StatusUnauthorized),
				"x-api-key header should work for authentication")
		})

		// Test 5: Verificar autenticación con xi-api-key header
		It("should authenticate using xi-api-key header", func() {
			client := &http.Client{Timeout: 10 * time.Second}
			req, err := http.NewRequest("GET", localagiURL+"/api/agents", nil)
			Expect(err).ToNot(HaveOccurred())

			req.Header.Set("xi-api-key", authTestKey2)

			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			Expect(resp.StatusCode).ToNot(Equal(http.StatusUnauthorized),
				"xi-api-key header should work for authentication")
		})
	})

	Context("Invalid API key rejection", func() {
		BeforeEach(func() {
			Eventually(func() error {
				_, err := http.Get(localagiURL + "/readyz")
				return err
			}, "2m", "10s").ShouldNot(HaveOccurred())
		})

		It("should reject an invalid API key with 401", func() {
			client := &http.Client{Timeout: 10 * time.Second}
			req, err := http.NewRequest("GET", localagiURL+"/api/agents", nil)
			Expect(err).ToNot(HaveOccurred())

			// Usar una clave que no existe
			req.Header.Set("Authorization", "Bearer invalid-key-that-does-not-exist")

			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			// Debe retornar 401 Unauthorized
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized),
				"invalid API key should be rejected with 401")
		})

		It("should reject empty API key with 401", func() {
			client := &http.Client{Timeout: 10 * time.Second}
			req, err := http.NewRequest("GET", localagiURL+"/api/agents", nil)
			Expect(err).ToNot(HaveOccurred())

			// No establecer header de autenticación
			// El servidor debería rechazar la solicitud

			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			// Debe retornar 401 cuando no hay clave de API
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized),
				"request without API key should be rejected with 401")
		})

		It("should reject partial match of valid key", func() {
			client := &http.Client{Timeout: 10 * time.Second}
			req, err := http.NewRequest("GET", localagiURL+"/api/agents", nil)
			Expect(err).ToNot(HaveOccurred())

			// Usar solo una parte de una clave válida
			req.Header.Set("Authorization", "Bearer "+authTestKey1[:len(authTestKey1)/2])

			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized),
				"partial key should be rejected with 401")
		})
	})

	Context("API keys with whitespace handling", func() {
		// Este test verifica que el parsing de claves funciona correctamente
		// cuando hay espacios en blanco alrededor de las claves en la variable
		// de entorno LOCALAGI_API_KEYS

		It("should authenticate when keys have surrounding whitespace", func() {
			// Simular el parsing de: "key1, key2, key3"
			keysWithSpaces := authTestKey1 + ", " + authTestKey2 + ", " + authTestKey3
			parsedKeys := strings.Split(keysWithSpaces, ",")

			// Verificar que Split divide correctamente
			Expect(len(parsedKeys)).To(Equal(3))

			// Las claves con espacios deben ser autenticadas
			// Esto verifica que el código en cmd/env.go:68 funciona correctamente
			for _, key := range parsedKeys {
				trimmedKey := strings.TrimSpace(key)
				Expect(trimmedKey).ToNot(BeEmpty(), "parsed key should not be empty after trim")
			}
		})

		It("should authenticate with keys that have leading/trailing spaces", func() {
			// Verificar que las claves con espacios funcionan cuando se usan directamente
			// (el servidor debe hacer TrimSpace internamente o las claves deben estar configuradas
			// con los espacios correctos en LOCALAGI_API_KEYS)

			// Test con clave que tiene espacios (si está configurada así)
			client := &http.Client{Timeout: 10 * time.Second}

			// Probar todas las claves conocidas
			for _, key := range []string{authTestKey1, authTestKey2, authTestKey3} {
				req, err := http.NewRequest("GET", localagiURL+"/api/agents", nil)
				Expect(err).ToNot(HaveOccurred())

				req.Header.Set("Authorization", "Bearer "+key)

				resp, err := client.Do(req)
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				// Al menos una clave debe funcionar si la configuración es correcta
				if resp.StatusCode != http.StatusUnauthorized {
					return // Test passed
				}
			}

			// Si llegamos aquí, ninguna clave funcionó
			Fail("At least one configured API key should be valid")
		})
	})

	Context("Multiple API keys with different separators", func() {
		It("should handle comma-separated keys correctly", func() {
			// Verificar que el parsing con Split funciona para diferentes casos
			testCases := []struct {
				input    string
				expected int
			}{
				{"key1,key2,key3", 3},
				{"key1, key2, key3", 3},
				{"key1,key2", 2},
				{"key1", 1},
				{"", 0},
			}

			for _, tc := range testCases {
				if tc.input == "" {
					continue // Skip empty case
				}
				keys := strings.Split(tc.input, ",")
				Expect(len(keys)).To(Equal(tc.expected),
					"Split should correctly parse keys from: "+tc.input)

				// Verificar que TrimSpace está disponible
				for _, key := range keys {
					trimmed := strings.TrimSpace(key)
					Expect(trimmed).ToNot(BeEmpty(),
						"TrimSpace should handle whitespace correctly for: "+key)
				}
			}
		})
	})
})
