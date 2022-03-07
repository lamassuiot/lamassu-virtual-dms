package store

type File interface {
	InsertCSR(id string, data []byte) error
	InsertCERT(id string, data []byte) error
	InsertCERTReenroll(id string, data []byte) error
	InsertKEY(id string, data []byte) error
	SelectByID(id string) ([]byte, error)
	Delete(id string) error
}
