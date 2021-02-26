package json

type GetUsersRequest struct {
}

type EventMeta struct {
	EventType string `json:"type"`
}

type CreateUserRequest struct {
	ID      string `json:id`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

type UpdateUserRequest struct {
	ID      string `json:id`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

type DeleteUserRequest struct {
	ID string `json:id`
}

type CreateTxsRequest struct {
	From  string `json:from`
	To    string `json:to`
	Money int    `json:money`
}
