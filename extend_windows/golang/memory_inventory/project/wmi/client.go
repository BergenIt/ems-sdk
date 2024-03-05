package wmi

import (
	"fmt"
	"time"

	"github.com/masterzen/winrm"
)

func newWinRMClient(ip, login, pass string, port int) (*winrm.Client, error) {
	login1251, err := encodeWindows1251(login)
	if err != nil {
		return nil, fmt.Errorf("encoding login [%s] error [%s]", login, err)
	}

	pass1251, err := encodeWindows1251(pass)
	if err != nil {
		return nil, fmt.Errorf("encoding password for login [%s] error [%s]", login, err)
	}

	endpoint := winrm.NewEndpoint(ip, port, false, false, nil, nil, nil, 10*time.Second)
	client, err := winrm.NewClient(endpoint, login1251, pass1251)
	if err != nil {
		return nil, fmt.Errorf("create WMI client error: %s", err)
	}

	return client, nil
}
