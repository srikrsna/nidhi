package activitylogs_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"

	activitylogs "github.com/srikrsna/nidhi/metadata/activity-logs"
)

func TestMetadata_MarshalMetadata(t *testing.T) {
	md := &activitylogs.Metadata{Created: &activitylogs.ActivityLog{}}

	w := jsoniter.ConfigDefault.BorrowStream(nil)
	defer jsoniter.ConfigDefault.ReturnStream(w)
	w.WriteObjectStart()
	if err := md.MarshalMetadata(w); err != nil {
		t.Fatal(err)
	}
	w.WriteObjectEnd()

	if !json.Valid(w.Buffer()) {
		t.Fatalf("invalid json: %s", string(w.Buffer()))
	}
}

func TestActivityLog_MarshalDocument(t *testing.T) {
	log := &activitylogs.ActivityLog{
		On: time.Now(),
		By: uuid.New().String(),
	}

	w := jsoniter.ConfigDefault.BorrowStream(nil)
	defer jsoniter.ConfigDefault.ReturnStream(w)

	if err := log.MarshalDocument(w); err != nil {
		t.Error(err)
		return
	}

	var act activitylogs.ActivityLog
	if err := json.Unmarshal(w.Buffer(), &act); err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(log, &act) {
		t.Errorf("output mismatch, exp: %v, act: %v", log, &act)
	}
}

func TestActivityLog_UnmarshalDocument(t *testing.T) {
	buf, _ := json.Marshal(&activitylogs.ActivityLog{
		On: time.Now(),
		By: uuid.New().String(),
	})

	var (
		exp activitylogs.ActivityLog
		act activitylogs.ActivityLog
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

	if !cmp.Equal(exp, act) {
		t.Fatalf("output mismatch, exp: %v, act: %v", exp, act)
	}
}
