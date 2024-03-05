package wmi

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func SendWinRMCommand(
	ctx context.Context,
	ip, login, pass string,
	port int,
	cmd string,
) (string, error) {
	client, err := newWinRMClient(ip, login, pass, port)
	if err != nil {
		return "", fmt.Errorf("create WMI client error: %s", err)
	}

	var stdout, stderr bytes.Buffer
	_, err = client.RunWithContext(ctx, "chcp 866 | "+cmd, &stdout, &stderr)
	if err != nil {
		return "", fmt.Errorf("cmd [%s] WinRM error: %s", cmd, err)
	}

	stderrStr := stderr.String()
	stdoutStr := stdout.String()
	if stderrStr != "" {
		return "", fmt.Errorf("cmd [%s] error: %s", cmd, stderrStr)
	}

	reader := transform.NewReader(bytes.NewReader([]byte(stdoutStr)), charmap.CodePage866.NewDecoder())
	d, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("encoding stdout [%s] error [%s]", stdoutStr, err)
	}

	return string(d), nil
}

func encodeWindows1251(in string) (string, error) {
	reader := transform.NewReader(bytes.NewReader([]byte(in)), charmap.Windows1251.NewEncoder())
	d, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(d), nil
}
