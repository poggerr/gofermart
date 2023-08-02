package storage

var schema = `
CREATE TABLE IF NOT EXISTS main_user (
    id UUID UNIQUE,
    username text,
    password text,
    balance float,
    withdrawn float
);

CREATE TABLE IF NOT EXISTS orders (
    id UUID UNIQUE,
    order_number bigint UNIQUE,
    order_user UUID,
    uploaded_at date,
    accrual_service float,
    status text
);

CREATE TABLE IF NOT EXISTS withdrawals (
    id UUID UNIQUE,
    order_number bigint UNIQUE,
    order_user UUID,
    sum float,
    processed_at date
)
`

func (strg *Storage) RestoreDB() {
	strg.Db.MustExec(schema)
}
