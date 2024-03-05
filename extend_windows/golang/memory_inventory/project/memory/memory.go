package memory

import (
	"context"

	pb "windows-handler/gen/cluster-contract"
)

// GetMemoryInv - Получить инвентарные данные по ОЗУ
func GetMemoryInv(
	ctx context.Context,
	address string,
	wmiCreds *pb.Credential,
) ([]*pb.MemoryCard, error) {
	ramInvCmd := composeWmiCmd(model.WmiCommandTemplate, model.WmiMemoryInvClass, model.WmiMemoryInvFileds)

	res, err := wmi.SendWinRMCommand(ctx, wmiCreds, ramInvCmd, timeout)
	if err != nil {
		return nil, err
	}

	return parseRAMInvInfo(res), nil
}
