package testutil

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"palabrario/core-service/internal/database"
)

const (
	// DefaultTestAdminURL apunta a la base administrativa 'postgres' del contenedor efímero (:5433)
	// para poder clonar 'palabrario_test' mediante TEMPLATE sin mantener conexiones abiertas sobre la plantilla.
	DefaultTestAdminURL = "postgresql://postgres:postgres@localhost:5433/postgres?sslmode=disable"
	TemplateDBName      = "palabrario_test"

	// templateAdvisoryLockID es una llave global en el motor de PostgreSQL 17 que serializa
	// la operación CREATE DATABASE ... TEMPLATE entre distintos procesos 'go test' concurrentes
	// (ej. cuando 'go test ./core-service/...' ejecuta los binarios de 'internal/api' e 'internal/service' al mismo tiempo).
	templateAdvisoryLockID int64 = 84920174
)

// SetupTestDB crea un clon instantáneo e independiente de la base de datos semilla ('palabrario_test')
// utilizando el mecanismo nativo CREATE DATABASE ... TEMPLATE de PostgreSQL sobre RAM (tmpfs).
//
// Ventajas arquitectónicas:
//  1. 100% agnóstico al esquema: No contiene nombres de tablas ni IDs quemados; cualquier cambio futuro
//     en SCHEMA_FINAL.sql es heredado automáticamente.
//  2. Aislamiento total de escrituras, UPDATEs y DELETEs sobre datos semilla (Silas, Gael, popote, etc.).
//  3. Sincronización inter-proceso con pg_advisory_lock: Permite que 'go test ./core-service/...' ejecute
//     múltiples paquetes simultáneamente sin el error SQLSTATE 55006 ("source database is being accessed by other users").
//  4. Prevención de envenenamiento del caché de Go: Si el contenedor no está activo, falla con t.Fatalf
//     (en vez de t.Skipf) porque en Go un 'SKIP' devuelve código 0 y se queda pegado en '(cached)' cuando se corre
//     'go test' sin el flag '-count=1'.
func SetupTestDB(t *testing.T) *database.DB {
	t.Helper()

	baseURL := os.Getenv("TEST_DATABASE_URL")
	if baseURL == "" {
		baseURL = DefaultTestAdminURL
	}

	adminURL, err := replaceDBNameInURL(baseURL, "postgres")
	if err != nil {
		t.Fatalf("URL de base de datos de prueba inválida (%s): %v", baseURL, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	adminPool, err := pgxpool.New(ctx, adminURL)
	if err != nil {
		t.Fatalf("No se pudo inicializar pool admin en %s: %v", adminURL, err)
		return nil
	}

	if err := adminPool.Ping(ctx); err != nil {
		adminPool.Close()
		t.Fatalf("Contenedor postgres-test (:5433) no disponible (%v). Levántalo con 'task test:up' o ejecuta 'task test'", err)
		return nil
	}

	rawID := strings.ReplaceAll(uuid.NewString(), "-", "")
	testDBName := fmt.Sprintf("test_%s", rawID[:12])

	// Adquirimos una conexión dedicada para tomar el candado inter-proceso (pg_advisory_lock) en PostgreSQL
	conn, err := adminPool.Acquire(ctx)
	if err != nil {
		adminPool.Close()
		t.Fatalf("Fallo adquiriendo conexión admin para clonar DB: %v", err)
	}

	_, _ = conn.Exec(ctx, "SELECT pg_advisory_lock($1)", templateAdvisoryLockID)

	createSQL := fmt.Sprintf("CREATE DATABASE %s TEMPLATE %s", testDBName, TemplateDBName)
	_, createErr := conn.Exec(ctx, createSQL)
	if createErr != nil {
		// Fallback autocurativo: SOLO si un cliente externo (ej. DBeaver) dejó una sesión colgada
		// sobre la plantilla 'palabrario_test', la desconectamos quirúrgicamente y reintentamos una vez.
		_, _ = conn.Exec(ctx, `
			SELECT pg_terminate_backend(pid)
			FROM pg_stat_activity
			WHERE datname = $1 AND pid <> pg_backend_pid()
		`, TemplateDBName)
		_, createErr = conn.Exec(ctx, createSQL)
	}

	_, _ = conn.Exec(ctx, "SELECT pg_advisory_unlock($1)", templateAdvisoryLockID)
	conn.Release()

	if createErr != nil {
		adminPool.Close()
		t.Fatalf("Fallo al clonar base de datos efímera desde plantilla '%s': %v", TemplateDBName, createErr)
	}

	testDBURL, err := replaceDBNameInURL(baseURL, testDBName)
	if err != nil {
		adminPool.Close()
		t.Fatalf("Fallo construyendo URL para DB efímera '%s': %v", testDBName, err)
	}

	testDB, err := database.Connect(ctx, testDBURL)
	if err != nil {
		_, _ = adminPool.Exec(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS %s WITH (FORCE)", testDBName))
		adminPool.Close()
		t.Fatalf("Fallo conectando a la base de datos efímera '%s': %v", testDBName, err)
	}

	t.Cleanup(func() {
		testDB.Close()

		dropCtx, dropCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer dropCancel()

		dropSQL := fmt.Sprintf("DROP DATABASE IF EXISTS %s WITH (FORCE)", testDBName)
		_, _ = adminPool.Exec(dropCtx, dropSQL)
		adminPool.Close()
	})

	return testDB
}

func replaceDBNameInURL(rawURL, newDBName string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	parsed.Path = "/" + newDBName
	return parsed.String(), nil
}
