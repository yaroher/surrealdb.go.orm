package migrator

import "io/fs"

// Mode defines migration strictness.
type Mode string

const (
	ModeStrict Mode = "strict"
	ModeLax    Mode = "lax"
)

// Config holds migrator configuration.
type Config struct {
	Dir            string
	Mode           Mode
	DSN            string
	NS             string
	DB             string
	Username       string
	Password       string
	Force          bool
	TwoWay         bool
	RenameStrategy string
	RenameExpr     string
	GrantsAlways   bool
}

// Source provides access to migration files.
type Source interface {
	Migrations() ([]Migration, error)
}

// Migration is a single migration pair.
type Migration struct {
	ID       string
	UpSQL    string
	DownSQL  string
	Checksum string
}

// MigrationRecord represents an applied migration row.
type MigrationRecord struct {
	ID        string
	Checksum  string
	AppliedAt string
}

// MigrationStatus represents file migration status.
type MigrationStatus struct {
	Migration Migration
	Applied   bool
	AppliedAt string
	Checksum  string
}

// EmbeddedSource reads migrations from an embedded filesystem.
type EmbeddedSource struct {
	FS  fs.FS
	Dir string
}
