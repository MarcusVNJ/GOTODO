package response

type MessagePayload struct {
	Message string `json:"message" example:"Operação realizada com sucesso"`
	Id      string `json:"id,omitempty" example:"123456789"`
}
