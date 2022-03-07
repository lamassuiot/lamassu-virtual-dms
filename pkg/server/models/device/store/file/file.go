package files

import (
	"encoding/pem"
	"io/ioutil"
	"os"

	"github.com/lamassuiot/lamassu-default-dms/pkg/server/models/device/store"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type File struct {
	dirPath string
	logger  log.Logger
}

func NewFile(dirPath string, logger log.Logger) store.File {
	return &File{dirPath: dirPath, logger: logger}
}

const (
	csrPerm = 5555
)

func (f *File) InsertCSR(id string, data []byte) error {
	name := f.dirPath + "/device-" + id + ".csr"
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, csrPerm)
	if err != nil {
		level.Error(f.logger).Log("err", err, "msg", "Could not insert CSR with ID "+id+" in filesystem")
		return err
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		level.Error(f.logger).Log("err", err, "msg", "Error encoding bytas as CSR")
		os.Remove(name)
		return err
	}
	level.Info(f.logger).Log("msg", "CSR with ID "+id+" inserted in file system")
	return nil
}

func (f *File) InsertCERTReenroll(id string, data []byte) error {
	name := f.dirPath + "/device-reenrolled-" + id + ".crt"
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, csrPerm)
	if err != nil {
		level.Error(f.logger).Log("err", err, "msg", "Could not insert CERT with ID "+id+" in filesystem")
		return err
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		level.Error(f.logger).Log("err", err, "msg", "Error encoding bytas as CSR")
		os.Remove(name)
		return err
	}
	level.Info(f.logger).Log("msg", "CERT with ID "+id+" inserted in file system")
	return nil
}

func (f *File) InsertCERT(id string, data []byte) error {
	name := f.dirPath + "/device-" + id + ".crt"
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, csrPerm)
	if err != nil {
		level.Error(f.logger).Log("err", err, "msg", "Could not insert CERT with ID "+id+" in filesystem")
		return err
	}
	defer file.Close()

	if err := pem.Encode(file, &pem.Block{Type: "CERTIFICATE", Bytes: data}); err != nil {
		level.Error(f.logger).Log("err", err.Error, "msg", "Error encoding bytes as a certificate for certficate with ID "+id)
		os.Remove(name)
	}
	level.Info(f.logger).Log("msg", "Certificate with ID "+id+" inserted in file system")
	return nil
}

func (f *File) InsertKEY(id string, data []byte) error {
	name := f.dirPath + "/device-" + id + ".key"
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, csrPerm)
	if err != nil {
		level.Error(f.logger).Log("err", err, "msg", "Could not insert KEY with ID "+id+" in filesystem")
		return err
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		level.Error(f.logger).Log("err", err, "msg", "Error encoding bytas as CSR")
		os.Remove(name)
		return err
	}
	level.Info(f.logger).Log("msg", "KEY with ID "+id+" inserted in file system")
	return nil
}

func (f *File) SelectByID(id string) ([]byte, error) {
	name := f.dirPath + "/device-" + id + ".csr"
	data, err := ioutil.ReadFile(name)
	if err != nil {
		level.Error(f.logger).Log("err", err, "msg", "Could not obtain CSR with ID "+id+" from filesystem")
		return nil, err
	}
	level.Info(f.logger).Log("msg", "CSR with ID "+id+" obtained from file system")
	return data, nil
}

func (f *File) Delete(id string) error {
	name := f.dirPath + "/" + id + ".csr"
	err := os.Remove(name)
	if err != nil {
		level.Error(f.logger).Log("err", err, "msg", "Could not delete CSR with ID "+id+" from filesystem")
		return err
	}
	level.Info(f.logger).Log("msg", "CSR with ID "+id+" deleted from file system")
	return nil
}
