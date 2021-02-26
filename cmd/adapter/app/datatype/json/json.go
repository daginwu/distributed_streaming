package json

type GetUsersRequest struct {
}

type CreateUserRequest struct {
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

type GetUserRequest struct {
}

type UpdateUserRequest struct {
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

type DeleteUserRequest struct {
}

type GetTxsRequest struct {
}

type CreateTxRequest struct {
	From  string `json:from`
	To    string `json:to`
	Money int    `json:money`
}
