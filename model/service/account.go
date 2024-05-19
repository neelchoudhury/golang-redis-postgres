package service

type Account struct {
	Id      int64   `pg:"type:serial,pk"`
	Name    string  `json:"name" pg:"type:varchar(50),pk,notnull,notnull"`
	Balance float64 `json:"balance" pg:"type:float8,notnull"`
}
