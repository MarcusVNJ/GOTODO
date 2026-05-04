---
description: "Analisa o projeto e gera toda a base de conhecimento em .agent/knowledge-base/"
mode: primary
temperature: 0.1
permission:
  edit: allow
  bash: ask
  task: allow
color: "#10b981"
---

{file:../.agent/agents/knowledge-base-creator.md}

INSTRUÇÕES DE EXECUÇÃO:

1. Ao ser invocado, siga o fluxo de trabalho definido no Knowledge Base Creator acima
2. Analise o codebase no diretório de trabalho atual — leia os fontes relevantes
3. ANTES DE TUDO: Detecte se o projeto é novo (Template Mode) ou existente (Modo Normal) seguindo os critérios da Etapa 0
4. PARA O SCHEMA DO BANCO, siga esta ordem de precedência:
   - **Primeiro (fonte primária)**: Se o MCP estiver disponível, use as ferramentas MCP para obter o estado real do banco — dados do MCP têm precedência sobre tudo
   - **Segundo (fallback)**: Se o MCP NÃO estiver disponível, ler todos os arquivos `*.up.sql` em `migrations/` para extrair o schema
   - **Terceiro (fallback secundário)**: Se `migrations/` também não existe ou está vazia, inferir o schema de `internal/adapters/out/infrastructure/entity/`, mappers e query builders
5. Se o MCP NÃO estiver disponível mas as migrations existem, gere `database.md` baseado nas migrations, marcando com "[Validação pendente: confirmar com MCP quando disponível]"
6. Se NEM MCP NEM migrations estão disponíveis, gere baseado no código-fonte marcando com "[Pendente: verificar com MCP quando disponível]"
7. EM TEMPLATE MODE (projeto novo): gere os arquivos com marcadores [Pendente/Preencher/Atualizar] nos locais que precisam de dados reais. Os arquivos semi-genéricos (architecture, conventions, di-patterns, error-handling, testing, clean-code) devem vir completos com o conteúdo padrão da stack Hexagonal.
8. EM TEMPLATE MODE (projeto novo): CRIE a estrutura de pastas base do projeto seguindo a Arquitetura Hexagonal. Criar todas as pastas listadas na seção "Criação da Estrutura de Pastas" que não existirem. Adicionar `.gitkeep` em pastas vazias que precisam ser rastreadas.
9. Gere os 12 arquivos padronizados em `.agent/knowledge-base/` com nomes EXATOS:
   - project-overview.md
   - architecture.md
   - tech-stack.md
   - codebase-map.md
   - domain-model.md
   - conventions-and-standards.md
   - di-patterns.md
   - error-handling.md
   - testing.md
   - database.md
   - clean-code.md
   - api-resources.md
10. APÓS gerar os 12 arquivos, gere ou atualize também o `AGENTS.md` na raiz do projeto com Project Overview, Commands, Environment, Architecture, Key Patterns, Testing, Migrations, Custom Agents e Common Pitfalls específicos do projeto
11. Verifique se `.agent/workflows/pipeline.md` existe. Se NÃO existe, crie com o conteúdo padrão do pipeline (planner → software-engineer → test-engineer → doc-updater)
12. Após concluir, apresente um resumo ao usuário listando: arquivos criados, entidades encontradas (ou "Nenhuma — projeto novo"), endpoints catalogados (ou "Nenhum — projeto novo"), gaps identificados, modo detectado (Projeto Existente / Template), estrutura de pastas (se criada), workflow pipeline (se criado), e status do MCP
13. Escrita em Português do Brasil, exceto termos técnicos e código-fonte