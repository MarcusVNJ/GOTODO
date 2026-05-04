---
description: "Implementa tarefas do Planner seguindo Arquitetura Hexagonal e bases de conhecimento"
mode: primary
temperature: 0.2
permission:
  edit: allow
  bash: allow
  task: allow
color: "#f59e0b"
---

{file:../.agent/agents/software-engineer.md}

CONTEXTO DO PROJETO (sempre consultar antes de implementar):

{file:../.agent/knowledge-base/architecture.md}

{file:../.agent/knowledge-base/codebase-map.md}

{file:../.agent/knowledge-base/conventions-and-standards.md}

{file:../.agent/knowledge-base/di-patterns.md}

{file:../.agent/knowledge-base/domain-model.md}

{file:../.agent/knowledge-base/error-handling.md}

{file:../.agent/knowledge-base/database.md}

{file:../.agent/knowledge-base/api-resources.md}

{file:../.agent/knowledge-base/tech-stack.md}

{file:../.agent/knowledge-base/testing.md}

INSTRUÇÕES DE EXECUÇÃO:

1. Ao ser invocado, identifique qual tarefa implementar (especificada pelo usuário ou a mais antiga em `.agent/tasks/output/`)
2. Leia o plano completo (ex: `task01_plan.md`) e extraia: tipo, objetivo, regras de negócio, especificações técnicas e ordem de implementação
3. Consulte as bases de conhecimento acima antes de escrever qualquer código
4. Para o schema do banco, siga a ordem de precedência: MCP (primário) > Migrations (fallback) > Código-fonte (fallback secundário)
5. Siga a Ordem de Implementação do plano: Core → Outbound → Inbound → Orchestrator
6. Se o plano incluir migrations, crie os arquivos `.up.sql` e `.down.sql` em `migrations/` e execute `migrate -database "$(grep DATABASE_URL .env | cut -d= -f2-)" -path migrations up`
7. Após implementar, execute `go build ./...` e `go test ./...` para verificar
8. Se houver falhas de compilação ou testes, corrija antes de considerar a tarefa concluída
9. NUNCA importe bibliotecas de infraestrutura em `internal/core/` (net/http, pgx, json, fx, etc.)
10. NUNCA defina HTTP status codes em handlers — deixe HandlerException classificar
11. SEMPRE use soft delete (UPDATE SET deleted_at) — nunca DELETE FROM
12. Após concluir, recomende ao usuário executar o `knowledge-base-creator` para atualizar a base de conhecimento