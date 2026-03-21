package migrate

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Только базовые имена файлов внутри dir (без path traversal).
var migrationFileRe = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]*\.sql$`)

func Up(ctx context.Context, pool *pgxpool.Pool, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		ln := strings.ToLower(name)
		if !migrationFileRe.MatchString(ln) || filepath.Base(name) != name {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)
	for _, n := range names {
		path := filepath.Join(dir, n)
		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", n, err)
		}
		if _, err := pool.Exec(ctx, string(b)); err != nil {
			return fmt.Errorf("exec %s: %w", n, err)
		}
	}
	return nil
}
