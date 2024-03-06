package memory

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	pb "windows-handler/gen/cluster-contract"
	"windows-handler/wmi"
)

const (
	wmiCommandTemplate = "WMIC PATH @Class GET @Fields /format:list"
	wmiClassTemplate   = "@Class"
	wmiFieldsTemplate  = "@Fields"

	wmiMemoryInvClass = "Win32_PhysicalMemory"

	wmiMemoryInvFields = "BankLabel,TypeDetail,MemoryType,SMBIOSMemoryType,FormFactor," +
		"Manufacturer,Capacity,PartNumber,SerialNumber,Status,Speed,DeviceLocator"
)

// GetMemoryInv - Получить инвентарные данные по ОЗУ
func GetMemoryInv(
	ctx context.Context,
	address string,
	wmiCreds *pb.Credential,
) ([]*pb.MemoryCard, error) {
	ramInvCmd := composeWmiCmd(wmiCommandTemplate, wmiMemoryInvClass, wmiMemoryInvFields)

	res, err := wmi.SendWinRMCommand(ctx, address, wmiCreds.Login, wmiCreds.Password, int(wmiCreds.Port), ramInvCmd)
	if err != nil {
		return nil, err
	}

	return parseRAMInvInfo(res), nil
}

// Формирование команды для отправки по WMI
func composeWmiCmd(cmd, class, fields string) string {
	wmiCmd := strings.ReplaceAll(cmd, wmiClassTemplate, class)
	wmiCmd = strings.ReplaceAll(wmiCmd, wmiFieldsTemplate, fields)
	return wmiCmd
}

// Парсинг ответа от WMI для вытягивания информации по ОЗУ
func parseRAMInvInfo(stdout string) []*pb.MemoryCard {
	var ramInvInfos []*pb.MemoryCard

	metricInfos := handleMetricStdout(stdout)
	for _, metricInfo := range metricInfos {
		ramInvInf := &pb.MemoryCard{}
		var (
			memoryDeviceType uint16
			smBIOSMemoryType uint32
		)

		for key, value := range metricInfo {
			var i64 int64

			switch key {
			case "BankLabel":
				i64, _ = parseStrToInt64(key, strings.ReplaceAll(value, "BANK ", ""))
				ramInvInf.Slot = int32(i64)
			case "Manufacturer":
				ramInvInf.Vendor = value
			case "PartNumber":
				ramInvInf.PartNumber = value
			case "SerialNumber":
				ramInvInf.SerialNumber = value
			case "Status":
				ramInvInf.State = memStateInternal(value)
			case "TypeDetail":
				i64, _ = parseStrToInt64(key, value)
				ramInvInf.MemoryType = memTypeInternal(uint16(i64))
			case "MemoryType":
				i64, _ = parseStrToInt64(key, value)
				memoryDeviceType = uint16(i64)
			case "SMBIOSMemoryType":
				i64, _ = parseStrToInt64(key, value)
				smBIOSMemoryType = uint32(i64)
			case "FormFactor":
				i64, _ = parseStrToInt64(key, value)
				ramInvInf.BaseModuleType = memBaseModuleTypeInternal(uint16(i64))
			case "Capacity":
				i64, _ = parseStrToInt64(key, value)
				if i64 != 0 {
					ramInvInf.Size = int32(i64 / 1024 / 1024)
				}
			case "Speed":
				i64, _ = parseStrToInt64(key, value)
				ramInvInf.SpeedMhz = int32(i64)
			case "DeviceLocator":
				ramInvInf.Location = value
			default:
				continue
			}
		}

		ramInvInf.MemoryDeviceType = memDeviceTypeInternal(memoryDeviceType, smBIOSMemoryType)
		if ramInvInf != (&pb.MemoryCard{}) {
			ramInvInfos = append(ramInvInfos, ramInvInf)
		}
	}

	return ramInvInfos
}

func memDeviceTypeInternal(memType uint16, smbiosMemType uint32) pb.MemoryDeviceType {
	var checkedMemType uint32
	if memType != 0 {
		checkedMemType = uint32(memType)
	} else {
		checkedMemType = smbiosMemType
	}

	switch checkedMemType {
	case 5:
		return pb.MemoryDeviceType_MEMORY_DEVICE_TYPE_EDO
	case 6:
		return pb.MemoryDeviceType_MEMORY_DEVICE_TYPE_EDRAM
	case 7:
		return pb.MemoryDeviceType_MEMORY_DEVICE_TYPE_VRAM
	case 9:
		return pb.MemoryDeviceType_MEMORY_DEVICE_TYPE_RAM
	case 10:
		return pb.MemoryDeviceType_MEMORY_DEVICE_TYPE_ROM
	case 12:
		return pb.MemoryDeviceType_MEMORY_DEVICE_TYPE_EEPROM
	case 13:
		return pb.MemoryDeviceType_MEMORY_DEVICE_TYPE_FEPROM
	case 14:
		return pb.MemoryDeviceType_MEMORY_DEVICE_TYPE_EPROM
	case 15:
		return pb.MemoryDeviceType_MEMORY_DEVICE_TYPE_CDRAM
	case 17:
		return pb.MemoryDeviceType_MEMORY_DEVICE_TYPE_SDRAM
	case 18:
		return pb.MemoryDeviceType_MEMORY_DEVICE_TYPE_DDR_SGRAM
	case 20:
		return pb.MemoryDeviceType_MEMORY_DEVICE_TYPE_DDR
	case 21:
		return pb.MemoryDeviceType_MEMORY_DEVICE_TYPE_DDR2
	case 22:
		return pb.MemoryDeviceType_MEMORY_DEVICE_TYPE_DDR2_SDRAM_FB_DIMM
	case 24:
		return pb.MemoryDeviceType_MEMORY_DEVICE_TYPE_DDR3
	case 26:
		return pb.MemoryDeviceType_MEMORY_DEVICE_TYPE_DDR4
	default:
		return pb.MemoryDeviceType_MEMORY_DEVICE_TYPE_UNSPECIFIED
	}
}

func memTypeInternal(typeDetail uint16) pb.MemoryType {
	switch typeDetail {
	case 8:
		return pb.MemoryType_MEMORY_TYPE_FPRAM
	case 16:
		return pb.MemoryType_MEMORY_TYPE_SRAM
	case 32:
		return pb.MemoryType_MEMORY_TYPE_PSRAM
	case 64:
		return pb.MemoryType_MEMORY_TYPE_RAMBUS
	case 128:
		return pb.MemoryType_MEMORY_TYPE_S_DRAM
	case 256:
		return pb.MemoryType_MEMORY_TYPE_CMOS
	case 512:
		return pb.MemoryType_MEMORY_TYPE_EDO_RAM
	case 1024:
		return pb.MemoryType_MEMORY_TYPE_WIN_DRAM
	case 2048:
		return pb.MemoryType_MEMORY_TYPE_CACHE_DRAM
	case 4096:
		return pb.MemoryType_MEMORY_TYPE_NVRAM
	default:
		return pb.MemoryType_MEMORY_TYPE_UNSPECIFIED
	}
}

func memBaseModuleTypeInternal(formFactor uint16) pb.BaseModuleType {
	switch formFactor {
	case 2:
		return pb.BaseModuleType_BASE_MODULE_TYPE_SIP
	case 3:
		return pb.BaseModuleType_BASE_MODULE_TYPE_DIP
	case 4:
		return pb.BaseModuleType_BASE_MODULE_TYPE_ZIP
	case 5:
		return pb.BaseModuleType_BASE_MODULE_TYPE_SOJ
	case 6:
		return pb.BaseModuleType_BASE_MODULE_TYPE_PROPRIETARY
	case 7:
		return pb.BaseModuleType_BASE_MODULE_TYPE_SIMM
	case 8:
		return pb.BaseModuleType_BASE_MODULE_TYPE_DIMM
	case 9:
		return pb.BaseModuleType_BASE_MODULE_TYPE_TSOP
	case 10:
		return pb.BaseModuleType_BASE_MODULE_TYPE_PGA
	case 11:
		return pb.BaseModuleType_BASE_MODULE_TYPE_RIMM
	case 12:
		return pb.BaseModuleType_BASE_MODULE_TYPE_SO_DIMM
	case 13:
		return pb.BaseModuleType_BASE_MODULE_TYPE_SRIMM
	case 14:
		return pb.BaseModuleType_BASE_MODULE_TYPE_SMD
	case 15:
		return pb.BaseModuleType_BASE_MODULE_TYPE_SSMP
	case 16:
		return pb.BaseModuleType_BASE_MODULE_TYPE_QFP
	case 17:
		return pb.BaseModuleType_BASE_MODULE_TYPE_TQFP
	case 18:
		return pb.BaseModuleType_BASE_MODULE_TYPE_SOIC
	case 19:
		return pb.BaseModuleType_BASE_MODULE_TYPE_LCC
	case 20:
		return pb.BaseModuleType_BASE_MODULE_TYPE_PLCC
	case 21:
		return pb.BaseModuleType_BASE_MODULE_TYPE_BGA
	case 22:
		return pb.BaseModuleType_BASE_MODULE_TYPE_FPBGA
	case 23:
		return pb.BaseModuleType_BASE_MODULE_TYPE_LGA
	default:
		return pb.BaseModuleType_BASE_MODULE_TYPE_UNSPECIFIED
	}
}

func memStateInternal(memState string) pb.MemoryState {
	switch memState {
	case "OK", "Degraded", "Pred Fail":
		return pb.MemoryState_MEMORY_STATE_OK
	case "Error", "Service", "Starting", "Stopping":
		return pb.MemoryState_MEMORY_STATE_CRITICAL
	default:
		return pb.MemoryState_MEMORY_STATE_UNKNOWN
	}
}

// HandleMetricStdout - Обработка потока вывода метрики
func handleMetricStdout(stdout string) []map[string]string {
	strList := strings.Trim(stdout, "\r\n")
	strList = strings.ReplaceAll(strList, "\r", "")
	metricInfosRaw := strings.Split(strList, "\n\n\n")

	metricInfos := getMapFromMetric(metricInfosRaw)

	return metricInfos
}

func getMapFromMetric(rawMetricInfos []string) []map[string]string {
	metricInfos := []map[string]string{}

	for _, ramStatInfoRaw := range rawMetricInfos {
		ramStatInfoItems := strings.Split(ramStatInfoRaw, "\n")
		metricInfo := make(map[string]string)
		for _, ramStatInfoItem := range ramStatInfoItems {
			itemAndValue := strings.Split(ramStatInfoItem, "=")
			if len(itemAndValue) < 2 || itemAndValue[1] == "" {
				continue
			}
			metricInfo[itemAndValue[0]] = itemAndValue[1]
		}
		metricInfos = append(metricInfos, metricInfo)
	}

	return metricInfos
}

func parseStrToInt64(field, value string) (int64, error) {
	i64, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse %s [%s] failed: %s", field, value, err)
	}
	return i64, nil
}
