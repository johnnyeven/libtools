package database

type Query struct {
	Type               string        `json:"type"`
	ExpectedSQLKeyWord string        `json:"sqlKeyword,omitempty"`
	WithArgs           []interface{} `json:"args,omitempty"`
	ReturnResult       *ReturnResult `json:"result,omitempty"`
	ReturnRows         *ReturnRows   `json:"rows,omitempty"`
	ReturnError        string        `json:"error,omitempty"`
}

type ReturnResult struct {
	LastInsertID int64 `json:"lastInsertID"`
	RowsEffected int64 `json:"rowsEffected"`
}

type ReturnRows struct {
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
}
