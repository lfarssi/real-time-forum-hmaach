package models

type Message struct {
	SenderID   int    `json:"sender_id"`
	Sender     string `json:"sender"`
	ReceiverID int    `json:"receiver_id"`
	Receiver   string `json:"receiver,omitempty"`
	Type       string `json:"type,omitempty"`
	Content    string `json:"content"`
	SentAt     string `json:"sent_at,omitempty"`
}

func GetMessages(receiver, sender, limit, page int) ([]Message, error) {
	var (
		messages []Message
		offset   = page * limit
	)
	query := `
        SELECT 
            m.message,
            receiver.nickname as receiver_name,
            sender.nickname as sender_name,
            m.sent_at,
            m.sender as sender_id,
            m.receiver as receiver_id
        FROM messages m
        LEFT JOIN users receiver ON m.receiver = receiver.id
        LEFT JOIN users sender ON m.sender = sender.id
        WHERE (m.sender = ? AND m.receiver = ?) OR (m.sender = ? AND m.receiver = ?)
        ORDER BY m.sent_at DESC
        LIMIT ? OFFSET ?`

	rows, err := DB.Query(query, sender, receiver, receiver, sender, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var message Message
		err := rows.Scan(
			&message.Content,
			&message.Receiver,
			&message.Sender,
			&message.SentAt,
			&message.SenderID,
			&message.ReceiverID,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		if len(messages) == 0 {
			return []Message{}, nil
		}
		return nil, err
	}

	return messages, nil
}

func StoreMessage(message Message) error {
	query := `
        INSERT INTO messages (sender, receiver, message) 
        VALUES (?, ?, ?)`

	_, err := DB.Exec(query, message.SenderID, message.ReceiverID, message.Content)
	if err != nil {
		return err
	}
	return nil
}
