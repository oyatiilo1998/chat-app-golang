package storage

import (
	"Study/websocket/models"
	"Study/websocket/pkg/security"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *PostgresStorage {
	return &PostgresStorage{
		db: db,
	}
}

func (ps *PostgresStorage) Login(req models.Login) (res models.LoginResp, err error) {
	var hashedPassword string
	chechPasswordQuery := `
		SELECT
			password
		FROM
			users
		WHERE
			username = $1
	`
	
	getIdQuery := `
		SELECT
			id
		FROM
			users
		WHERE
			username = $1
	`

	userRow := ps.db.QueryRow(chechPasswordQuery,req.Username)

	err = userRow.Scan(&hashedPassword)
	if err != nil {
		return res, err
	}

	doesMatch, err := security.ComparePassword(hashedPassword, req.Password)
	if err != nil {
		return res, err
	}

	if doesMatch {
		userRow := ps.db.QueryRow(getIdQuery,req.Username)

		err = userRow.Scan(&res.ID)
		if err != nil {
			return res, err
		}
	}

	return res, nil
}


func (ps *PostgresStorage) GetChatUsers(ID string) (response models.GetChatUsersResponse, err error) {	
	query  := `
	SELECT 
		"from", "to", text, created_at, messages.id, users.username, users.id, result.unread
    FROM 
		messages 
	JOIN 
		users
	ON 
		("to"=users.id or "from"=users.id) and users.id<>$1,
	(
		SELECT
			LEAST("from", "to"), GREATEST("from", "to"), MAX(created_at), count ($1) filter (where read=false and "from"<>$1) as unread
		FROM
			messages
		GROUP BY
			LEAST("from", "to"), GREATEST("from", "to")
	) as result
    WHERE
		messages.created_at=result.max and (messages.from=$1 or messages.to=$1)
	ORDER BY 
		result.max DESC
	`
	rows, err := ps.db.Query(query, ID)
	if err != nil {
		fmt.Println(err)
		return response, err
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User

		err = rows.Scan(
			&user.LastMessage.From,
			&user.LastMessage.To,
			&user.LastMessage.Text,
			&user.LastMessage.CreatedAt,
			&user.LastMessage.ID,
			&user.Username,
			&user.ID,
			&user.UnreadMessages,
		)

		response.Users = append(response.Users, user)
	}

	return
}

func (ps *PostgresStorage) GetChatHistory(currentUserId, peerId string) (messagesArray []models.DateMessage, err error) {
	messages := make(map[string][]models.Message)
	query := `
		SELECT
			id,
			"from",
			"to",
			text,
			created_at,
			read
		FROM
			messages
		WHERE
			("from"=$1 and "to"=$2) or ("to"=$1 and "from"=$2)
		ORDER BY
			created_at
		ASC
	`

	rows, err := ps.db.Query(query, currentUserId, peerId)
	if err != nil {
		return messagesArray, err
	}

	for rows.Next() {
		var msg models.Message

		err = rows.Scan(
			&msg.ID,
			&msg.From,
			&msg.To,
			&msg.Text,
			&msg.CreatedAt,
			&msg.Read,
		)
		if err != nil {
			return messagesArray, err
		}
		messages[msg.CreatedAt[:10]] = append(messages[msg.CreatedAt[:10]], msg)
	}

	for key, val := range messages {
		messagesArray = append(messagesArray, models.DateMessage{
			Date: key,
			Messages: val,
		})
	}

	return
}

func (ps *PostgresStorage) CreateMessage(req models.Message) (res models.Message, err error) {
	req.CreatedAt = time.Now().Format("2006-01-02T15:04:05Z07:00")
	query := `
		INSERT INTO
			messages
		(
			id,
			"from",
			"to",
			text,
			created_at,
			read
		)
		VALUES
		(
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		)
	`
	_, err = ps.db.Exec(query,
			&req.ID,
			&req.From,
			&req.To,
			&req.Text,
			&req.CreatedAt,
			&req.Read,			
		)
	if err != nil {
		return req, err
	}
	
	return req, nil
}

func (ps *PostgresStorage) UpdateReadStatus(req []string) error {
	params := make(map[string]interface{})

	params["message_id"] = pq.Array(req)

	updateQuery := `
		UPDATE
			messages
		SET
			read=true
		WHERE
			id=ANY(:message_id)
	`

	_, err := ps.db.NamedExec(updateQuery, params)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}