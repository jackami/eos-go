package eos

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestGetTransactionCustom(t *testing.T) {
	client := New("https://history.meet.one")
	resp, err := client.GetTransactionCustom("ac9e5cb24902bee691551a6a294461016b5240a14ac8e0d00592b7958cf4664f")
	if err != nil {

	} else {
		fmt.Println(resp)
	}
}

func TestReadTraceLog(t *testing.T) {

	var maxNum int64 = 0
	logFile, err := os.OpenFile("/data/state-data/trace_history.log", os.O_RDONLY, 0755)
	if err != nil {
		t.Error("blocks.log file open failure.")
	}
	idxFile, err := os.OpenFile("/data/state-data/trace_history.index", os.O_RDONLY, 0755)
	if err != nil {
		t.Error("blocks.index file open failure.")
	}
	defer logFile.Close()
	defer idxFile.Close()

	state, err := idxFile.Stat()
	if err != nil {
		t.Error(err)
	} else {
		size := state.Size()
		if size == 0 {
			t.Error("blocks.index file is empty.")
		}
		if size%8 != 0 {
			t.Error("blocks.index file size not correct.")
		}
		maxNum = size / 8
	}

	t.Log(maxNum)

	var num int64 = 71211 // 1227211

	offsetBuf := make([]byte, 16)
	n, err := idxFile.ReadAt(offsetBuf, (num-1)*8)
	if err != nil || n != 16 {
		t.Error(err)
	}
	start := binary.LittleEndian.Uint64(offsetBuf[:8])
	end := binary.LittleEndian.Uint64(offsetBuf[8:])

	dataSize := end - start
	traceBytes := make([]byte, dataSize)
	n, err = logFile.ReadAt(traceBytes, int64(start))
	if err != nil || n != int(dataSize) {
		t.Error(err)
	} else {
		t.Log(traceBytes)
	}

	header := &EntryHeader{}

	traceHeader := NewDecoder(traceBytes[:56])

	err = traceHeader.Decode(&header)

	//pos := binary.LittleEndian.Uint64(traceBytes[dataSize-8:])
	//t.Log(pos)

	//payload := &EntryPayload{}
	//tracePayload := NewDecoder(traceBytes[56 : dataSize-8])
	//err = tracePayload.Decode(payload)

	var payloadReader io.Reader
	payloadReader = bytes.NewBuffer(traceBytes[60 : dataSize-8])

	payloadReader, err = zlib.NewReader(payloadReader)
	//if err != nil {
	//	fmt.Errorf("new reader for tx, %s\n", err)
	//}
	payloadData, err := ioutil.ReadAll(payloadReader)
	//if err != nil {
	//	fmt.Errorf("unpack read all, %s", err)
	//}

	//list := []TrxTrace{}
	atraces := &TrxTraces{}

	trxTraceDec := NewDecoder(payloadData)
	err = trxTraceDec.Decode(atraces)

	t.Log(atraces.TrxTraces[0].TrxId.String())

	bufStr, _ := json.Marshal(atraces)
	t.Log(string(bufStr))
	//t.Log(traces.TrxTraces[0].TrxTraces[0].ReceiptActionDigest.String())
	t.Log(header)
	t.Log(hex.EncodeToString(traceBytes[8:40]))

}

func TestBlockTime(t *testing.T) {
	blk2Time, err := time.Parse("2006-01-02 15:04:05.000000000", "2018-06-09 19:56:30.000000000")
	if err != nil {
		t.Error(err)
	}

	var blkNum int64 = 1333

	x, e := time.ParseDuration(fmt.Sprintf("%dms", (blkNum-2)*500))
	if e != nil {
		t.Error(e)
	}

	blkTime := blk2Time.Add(x)

	t.Log(blk2Time.UTC().String())
	t.Log(blkTime.UTC().String())
}
