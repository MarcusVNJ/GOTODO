Arquitetura do Projeto GOTODO

1. Visão Geral

O GOTODO é uma API RESTful para o gerenciamento de tarefas (estilo Kanban) desenvolvida em Go. Para garantir que o
sistema seja escalável, altamente testável e resistente à deterioração do código ao longo do tempo, adotamos a
Arquitetura Hexagonal (também conhecida como Ports and Adapters).

O princípio fundamental desta arquitetura é que as regras de negócio devem ser o centro da aplicação. As dependências
devem sempre apontar "de fora para dentro". O domínio (Core) não sabe que existe um banco de dados PostgreSQL ou que
está respondendo a requisições HTTP; ele se comunica com o mundo externo puramente através de interfaces (Ports).

2. Estrutura de Diretórios e Fronteiras

A separação de responsabilidades reflete rigorosamente nossa árvore de diretórios:

    cmd/api/: O ponto de entrada da aplicação e nossa Raiz de Composição. O main.go e o subpacote router detêm a responsabilidade exclusiva de inicializar dependências e acoplar os adaptadores ao Core.

    internal/config/: Lida com a leitura tipada de variáveis de ambiente (12-Factor App).

    internal/core/ (O Hexágono): O coração do sistema. Contém a lógica pura.

        base/: Estruturas genéricas de base (ex: abstrações de UseCases).

        enums/: Definições estáticas do negócio (ex: status das tarefas).

        models/: Entidades de negócio ricas (ex: Task), contendo invariantes e regras de transição de status do Kanban.

        ports/: Interfaces que definem os contratos de entrada (UseCases) e saída (Repositories).

        usecase/: Orquestradores da lógica da aplicação.

        exceptions/: Erros de negócio tipados (ex: BusinessException), totalmente agnósticos ao protocolo HTTP.

    internal/adapters/ (A Infraestrutura): Tradutores entre o mundo externo e o Core.

        in/http/: Adaptadores primários (Handlers/Resources, DTOs de request/response e Middlewares).

        out/infrastructure/: Adaptadores secundários. Contém entity (modelos exclusivos de DB), mappers (tradutores de Domain <-> Entity), query_builder (geradores de SQL) e implementações do repositório PostgreSQL.

    migrations/: Scripts SQL puros (Up/Down) processados pelo golang-migrate, mantidos na raiz para garantir que a infraestrutura de banco de dados não polua o código interno.

3. Padrões de Engenharia e Design
   3.1. Tratamento Centralizado de Erros (Exception Handler)

Implementamos um adaptador de manipulador customizado (ExceptionHandler). Nossos Resources implementam a assinatura
func(w, r) error. O middleware intercepta o erro retornado:

    Erros do Cliente (4xx): Falhas de conversão de JSON (ex: string enviada em campo int) ou erros de negócio gerados pelo Core são interceptados, convertidos para o Status HTTP correspondente (400, 404, 409) e retornados como JSON estruturado. Não geram logs de nível Error, evitando falsos-positivos na observabilidade.

    Erros de Servidor (5xx): Falhas técnicas são logadas pelo slog com todo o contexto (stack trace, path, método) e o cliente recebe apenas uma mensagem sanitizada genérica.

3.2. Segregação de Modelos (DTO -> Domain -> Entity)

Para proteger o encapsulamento, utilizamos modelos diferentes para cada fase da jornada dos dados:

    DTO (adapters/in/http/dto): Usado unicamente para parsing e validação de payloads HTTP JSON.

    Domain (core/models): Usado para a lógica de negócio pura. Não possui tags de serialização (json ou db).

    Entity (adapters/out/infrastructure/entity): Usado pelo repositório e pelo query_builder para espelhar as tabelas do banco de dados (contendo tags db). A conversão ocorre no pacote mappers.

4. O Fluxo de uma Requisição

Ao executar um fluxo (ex: Criar Tarefa), a requisição atravessa as camadas da seguinte forma:

    Entrada: A requisição atinge o framework go-chi, passando por middlewares globais e caindo no nosso AppRouter customizado.

    Transporte (adapters/in): O Resource executa o decode do JSON para o DTO. Falhas de parse retornam imediatamente. Ocorrendo sucesso, o DTO é mapeado para o modelo do Core, e o UseCase é acionado.

    Lógica (core/usecase): O UseCase orquestra as regras. Se todas as validações do modelo de negócio passarem, ele invoca a interface do Repositório (Porta de saída).

    Persistência (adapters/out): A implementação concreta (task_repository_postgres_impl.go) recebe o modelo de domínio, usa o Mapper para convertê-lo em TaskEntity, aciona o TaskQueryBuilder para montar o SQL e executa no pool do Postgres.

    Retorno: O sucesso (ou falha) borbulha de volta pela pilha até o Resource ou o ExceptionHandler formartar a resposta HTTP correspondente.

5. 🛡️ Tratamento Centralizado de Erros no GOTODO

O projeto GOTODO utiliza o padrão de Handler Adapter atrelado a um roteador customizado para garantir um tratamento de
exceções centralizado, uniforme e seguro em todas as rotas web da aplicação. Essa arquitetura elimina a duplicação de
código, garante que as respostas da API sejam sempre consistentes e mitiga totalmente o risco de vazamento de dados
sensíveis da infraestrutura para o usuário final.

A mecânica rigorosa de tratamento de erros funciona da seguinte forma:

    Roteamento Seguro por Design (AppRouter): Para garantir que nenhuma rota fuja da padronização de erros, implementamos um wrapper (Decorator) sobre o roteador do framework chi. Nós sobrescrevemos os verbos HTTP principais (Get, Post, Delete, etc.) para aceitarem exclusivamente a nossa assinatura customizada ResourceHandler: func(w http.ResponseWriter, r *http.Request) error. Isso envelopa automaticamente todas as requisições com o nosso middleware interceptador, tornando impossível que um desenvolvedor esqueça de aplicar o tratamento de erros em uma nova rota.

    Assinatura Simplificada nos Handlers: Graças a esse roteador customizado, os manipuladores de rota (Resources) não precisam se preocupar em formatar a resposta HTTP em caso de erro. Eles simplesmente retornam a interface error e delegam a responsabilidade.

    O Middleware Interceptador (ExceptionHandler): Este adaptador central atua envolvendo todas as rotas da aplicação. Sua única responsabilidade é capturar qualquer erro devolvido pelos Handlers, analisar a sua natureza (usando errors.As) e dar o destino correto a ele, dividindo-os em duas categorias estritas:

1. Erros do Cliente e Regras de Negócio (Status 4xx)

   Este fluxo trata falhas de requisição, como problemas no parsing do JSON (por exemplo, enviar uma string em um campo
   numérico, gerando um json.UnmarshalTypeError), ou exceções geradas pelo Core (nossa BusinessException / códigos como
   CodeInvalidData).

   O middleware intercepta essas falhas e as traduz para um JSON padronizado e estruturado, acompanhado do Status HTTP
   adequado (como 400 Bad Request, 404 Not Found ou 409 Conflict).

   Observabilidade Limpa: Para garantir que as nossas ferramentas de monitoramento permaneçam úteis, esses erros (que
   são desvios de fluxo normais e esperados da regra de negócio) não geram logs de nível Error, evitando alertas de
   falsos-positivos para a equipe.

2. Erros Técnicos e Falhas de Infraestrutura (Status 5xx)

   Quando ocorrem falhas não previstas ou problemas estruturais (como falhas de conexão com o banco de dados Postgres
   encapsuladas pela nossa UnexpectedException), o sistema prioriza a segurança (Fail-Safe).

   O cliente recebe um erro genérico de Status 500, com uma mensagem sanitizada (ex: "Erro interno no servidor"),
   garantindo que nenhum detalhe técnico interno ou stack trace seja vazado.

   Simultaneamente, o sistema utiliza o pacote nativo log/slog para registrar a falha criticamente no servidor. Este log
   contém o contexto técnico completo e estruturado para a equipe de engenharia, incluindo a stack trace rica gerada
   pela biblioteca oops, o path da URL e o método HTTP utilizado.

   Malha Fina de Segurança: Se um erro cru (não mapeado como Business ou Unexpected) chegar ao handler, o middleware
   aciona um alerta crítico de que o fluxo de erro padrão foi violado.