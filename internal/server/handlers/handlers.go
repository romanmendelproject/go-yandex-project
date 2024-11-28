package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/romanmendelproject/go-yandex-project/internal/server/config"
	"github.com/romanmendelproject/go-yandex-project/internal/server/user"
	"github.com/romanmendelproject/go-yandex-project/internal/types"
	log "github.com/sirupsen/logrus"

	"github.com/romanmendelproject/go-yandex-project/internal/server/jwt"
)

type Storage interface {
	GetCred(ctx context.Context, name string, userID int) (*types.CredType, error)
	SetCred(ctx context.Context, value types.CredType, userID int) error
	GetText(ctx context.Context, name string) (*types.TextType, error)
	SetText(ctx context.Context, value types.TextType) error
	GetByte(ctx context.Context, name string) (*types.ByteType, error)
	SetByte(ctx context.Context, value types.ByteType) error
	GetCard(ctx context.Context, name string) (*types.CardType, error)
	SetCard(ctx context.Context, value types.CardType) error
	Ping(ctx context.Context) error
}

type RegisterReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func HandleBadRequest(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
}

func HandleStatusNotFound(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotFound)
}

func handleError(res http.ResponseWriter, err error, status int) {
	log.Error(err)
	http.Error(res, err.Error(), status)
}

func customError(res http.ResponseWriter, err string, status int) {
	log.Error(err)
	http.Error(res, err, status)
}

type ServiceHandlers struct {
	cfg     config.Config
	storage Storage
	token   *jwt.JWT
	user    *user.User
}

func NewHandlers(cfg config.Config, storage Storage, token *jwt.JWT, userData *user.User) *ServiceHandlers {
	return &ServiceHandlers{
		cfg:     cfg,
		storage: storage,
		token:   token,
		user:    userData,
	}
}

func (h *ServiceHandlers) GetCredValue(res http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")

	userID, err := h.getUserID(req)
	if err != nil {
		handleError(res, err, http.StatusUnauthorized)
		return
	}

	value, err := h.storage.GetCred(req.Context(), name, userID)
	if err != nil {
		customError(res, "there is no record named test", http.StatusNotFound)
		return
	}
	var requestData types.CredType

	requestData.Name = name
	requestData.Username = value.Username
	requestData.Password = value.Password
	requestData.Meta = value.Meta

	resp, err := json.Marshal(requestData)
	if err != nil {
		handleError(res, err, http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(resp)
}

func (h *ServiceHandlers) SetCredValue(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var request types.CredType

	userID, err := h.getUserID(req)
	if err != nil {
		handleError(res, err, http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		res.Write([]byte(err.Error()))
		handleError(res, err, http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	err = h.storage.SetCred(ctx, request, userID)
	if err != nil {
		res.Write([]byte(err.Error()))
		handleError(res, err, http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func (h *ServiceHandlers) GetTextValue(res http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")

	value, err := h.storage.GetText(req.Context(), name)
	if err != nil {
		handleError(res, err, http.StatusNotFound)
		return
	}
	var requestData types.TextType

	requestData.Name = name
	requestData.Data = value.Data
	requestData.Meta = value.Meta

	resp, err := json.Marshal(requestData)
	if err != nil {
		handleError(res, err, http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(resp)
}

func (h *ServiceHandlers) SetTextValue(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var request types.TextType
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		res.Write([]byte(err.Error()))
		handleError(res, err, http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	err := h.storage.SetText(ctx, request)
	if err != nil {
		res.Write([]byte(err.Error()))
		handleError(res, err, http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func (h *ServiceHandlers) GetByteValue(res http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")

	value, err := h.storage.GetByte(req.Context(), name)
	if err != nil {
		handleError(res, err, http.StatusNotFound)
		return
	}
	var requestData types.ByteType

	requestData.Name = name
	requestData.Data = value.Data
	requestData.Meta = value.Meta

	resp, err := json.Marshal(requestData)
	if err != nil {
		handleError(res, err, http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(resp)
}

func (h *ServiceHandlers) SetByteValue(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var request types.ByteType
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		res.Write([]byte(err.Error()))
		handleError(res, err, http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	err := h.storage.SetByte(ctx, request)
	if err != nil {
		res.Write([]byte(err.Error()))
		handleError(res, err, http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func (h *ServiceHandlers) GetCardValue(res http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")

	value, err := h.storage.GetCard(req.Context(), name)
	if err != nil {
		handleError(res, err, http.StatusNotFound)
		return
	}
	var requestData types.CardType

	requestData.Name = name
	requestData.Data = value.Data
	requestData.Meta = value.Meta

	resp, err := json.Marshal(requestData)
	if err != nil {
		handleError(res, err, http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(resp)
}

func (h *ServiceHandlers) SetCardValue(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var request types.CardType
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		res.Write([]byte(err.Error()))
		handleError(res, err, http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	err := h.storage.SetCard(ctx, request)
	if err != nil {
		res.Write([]byte(err.Error()))
		handleError(res, err, http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func (h *ServiceHandlers) RegisterUser(res http.ResponseWriter, req *http.Request) {
	var request RegisterReq

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		customError(res, "error while decoding body", http.StatusBadRequest)
		return
	}

	defer func() {
		if err := req.Body.Close(); err != nil {
			log.Error("error closing body", "error", err)
		}
	}()

	if !checkEmpty(request.Login, request.Password) {
		customError(res, "empty login or password", http.StatusBadRequest)
		return
	}

	token, err := h.user.RegisterUser(req.Context(), request.Login, request.Password)
	if err != nil {
		handleError(res, err, http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:  "Token",
		Value: token,
		Path:  "/",
	})
}

func (h *ServiceHandlers) LoginUser(res http.ResponseWriter, req *http.Request) {
	var request RegisterReq

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		customError(res, "error while decoding body", http.StatusBadRequest)
		return
	}
	defer func() {
		if err := req.Body.Close(); err != nil {
			log.Error("error closing body", "error", err)
		}
	}()

	if !checkEmpty(request.Login, request.Password) {
		customError(res, "empty login or password", http.StatusBadRequest)
		return
	}

	token, err := h.user.LoginUser(req.Context(), request.Login, request.Password)
	if err != nil {
		handleError(res, err, http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:  "Token",
		Value: token,
		Path:  "/",
	})
}

func checkEmpty(login, password string) bool {
	if login == "" || password == "" {
		log.Error("empty login or password", "login", login, "password", password)
		return false
	}

	return true
}

func (h *ServiceHandlers) Ping(res http.ResponseWriter, req *http.Request) {
	err := h.storage.Ping(req.Context())
	if err != nil {
		handleError(res, err, http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func (h *ServiceHandlers) getUserID(req *http.Request) (int, error) {
	reqToken, err := req.Cookie("Token")
	if err != nil {
		return 0, fmt.Errorf("error getting token", "error", err)
	}

	userID, err := h.token.ParseToken(reqToken.Value)
	if err != nil || userID == 0 {
		return 0, fmt.Errorf("Unauthorized", "error", err)
	}
	return userID, nil
}
