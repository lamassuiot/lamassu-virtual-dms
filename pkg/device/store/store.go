package store

type File interface {
	InsertCSR(id string, data []byte, types string, dmsname string) error
	InsertCERT(id string, data []byte, types string, dmsname string, serialnumber string) error
	InsertCERTReenroll(id string, data []byte) error
	InsertKEY(id string, data []byte, types string, dmsname string) error
	SelectByID(id string) ([]byte, error)
	Delete(id string) error
}
