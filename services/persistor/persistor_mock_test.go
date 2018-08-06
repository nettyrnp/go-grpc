package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	myproto "github.com/nettyrnp/go-grpc/proto"
	"golang.org/x/net/context"
)

// rpcMsg implements the gomock.Matcher interface
type rpcMsg struct {
	msg proto.Message
}

func (r *rpcMsg) Matches(msg interface{}) bool {
	m, ok := msg.(proto.Message)
	if !ok {
		return false
	}
	return proto.Equal(m, r.msg)
}

func (r *rpcMsg) String() string {
	return fmt.Sprintf("is %s", r.msg)
}

func TestSavePersons(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPersistor := NewMockPersistor(ctrl)
	req := &myproto.PersonsReq{
		Persons: []*myproto.Person{
			{
				Id:           123,
				Name:         "name1",
				Email:        "ee@ee.ee",
				MobileNumber: "08730948608",
			},
			{
				Id:           234,
				Name:         "name2",
				Email:        "ee@ee.ch",
				MobileNumber: "09730948608",
			},
		},
	}

	mockPersistor.EXPECT().SavePersons(
		gomock.Any(),
		&rpcMsg{msg: req},
	).Return(&myproto.PersonsReply{
		CreatedCount: 22,
		UpdatedCount: 4,
	}, nil)
	testSavePersons(t, mockPersistor)
}

func testSavePersons(t *testing.T, client myproto.PersistorClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.SavePersons(ctx, &myproto.PersonsReq{
		Persons: []*myproto.Person{
			{
				Id:           123,
				Name:         "name1",
				Email:        "ee@ee.ee",
				MobileNumber: "08730948608",
			},
			{
				Id:           234,
				Name:         "name2",
				Email:        "ee@ee.ch",
				MobileNumber: "09730948608",
			},
		}})
	if err != nil || r.CreatedCount != 22 || r.UpdatedCount != 4 {
		t.Errorf("mocking failed")
	}
	t.Log("Reply : ", r)
}
