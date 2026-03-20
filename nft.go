package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/jackz-jones/nft-contract-go/code"
	_const "github.com/jackz-jones/nft-contract-go/const"
	"github.com/jackz-jones/nft-contract-go/event"
	"github.com/jackz-jones/nft-contract-go/types"
	"github.com/jackz-jones/nft-contract-go/util"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

type EblNFT struct{}

func (e *EblNFT) InitContract() protogo.Response {

	// 记录合约创建者
	errResp := util.StoreContractCreator()
	if errResp != nil {
		return *errResp
	}

	// 记录合约版本
	errResp = util.StoreCurrentVersion()
	if errResp != nil {
		return *errResp
	}

	return sdk.Success([]byte("Init success"))
}

func (e *EblNFT) UpgradeContract() protogo.Response {

	// 更新当前合约版本
	errResp := util.StoreCurrentVersion()
	if errResp != nil {
		return *errResp
	}

	return sdk.Success([]byte("Upgrade success"))
}

func (e *EblNFT) InvokeContract(method string) (resp protogo.Response) {
	args := sdk.Instance.GetArgs()
	if len(method) == 0 {
		return sdk.Error("method of param should not be empty")
	}

	defer func() {
		if resp.Status != sdk.OK {
			sdk.Instance.Warnf(resp.Message)
		}
	}()

	switch method {
	case _const.MethodMint:
		var ni types.NFTInfo
		nftInfoBytes := args[_const.ParamNFTInfo]
		err := json.Unmarshal(nftInfoBytes, &ni)
		if err != nil {
			return code.ErrorResp(code.ErrInternalJsonUnmarshalFailed,
				fmt.Sprintf(code.InternalJsonUnmarshalFailedStrModifier, string(nftInfoBytes), err))
		}

		return e.Mint(ni)

	case _const.MethodTransfer:
		tokenIdBytes := args[_const.ParamTokenId]
		fromBytes := args[_const.ParamFrom]
		toBytes := args[_const.ParamTo]
		return e.Transfer(string(tokenIdBytes), string(fromBytes), string(toBytes))

	case _const.MethodCrossChainTransfer:
		tokenIdBytes := args[_const.ParamTokenId]
		fromBytes := args[_const.ParamFrom]
		toBytes := args[_const.ParamTo]
		return e.CrossChainTransfer(string(tokenIdBytes), string(fromBytes), string(toBytes))

	case _const.MethodCrossChainMint:
		var ni types.NFTInfo
		nftInfoBytes := args[_const.ParamNFTInfo]
		err := json.Unmarshal(nftInfoBytes, &ni)
		if err != nil {
			return code.ErrorResp(code.ErrInternalJsonUnmarshalFailed,
				fmt.Sprintf(code.InternalJsonUnmarshalFailedStrModifier, string(nftInfoBytes), err))
		}

		return e.CrossChainMint(ni)

	case _const.MethodUpdateCrossChainStatus:
		tokenIdBytes := args[_const.ParamTokenId]
		stateBytes := args[_const.ParamState]
		state, err := strconv.Atoi(string(stateBytes))
		if err != nil {
			return code.ErrorResp(code.ErrInternalAtoiFailed, fmt.Sprintf(code.InternalAtoiFailedStrModifier, _const.ParamState, string(stateBytes), err))
		}

		return e.UpdateCrossChainStatus(string(tokenIdBytes), types.CrossState(state))

	case _const.MethodQueryContractVersion:
		return e.QueryContractVersion()

	default:
		return code.ErrorResp(code.ErrInvalidMethod, fmt.Sprintf("unknown method %s", method))
	}
}

// Mint 出口企业首先在本地 mint 一个 nft，当前权属于出口企业
func (e *EblNFT) Mint(ni types.NFTInfo) protogo.Response {
	sender, err := sdk.Instance.Sender()
	if err != nil {
		return code.ErrorResp(code.ErrInternalGetSenderFailed, "")
	}

	// set sender
	ni.Sender = sender

	// init state unlock
	ni.State = types.NFTState_Unlock

	// gen token id
	ni.TokenId = util.GenTokenId(ni.OriginHash, ni.Data)

	// store nft info
	errResp := util.StoreNFTInfo(ni)
	if errResp != nil {
		return *errResp
	}

	// store nft ownership
	errResp = util.StoreNFTOwnership(ni.TokenId, ni.Owner)
	if errResp != nil {
		return *errResp
	}

	// emit event
	err = event.EmitMintEvent(ni)
	if err != nil {
		return code.ErrorResp(code.ErrInternalEmitEvent, fmt.Sprintf(code.InternalEmitEventFailedStrModifier,
			event.MintEvent, err))
	}

	return util.SuccessNormal()
}

// Transfer 出口企业在本侧链上转移 nft 的所有权到进口企业，当前是 lock 状态
func (e *EblNFT) Transfer(tokenId, from, to string) protogo.Response {

	// get nftInfo
	ni, errResp := util.GetNFTInfo(tokenId)
	if errResp != nil {
		return *errResp
	}

	// check whether nft can be transferred
	errResp = util.CanTransfer(*ni, from)
	if errResp != nil {
		return *errResp
	}

	// get sender
	sender, err := sdk.Instance.Sender()
	if err != nil {
		return code.ErrorResp(code.ErrInternalGetSenderFailed, "")
	}

	// set operation sender
	ni.Sender = sender

	// set state lock
	ni.State = types.NFTState_Lock

	// store nft info
	errResp = util.StoreNFTInfo(*ni)
	if errResp != nil {
		return *errResp
	}

	// emit event
	err = event.EmitTransferEvent(*ni, from, to)
	if err != nil {
		return code.ErrorResp(code.ErrInternalEmitEvent, fmt.Sprintf(code.InternalEmitEventFailedStrModifier,
			event.TransferEvent, err))
	}

	return util.SuccessNormal()
}

// CrossChainTransfer 本地侧发起跨链通知，通知跨链服务处理
func (e *EblNFT) CrossChainTransfer(tokenId, from, to string) protogo.Response {

	// get nftInfo
	ni, errResp := util.GetNFTInfo(tokenId)
	if errResp != nil {
		return *errResp
	}

	// get sender
	sender, err := sdk.Instance.Sender()
	if err != nil {
		return code.ErrorResp(code.ErrInternalGetSenderFailed, "")
	}

	// set operation sender
	ni.Sender = sender

	// store nft info
	errResp = util.StoreNFTInfo(*ni)
	if errResp != nil {
		return *errResp
	}

	// emit event
	err = event.EmitCrossChainTransferEvent(*ni, from, to)
	if err != nil {
		return code.ErrorResp(code.ErrInternalEmitEvent, fmt.Sprintf(code.InternalEmitEventFailedStrModifier,
			event.CrossChainTransferEvent, err))
	}

	return util.SuccessNormal()
}

// CrossChainMint 跨链服务收到 CrossChainTransfer 事件后在目标链上发起，这时的 owner 应该传 to 地址了
func (e *EblNFT) CrossChainMint(ni types.NFTInfo) protogo.Response {
	sender, err := sdk.Instance.Sender()
	if err != nil {
		return code.ErrorResp(code.ErrInternalGetSenderFailed, "")
	}

	// set sender
	ni.Sender = sender

	// init state unlock
	ni.State = types.NFTState_Unlock

	// gen token id
	ni.TokenId = util.GenTokenId(ni.OriginHash, ni.Data)

	// store nft info
	errResp := util.StoreNFTInfo(ni)
	if errResp != nil {
		return *errResp
	}

	// store nft ownership
	errResp = util.StoreNFTOwnership(ni.TokenId, ni.Owner)
	if errResp != nil {
		return *errResp
	}

	// emit event
	err = event.EmitCrossChainMintEvent(ni)
	if err != nil {
		return code.ErrorResp(code.ErrInternalEmitEvent, fmt.Sprintf(code.InternalEmitEventFailedStrModifier,
			event.CrossChainMintEvent, err))
	}

	return util.SuccessNormal()
}

// UpdateCrossChainStatus 更新跨链状态，跨链服务调用，本侧收到状态通知更新本地的 nft 信息
func (e *EblNFT) UpdateCrossChainStatus(tokenId string, state types.CrossState) protogo.Response {
	sender, err := sdk.Instance.Sender()
	if err != nil {
		return code.ErrorResp(code.ErrInternalGetSenderFailed, "")
	}

	// get nftInfo
	ni, errResp := util.GetNFTInfo(tokenId)
	if errResp != nil {
		return *errResp
	}

	// set sender
	ni.Sender = sender

	// 跨链成功则把之前 lock 住的 nft ownership 改成零地址，否则 unlock
	switch state {
	case types.CrossState_Success:

		// set nft info owner to zero address
		ni.Owner = _const.ZeroAddress

		// set nft ownership to zero address
		errResp = util.SetNFTOwnershipToZero(tokenId)
		if errResp != nil {
			return *errResp
		}

	case types.CrossState_Failed:
		ni.State = types.NFTState_Unlock

	default:
		return code.ErrorResp(code.ErrUnknownState, fmt.Sprintf("invalid cross state: %d", state))
	}

	// store nft info
	errResp = util.StoreNFTInfo(*ni)
	if errResp != nil {
		return *errResp
	}

	// emit event
	err = event.EmitUpdateCrossChainStatusEvent(*ni, state)
	if err != nil {
		return code.ErrorResp(code.ErrInternalEmitEvent, fmt.Sprintf(code.InternalEmitEventFailedStrModifier,
			event.UpdateCrossChainStatusEvent, err))
	}

	return util.SuccessNormal()
}

func (e *EblNFT) QueryContractVersion() protogo.Response {
	verBytes, err := sdk.Instance.GetStateByte(_const.KeyNFTContractVersion, "")
	if err != nil {
		return code.ErrorResp(code.ErrInternalDataReadFailed, fmt.Sprintf(code.InternalDataReadFailedStrModifier,
			_const.KeyNFTContractVersion, "", err))
	}

	return sdk.Success(verBytes)
}

func main() {
	err := sandbox.Start(new(EblNFT))
	if err != nil {
		log.Fatal(err)
	}
}
