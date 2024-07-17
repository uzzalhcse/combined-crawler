package models

type GlobalSecret struct {
}

func (c *GlobalSecret) GetTableName() string {
	return "global_secrets"
}
