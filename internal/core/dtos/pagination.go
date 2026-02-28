package dtos

// PaginationRequest se usa en cualquier controlador que devuelva listas
type PaginationRequest struct {
	Limit  int `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset int `form:"offset" binding:"omitempty,min=0"`
}

// PaginatedResponse estandariza cómo la API devuelve listas en toda la notaría
type PaginatedResponse struct {
	Total int         `json:"total"`
	Data  interface{} `json:"data"` // Aquí puedes meter una lista de Usuarios, Clientes o Trabajos
}