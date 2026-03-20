package types

import (
	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
)

// NFT 数字资产接口定义
type NFT interface {

	// Mint 创建本地 NFT，业务侧调用
	Mint(ni NFTInfo) protogo.Response

	// Transfer 本地发起转移 EBL NFT，业务侧调用
	Transfer(tokenId, from, to string) protogo.Response

	// CrossChainTransfer 跨链转移 EBL NFT，业务侧调用
	CrossChainTransfer(tokenId, from, to string) protogo.Response

	// CrossChainMint 监听到 CrossChainTransfer 事件后，创建跨链 NFT，跨链服务调用
	CrossChainMint(ni NFTInfo) protogo.Response

	// UpdateCrossChainStatus 更新跨链转移状态，跨链服务调用
	UpdateCrossChainStatus(tokenId string, state CrossState) protogo.Response

	// QueryContractVersion 查询合约版本
	QueryContractVersion() protogo.Response
}
