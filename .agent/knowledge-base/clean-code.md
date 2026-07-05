# Clean Code — Princípios e Práticas

Baseado nos conceitos de **Clean Code** (Robert C. Martin) e **Refactoring** (Martin Fowler).

> **Nota**: Os exemplos de código usam nomes como `Task`, `CreateTaskUC`, `PostgresTaskRepository` como **ilustração**. Substitua pela entidade e nomes específicos do domínio da aplicação em que está trabalhando. Os princípios são universais e se aplicam a qualquer projeto.

---

## 1. Nomes Significativos

### Regras
- **Revelem intenção**: O nome deve responder *o que*, *por que* e *como*. Evite nomes genéricos como `data`, `info`, `temp`.
- **Evite desinformação**: Não chame lista de `groupList`. Não use nomes que mascarem a intenção.
- **Faça distinções meaningfuls**: Não sufixos redundantes (`NameString`, `CustomerObject`). Se dois nomes são diferentes, devem significar coisas diferentes.
- **Pronunciáveis e buscáveis**: `generationTimestamp` > `genTs`. Nomes buscáveis evitam magic numbers.
- **Evite codificações**: Sem notação húngara (`strName`, `iCount`). O tipo já está no código.
- **Nomes de classes/structs**: Substantivos (`Customer`, `Account`, `Order`). Evite `Manager`, `Processor`, `Data`, `Info` como sufixo genérico.
- **Nomes de métodos**: Verbos ou frases (`Save`, `FindByID`, `Delete`). Privados podem ser mais curtos se o contexto é claro.

### Exemplos ilustrativos
| Ruim | Bom | Por quê |
|---|---|---|
| `Get<Entity>ById(id)` | `FindByID(id)` | Repositories *encontram*, não *obtêm* |
| `<Entity>Data` | `<Entity>` | Sem sufixos redundantes |
| `procData()` | `create<Entity>()` | Revela a intenção |
| `m *Mock<Entity>Repository` | `mockRepo := new(Mock<Entity>Repository)` | Curto e claro no contexto do teste |
| `err1`, `err2` | `expectedErr`, `dbErr` | Nomes distinguem significado |

---

## 2. Funções

### Regras
- **Pequenas**: Funções devem ser menores que 20 linhas. Quanto menor, melhor.
- **Faça uma coisa**: Cada função faz exatamente uma coisa. Se pode extrair uma sub-função com nome descritivo, ela faz mais de uma.
- **Um nível de abstração por função**: Não misture conceitos de alto nível (`Create<Entity>`) com detalhes de baixo (`sql.Exec`).
- **Switch statements**: Usem polimorfismo em vez de `switch/if-else` encadeados para seleção de tipo. Extraia para factory ou strategy.
- **Poucos argumentos**: `2-3 argumentos` é o máximo ideal. Use structs/objects para agrupar parâmetros relacionados.
- **Sem efeitos colaterais**: Funções que prometem uma coisa mas fazem outra violam a confiança. `CheckPassword` não deve modificar estado.
- **Separação Comando/Consulta**: Funções que retornam valor **não** devem modificar estado. `Save()` retorna `error`. `ExistByID()` retorna `(bool, error)` — não salva nada.
- **Prefira erros tipados a códigos de erro**: Retorne exceções tipadas (`BusinessException`) em vez de códigos de status numéricos. Erros se propagam via call stack, códigos não.

### Exemplo ilustrativo
```go
// RUIM: Função faz 3 coisas (validar, criar, salvar)
func Process<Entity>(title, desc string, priority int) error { ... }

// BOM: Cada uma faz uma coisa
func New<Entity>(title, desc string, priority int) (*Entity, error) { ... }  // valida e cria
func (uc *Create<Entity>UC) create<Entity>(ctx, cmd) (Result, error) { ... }  // orquestra
func (r *Postgres<Entity>Repository) Save(ctx, entity) error { ... }        // persiste
```

---

## 3. Comentários

### Regras
- **Comentários não compensam código ruim**: Em vez de comentar código confuso, reescreva-o.
- **Explique a intenção, não a mecânica**: O *porquê*, não o *o quê*.
- **Comentários bons**: Legalidade, intenção, esclarecimento, consequência (TODO), amplificação.
- **Comentários ruins**: Mumbling, redundantes, mandados, journal, fechar-chave, atribuições, HTML, informação demais, posição, código comentado.

### Exemplo ilustrativo
```go
// RUIM: Repete o que o código já diz
// Save persists the entity in the database
func (r *PostgresRepository) Save(ctx context.Context, entity *Entity) error { ... }

// BOM: Explica o porquê
// Soft delete preserves audit trail — physical deletion violates compliance requirements
func (r *PostgresRepository) Delete(ctx context.Context, id string) error { ... }

// RUIM: Código comentado
// result, err := oldRepository.Find(id)
// if err != nil { return err }

// BOM: Nomes expressivos eliminam a necessidade de comentários
func (uc *Create<Entity>UC) Execute(ctx context.Context, entity *Entity) (struct{}, error) {
    return uc.Call(ctx, entity, uc.createEntity)
}
```

---

## 4. Formatação

### Regras
- **Formatação vertical**: Arquivos de cima a baixo devem contar uma história. Alto nível → baixo nível. Funções chamadas devem estar abaixo das que chamam.
- **Distância vertical**: Funções relacionadas devem estar próximas. Variáveis locais declaradas logo antes do uso.
- **Formatação horizontal**: Linhas não devem ultrapassar 100-120 caracteres. Alinhamento não enfatiza estrutura.
- **Agrupamento conceitual**: Separe conceitos com linhas em branco. Agrupe lógica relacionada.
- **Consistência**: Siga o estilo do time. Um estilo consistente vence qualquer estilo "perfeito" 100% das vezes.

### No contexto do projeto
- Imports agrupados: stdlib, terceiros, interno
- Constantes e tipos antes de funções
- Structs com campos agrupados logicamente
- Métodos de interface antes de implementação
- Factory functions: `New`, `NewWithoutAudit`, `NewInit` em sequência

---

## 5. Tratamento de Erros

### Regras (Martin Fowler — Refactoring)
- **Replace Error Code with Exception**: Em vez de retornar códigos mágicos, retornem erros tipados.
- **Replace Conditional with Polymorphism**: Diferentes comportamentos para diferentes tipos → use interfaces/polimorfismo em vez de `if-else` encadeado.
- **Encapsule condicionais**: `if (entity.Cancomplete())` > `if (entity.status == COMPLETED)`. O chamador não precisa conhecer invariantes internas.

### Exemplo ilustrativo
```go
// RUIM: Chamador conhece invariantes internas
if entity.status == "COMPLETED" {
    return BusinessException(EntityAlreadyDone)
}

// BOM: Objeto protege seus invariantes
func (e *Entity) Complete() error {
    if e.status == Completed {
        return BusinessException(EntityAlreadyDone)
    }
    e.status = Completed
    return nil
}
```

- **Erros de negócio** (4xx): `BusinessException` — o cliente deve saber o que aconteceu.
- **Erros técnicos** (5xx): `UnexpectedException` — o cliente recebe mensagem genérica, detalhes ficam nos logs.
- **Erros são transparentes na borda**: Handlers retornam `nil, err`. O `HandlerException` classifica e formata.
- **Enriquecimento de contexto**: No repositório, erros são envolvidos com `oops` para preservar stack trace.

---

## 6. Objetos e Estruturas de Dados

### Regras
- **Ocultação de dados**: Objetos ocultam dados atrás de abstrações. Estruturas de dados expõem dados sem comportamento. Não misture os dois.
- **Lei de Demeter**: Um módulo não deve saber sobre as entranhas dos objetos que manipula. `a.GetB().GetC().DoD()` viola Demeter — fale com o objeto direto.
- **Data Transfer Objects (DTOs)**: Estruturas de dados puras sem comportamento, usadas para transferência entre camadas. Sem lógica.

### No contexto do projeto
| Tipo | Ocultação | Comportamento | Uso |
|---|---|---|---|
| **Models** (Domain) | Campos privados + getters | Validadores, transições de estado | Core — regras de negócio |
| **Commands/Queries** (App) | Campos públicos | Nenhum | Transporte Handler → UseCase |
| **Entities** (DB) | Campos públicos | Nenhum | Mapeamento de colunas DB |
| **Request DTOs** (HTTP) | Campos públicos | Nenhum (só tags Huma) | Validação de input HTTP |
| **Response DTOs** (HTTP) | Campos públicos | Nenhum (só tags JSON) | Serialização de output HTTP |

---

## 7. Limites (Boundaries)

### Regras
- **Código de borda é diferente**: Código que interage com terceiros (DB, APIs) precisa ser isolado. Adapters são a fronteira.
- **Não envolva o que não precisa**: Envolva apenas o que é volátil ou pode mudar.
- **Clean boundaries**: Interfaces bem definidas entre camadas permitem trocar implementações sem tocar no domínio.

### Padrões de Refactoring de Fronteira (Martin Fowler)
- **Introduce Parameter Object**: Parâmetros relacionados devem ser um objeto.
- **Replace Subclass with Fields**: Se subclasses só diferem em dados, use configuração em vez de hierarquia.
- **Change Reference to Value**: Se o objeto é imutável e comparável por valor, trate como valor.

### No contexto do projeto
```
Core ← (port) → Adapter.Out
Core ← (IUsecase) → App (UseCases)
App ← (Command/Query) → Adapter.In
```
- Core define os contratos (Ports/IUsecase). App e Adapters implementam.
- App recebe Commands/Queries puros (sem tags de infra) e cria modelos de domínio.
- Trocar PostgreSQL por MongoDB = trocar `repository_impl`. Core e App permanecem intocados.
- Trocar Chi/Huma por outro router = trocar `http/server` + handlers. Core e App permanecem intocados.

---

## 8. Princípios SOLID

### Single Responsibility Principle (SRP)
Cada módulo/classe tem **um motivo para mudar**.

- `models/<entity>.go` muda quando regras de negócio mudam.
- `repository/<entity>_repository_postgres_impl.go` muda quando esquema DB muda.
- handlers mudam quando API muda. Nunca o contrário.

### Open/Closed Principle (OCP)
Aberto para extensão, fechado para modificação. Adicionar uma nova entidade = criar novos arquivos, não modificar os existentes.

- Value Groups do FX: novos handlers se registram sem alterar o servidor.
- `fx.As` permite trocar implementação sem mudar quem consome.

### Liskov Substitution Principle (LSP)
Subtipos devem ser substituíveis por seus tipos base. Qualquer implementação de `<Entity>Repository` (Postgres, InMemory, Mock) deve funcionar com qualquer `UseCase`.

### Interface Segregation Principle (ISP)
Clientes não devem depender de interfaces que não usam. Interfaces de Port são enxutas e específicas por domínio.

### Dependency Inversion Principle (DIP)
Módulos de alto nível não dependem de módulos de baixo nível. Ambos dependem de abstrações.

- UseCases dependem de `ports.<Entity>Repository` (interface), não da implementação concreta.
- `fx.As` garante que o contêiner injeta a interface, não a struct.

---

## 9. Classes e Sistema de Tipos

### Regras
- **Sistema de tipos como documentação**: Tipos nomeados (enums, BusinessException) documentam melhor que comentários.
- **Evite tipos primitivos para conceitos de domínio**: `type Priority int` > `int`. `type EntityID string` > `string`.
- **Prefira enums a strings mágicas**: Use tipos definidos com constantes, não strings soltas.

### Exemplo ilustrativo
```go
// RUIM: String solta não documenta intenção
func Create(title string, status string) {}

// BOM: Tipos nomeados documentam e constraint
func Create(title string, status enums.Status) {}
```

---

## 10. Testes

### Regras (Martin Fowler — Refactoring)
- **Testes como proteção**: Antes de refatorar, tenha testes sólidos. Testes te dão confiança para mudar.
- **Self-Verifying**: Testes devem ser determinísticos. Mesmo input, mesmo resultado. Sem depender de estado externo.
- **Um assertion lógico por teste**: Se um teste verifica 3 coisas, deveria ser 3 testes. Mas mock assertions são exceção aceitável.
- **FIRST**: Fast (rápidos), Independent (independentes), Repeatable (repetíveis), Self-validating (auto-verificáveis), Timely (escritos no momento certo).

### Padrões de Refactoring aplicados a Testes
- **Extract Method**: Testes longos com setup repetido → extrair helpers.
- **Introduce Parameter Object**: Muitos parâmetros no Arrange → criar struct de teste.
- **Replace Magic Literal**: Hardcode `1` → constantes nomeadas `ValidPriority`.

### No contexto do projeto
```
Arrange → mock com testify/mock
Act → chamar UseCase.Execute()
Assert → testify/assert + mock.AssertExpectations()
```
- UseCases são 100% testáveis sem infraestrutura
- Mocks substituem repositories concretos
- Padrão AAA aplicado consistentemente

---

## 11. Cheiros de Código (Martin Fowler — Refactoring)

| Cheiro | Refactoring |
|---|---|
| Função longa | Extract Method, Replace Conditional with Polymorphism |
| Classe grande | Extract Class, Introduce Parameter Object |
| Lista longa de parâmetros | Introduce Parameter Object, Preserve Whole Object |
| Mudanças divergentes (SRP) | Extract Class |
| Shotgun Surgery (uma mudança, muitos arquivos) | Move Method, Move Field, Inline Class |
| Feature Envy (método usa mais dados de outra classe) | Move Method |
| Data Clumps (parâmetros que andam juntos) | Introduce Parameter Object |
| Switch Statements | Replace Conditional with Polymorphism |
| Speculative Generality (código "caso um dia precise") | Remove Dead Code, Collapse Hierarchy |
| Temporary Field (campos usados só às vezes) | Extract Class, Introduce Null Object |
| Message Chains (`a.getB().getC()`) | Hide Delegate, Extract Method |
| Middle Man (classe que só delega) | Remove Middle Man, Inline Method |
| Refused Bequest (subclasse que ignora herança) | Replace Inheritance with Delegation |
| Comments (comentários que explicam código ruim) | Extract Method, Rename Method, Introduce Assertion |

---

## 12. Refactoring Core (Martin Fowler)

### Quando Refatorar
- **Regra dos Três**: Primeira vez, faça. Segunda vez, tolera. Terceira vez, refatore.
- **Preparar**: Antes de adicionar funcionalidade, refatore o código para que seja fácil adicionar.
- **Compreender**: Ao entender código confuso, refatore para deixar seu entendimento claro.
- **Revisão**: Durante code review, sugira refatorações que melhorem legibilidade.

### Catálogo Essencial de Refactorings

| Refactoring | Quando Usar |
|---|---|
| **Extract Method** | Função faz mais de uma coisa, ou trecho pode ser nomeado |
| **Inline Method** | Corpo do método é tão óbvio quanto o nome |
| **Extract Variable** | Expressão complexa precisa de nome explicativo |
| **Rename Variable/Method** | Nome não revela intenção |
| **Move Method** | Método usa mais dados de outra classe que da própria |
| **Replace Conditional with Polymorphism** | Switch/if-else baseado em tipo |
| **Replace Error Code with Exception** | Retorno de códigos de erro em vez de propagação |
| **Introduce Parameter Object** | Parâmetros que sempre aparecem juntos |
| **Decompose Conditional** | Condicional complexa com lógica de branch em métodos nomeados |
| **Preserve Whole Object** | Passando vários campos de um objeto → passe o objeto |

### Princípios de Refactoring Seguro
1. **Testes primeiro**: Nunca refatore sem testes que verifiquem o comportamento atual.
2. **Passos pequenos**: Uma mudança por vez. Compile e teste entre passos.
3. **Não adicione funcionalidade**: Refatorar muda estrutura, não comportamento.
4. **Use controle de versão**: Commit antes de cada refactoring para rollback fácil.