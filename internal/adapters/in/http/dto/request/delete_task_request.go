package request

type DeleteTaskRequest struct {
	ID string `path:"id" minLength:"1" maxLength:"36" description:"ID da task a ser excluída" example:"123e4567-e89b-12d3-a456-426614174000"`
}
