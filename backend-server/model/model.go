package model

// Status représente un enregistrement de status dans la base de données.
type Status struct {
	ID          int    `json:"id"`
	Designation string `json:"designation"`
	State       string `json:"state"`
}

// CrashResponse représente la réponse JSON envoyée lors de l'appel à la route /crash.
type CrashResponse struct {
	Message string `json:"message"`
}
