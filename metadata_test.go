package nidhi_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/srikrsna/nidhi"
)

func TestActivityLog_MarshalDocument(t *testing.T) {
	log := &nidhi.ActivityLog{
		On: time.Now(),
		By: uuid.New().String(),
	}

	exp, _ := json.Marshal(log)

	w := jsoniter.ConfigDefault.BorrowStream(nil)
	defer jsoniter.ConfigDefault.ReturnStream(w)

	if err := log.MarshalDocument(w); err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(exp, w.Buffer()) {
		t.Errorf("output mismatch, exp: %s, act: %s", exp, w.Buffer())
	}
}

func TestActivityLog_UnmarshalDocument(t *testing.T) {

	buf, _ := json.Marshal(&nidhi.ActivityLog{
		On: time.Now(),
		By: uuid.New().String(),
	})

	var (
		exp nidhi.ActivityLog
		act nidhi.ActivityLog
	)

	if err := json.Unmarshal(buf, &exp); err != nil {
		t.Fatalf("error in json: %v", err)
	}

	r := jsoniter.ConfigDefault.BorrowIterator(buf)
	defer jsoniter.ConfigDefault.ReturnIterator(r)

	if err := act.UnmarshalDocument(r); err != nil {
		t.Fatal(err)
		return
	}

	if !reflect.DeepEqual(exp, act) {
		t.Fatalf("output mismatch, exp: %v, act: %v", exp, act)
	}
}

