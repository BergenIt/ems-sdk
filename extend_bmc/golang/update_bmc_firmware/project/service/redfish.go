package redfish

import (
	rfish "bmc-handler/adapter"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	insecureHTTP = "http://"
	secureHTTPS  = "https://"
	SFTP         = "sftp://"
)

const AuthHeader = "X-Auth-Token"

type Host struct {
	Headers  map[string]string
	IP       string
	User     string
	Password string
	IsHTTPS  bool
}

type redfishService struct {
	adapter rfish.RedfishAdapter
}

type RedfishCFG struct {
	Vendor       string       `yaml:"vendor"`
	Model        string       `yaml:"model"`
	SimpleUpdate SimpleUpdate `yaml:"simple_update,omitempty"`
	Optional     Optional     `yaml:"optional,omitempty"`
}

type SimpleUpdate struct {
	UpdateURL     string `yaml:"url"`
	UpdateRequest string `yaml:"request"`
}

type Optional struct {
	BmcServiceURL     string `yaml:"bmc_service_url,omitempty"`
	BmcUpdateRequest  string `yaml:"bmc_update_request,omitempty"`
	BiosServiceURL    string `yaml:"bios_service_url,omitempty"`
	BiosUpdateRequest string `yaml:"bios_update_request,omitempty"`
	BiosReqpath       string `yaml:"bios_reqpath,omitempty"`
}

var (
	ErrTaskURINotFound     = errors.New("can't find task monitoring uri in response from device")
	ErrUnexpectedModel     = errors.New("unexpected vendor model")
	ErrHuaweiWrongFilePath = errors.New("filePath should be with sftp protocol")
	ErrEmptyAuthToken      = errors.New("can't find auth token in headers after creating session")
	ErrUnknownVendor       = errors.New("can't update firmware for unknown vendor")
	ErrUnknownNetworkCause = errors.New("can't suspect network cause of timeout")
)

type RedfishRequestErr struct {
	Origin error
}

func NewRedfish(ra rfish.RedfishAdapter) *redfishService {
	return &redfishService{
		adapter: ra,
	}
}

func (r *RedfishRequestErr) Error() string {
	return fmt.Sprint("redfish request error: " + r.Origin.Error())
}

func (r *redfishService) SimpleUpdate(filePath string, cfg RedfishCFG, host Host) error {
	redfish := r.adapter

	if err := r.connectAdapter(host, &redfish); err != nil {
		return fmt.Errorf("connectAdapter failed: %w", err)
	}

	if err := r.huaweiUpdateFirmware(filePath, host, cfg, &redfish); err != nil {
		return fmt.Errorf("huaweiUpdateFirmware failed: %w", err)
	}

	return nil
}

func (r *redfishService) huaweiUpdateFirmware(
	filePath string,
	host Host,
	cfg RedfishCFG,
	redfish *rfish.RedfishAdapter,
) error {
	if !strings.HasPrefix(filePath, SFTP) {
		return ErrHuaweiWrongFilePath
	}

	h, id, err := redfish.CreateSession(host.Headers, host.User, host.Password)
	if err != nil {
		return fmt.Errorf("CreateSession failed: %w", err)
	}

	_, ok := h[AuthHeader]
	if !ok {
		return ErrEmptyAuthToken
	}

	if err := r.updateFirmware(h, filePath, "@odata.id", "TaskState", cfg, redfish); err != nil {
		return fmt.Errorf("updateFirmware failed: %w", err)
	}

	if err := redfish.DeleteSession(h, id); err != nil {
		return fmt.Errorf("DeleteSession failed: X-Auth-Token: %s: Id: %q", h[AuthHeader], id)
	}

	return nil
}

func (r *redfishService) updateFirmware(
	headers map[string]string,
	filePath, taskField, stateField string,
	cfg RedfishCFG,
	redfish *rfish.RedfishAdapter) error {

	taskMonitor, err := r.prepareUpdFirmware(
		headers,
		cfg.SimpleUpdate.UpdateURL,
		cfg.SimpleUpdate.UpdateRequest,
		filePath,
		taskField,
		redfish,
	)
	if err != nil {
		return fmt.Errorf("prepareUpdFirmware failed: %w", err)
	}

	if err := r.monitoringUpdFirmware(taskMonitor, stateField, redfish); err != nil {
		return fmt.Errorf("monitoringUpdFirmware failed: %w", err)
	}

	return nil
}

func (r *redfishService) prepareUpdFirmware(
	headers map[string]string,
	updateURL, request, filePath, taskField string,
	redfish *rfish.RedfishAdapter,
) (string, error) {
	// формируем тело запроса с подстановкой пути файла вместо {filePath}
	// и делаем маршл в json формат
	body := strings.ReplaceAll(request, "{filePath}", filePath)

	// для Huawei учитываем возожность транспортировки через sftp
	// если нет возможности трансфера файла через https
	if strings.Contains(body, "TransferProtocol") && strings.Contains(body, "{protocol}") {
		switch {
		case strings.HasPrefix(filePath, secureHTTPS):
			protocol := "HTTPS"
			body = strings.ReplaceAll(body, "{protocol}", protocol)
		case strings.HasPrefix(filePath, insecureHTTP):
			protocol := "HTTP"
			body = strings.ReplaceAll(body, "{protocol}", protocol)
		case strings.HasPrefix(filePath, SFTP):
			protocol := "SFTP"
			body = strings.ReplaceAll(body, "{protocol}", protocol)
		}
	}

	resp, err := redfish.PostWithHeaders(headers, updateURL, body)
	if err != nil {
		return "", fmt.Errorf("PostWithHeaders failed: %w", err)
	}
	defer resp.Body.Close()

	var taskMonitor string
	path := resp.Header.Get(taskField)
	if path != "" {
		taskMonitor = path

		return taskMonitor, nil
	}

	if resp.ContentLength != 0 {
		// результат парсим
		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return "", fmt.Errorf("unmarshal response: %w", err)
		}

		// достаем путь мониторинга задачи
		path, ok := result[taskField]
		if !ok {
			return "", ErrTaskURINotFound
		} else {
			taskMonitor = fmt.Sprint(path)
		}
	}

	// возвращаем путь мониторинга для следующего этапа наката обновления
	return taskMonitor, nil
}

// monitoringUpdFirmware метод для мониторинга задачи по обновлению прошивки
func (r *redfishService) monitoringUpdFirmware(taskMonitor, state string, redfish *rfish.RedfishAdapter) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	var getTaskErr error
	for {
		select {
		default:
			// каждый раз делаем запрос по пути мониторинга для получения свежей информации
			resp, err := redfish.GetTaskMonitor(taskMonitor)
			if err != nil {
				if getTaskErr == nil {
					getTaskErr = &RedfishRequestErr{Origin: err}
				}
				continue
			}

			switch resp.StatusCode {
			case http.StatusNoContent, http.StatusOK:
				if resp.ContentLength <= 0 {
					return nil
				}
			case http.StatusAccepted:
				continue
			}

			var result map[string]any
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return fmt.Errorf("unmarshal response from task service: %w", err)
			}
			resp.Body.Close()

			states, ok := result[state].(string)
			if !ok {
				return fmt.Errorf(
					"failed assert states: states value: [%s] status: [%v] content lenght: [%v]",
					states,
					resp.StatusCode,
					resp.ContentLength,
				)
			}
			states = strings.ToLower(states)
			switch states {
			case "completed", "ok":
				return nil
			case "exception", "error", "critical":
				return fmt.Errorf("failed task for updating source: [%s]", result)
			}

		case <-ctx.Done():
			return fmt.Errorf("%w: %s", ctx.Err(), getTaskErr.Error())
		}
	}
}

func (r *redfishService) connectAdapter(host Host, redfish *rfish.RedfishAdapter) error {
	baseURL := composeBaseURL(host.IP, host.IsHTTPS)

	opts := []rfish.Option{
		rfish.WithBasicAuth(true),
		rfish.WithCredentials(host.User, host.Password),
		rfish.WithEndpoint(baseURL),
		rfish.WithInsecure(true),
	}

	if err := redfish.Connect(opts...); err != nil {
		return err
	}

	return nil
}

func composeBaseURL(ip string, isHTTPS bool) string {
	if isHTTPS {
		return secureHTTPS + ip
	}
	return insecureHTTP + ip
}
