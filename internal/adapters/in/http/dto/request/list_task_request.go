package request

type ListTaskRequest struct {
	Status      string `query:"status" description:"Filtrar por status (PENDING, IN_PROCESS, COMPLETED, CANCELLED)"`
	MinPriority int    `query:"min_priority" description:"Filtrar por prioridade mínima (1-5)"`
}