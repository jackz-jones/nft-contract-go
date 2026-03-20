package types

import "time"

// CrossState 跨链状态
type CrossState int

const (
	CrossState_Success CrossState = iota // 成功状态
	CrossState_Failed                    // 失败状态
)

// NFTState NFT 状态
type NFTState int

const (
	NFTState_Unlock NFTState = iota // 解锁状态
	NFTState_Lock                   // 锁住状态
)

// NFTInfo NFT 信息
type NFTInfo struct {
	ID         string    `json:"id"`         // 电子单源文件id
	Owner      string    `json:"owner"`      // 出口或进口方地址
	Holder     string    `json:"holder"`     // 承运方地址
	Sender     string    `json:"sender"`     // 目标链交易签名者地址
	TokenId    string    `json:"tokenId"`    // 跨链资产id，hash(OriginHash+hash(data))，合约里计算返回
	OriginHash string    `json:"originHash"` // 业务上链数据的 hash 值
	State      NFTState  `json:"state"`      // 状态，0-解锁状态，1-锁住状态
	Data       string    `json:"data"`       // 附加信息
	CreatedAt  time.Time `json:"createdAt"`  // 创建时间
}

// MessageType 消息类型
type MessageType int

const (
	MessageType_CrossChainTransferEvent MessageType = 500 + iota // 跨链转移消息
	MessageType_CrossChainMintEvent                              // 跨链增发消息
)
