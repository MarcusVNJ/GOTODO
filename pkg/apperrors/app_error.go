package apperrors

type AppError struct {
    // Err é o erro original (causa raiz). Não é exportado para JSON.
    Err error `json:"-"`

    // HTTPStatus define o código de status HTTP sugerido (ex: 404, 500).
    HTTPStatus int `json:"-"`

    // Code é um identificador estável para o erro (ex: "RESOURCE_NOT_FOUND").
    // Útil para clientes programáticos (frontend/mobile) reagirem a erros específicos.
    Code string `json:"code"`

    // Message é a mensagem amigável para o usuário final.
    Message string `json:"message"`

    // RequestID permite correlação em sistemas distribuídos.
    RequestID string `json:"request_id,omitempty"`

    // Details fornece contexto adicional, útil para erros de validação.
    Details map[string]interface{} `json:"details,omitempty"`
}
