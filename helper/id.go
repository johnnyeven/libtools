package helper

import (
	"github.com/johnnyeven/libtools/clients/client_id"
	"github.com/sirupsen/logrus"
)

func NewUniqueID(client client_id.ClientIDInterface) (uint64, error) {
	resp, err := client.GetNewId()
	if err != nil {
		logrus.Errorf("ClientID.GetNewId err: %v", err)
		return 0, err
	}

	return resp.Body.ID, nil
}
