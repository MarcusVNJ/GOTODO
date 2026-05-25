package request

type CreateTaskPayload struct {
	Title       string `json:"title" minLength:"1" maxLength:"150" description:"Título da tarefa"`
	Description string `json:"description" maxLength:"500" description:"Uma breve descrição da tarefa"`
	Priority    int    `json:"priority" minimum:"1" maximum:"5" description:"Prioridade de 1 a 5"`
}

type CreateTaskRequest struct {
	Body CreateTaskPayload
}
