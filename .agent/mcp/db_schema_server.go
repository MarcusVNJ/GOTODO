package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func getDBConfig() DBConfig {
	return DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "root"),
		DBName:   getEnv("DB_NAME", "todo_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}

func connectDB(cfg DBConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}

func main() {
	cfg := getDBConfig()

	s := server.NewMCPServer(
		"GOTODO Database MCP",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Tool: List all tables
	s.AddTool(mcp.NewTool("list_tables",
		mcp.WithDescription("Lista todas as tabelas do banco de dados com seus schemas"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		db, err := connectDB(cfg)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Erro ao conectar ao banco: %v", err)), nil
		}
		defer db.Close()

		rows, err := db.Query(`
			SELECT table_schema, table_name, table_type
			FROM information_schema.tables
			WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
			ORDER BY table_schema, table_name
		`)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Erro ao listar tabelas: %v", err)), nil
		}
		defer rows.Close()

		type TableInfo struct {
			Schema string `json:"schema"`
			Name   string `json:"name"`
			Type   string `json:"type"`
		}

		var tables []TableInfo
		for rows.Next() {
			var t TableInfo
			if err := rows.Scan(&t.Schema, &t.Name, &t.Type); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Erro ao escanear tabela: %v", err)), nil
			}
			tables = append(tables, t)
		}

		jsonBytes, err := json.MarshalIndent(tables, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Erro ao serializar tabelas: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonBytes)), nil
	})

	// Tool: Describe table structure
	s.AddTool(mcp.NewTool("describe_table",
		mcp.WithDescription("Descreve a estrutura completa de uma tabela: colunas, tipos, constraints, defaults,.nullable"),
		mcp.WithString("table_name",
			mcp.Required(),
			mcp.Description("Nome da tabela (ex: tasks)"),
		),
		mcp.WithString("schema",
			mcp.Description("Schema da tabela (default: public)"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		tableName, ok := args["table_name"].(string)
		if !ok {
			return mcp.NewToolResultError("table_name é obrigatório"), nil
		}
		schema := "public"
		if s, ok := args["schema"].(string); ok && s != "" {
			schema = s
		}

		db, err := connectDB(cfg)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Erro ao conectar ao banco: %v", err)), nil
		}
		defer db.Close()

		type ColumnInfo struct {
			Name         string `json:"name"`
			DataType     string `json:"data_type"`
			Nullable      string `json:"nullable"`
			DefaultValue  string `json:"default_value"`
			IsPrimaryKey bool   `json:"is_primary_key"`
			MaxLength     string `json:"max_length,omitempty"`
		}

		rows, err := db.Query(`
			SELECT
				c.column_name,
				c.data_type,
				c.udt_name,
				c.is_nullable,
				COALESCE(c.column_default, ''),
				CASE WHEN pk.column_name IS NOT NULL THEN true ELSE false END,
				COALESCE(c.character_maximum_length::text, '')
			FROM information_schema.columns c
			LEFT JOIN (
				SELECT kcu.column_name, kcu.table_name, kcu.table_schema
				FROM information_schema.table_constraints tc
				JOIN information_schema.key_column_usage kcu
					ON tc.constraint_name = kcu.constraint_name
					AND tc.table_schema = kcu.table_schema
				WHERE tc.constraint_type = 'PRIMARY KEY'
			) pk ON pk.column_name = c.column_name
				AND pk.table_name = c.table_name
				AND pk.table_schema = c.table_schema
			WHERE c.table_name = $1 AND c.table_schema = $2
			ORDER BY c.ordinal_position
		`, tableName, schema)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Erro ao descrever tabela: %v", err)), nil
		}
		defer rows.Close()

		var columns []ColumnInfo
		for rows.Next() {
			var col ColumnInfo
			var udtName, maxLength string
			if err := rows.Scan(&col.Name, &col.DataType, &udtName, &col.Nullable, &col.DefaultValue, &col.IsPrimaryKey, &maxLength); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Erro ao escanear coluna: %v", err)), nil
			}
			if strings.HasPrefix(udtName, "task_") || col.DataType == "USER-DEFINED" {
				col.DataType = udtName
			}
			if maxLength != "" {
				col.MaxLength = maxLength
			}
			columns = append(columns, col)
		}

		if len(columns) == 0 {
			return mcp.NewToolResultError(fmt.Sprintf("Tabela '%s' não encontrada no schema '%s'", tableName, schema)), nil
		}

		jsonBytes, err := json.MarshalIndent(columns, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Erro ao serializar colunas: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonBytes)), nil
	})

	// Tool: List enum types
	s.AddTool(mcp.NewTool("list_enums",
		mcp.WithDescription("Lista todos os tipos ENUM do PostgreSQL com seus valores"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		db, err := connectDB(cfg)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Erro ao conectar ao banco: %v", err)), nil
		}
		defer db.Close()

		type EnumInfo struct {
			Name   string   `json:"name"`
			Values []string `json:"values"`
		}

		rows, err := db.Query(`
			SELECT t.typname, e.enumlabel
			FROM pg_type t
			JOIN pg_enum e ON t.oid = e.enumtypid
			JOIN pg_namespace n ON n.oid = t.typnamespace
			WHERE n.nspname = 'public'
			ORDER BY t.typname, e.enumsortorder
		`)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Erro ao listar enums: %v", err)), nil
		}
		defer rows.Close()

		enumMap := make(map[string][]string)
		for rows.Next() {
			var typeName, enumLabel string
			if err := rows.Scan(&typeName, &enumLabel); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Erro ao escanear enum: %v", err)), nil
			}
			enumMap[typeName] = append(enumMap[typeName], enumLabel)
		}

		var enums []EnumInfo
		for name, values := range enumMap {
			enums = append(enums, EnumInfo{Name: name, Values: values})
		}

		jsonBytes, err := json.MarshalIndent(enums, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Erro ao serializar enums: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonBytes)), nil
	})

	// Tool: List indexes
	s.AddTool(mcp.NewTool("list_indexes",
		mcp.WithDescription("Lista todos os indexes de uma tabela com suas colunas e tipo (unique ou não)"),
		mcp.WithString("table_name",
			mcp.Required(),
			mcp.Description("Nome da tabela (ex: tasks)"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		tableName, ok := args["table_name"].(string)
		if !ok {
			return mcp.NewToolResultError("table_name é obrigatório"), nil
		}

		db, err := connectDB(cfg)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Erro ao conectar ao banco: %v", err)), nil
		}
		defer db.Close()

		type IndexInfo struct {
			Name     string `json:"name"`
			Unique   bool   `json:"unique"`
			Columns  string `json:"columns"`
		}

		rows, err := db.Query(`
			SELECT
				i.relname AS index_name,
				ix.indisunique AS is_unique,
				string_agg(a.attname, ', ' ORDER BY array_position(ix.indkey, a.attnum)) AS columns
			FROM pg_class t
			JOIN pg_index ix ON t.oid = ix.indrelid
			JOIN pg_class i ON i.oid = ix.indexrelid
			JOIN pg_attribute a ON a.attrelid = t.oid AND a.attnum = ANY(ix.indkey)
			JOIN pg_namespace n ON n.oid = t.relnamespace
			WHERE t.relname = $1 AND n.nspname = 'public'
			GROUP BY i.relname, ix.indisunique
			ORDER BY i.relname
		`, tableName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Erro ao listar indexes: %v", err)), nil
		}
		defer rows.Close()

		var indexes []IndexInfo
		for rows.Next() {
			var idx IndexInfo
			if err := rows.Scan(&idx.Name, &idx.Unique, &idx.Columns); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Erro ao escanear index: %v", err)), nil
			}
			indexes = append(indexes, idx)
		}

		if len(indexes) == 0 {
			return mcp.NewToolResultText(fmt.Sprintf("Nenhum index encontrado para a tabela '%s'", tableName)), nil
		}

		jsonBytes, err := json.MarshalIndent(indexes, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Erro ao serializar indexes: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonBytes)), nil
	})

	// Tool: List foreign keys
	s.AddTool(mcp.NewTool("list_foreign_keys",
		mcp.WithDescription("Lista todas as foreign keys de uma tabela com a tabela e coluna referenciada"),
		mcp.WithString("table_name",
			mcp.Required(),
			mcp.Description("Nome da tabela (ex: tasks)"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		tableName, ok := args["table_name"].(string)
		if !ok {
			return mcp.NewToolResultError("table_name é obrigatório"), nil
		}

		db, err := connectDB(cfg)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Erro ao conectar ao banco: %v", err)), nil
		}
		defer db.Close()

		type FKInfo struct {
			ConstraintName    string `json:"constraint_name"`
			ColumnName        string `json:"column_name"`
			ReferencedTable   string `json:"referenced_table"`
			ReferencedColumn  string `json:"referenced_column"`
		}

		rows, err := db.Query(`
			SELECT
				kcu.constraint_name,
				kcu.column_name,
				ccu.table_name AS referenced_table,
				ccu.column_name AS referenced_column
			FROM information_schema.table_constraints tc
			JOIN information_schema.key_column_usage kcu
				ON tc.constraint_name = kcu.constraint_name
				AND tc.table_schema = kcu.table_schema
			JOIN information_schema.constraint_column_usage ccu
				ON tc.constraint_name = ccu.constraint_name
				AND tc.table_schema = ccu.table_schema
			WHERE tc.constraint_type = 'FOREIGN KEY'
				AND tc.table_name = $1
				AND tc.table_schema = 'public'
		`, tableName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Erro ao listar foreign keys: %v", err)), nil
		}
		defer rows.Close()

		var fks []FKInfo
		for rows.Next() {
			var fk FKInfo
			if err := rows.Scan(&fk.ConstraintName, &fk.ColumnName, &fk.ReferencedTable, &fk.ReferencedColumn); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Erro ao escanear FK: %v", err)), nil
			}
			fks = append(fks, fk)
		}

		if len(fks) == 0 {
			return mcp.NewToolResultText(fmt.Sprintf("Nenhuma foreign key encontrada para a tabela '%s'", tableName)), nil
		}

		jsonBytes, err := json.MarshalIndent(fks, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Erro ao serializar FKs: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonBytes)), nil
	})

	// Tool: Read existing domain models
	s.AddTool(mcp.NewTool("read_domain_models",
		mcp.WithDescription("Lê todos os modelos de domínio existentes em internal/core/models/ e retorna a estrutura de cada um (campos, tipos, métodos, factories)"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		modelsDir := "internal/core/models"

		entries, err := os.ReadDir(modelsDir)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Erro ao ler diretório de modelos: %v", err)), nil
		}

		type ModelInfo struct {
			File    string `json:"file"`
			Content string `json:"content"`
		}

		var models []ModelInfo
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
				continue
			}

			filePath := modelsDir + "/" + entry.Name()
			content, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}

			models = append(models, ModelInfo{
				File:    entry.Name(),
				Content: string(content),
			})
		}

		if len(models) == 0 {
			return mcp.NewToolResultText("Nenhum modelo de domínio encontrado em internal/core/models/"), nil
		}

		jsonBytes, err := json.MarshalIndent(models, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Erro ao serializar modelos: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonBytes)), nil
	})

	// Tool: Read existing enums
	s.AddTool(mcp.NewTool("read_domain_enums",
		mcp.WithDescription("Lê todos os enums existentes em internal/core/enums/ e retorna seus valores"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		enumsDir := "internal/core/enums"

		entries, err := os.ReadDir(enumsDir)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Erro ao ler diretório de enums: %v", err)), nil
		}

		type EnumInfo struct {
			File    string `json:"file"`
			Content string `json:"content"`
		}

		var enums []EnumInfo
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
				continue
			}

			filePath := enumsDir + "/" + entry.Name()
			content, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}

			enums = append(enums, EnumInfo{
				File:    entry.Name(),
				Content: string(content),
			})
		}

		if len(enums) == 0 {
			return mcp.NewToolResultText("Nenhum enum encontrado em internal/core/enums/"), nil
		}

		jsonBytes, err := json.MarshalIndent(enums, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Erro ao serializar enums: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonBytes)), nil
	})

	// Tool: Read existing error codes
	s.AddTool(mcp.NewTool("read_error_codes",
		mcp.WithDescription("Lê todos os códigos de erro existentes em internal/core/exceptions/codes/ e retorna seus valores"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		codesDir := "internal/core/exceptions/codes"

		entries, err := os.ReadDir(codesDir)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Erro ao ler diretório de códigos de erro: %v", err)), nil
		}

		type CodeInfo struct {
			File    string `json:"file"`
			Content string `json:"content"`
		}

		var codes []CodeInfo
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
				continue
			}

			filePath := codesDir + "/" + entry.Name()
			content, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}

			codes = append(codes, CodeInfo{
				File:    entry.Name(),
				Content: string(content),
			})
		}

		if len(codes) == 0 {
			return mcp.NewToolResultText("Nenhum código de erro encontrado em internal/core/exceptions/codes/"), nil
		}

		jsonBytes, err := json.MarshalIndent(codes, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Erro ao serializar códigos: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonBytes)), nil
	})

	// Tool: Full database schema overview
	s.AddTool(mcp.NewTool("database_schema_overview",
		mcp.WithDescription("Retorna uma visão geral completa do schema do banco: todas as tabelas com colunas, tipos, enums, indexes, foreign keys e soft delete status"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		db, err := connectDB(cfg)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Erro ao conectar ao banco: %v", err)), nil
		}
		defer db.Close()

		// Get all tables
		tableRows, err := db.Query(`
			SELECT table_name
			FROM information_schema.tables
			WHERE table_schema = 'public' AND table_type = 'BASE TABLE'
			ORDER BY table_name
		`)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Erro ao listar tabelas: %v", err)), nil
		}
		defer tableRows.Close()

		var tableNames []string
		for tableRows.Next() {
			var name string
			if err := tableRows.Scan(&name); err != nil {
				continue
			}
			tableNames = append(tableNames, name)
		}

		type ColumnInfo struct {
			Name         string `json:"name"`
			DataType     string `json:"data_type"`
			Nullable      string `json:"nullable"`
			DefaultValue  string `json:"default_value,omitempty"`
			IsPrimaryKey bool   `json:"is_primary_key"`
		}

		type TableSchema struct {
			Name      string        `json:"table_name"`
			Columns   []ColumnInfo  `json:"columns"`
			HasSoftDelete bool      `json:"has_soft_delete"`
		}

		var schemas []TableSchema
		for _, tableName := range tableNames {
			colRows, err := db.Query(`
				SELECT
					c.column_name,
					COALESCE(c.udt_name, c.data_type),
					c.is_nullable,
					COALESCE(c.column_default, ''),
					CASE WHEN pk.column_name IS NOT NULL THEN true ELSE false END
				FROM information_schema.columns c
				LEFT JOIN (
					SELECT kcu.column_name, kcu.table_name
					FROM information_schema.table_constraints tc
					JOIN information_schema.key_column_usage kcu
						ON tc.constraint_name = kcu.constraint_name
					WHERE tc.constraint_type = 'PRIMARY KEY'
				) pk ON pk.column_name = c.column_name AND pk.table_name = c.table_name
				WHERE c.table_name = $1 AND c.table_schema = 'public'
				ORDER BY c.ordinal_position
			`, tableName)
			if err != nil {
				continue
			}

			var columns []ColumnInfo
			hasSoftDelete := false
			for colRows.Next() {
				var col ColumnInfo
				if err := colRows.Scan(&col.Name, &col.DataType, &col.Nullable, &col.DefaultValue, &col.IsPrimaryKey); err != nil {
					continue
				}
				if col.Name == "deleted_at" {
					hasSoftDelete = true
				}
				columns = append(columns, col)
			}
			colRows.Close()

			schemas = append(schemas, TableSchema{
				Name:           tableName,
				Columns:        columns,
				HasSoftDelete:  hasSoftDelete,
			})
		}

		// Get enums
		enumRows, err := db.Query(`
			SELECT t.typname, e.enumlabel
			FROM pg_type t
			JOIN pg_enum e ON t.oid = e.enumtypid
			JOIN pg_namespace n ON n.oid = t.typnamespace
			WHERE n.nspname = 'public'
			ORDER BY t.typname, e.enumsortorder
		`)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Erro ao listar enums: %v", err)), nil
		}

		type EnumDef struct {
			Name   string   `json:"name"`
			Values []string `json:"values"`
		}

		enumMap := make(map[string][]string)
		for enumRows.Next() {
			var typeName, enumLabel string
			if err := enumRows.Scan(&typeName, &enumLabel); err != nil {
				continue
			}
			enumMap[typeName] = append(enumMap[typeName], enumLabel)
		}
		enumRows.Close()

		var enums []EnumDef
		for name, values := range enumMap {
			enums = append(enums, EnumDef{Name: name, Values: values})
		}

		type Overview struct {
			Tables []TableSchema `json:"tables"`
			Enums  []EnumDef     `json:"enums"`
		}

		overview := Overview{
			Tables: schemas,
			Enums:  enums,
		}

		jsonBytes, err := json.MarshalIndent(overview, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Erro ao serializar overview: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonBytes)), nil
	})

	log.Println("Starting GOTODO Database MCP server on stdio...")
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}