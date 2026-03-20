package event

import (
	"encoding/json"
	"fmt"

	"github.com/jackz-jones/nft-contract-go/types"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

const (
	MintEvent                   = "MintEvent"
	TransferEvent               = "TransferEvent"
	CrossChainTransferEvent     = "CrossChainTransferEvent"
	CrossChainMintEvent         = "CrossChainMintEvent"
	UpdateCrossChainStatusEvent = "UpdateCrossChainStatusEvent"
)

// EmitMintEvent NFT Mint 事件
func EmitMintEvent(ni types.NFTInfo) error {
	nftInfoBytes, err := json.Marshal(ni)
	if err != nil {
		return err
	}

	sdk.Instance.EmitEvent(MintEvent, []string{string(nftInfoBytes)})
	return nil
}

// EmitTransferEvent NFT Transfer 事件
func EmitTransferEvent(ni types.NFTInfo, from, to string) error {
	nftInfoBytes, err := json.Marshal(ni)
	if err != nil {
		return err
	}

	sdk.Instance.EmitEvent(TransferEvent, []string{string(nftInfoBytes), from, to})
	return nil
}

// EmitCrossChainTransferEvent CrossChainTransfer 事件
func EmitCrossChainTransferEvent(ni types.NFTInfo, from, to string) error {
	nftInfoBytes, err := json.Marshal(ni)
	if err != nil {
		return err
	}

	sdk.Instance.EmitEvent(CrossChainTransferEvent, []string{string(nftInfoBytes), from, to})
	return nil
}

// EmitCrossChainMintEvent CrossChainMint 事件
func EmitCrossChainMintEvent(ni types.NFTInfo) error {
	nftInfoBytes, err := json.Marshal(ni)
	if err != nil {
		return err
	}

	sdk.Instance.EmitEvent(CrossChainMintEvent, []string{string(nftInfoBytes)})
	return nil
}

// EmitUpdateCrossChainStatusEvent UpdateCrossChainStatus 事件
func EmitUpdateCrossChainStatusEvent(ni types.NFTInfo, crossState types.CrossState) error {
	nftInfoBytes, err := json.Marshal(ni)
	if err != nil {
		return err
	}

	sdk.Instance.EmitEvent(UpdateCrossChainStatusEvent, []string{string(nftInfoBytes), fmt.Sprintf("%d", crossState)})
	return nil
}
