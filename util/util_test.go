package util

import "testing"

func TestGetContractId(t *testing.T) {
	originHash := "asasdfgfafafg"
	data := "hello"
	tokenId := GenTokenId(originHash, data)
	if tokenId != "2a10acdc66511d6a09c25c62311f93cdcdd53d03b9e14da98546e9b9fe34b846" {
		t.Errorf("token id is wrong, expected: %s, actual: %s", "1615623753565630000", tokenId)
	} else {
		t.Logf("token id: %s", tokenId)
	}
}
