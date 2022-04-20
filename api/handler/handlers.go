package handler

import (
	"Study/websocket/config"
	"Study/websocket/models"
	"Study/websocket/pkg/security"
	"Study/websocket/storage"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type handler struct {
	storage storage.PostgresStorage
	cfg config.Config
}

func New(db *sqlx.DB, cfg config.Config) *handler {
	return &handler{
		storage: *storage.New(db),
		cfg: cfg,
	}
}

func (h *handler) Login(c *gin.Context) {
	var Credentials models.Login

	err := c.ShouldBindJSON(&Credentials)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := h.storage.Login(Credentials)
	
	if err != nil {
		h.handleErrorResponse(c, 400, "invalid user credentials", err)
		return
	}

	m := map[string]interface{}{
		"id":                 resp.ID,
	}

	accessToken, err := security.GenerateJWT(m, (7 * time.Hour * 24), h.cfg.SecretKey)
	if err != nil {
		h.handleErrorResponse(c, 500, "server error", err.Error())
		return
	}

	c.JSON(http.StatusOK, models.LoginTokenResp{AccessToken: accessToken, UserId: resp.ID, Code: "200"})
}

func (h *handler) GetChatUsers(c *gin.Context) {
	id := c.Param("id")

	res, err :=h.storage.GetChatUsers(id)
	if err  != nil {
		h.handleErrorResponse(c, http.StatusBadRequest, "something went wrong" , err)
	}

	c.JSON(http.StatusOK, res)
}


func (h *handler) GetChatHistory(c *gin.Context) {
	var users models.GetChatHistoryRequest
	
	users.UserId = c.DefaultQuery("user_id", "")
	users.PeerId = c.DefaultQuery("peer_id", "")

	res, err :=h.storage.GetChatHistory(users.UserId, users.PeerId)
	if err  != nil {
		h.handleErrorResponse(c, http.StatusBadRequest, "something went wrong" , err)
	}

	c.JSON(http.StatusOK, res)
}

func (h *handler) CreateMessage(c *gin.Context) {
	var message models.Message
	
	err := c.ShouldBindJSON(&message)
	if err != nil {
		fmt.Println(err)
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		h.handleErrorResponse(c, http.StatusBadRequest, "something went wrong while creating id" , err)
	}
	message.ID = id.String()

	_, err =h.storage.CreateMessage(message)
	if err  != nil {
		h.handleErrorResponse(c, http.StatusBadRequest, "something went wrong" , err)
	}

	c.JSON(http.StatusOK, message)
}

func (h *handler) UpdateReadStatus(c *gin.Context) {
	var messages []string
	
	err := c.ShouldBindJSON(&messages)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(messages)

	err =h.storage.UpdateReadStatus(messages)
	if err  != nil {
		h.handleErrorResponse(c, http.StatusBadRequest, "something went wrong" , err)
		return
	}

	c.JSON(http.StatusOK, messages,)
}

