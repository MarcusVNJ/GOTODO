Guia de Integração para Desenvolvedores - Projeto GOTODO

Bem-vindo ao projeto GOTODO! Este é o nosso sistema de gerenciamento de tarefas (estilo Kanban). Para garantir que o código permaneça escalável, limpo e de fácil manutenção, nós adotamos uma série de padrões avançados de engenharia de software e recursos específicos da linguagem Go.

Este documento foi criado para acelerar o seu entendimento sobre como o Go funciona por debaixo dos panos, como estruturamos nossa Arquitetura Hexagonal e como lidamos de forma elegante com os fluxos HTTP, banco de dados e erros.
1. Fundamentos Computacionais em Go: Memória e Desempenho

Diferente de linguagens interpretadas ou que rodam em máquinas virtuais densas, Go é compilado para código de máquina e nos dá grande controle sobre a memória. Para escrever um código eficiente aqui, você precisa entender como o Go gerencia a alocação de dados.
1.1. Stack vs. Heap

Toda vez que uma variável é criada, o Go precisa decidir onde guardá-la na memória RAM do servidor. Ele tem dois lugares principais: a Stack (Pilha) e o Heap (Monte).

    A Stack (A Mochila): É uma região de memória incrivelmente rápida que opera no formato LIFO (Last In, First Out). É como sua mochila para uma viagem de um dia: rápida para acessar, mas com tamanho limitado e esvaziada automaticamente assim que a função que você estava executando termina.

    O Heap (A Mala de Porão): É uma região de memória maior, usada para armazenamento dinâmico. É como uma mala grande que fica guardada mesmo depois que a "viagem" da função acaba. É mais lento e exige que o Garbage Collector (GC) do Go trabalhe periodicamente para limpar o que não está mais sendo usado.

1.2. Escape Analysis (Análise de Fuga)

A regra de ouro do Go é: o tamanho da variável não importa tanto quanto o tempo de vida dela. Durante a compilação, o Go rastreia o caminho de cada variável (Escape Analysis). Se o compilador provar que a variável só será usada dentro daquela função, ela fica na Stack. Mas se a variável "escapar" (ex: você retornar um ponteiro para ela, passá-la para uma goroutine ou para uma interface vazia any), o Go a joga para o Heap.
1.3. Ponteiros vs. Valores (Quando usar o quê?)

É comum que novos desenvolvedores Go passem tudo como ponteiros (*models.Task) achando que isso economiza memória. Isso é um mito.

    Passe por Valor: Para structs pequenas ou médias (como DTOs), o custo de copiar os dados é quase zero e mantém a variável na Stack, aliviando o trabalho do Garbage Collector. Nós preferimos a passagem por valor por padrão.

    Passe por Ponteiro: Use ponteiros apenas quando precisar modificar/mutar a estrutura original ou quando a struct for excepcionalmente gigantesca.

2. A Arquitetura Hexagonal (Ports and Adapters)

O projeto GOTODO utiliza a Arquitetura Hexagonal. O objetivo é garantir que as nossas regras de negócio nunca dependam de detalhes técnicos como bancos de dados (Postgres) ou frameworks web (Chi).
2.1. O Core (Domínio e UseCases)

É o coração da nossa aplicação, localizado em internal/core.

    Regra Absoluta: Nenhuma interface, struct ou função dentro do core pode importar bibliotecas de infraestrutura (como net/http, pgx ou json). O domínio é cego para o mundo externo.

    Possuímos pacotes como enums (para os status Kanban), models (entidades ricas com regras de negócio) e base (interfaces genéricas de caso de uso).

2.2. Ports (As Portas)

Para que o core consiga conversar com o banco de dados sem conhecê-lo, ele define Interfaces em internal/core/ports.

    O domínio apenas diz: "Eu preciso de alguém que implemente a função Save(task *models.Task)". Ele não se importa se quem vai salvar é o Postgres ou a memória.

2.3. Adapters (Os Adaptadores)

Ficam em internal/adapters. Eles são os tradutores entre o mundo externo e o nosso domínio.

    Inbound (in/http): Nossos Handlers utilizam o ecossistema Huma. Não lidamos com decodificação de JSON manual (`json.Unmarshal`). Definimos structs `Input` e `Output`, e o Huma valida magicamente via tags (ex: `required:"true"`). Se for válido, o Handler converte o DTO para o modelo puro e chama o UseCase.

    Outbound (out/infrastructure): Nossa implementação do Postgres. Aqui usamos mappers para converter a entidade de domínio em uma TaskEntity e o query_builder para gerar o SQL.

3. Padrões Específicos do Projeto GOTODO
3.1. Injeção de Dependência e Raiz de Composição (uber-go/fx)

Não instanciamos repositórios ou casos de uso manualmente dentro dos Handlers (ex: não fazemos `new(...)` espalhado pelo código). Usamos a biblioteca `uber-go/fx` para o gerenciamento de Injeção de Dependência (DI) e controle do ciclo de vida da aplicação.

Se você olhar pelo código, verá pacotes exportando uma variável `var Module = fx.Module(...)`. Entenda como o FX funciona:

**a) Avaliação Preguiçosa vs Imediata (`fx.Provide` vs `fx.Invoke`)**
- `fx.Provide(Construtor)`: Ensina o FX a criar um objeto (ex: `LoadConfig` ou `InitDB`). O FX é **preguiçoso (lazy)**: ele registra a "receita" de como criar a conexão com o banco, mas só a executará se alguém na aplicação pedir por um banco de dados.
- `fx.Invoke(Função)`: É o gatilho **imediato (eager)**. Diz para o FX: "Execute isso obrigatoriamente no momento da inicialização". É usado para iniciar o Logger ou subir o Servidor HTTP (`StartHTTPServer`). O FX lerá a assinatura dessa função, verá as dependências necessárias e executará a cadeia de construtores do `fx.Provide` para resolvê-la.

**b) Ocultando Código Concreto sob Interfaces (`fx.Annotate` e `fx.As`)**
Na Arquitetura Hexagonal, o Core exige interfaces (ex: `ports.TaskRepository`), mas o adaptador fornece uma struct concreta (`*PostgresTaskRepository`). Para ensinar o FX a disfarçar o construtor, usamos:
```go
fx.Annotate(NewPostgresTaskRepository, fx.As(new(repository.TaskRepository)))
```
Isso garante o princípio de Inversão de Dependência: o Core não acopla com o banco, acopla com a Interface.

**c) Colisão de Tipos em Generics e Named Instances (`fx.ResultTags` e `fx.ParamTags`)**
Um problema complexo ocorre quando duas funções fornecem o mesmo tipo de retorno. Exemplo: a criação e a atualização de uma tarefa retornam a mesmíssima interface genérica: `usecase.IUsecase[*models.Task, struct{}]`. Se nada for feito, o FX quebra na inicialização avisando sobre *ambiguidade* (não sabe qual injetar).
Para resolver isso, usamos "Instâncias Nomeadas":
1. O provedor cola uma etiqueta: `fx.ResultTags("name:\"createTaskUC\"")`.
2. O consumidor exige a etiqueta: `fx.ParamTags("name:\"createTaskUC\"")`.

**d) Registro Dinâmico de Rotas (Value Groups)**
Não registramos rotas uma por uma no arquivo do servidor. Os Handlers de rota usam a função `router.AsRoute(...)`, que embala o resource com a anotação `fx.ResultTags("group:\"routes\"")`. 
Essa instrução não dá um nome para o objeto, mas sim o joga dentro de uma "caixa" (array/slice) chamada `"routes"`. 
Lá no Módulo de Servidor (`server.go`), o FX injeta todo o grupo em `Routes []router.RouteRegister` e o servidor simplesmente itera o *array* registrando todos na Huma.API.

**e) Pureza do Core**
Para manter o `internal/core` isolado, a orquestração do FX para os UseCases não fica na pasta Core. Ela é "terceirizada" e mantida na borda do sistema em `cmd/api/di/usecases.go`.

Nossos Handlers de rota são super limpos, retornando apenas structs tipadas e um `error` (ex: `func(ctx context.Context, input *Input) (*Output, error)`).
Nós criamos um wrapper chamado `middlewares.HandlerException` que envolve todas as rotas no momento do registro. O seu único trabalho na Resource é: deu erro? Faça `return nil, err`!
O wrapper captura esse erro, converte para o padrão RFC 7807 problem details (Status 400, 404, etc.) e caso seja um erro crítico não mapeado (500), ele responde com uma mensagem sanitizada e cria uma *goroutine* separada para fazer o log completo (com a stack trace do `oops`) via `slog` no console, sem travar a requisição do usuário.
3.3. Configuração Orientada a Ambiente

Lemos configurações via variáveis de ambiente usando o pacote internal/config, populando uma struct fortemente tipada (AppConfig). Se faltar uma variável obrigatória, o sistema aplica o padrão Fail-Fast e não inicia.
3.4. Versionamento de Banco de Dados (Migrations)

Nunca criamos tabelas manualmente. Utilizamos a biblioteca golang-migrate/migrate.

    Onde ficam: Todos os scripts SQL ficam na pasta migrations/ na raiz do projeto, fora do internal (já que o Core não deve saber sobre SQL).

    Up e Down: Para cada alteração, há um arquivo .up.sql (aplica a mudança) e um .down.sql (reverte a mudança). Isso garante que o esquema do banco evolua de forma segura e rastreável junto com a API.

3.5. Documentação Automática (Swagger/OpenAPI)

O projeto se auto-documenta dinamicamente via Huma. Se você quiser ver quais rotas estão disponíveis e testá-las sem precisar abrir o Postman:
Acesse: `http://localhost:8080/docs` (garanta que a env `ENABLE_DOCS=true` está configurada no seu `.env`).

4. Escrevendo Testes de Unidade

A nossa regra de ouro é: **A lógica de negócio deve ser 100% testável de forma rápida e sem depender de infraestrutura externa (banco de dados).**

Utilizamos a biblioteca `testify` para facilitar nossos testes e aplicamos o padrão **AAA (Arrange, Act, Assert)**:

1. **Arrange (Preparação):** Usamos o pacote `testify/mock` para criar "Repositórios Falsos" que obedecem o que o teste mandar. Injetamos esse Mock no UseCase.
2. **Act (Ação):** Acionamos a função principal do UseCase (ex: `Execute(...)`).
3. **Assert (Verificação):** Usamos o pacote `testify/assert` para garantir que a resposta está correta e `mock.AssertExpectations` para garantir que o fluxo chamou as funções de banco apenas quando deveria.

**Exemplo Prático (Simulando um erro no Banco de Dados):**
```go
func Test_Execute_RepositoryError(t *testing.T) {
	// ARRANGE
	mockRepo := new(MockTaskRepository)
	uc := usecase.NewCreateTaskUC(mockRepo)
	
	task := models.NewTask("Título", "Desc", 1)
	expectedErr := errors.New("banco fora do ar")
	
	// Dita a regra: "Quando tentarem Salvar essa task, retorne erro!"
	mockRepo.On("Save", mock.Anything, task).Return(expectedErr)

	// ACT
	_, err := uc.Execute(context.Background(), task)

	// ASSERT
	assert.ErrorIs(t, err, expectedErr) // Garantimos que falhou pelo motivo certo
	mockRepo.AssertExpectations(t)      // Garantimos que a função Save realmente foi chamada
}
```

Sempre cubra o **Caminho Feliz**, os caminhos **Alternativos**, e os de **Lançamento de Exceções** (ex: `BusinessException`).

5. Primeiros Passos

    Abra a pasta internal/core/models e analise as regras de negócio do Kanban.

    Veja a configuração dos módulos (`var Module = fx.Module(...)`) nas camadas de adapters.

    Observe como o `cmd/api/di` foi criado exclusivamente para não poluir o domínio com a injeção do FX.

    Observe o `internal/adapters/in/http/handlers` para ver como o fluxo HTTP é limpo graças ao nosso ExceptionHandler.