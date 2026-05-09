// users/adapters/http/handler.go
type UserHandler struct {
	registerUser usecases.RegisterUser // interface, não o struct concreto
	login        usecases.Login
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"bad request"}`, 400)
		return
	}

	out, err := h.registerUser.Execute(r.Context(), usecases.RegisterUserInput{
		Name:     body.Name,
		Email:    body.Email,
		Password: body.Password,
	})
	if errors.Is(err, domain.ErrEmailAlreadyInUse) {
		http.Error(w, `{"error":"email already in use"}`, 409)
		return
	}
	if err != nil {
		http.Error(w, `{"error":"internal server error"}`, 500)
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]string{"user_id": out.UserID})
}