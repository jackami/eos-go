package eos

import (
	"fmt"
	"testing"
)

func TestGetTransactionCustom(t *testing.T) {
	client := New("https://history.meet.one")
	resp, err := client.GetTransactionCustom("ac9e5cb24902bee691551a6a294461016b5240a14ac8e0d00592b7958cf4664f")
	if err != nil {

	} else {
		fmt.Println(resp)
	}
}