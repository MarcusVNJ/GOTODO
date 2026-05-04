---
description: "Cria e mantém testes de unidade focados em regras de negócio"
mode: primary
temperature: 0.1
permission:
  edit: allow
  bash: allow
  task: allow
color: "#8b5cf6"
---

{file:../.agent/agents/test-engineer.md}

CONTEXTO DO PROJETO (sempre consultar antes de criar testes):

{file:../.agent/knowledge-base/testing.md}

{file:../.agent/knowledge-base/api-resources.md}

{file:../.agent/knowledge-base/domain-model.md}

{file:../.agent/knowledge-base/conventions-and-standards.md}

{file:../.agent/knowledge-base/architecture.md}

{file:../.agent/knowledge-base/error-handling.md}

{file:../.agent/knowledge-base/codebase-map.md}

{file:../.agent/knowledge-base/di-patterns.md}

{file:../.agent/knowledge-base/tech-stack.md}

INSTRUÇÕES DE EXECUÇÃO:

1. Ao ser invocado, identificar qual tarefa/UseCase testar (especificado pelo usuário ou o mais recente em `.agent/tasks/output/`)
2. Ler o plano da tarefa e extrair regras de negócio, validações e error codes
3. Ler o código-fonte dos Models e UseCases relevantes para identificar TODAS as regras de negócio implementadas
4. Consultar as bases de conhecimento acima (especialmente `testing.md`, `domain-model.md`, `api-resources.md`)
5. Verificar se o Mock do Repository está completo e atualizado — atualizar se necessário
6. Criar os testes seguindo o padrão AAA com `stretchr/testify`
7. Usar `package <entity>_test` para black-box testing
8. NOMEAR testes com padrão `Test<UseCase>_<Cenário>`
9. Testar TODAS as regras de negócio identificadas — não apenas caminho feliz
10. Verificar error codes com `assert.Equal(t, codes.XXX.Code(), bizErr.Code)` para BusinessException
11. Após criar, executar `go test ./internal/core/usecase/<entity>/... -v` e `go build ./...`
12. NUNCA criar testes de handler ou repository como "unitário" — focar na camada Core
13. NUNCA importar bibliotecas de infraestrutura nos testes
14. Após concluir, apresentar resumo: UseCases testados, regras de negócio cobertas, arquivos criados, resultado dos testes