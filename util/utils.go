package util

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/jackz-jones/nft-contract-go/code"
	_const "github.com/jackz-jones/nft-contract-go/const"
	"github.com/jackz-jones/nft-contract-go/types"
	"github.com/jackz-jones/nft-contract-go/version"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

func SuccessNormal() protogo.Response {
	return sdk.Success([]byte("ok"))
}

func StoreCurrentVersion() *protogo.Response {

	// 获取当前版本信息
	ver := version.VerInfo{Version: version.Version, CommitID: version.CommitID, BuildTime: version.BuildTime}
	verBytes, err := json.Marshal(&ver)
	if err != nil {
		errResp := code.ErrorResp(code.ErrInternalJsonMarshalFailed,
			fmt.Sprintf(code.InternalJsonMarshalFailedStrModifier, ver, err))
		return &errResp
	}

	// 记录合约版本
	if err = sdk.Instance.PutStateByte(_const.KeyNFTContractVersion, "", verBytes); err != nil {
		errResp := code.ErrorResp(code.ErrInternalDataStoreFailed, fmt.Sprintf(code.InternalDataStoreFailedStrModifier,
			_const.KeyNFTContractVersion, "", err))
		return &errResp
	}

	return nil
}

func GenTokenId(originHash, data string) string {
	dataHash := sha256.Sum256([]byte(data))
	finalHash := sha256.Sum256([]byte(originHash + string(dataHash[:])))
	return fmt.Sprintf("%x", finalHash)
}

func CanTransfer(ni types.NFTInfo, from string) *protogo.Response {

	// check state, only unlock state can be transfer
	if ni.State == types.NFTState_Lock {
		errFResp := code.ErrorResp(code.ErrNFTLocked, "")
		return &errFResp
	}

	// get owner
	owner, errResp := GetNFTOwnership(ni.TokenId)
	if errResp != nil {
		return errResp
	}

	// check owner and from
	if owner != from {
		errFResp := code.ErrorResp(code.ErrInvalidParam, "from is not nft owner")
		return &errFResp
	}

	return nil
}

func GetNFTInfo(tokenId string) (*types.NFTInfo, *protogo.Response) {

	// get nftInfo
	niBytes, err := sdk.Instance.GetStateByte(_const.KeyNFTInfo, tokenId)
	if err != nil {
		errFResp := code.ErrorResp(code.ErrInternalDataReadFailed, fmt.Sprintf(code.InternalDataReadFailedStrModifier,
			_const.KeyNFTInfo, "", err))
		return nil, &errFResp
	}

	// unmarshal nftInfo
	var ni types.NFTInfo
	err = json.Unmarshal(niBytes, &ni)
	if err != nil {
		errFResp := code.ErrorResp(code.ErrInternalJsonUnmarshalFailed,
			fmt.Sprintf(code.InternalJsonUnmarshalFailedStrModifier, string(niBytes), err))
		return nil, &errFResp
	}

	return &ni, nil
}

func StoreNFTInfo(ni types.NFTInfo) *protogo.Response {

	// marshal
	newNiBytes, err := json.Marshal(ni)
	if err != nil {
		errFResp := code.ErrorResp(code.ErrInternalJsonMarshalFailed,
			fmt.Sprintf(code.InternalJsonMarshalFailedStrModifier, ni, err))
		return &errFResp
	}

	// update nft info
	err = sdk.Instance.PutStateByte(_const.KeyNFTInfo, ni.TokenId, newNiBytes)
	if err != nil {
		errFResp := code.ErrorResp(code.ErrInternalDataStoreFailed, fmt.Sprintf(code.InternalDataStoreFailedStrModifier,
			_const.KeyNFTInfo, ni.TokenId, err))
		return &errFResp
	}

	return nil
}

func GetNFTOwnership(tokenId string) (string, *protogo.Response) {
	owner, err := sdk.Instance.GetState(_const.KeyNFTOwnership, tokenId)
	if err != nil {
		errFResp := code.ErrorResp(code.ErrInternalDataReadFailed, fmt.Sprintf(code.InternalDataReadFailedStrModifier,
			_const.KeyNFTOwnership, tokenId, err))
		return "", &errFResp
	}

	return owner, nil
}

func StoreNFTOwnership(tokenId, owner string) *protogo.Response {
	err := sdk.Instance.PutStateByte(_const.KeyNFTOwnership, tokenId, []byte(owner))
	if err != nil {
		errFResp := code.ErrorResp(code.ErrInternalDataStoreFailed, fmt.Sprintf(code.InternalDataStoreFailedStrModifier,
			_const.KeyNFTOwnership, tokenId, err))
		return &errFResp
	}

	return nil
}

func SetNFTOwnershipToZero(tokenId string) *protogo.Response {
	err := sdk.Instance.PutStateByte(_const.KeyNFTOwnership, tokenId, []byte(_const.ZeroAddress))
	if err != nil {
		errFResp := code.ErrorResp(code.ErrInternalDataStoreFailed, fmt.Sprintf(code.InternalDataStoreFailedStrModifier,
			_const.KeyNFTOwnership, tokenId, err))
		return &errFResp
	}

	return nil
}

func StoreContractCreator() *protogo.Response {

	// 获取合约创建者
	origin, err := sdk.Instance.Origin()
	if err != nil {
		errFResp := code.ErrorResp(code.ErrGetOriginFailed, "")
		return &errFResp
	}

	// 记录合约创建者
	if err = sdk.Instance.PutStateByte(_const.KeyNFTContractCreator, "", []byte(origin)); err != nil {
		errFResp := code.ErrorResp(code.ErrInternalDataStoreFailed, fmt.Sprintf(code.InternalDataStoreFailedStrModifier,
			_const.KeyNFTContractCreator, "", err))
		return &errFResp
	}

	return nil
}
