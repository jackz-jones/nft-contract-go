package _const

const ZeroAddress = "0x0000000000000000000000000000000000000000"

// 一些参数名称定义
const (
	ParamNFTInfo = "nftInfo"
	ParamTokenId = "tokenId"
	ParamFrom    = "from"
	ParamTo      = "to"
	ParamState   = "state"
)

// 一些写状态的 key 定义
const (
	KeyNFTContractCreator = "NFTContractCreator" // 合约创建者
	KeyNFTContractVersion = "NFTContractVersion" // 合约版本
	KeyNFTInfo            = "NFTInfo"            // NFT 信息
	KeyNFTOwnership       = "NFTOwnership"       // NFT 所有权
)

// 一些合约方法名定义
const (
	MethodMint                   = "Mint"
	MethodTransfer               = "Transfer"
	MethodCrossChainTransfer     = "CrossChainTransfer"
	MethodCrossChainMint         = "CrossChainMint"
	MethodUpdateCrossChainStatus = "UpdateCrossChainStatus"
	MethodQueryContractVersion   = "QueryContractVersion"
)
