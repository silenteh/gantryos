// Code generated by protoc-gen-gogo.
// source: envelope.proto
// DO NOT EDIT!

/*
	Package proto is a generated protocol buffer package.

	It is generated from these files:
		envelope.proto
		gantryos.proto
		messages.proto

	It has these top-level messages:
		Envelope
*/
package proto

import proto1 "github.com/gogo/protobuf/proto"
import math "math"

// discarding unused import gogoproto "github.com/gogo/protobuf/gogoproto/gogo.pb"

import io2 "io"
import fmt8 "fmt"
import github_com_gogo_protobuf_proto4 "github.com/gogo/protobuf/proto"

import fmt9 "fmt"
import strings4 "strings"
import reflect4 "reflect"

import fmt10 "fmt"
import strings5 "strings"
import github_com_gogo_protobuf_proto5 "github.com/gogo/protobuf/proto"
import sort2 "sort"
import strconv2 "strconv"
import reflect5 "reflect"

import fmt11 "fmt"
import bytes2 "bytes"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = math.Inf

// Basic messages type so that we can match them easily on both side of the channel
type MessageType int32

const (
	// Send an heartbit from the slave to the master
	// if an heartbit is not received for N seconds, then we consider the slave gone
	MessageType_HEARTBIT MessageType = 0
	// Allow the master to acknowledge the slave
	MessageType_ACK_HEARTBIT MessageType = 1
	// The slave ask to be added to the pool
	MessageType_SLAVE_SUBSCRIBE_REQUEST MessageType = 2
	// Acknowledge that the slave was added
	MessageType_ACK_SLAVE_SUBSCRIBE_REQUEST MessageType = 3
	// the slave offers resources to the master
	MessageType_RESOURCE_OFFER MessageType = 4
	// the slave sends updates abut the resource usage to update the master stats ()
	MessageType_RESOURCE_USAGE MessageType = 5
	// this is the request sent from the master to start a task on a slave
	MessageType_TASK_REQUEST MessageType = 6
	// this is the status of the TASK send from the slaves to the master
	// this is sent everytime a task changes its status on the slaves
	MessageType_TASK_STATUS MessageType = 7
	// this gives information about the slave (hostname for ex)
	MessageType_SLAVE_INFO MessageType = 8
)

var MessageType_name = map[int32]string{
	0: "HEARTBIT",
	1: "ACK_HEARTBIT",
	2: "SLAVE_SUBSCRIBE_REQUEST",
	3: "ACK_SLAVE_SUBSCRIBE_REQUEST",
	4: "RESOURCE_OFFER",
	5: "RESOURCE_USAGE",
	6: "TASK_REQUEST",
	7: "TASK_STATUS",
	8: "SLAVE_INFO",
}
var MessageType_value = map[string]int32{
	"HEARTBIT":                    0,
	"ACK_HEARTBIT":                1,
	"SLAVE_SUBSCRIBE_REQUEST":     2,
	"ACK_SLAVE_SUBSCRIBE_REQUEST": 3,
	"RESOURCE_OFFER":              4,
	"RESOURCE_USAGE":              5,
	"TASK_REQUEST":                6,
	"TASK_STATUS":                 7,
	"SLAVE_INFO":                  8,
}

func (x MessageType) Enum() *MessageType {
	p := new(MessageType)
	*p = x
	return p
}
func (x MessageType) String() string {
	return proto1.EnumName(MessageType_name, int32(x))
}
func (x *MessageType) UnmarshalJSON(data []byte) error {
	value, err := proto1.UnmarshalJSONEnum(MessageType_value, data, "MessageType")
	if err != nil {
		return err
	}
	*x = MessageType(value)
	return nil
}

type Envelope struct {
	SenderId      *string        `protobuf:"bytes,1,opt,name=sender_id" json:"sender_id,omitempty"`
	DestinationId *string        `protobuf:"bytes,2,opt,name=destination_id" json:"destination_id,omitempty"`
	ResourceOffer *ResourceOffer `protobuf:"bytes,3,opt,name=resource_offer" json:"resource_offer,omitempty"`
	TaskInfo      *TaskInfo      `protobuf:"bytes,4,opt,name=task_info" json:"task_info,omitempty"`
	TaskStatus    *TaskStatus    `protobuf:"bytes,5,opt,name=task_status" json:"task_status,omitempty"`
	MasterInfo    *MasterInfo    `protobuf:"bytes,6,opt,name=master_info" json:"master_info,omitempty"`
	SlaveInfo     *SlaveInfo     `protobuf:"bytes,7,opt,name=slave_info" json:"slave_info,omitempty"`
	Request       *Request       `protobuf:"bytes,8,opt,name=request" json:"request,omitempty"`
	// Tasks
	RunTask  *RunTaskMessage  `protobuf:"bytes,9,opt,name=run_task" json:"run_task,omitempty"`
	KillTask *KillTaskMessage `protobuf:"bytes,10,opt,name=kill_task" json:"kill_task,omitempty"`
	// messages
	RegisterSlave     *RegisterSlaveMessage   `protobuf:"bytes,11,opt,name=register_slave" json:"register_slave,omitempty"`
	ReRegisterSlave   *ReregisterSlaveMessage `protobuf:"bytes,12,opt,name=re_register_slave" json:"re_register_slave,omitempty"`
	SlaveReRegistered *SlaveRegisteredMessage `protobuf:"bytes,13,opt,name=slave_re_registered" json:"slave_re_registered,omitempty"`
	UnregisterSlave   *UnregisterSlaveMessage `protobuf:"bytes,14,opt,name=unregister_slave" json:"unregister_slave,omitempty"`
	Heartbeat         *HeartbeatMessage       `protobuf:"bytes,15,opt,name=heartbeat" json:"heartbeat,omitempty"`
	ReconcileTasks    *ReconcileTasksMessage  `protobuf:"bytes,16,opt,name=reconcile_tasks" json:"reconcile_tasks,omitempty"`
	LostSlave         *LostSlaveMessage       `protobuf:"bytes,17,opt,name=lost_slave" json:"lost_slave,omitempty"`
	XXX_unrecognized  []byte                  `json:"-"`
}

func (m *Envelope) Reset()      { *m = Envelope{} }
func (*Envelope) ProtoMessage() {}

func (m *Envelope) GetSenderId() string {
	if m != nil && m.SenderId != nil {
		return *m.SenderId
	}
	return ""
}

func (m *Envelope) GetDestinationId() string {
	if m != nil && m.DestinationId != nil {
		return *m.DestinationId
	}
	return ""
}

func (m *Envelope) GetResourceOffer() *ResourceOffer {
	if m != nil {
		return m.ResourceOffer
	}
	return nil
}

func (m *Envelope) GetTaskInfo() *TaskInfo {
	if m != nil {
		return m.TaskInfo
	}
	return nil
}

func (m *Envelope) GetTaskStatus() *TaskStatus {
	if m != nil {
		return m.TaskStatus
	}
	return nil
}

func (m *Envelope) GetMasterInfo() *MasterInfo {
	if m != nil {
		return m.MasterInfo
	}
	return nil
}

func (m *Envelope) GetSlaveInfo() *SlaveInfo {
	if m != nil {
		return m.SlaveInfo
	}
	return nil
}

func (m *Envelope) GetRequest() *Request {
	if m != nil {
		return m.Request
	}
	return nil
}

func (m *Envelope) GetRunTask() *RunTaskMessage {
	if m != nil {
		return m.RunTask
	}
	return nil
}

func (m *Envelope) GetKillTask() *KillTaskMessage {
	if m != nil {
		return m.KillTask
	}
	return nil
}

func (m *Envelope) GetRegisterSlave() *RegisterSlaveMessage {
	if m != nil {
		return m.RegisterSlave
	}
	return nil
}

func (m *Envelope) GetReRegisterSlave() *ReregisterSlaveMessage {
	if m != nil {
		return m.ReRegisterSlave
	}
	return nil
}

func (m *Envelope) GetSlaveReRegistered() *SlaveRegisteredMessage {
	if m != nil {
		return m.SlaveReRegistered
	}
	return nil
}

func (m *Envelope) GetUnregisterSlave() *UnregisterSlaveMessage {
	if m != nil {
		return m.UnregisterSlave
	}
	return nil
}

func (m *Envelope) GetHeartbeat() *HeartbeatMessage {
	if m != nil {
		return m.Heartbeat
	}
	return nil
}

func (m *Envelope) GetReconcileTasks() *ReconcileTasksMessage {
	if m != nil {
		return m.ReconcileTasks
	}
	return nil
}

func (m *Envelope) GetLostSlave() *LostSlaveMessage {
	if m != nil {
		return m.LostSlave
	}
	return nil
}

func init() {
	proto1.RegisterEnum("proto.MessageType", MessageType_name, MessageType_value)
}
func (m *Envelope) Unmarshal(data []byte) error {
	l := len(data)
	index := 0
	for index < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if index >= l {
				return io2.ErrUnexpectedEOF
			}
			b := data[index]
			index++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt8.Errorf("proto: wrong wireType = %d for field SenderId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io2.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + int(stringLen)
			if postIndex > l {
				return io2.ErrUnexpectedEOF
			}
			s := string(data[index:postIndex])
			m.SenderId = &s
			index = postIndex
		case 2:
			if wireType != 2 {
				return fmt8.Errorf("proto: wrong wireType = %d for field DestinationId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io2.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + int(stringLen)
			if postIndex > l {
				return io2.ErrUnexpectedEOF
			}
			s := string(data[index:postIndex])
			m.DestinationId = &s
			index = postIndex
		case 3:
			if wireType != 2 {
				return fmt8.Errorf("proto: wrong wireType = %d for field ResourceOffer", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io2.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + msglen
			if postIndex > l {
				return io2.ErrUnexpectedEOF
			}
			if m.ResourceOffer == nil {
				m.ResourceOffer = &ResourceOffer{}
			}
			if err := m.ResourceOffer.Unmarshal(data[index:postIndex]); err != nil {
				return err
			}
			index = postIndex
		case 4:
			if wireType != 2 {
				return fmt8.Errorf("proto: wrong wireType = %d for field TaskInfo", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io2.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + msglen
			if postIndex > l {
				return io2.ErrUnexpectedEOF
			}
			if m.TaskInfo == nil {
				m.TaskInfo = &TaskInfo{}
			}
			if err := m.TaskInfo.Unmarshal(data[index:postIndex]); err != nil {
				return err
			}
			index = postIndex
		case 5:
			if wireType != 2 {
				return fmt8.Errorf("proto: wrong wireType = %d for field TaskStatus", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io2.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + msglen
			if postIndex > l {
				return io2.ErrUnexpectedEOF
			}
			if m.TaskStatus == nil {
				m.TaskStatus = &TaskStatus{}
			}
			if err := m.TaskStatus.Unmarshal(data[index:postIndex]); err != nil {
				return err
			}
			index = postIndex
		case 6:
			if wireType != 2 {
				return fmt8.Errorf("proto: wrong wireType = %d for field MasterInfo", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io2.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + msglen
			if postIndex > l {
				return io2.ErrUnexpectedEOF
			}
			if m.MasterInfo == nil {
				m.MasterInfo = &MasterInfo{}
			}
			if err := m.MasterInfo.Unmarshal(data[index:postIndex]); err != nil {
				return err
			}
			index = postIndex
		case 7:
			if wireType != 2 {
				return fmt8.Errorf("proto: wrong wireType = %d for field SlaveInfo", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io2.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + msglen
			if postIndex > l {
				return io2.ErrUnexpectedEOF
			}
			if m.SlaveInfo == nil {
				m.SlaveInfo = &SlaveInfo{}
			}
			if err := m.SlaveInfo.Unmarshal(data[index:postIndex]); err != nil {
				return err
			}
			index = postIndex
		case 8:
			if wireType != 2 {
				return fmt8.Errorf("proto: wrong wireType = %d for field Request", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io2.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + msglen
			if postIndex > l {
				return io2.ErrUnexpectedEOF
			}
			if m.Request == nil {
				m.Request = &Request{}
			}
			if err := m.Request.Unmarshal(data[index:postIndex]); err != nil {
				return err
			}
			index = postIndex
		case 9:
			if wireType != 2 {
				return fmt8.Errorf("proto: wrong wireType = %d for field RunTask", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io2.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + msglen
			if postIndex > l {
				return io2.ErrUnexpectedEOF
			}
			if m.RunTask == nil {
				m.RunTask = &RunTaskMessage{}
			}
			if err := m.RunTask.Unmarshal(data[index:postIndex]); err != nil {
				return err
			}
			index = postIndex
		case 10:
			if wireType != 2 {
				return fmt8.Errorf("proto: wrong wireType = %d for field KillTask", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io2.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + msglen
			if postIndex > l {
				return io2.ErrUnexpectedEOF
			}
			if m.KillTask == nil {
				m.KillTask = &KillTaskMessage{}
			}
			if err := m.KillTask.Unmarshal(data[index:postIndex]); err != nil {
				return err
			}
			index = postIndex
		case 11:
			if wireType != 2 {
				return fmt8.Errorf("proto: wrong wireType = %d for field RegisterSlave", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io2.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + msglen
			if postIndex > l {
				return io2.ErrUnexpectedEOF
			}
			if m.RegisterSlave == nil {
				m.RegisterSlave = &RegisterSlaveMessage{}
			}
			if err := m.RegisterSlave.Unmarshal(data[index:postIndex]); err != nil {
				return err
			}
			index = postIndex
		case 12:
			if wireType != 2 {
				return fmt8.Errorf("proto: wrong wireType = %d for field ReRegisterSlave", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io2.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + msglen
			if postIndex > l {
				return io2.ErrUnexpectedEOF
			}
			if m.ReRegisterSlave == nil {
				m.ReRegisterSlave = &ReregisterSlaveMessage{}
			}
			if err := m.ReRegisterSlave.Unmarshal(data[index:postIndex]); err != nil {
				return err
			}
			index = postIndex
		case 13:
			if wireType != 2 {
				return fmt8.Errorf("proto: wrong wireType = %d for field SlaveReRegistered", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io2.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + msglen
			if postIndex > l {
				return io2.ErrUnexpectedEOF
			}
			if m.SlaveReRegistered == nil {
				m.SlaveReRegistered = &SlaveRegisteredMessage{}
			}
			if err := m.SlaveReRegistered.Unmarshal(data[index:postIndex]); err != nil {
				return err
			}
			index = postIndex
		case 14:
			if wireType != 2 {
				return fmt8.Errorf("proto: wrong wireType = %d for field UnregisterSlave", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io2.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + msglen
			if postIndex > l {
				return io2.ErrUnexpectedEOF
			}
			if m.UnregisterSlave == nil {
				m.UnregisterSlave = &UnregisterSlaveMessage{}
			}
			if err := m.UnregisterSlave.Unmarshal(data[index:postIndex]); err != nil {
				return err
			}
			index = postIndex
		case 15:
			if wireType != 2 {
				return fmt8.Errorf("proto: wrong wireType = %d for field Heartbeat", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io2.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + msglen
			if postIndex > l {
				return io2.ErrUnexpectedEOF
			}
			if m.Heartbeat == nil {
				m.Heartbeat = &HeartbeatMessage{}
			}
			if err := m.Heartbeat.Unmarshal(data[index:postIndex]); err != nil {
				return err
			}
			index = postIndex
		case 16:
			if wireType != 2 {
				return fmt8.Errorf("proto: wrong wireType = %d for field ReconcileTasks", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io2.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + msglen
			if postIndex > l {
				return io2.ErrUnexpectedEOF
			}
			if m.ReconcileTasks == nil {
				m.ReconcileTasks = &ReconcileTasksMessage{}
			}
			if err := m.ReconcileTasks.Unmarshal(data[index:postIndex]); err != nil {
				return err
			}
			index = postIndex
		case 17:
			if wireType != 2 {
				return fmt8.Errorf("proto: wrong wireType = %d for field LostSlave", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io2.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + msglen
			if postIndex > l {
				return io2.ErrUnexpectedEOF
			}
			if m.LostSlave == nil {
				m.LostSlave = &LostSlaveMessage{}
			}
			if err := m.LostSlave.Unmarshal(data[index:postIndex]); err != nil {
				return err
			}
			index = postIndex
		default:
			var sizeOfWire int
			for {
				sizeOfWire++
				wire >>= 7
				if wire == 0 {
					break
				}
			}
			index -= sizeOfWire
			skippy, err := github_com_gogo_protobuf_proto4.Skip(data[index:])
			if err != nil {
				return err
			}
			if (index + skippy) > l {
				return io2.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, data[index:index+skippy]...)
			index += skippy
		}
	}
	return nil
}
func (this *Envelope) String() string {
	if this == nil {
		return "nil"
	}
	s := strings4.Join([]string{`&Envelope{`,
		`SenderId:` + valueToStringEnvelope(this.SenderId) + `,`,
		`DestinationId:` + valueToStringEnvelope(this.DestinationId) + `,`,
		`ResourceOffer:` + strings4.Replace(fmt9.Sprintf("%v", this.ResourceOffer), "ResourceOffer", "ResourceOffer", 1) + `,`,
		`TaskInfo:` + strings4.Replace(fmt9.Sprintf("%v", this.TaskInfo), "TaskInfo", "TaskInfo", 1) + `,`,
		`TaskStatus:` + strings4.Replace(fmt9.Sprintf("%v", this.TaskStatus), "TaskStatus", "TaskStatus", 1) + `,`,
		`MasterInfo:` + strings4.Replace(fmt9.Sprintf("%v", this.MasterInfo), "MasterInfo", "MasterInfo", 1) + `,`,
		`SlaveInfo:` + strings4.Replace(fmt9.Sprintf("%v", this.SlaveInfo), "SlaveInfo", "SlaveInfo", 1) + `,`,
		`Request:` + strings4.Replace(fmt9.Sprintf("%v", this.Request), "Request", "Request", 1) + `,`,
		`RunTask:` + strings4.Replace(fmt9.Sprintf("%v", this.RunTask), "RunTaskMessage", "RunTaskMessage", 1) + `,`,
		`KillTask:` + strings4.Replace(fmt9.Sprintf("%v", this.KillTask), "KillTaskMessage", "KillTaskMessage", 1) + `,`,
		`RegisterSlave:` + strings4.Replace(fmt9.Sprintf("%v", this.RegisterSlave), "RegisterSlaveMessage", "RegisterSlaveMessage", 1) + `,`,
		`ReRegisterSlave:` + strings4.Replace(fmt9.Sprintf("%v", this.ReRegisterSlave), "ReregisterSlaveMessage", "ReregisterSlaveMessage", 1) + `,`,
		`SlaveReRegistered:` + strings4.Replace(fmt9.Sprintf("%v", this.SlaveReRegistered), "SlaveRegisteredMessage", "SlaveRegisteredMessage", 1) + `,`,
		`UnregisterSlave:` + strings4.Replace(fmt9.Sprintf("%v", this.UnregisterSlave), "UnregisterSlaveMessage", "UnregisterSlaveMessage", 1) + `,`,
		`Heartbeat:` + strings4.Replace(fmt9.Sprintf("%v", this.Heartbeat), "HeartbeatMessage", "HeartbeatMessage", 1) + `,`,
		`ReconcileTasks:` + strings4.Replace(fmt9.Sprintf("%v", this.ReconcileTasks), "ReconcileTasksMessage", "ReconcileTasksMessage", 1) + `,`,
		`LostSlave:` + strings4.Replace(fmt9.Sprintf("%v", this.LostSlave), "LostSlaveMessage", "LostSlaveMessage", 1) + `,`,
		`XXX_unrecognized:` + fmt9.Sprintf("%v", this.XXX_unrecognized) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringEnvelope(v interface{}) string {
	rv := reflect4.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect4.Indirect(rv).Interface()
	return fmt9.Sprintf("*%v", pv)
}
func (m *Envelope) Size() (n int) {
	var l int
	_ = l
	if m.SenderId != nil {
		l = len(*m.SenderId)
		n += 1 + l + sovEnvelope(uint64(l))
	}
	if m.DestinationId != nil {
		l = len(*m.DestinationId)
		n += 1 + l + sovEnvelope(uint64(l))
	}
	if m.ResourceOffer != nil {
		l = m.ResourceOffer.Size()
		n += 1 + l + sovEnvelope(uint64(l))
	}
	if m.TaskInfo != nil {
		l = m.TaskInfo.Size()
		n += 1 + l + sovEnvelope(uint64(l))
	}
	if m.TaskStatus != nil {
		l = m.TaskStatus.Size()
		n += 1 + l + sovEnvelope(uint64(l))
	}
	if m.MasterInfo != nil {
		l = m.MasterInfo.Size()
		n += 1 + l + sovEnvelope(uint64(l))
	}
	if m.SlaveInfo != nil {
		l = m.SlaveInfo.Size()
		n += 1 + l + sovEnvelope(uint64(l))
	}
	if m.Request != nil {
		l = m.Request.Size()
		n += 1 + l + sovEnvelope(uint64(l))
	}
	if m.RunTask != nil {
		l = m.RunTask.Size()
		n += 1 + l + sovEnvelope(uint64(l))
	}
	if m.KillTask != nil {
		l = m.KillTask.Size()
		n += 1 + l + sovEnvelope(uint64(l))
	}
	if m.RegisterSlave != nil {
		l = m.RegisterSlave.Size()
		n += 1 + l + sovEnvelope(uint64(l))
	}
	if m.ReRegisterSlave != nil {
		l = m.ReRegisterSlave.Size()
		n += 1 + l + sovEnvelope(uint64(l))
	}
	if m.SlaveReRegistered != nil {
		l = m.SlaveReRegistered.Size()
		n += 1 + l + sovEnvelope(uint64(l))
	}
	if m.UnregisterSlave != nil {
		l = m.UnregisterSlave.Size()
		n += 1 + l + sovEnvelope(uint64(l))
	}
	if m.Heartbeat != nil {
		l = m.Heartbeat.Size()
		n += 1 + l + sovEnvelope(uint64(l))
	}
	if m.ReconcileTasks != nil {
		l = m.ReconcileTasks.Size()
		n += 2 + l + sovEnvelope(uint64(l))
	}
	if m.LostSlave != nil {
		l = m.LostSlave.Size()
		n += 2 + l + sovEnvelope(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovEnvelope(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozEnvelope(x uint64) (n int) {
	return sovEnvelope(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func NewPopulatedEnvelope(r randyEnvelope, easy bool) *Envelope {
	this := &Envelope{}
	if r.Intn(10) != 0 {
		v1 := randStringEnvelope(r)
		this.SenderId = &v1
	}
	if r.Intn(10) != 0 {
		v2 := randStringEnvelope(r)
		this.DestinationId = &v2
	}
	if r.Intn(10) != 0 {
		this.ResourceOffer = NewPopulatedResourceOffer(r, easy)
	}
	if r.Intn(10) != 0 {
		this.TaskInfo = NewPopulatedTaskInfo(r, easy)
	}
	if r.Intn(10) != 0 {
		this.TaskStatus = NewPopulatedTaskStatus(r, easy)
	}
	if r.Intn(10) != 0 {
		this.MasterInfo = NewPopulatedMasterInfo(r, easy)
	}
	if r.Intn(10) != 0 {
		this.SlaveInfo = NewPopulatedSlaveInfo(r, easy)
	}
	if r.Intn(10) != 0 {
		this.Request = NewPopulatedRequest(r, easy)
	}
	if r.Intn(10) != 0 {
		this.RunTask = NewPopulatedRunTaskMessage(r, easy)
	}
	if r.Intn(10) != 0 {
		this.KillTask = NewPopulatedKillTaskMessage(r, easy)
	}
	if r.Intn(10) != 0 {
		this.RegisterSlave = NewPopulatedRegisterSlaveMessage(r, easy)
	}
	if r.Intn(10) != 0 {
		this.ReRegisterSlave = NewPopulatedReregisterSlaveMessage(r, easy)
	}
	if r.Intn(10) != 0 {
		this.SlaveReRegistered = NewPopulatedSlaveRegisteredMessage(r, easy)
	}
	if r.Intn(10) != 0 {
		this.UnregisterSlave = NewPopulatedUnregisterSlaveMessage(r, easy)
	}
	if r.Intn(10) != 0 {
		this.Heartbeat = NewPopulatedHeartbeatMessage(r, easy)
	}
	if r.Intn(10) != 0 {
		this.ReconcileTasks = NewPopulatedReconcileTasksMessage(r, easy)
	}
	if r.Intn(10) != 0 {
		this.LostSlave = NewPopulatedLostSlaveMessage(r, easy)
	}
	if !easy && r.Intn(10) != 0 {
		this.XXX_unrecognized = randUnrecognizedEnvelope(r, 18)
	}
	return this
}

type randyEnvelope interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func randUTF8RuneEnvelope(r randyEnvelope) rune {
	res := rune(r.Uint32() % 1112064)
	if 55296 <= res {
		res += 2047
	}
	return res
}
func randStringEnvelope(r randyEnvelope) string {
	v3 := r.Intn(100)
	tmps := make([]rune, v3)
	for i := 0; i < v3; i++ {
		tmps[i] = randUTF8RuneEnvelope(r)
	}
	return string(tmps)
}
func randUnrecognizedEnvelope(r randyEnvelope, maxFieldNumber int) (data []byte) {
	l := r.Intn(5)
	for i := 0; i < l; i++ {
		wire := r.Intn(4)
		if wire == 3 {
			wire = 5
		}
		fieldNumber := maxFieldNumber + r.Intn(100)
		data = randFieldEnvelope(data, r, fieldNumber, wire)
	}
	return data
}
func randFieldEnvelope(data []byte, r randyEnvelope, fieldNumber int, wire int) []byte {
	key := uint32(fieldNumber)<<3 | uint32(wire)
	switch wire {
	case 0:
		data = encodeVarintPopulateEnvelope(data, uint64(key))
		v4 := r.Int63()
		if r.Intn(2) == 0 {
			v4 *= -1
		}
		data = encodeVarintPopulateEnvelope(data, uint64(v4))
	case 1:
		data = encodeVarintPopulateEnvelope(data, uint64(key))
		data = append(data, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	case 2:
		data = encodeVarintPopulateEnvelope(data, uint64(key))
		ll := r.Intn(100)
		data = encodeVarintPopulateEnvelope(data, uint64(ll))
		for j := 0; j < ll; j++ {
			data = append(data, byte(r.Intn(256)))
		}
	default:
		data = encodeVarintPopulateEnvelope(data, uint64(key))
		data = append(data, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	}
	return data
}
func encodeVarintPopulateEnvelope(data []byte, v uint64) []byte {
	for v >= 1<<7 {
		data = append(data, uint8(uint64(v)&0x7f|0x80))
		v >>= 7
	}
	data = append(data, uint8(v))
	return data
}
func (m *Envelope) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *Envelope) MarshalTo(data []byte) (n int, err error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.SenderId != nil {
		data[i] = 0xa
		i++
		i = encodeVarintEnvelope(data, i, uint64(len(*m.SenderId)))
		i += copy(data[i:], *m.SenderId)
	}
	if m.DestinationId != nil {
		data[i] = 0x12
		i++
		i = encodeVarintEnvelope(data, i, uint64(len(*m.DestinationId)))
		i += copy(data[i:], *m.DestinationId)
	}
	if m.ResourceOffer != nil {
		data[i] = 0x1a
		i++
		i = encodeVarintEnvelope(data, i, uint64(m.ResourceOffer.Size()))
		n1, err := m.ResourceOffer.MarshalTo(data[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.TaskInfo != nil {
		data[i] = 0x22
		i++
		i = encodeVarintEnvelope(data, i, uint64(m.TaskInfo.Size()))
		n2, err := m.TaskInfo.MarshalTo(data[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	if m.TaskStatus != nil {
		data[i] = 0x2a
		i++
		i = encodeVarintEnvelope(data, i, uint64(m.TaskStatus.Size()))
		n3, err := m.TaskStatus.MarshalTo(data[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	if m.MasterInfo != nil {
		data[i] = 0x32
		i++
		i = encodeVarintEnvelope(data, i, uint64(m.MasterInfo.Size()))
		n4, err := m.MasterInfo.MarshalTo(data[i:])
		if err != nil {
			return 0, err
		}
		i += n4
	}
	if m.SlaveInfo != nil {
		data[i] = 0x3a
		i++
		i = encodeVarintEnvelope(data, i, uint64(m.SlaveInfo.Size()))
		n5, err := m.SlaveInfo.MarshalTo(data[i:])
		if err != nil {
			return 0, err
		}
		i += n5
	}
	if m.Request != nil {
		data[i] = 0x42
		i++
		i = encodeVarintEnvelope(data, i, uint64(m.Request.Size()))
		n6, err := m.Request.MarshalTo(data[i:])
		if err != nil {
			return 0, err
		}
		i += n6
	}
	if m.RunTask != nil {
		data[i] = 0x4a
		i++
		i = encodeVarintEnvelope(data, i, uint64(m.RunTask.Size()))
		n7, err := m.RunTask.MarshalTo(data[i:])
		if err != nil {
			return 0, err
		}
		i += n7
	}
	if m.KillTask != nil {
		data[i] = 0x52
		i++
		i = encodeVarintEnvelope(data, i, uint64(m.KillTask.Size()))
		n8, err := m.KillTask.MarshalTo(data[i:])
		if err != nil {
			return 0, err
		}
		i += n8
	}
	if m.RegisterSlave != nil {
		data[i] = 0x5a
		i++
		i = encodeVarintEnvelope(data, i, uint64(m.RegisterSlave.Size()))
		n9, err := m.RegisterSlave.MarshalTo(data[i:])
		if err != nil {
			return 0, err
		}
		i += n9
	}
	if m.ReRegisterSlave != nil {
		data[i] = 0x62
		i++
		i = encodeVarintEnvelope(data, i, uint64(m.ReRegisterSlave.Size()))
		n10, err := m.ReRegisterSlave.MarshalTo(data[i:])
		if err != nil {
			return 0, err
		}
		i += n10
	}
	if m.SlaveReRegistered != nil {
		data[i] = 0x6a
		i++
		i = encodeVarintEnvelope(data, i, uint64(m.SlaveReRegistered.Size()))
		n11, err := m.SlaveReRegistered.MarshalTo(data[i:])
		if err != nil {
			return 0, err
		}
		i += n11
	}
	if m.UnregisterSlave != nil {
		data[i] = 0x72
		i++
		i = encodeVarintEnvelope(data, i, uint64(m.UnregisterSlave.Size()))
		n12, err := m.UnregisterSlave.MarshalTo(data[i:])
		if err != nil {
			return 0, err
		}
		i += n12
	}
	if m.Heartbeat != nil {
		data[i] = 0x7a
		i++
		i = encodeVarintEnvelope(data, i, uint64(m.Heartbeat.Size()))
		n13, err := m.Heartbeat.MarshalTo(data[i:])
		if err != nil {
			return 0, err
		}
		i += n13
	}
	if m.ReconcileTasks != nil {
		data[i] = 0x82
		i++
		data[i] = 0x1
		i++
		i = encodeVarintEnvelope(data, i, uint64(m.ReconcileTasks.Size()))
		n14, err := m.ReconcileTasks.MarshalTo(data[i:])
		if err != nil {
			return 0, err
		}
		i += n14
	}
	if m.LostSlave != nil {
		data[i] = 0x8a
		i++
		data[i] = 0x1
		i++
		i = encodeVarintEnvelope(data, i, uint64(m.LostSlave.Size()))
		n15, err := m.LostSlave.MarshalTo(data[i:])
		if err != nil {
			return 0, err
		}
		i += n15
	}
	if m.XXX_unrecognized != nil {
		i += copy(data[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeFixed64Envelope(data []byte, offset int, v uint64) int {
	data[offset] = uint8(v)
	data[offset+1] = uint8(v >> 8)
	data[offset+2] = uint8(v >> 16)
	data[offset+3] = uint8(v >> 24)
	data[offset+4] = uint8(v >> 32)
	data[offset+5] = uint8(v >> 40)
	data[offset+6] = uint8(v >> 48)
	data[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32Envelope(data []byte, offset int, v uint32) int {
	data[offset] = uint8(v)
	data[offset+1] = uint8(v >> 8)
	data[offset+2] = uint8(v >> 16)
	data[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintEnvelope(data []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		data[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	data[offset] = uint8(v)
	return offset + 1
}
func (this *Envelope) GoString() string {
	if this == nil {
		return "nil"
	}
	s := strings5.Join([]string{`&proto.Envelope{` +
		`SenderId:` + valueToGoStringEnvelope(this.SenderId, "string"),
		`DestinationId:` + valueToGoStringEnvelope(this.DestinationId, "string"),
		`ResourceOffer:` + fmt10.Sprintf("%#v", this.ResourceOffer),
		`TaskInfo:` + fmt10.Sprintf("%#v", this.TaskInfo),
		`TaskStatus:` + fmt10.Sprintf("%#v", this.TaskStatus),
		`MasterInfo:` + fmt10.Sprintf("%#v", this.MasterInfo),
		`SlaveInfo:` + fmt10.Sprintf("%#v", this.SlaveInfo),
		`Request:` + fmt10.Sprintf("%#v", this.Request),
		`RunTask:` + fmt10.Sprintf("%#v", this.RunTask),
		`KillTask:` + fmt10.Sprintf("%#v", this.KillTask),
		`RegisterSlave:` + fmt10.Sprintf("%#v", this.RegisterSlave),
		`ReRegisterSlave:` + fmt10.Sprintf("%#v", this.ReRegisterSlave),
		`SlaveReRegistered:` + fmt10.Sprintf("%#v", this.SlaveReRegistered),
		`UnregisterSlave:` + fmt10.Sprintf("%#v", this.UnregisterSlave),
		`Heartbeat:` + fmt10.Sprintf("%#v", this.Heartbeat),
		`ReconcileTasks:` + fmt10.Sprintf("%#v", this.ReconcileTasks),
		`LostSlave:` + fmt10.Sprintf("%#v", this.LostSlave),
		`XXX_unrecognized:` + fmt10.Sprintf("%#v", this.XXX_unrecognized) + `}`}, ", ")
	return s
}
func valueToGoStringEnvelope(v interface{}, typ string) string {
	rv := reflect5.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect5.Indirect(rv).Interface()
	return fmt10.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func extensionToGoStringEnvelope(e map[int32]github_com_gogo_protobuf_proto5.Extension) string {
	if e == nil {
		return "nil"
	}
	s := "map[int32]proto.Extension{"
	keys := make([]int, 0, len(e))
	for k := range e {
		keys = append(keys, int(k))
	}
	sort2.Ints(keys)
	ss := []string{}
	for _, k := range keys {
		ss = append(ss, strconv2.Itoa(k)+": "+e[int32(k)].GoString())
	}
	s += strings5.Join(ss, ",") + "}"
	return s
}
func (this *Envelope) VerboseEqual(that interface{}) error {
	if that == nil {
		if this == nil {
			return nil
		}
		return fmt11.Errorf("that == nil && this != nil")
	}

	that1, ok := that.(*Envelope)
	if !ok {
		return fmt11.Errorf("that is not of type *Envelope")
	}
	if that1 == nil {
		if this == nil {
			return nil
		}
		return fmt11.Errorf("that is type *Envelope but is nil && this != nil")
	} else if this == nil {
		return fmt11.Errorf("that is type *Envelopebut is not nil && this == nil")
	}
	if this.SenderId != nil && that1.SenderId != nil {
		if *this.SenderId != *that1.SenderId {
			return fmt11.Errorf("SenderId this(%v) Not Equal that(%v)", *this.SenderId, *that1.SenderId)
		}
	} else if this.SenderId != nil {
		return fmt11.Errorf("this.SenderId == nil && that.SenderId != nil")
	} else if that1.SenderId != nil {
		return fmt11.Errorf("SenderId this(%v) Not Equal that(%v)", this.SenderId, that1.SenderId)
	}
	if this.DestinationId != nil && that1.DestinationId != nil {
		if *this.DestinationId != *that1.DestinationId {
			return fmt11.Errorf("DestinationId this(%v) Not Equal that(%v)", *this.DestinationId, *that1.DestinationId)
		}
	} else if this.DestinationId != nil {
		return fmt11.Errorf("this.DestinationId == nil && that.DestinationId != nil")
	} else if that1.DestinationId != nil {
		return fmt11.Errorf("DestinationId this(%v) Not Equal that(%v)", this.DestinationId, that1.DestinationId)
	}
	if !this.ResourceOffer.Equal(that1.ResourceOffer) {
		return fmt11.Errorf("ResourceOffer this(%v) Not Equal that(%v)", this.ResourceOffer, that1.ResourceOffer)
	}
	if !this.TaskInfo.Equal(that1.TaskInfo) {
		return fmt11.Errorf("TaskInfo this(%v) Not Equal that(%v)", this.TaskInfo, that1.TaskInfo)
	}
	if !this.TaskStatus.Equal(that1.TaskStatus) {
		return fmt11.Errorf("TaskStatus this(%v) Not Equal that(%v)", this.TaskStatus, that1.TaskStatus)
	}
	if !this.MasterInfo.Equal(that1.MasterInfo) {
		return fmt11.Errorf("MasterInfo this(%v) Not Equal that(%v)", this.MasterInfo, that1.MasterInfo)
	}
	if !this.SlaveInfo.Equal(that1.SlaveInfo) {
		return fmt11.Errorf("SlaveInfo this(%v) Not Equal that(%v)", this.SlaveInfo, that1.SlaveInfo)
	}
	if !this.Request.Equal(that1.Request) {
		return fmt11.Errorf("Request this(%v) Not Equal that(%v)", this.Request, that1.Request)
	}
	if !this.RunTask.Equal(that1.RunTask) {
		return fmt11.Errorf("RunTask this(%v) Not Equal that(%v)", this.RunTask, that1.RunTask)
	}
	if !this.KillTask.Equal(that1.KillTask) {
		return fmt11.Errorf("KillTask this(%v) Not Equal that(%v)", this.KillTask, that1.KillTask)
	}
	if !this.RegisterSlave.Equal(that1.RegisterSlave) {
		return fmt11.Errorf("RegisterSlave this(%v) Not Equal that(%v)", this.RegisterSlave, that1.RegisterSlave)
	}
	if !this.ReRegisterSlave.Equal(that1.ReRegisterSlave) {
		return fmt11.Errorf("ReRegisterSlave this(%v) Not Equal that(%v)", this.ReRegisterSlave, that1.ReRegisterSlave)
	}
	if !this.SlaveReRegistered.Equal(that1.SlaveReRegistered) {
		return fmt11.Errorf("SlaveReRegistered this(%v) Not Equal that(%v)", this.SlaveReRegistered, that1.SlaveReRegistered)
	}
	if !this.UnregisterSlave.Equal(that1.UnregisterSlave) {
		return fmt11.Errorf("UnregisterSlave this(%v) Not Equal that(%v)", this.UnregisterSlave, that1.UnregisterSlave)
	}
	if !this.Heartbeat.Equal(that1.Heartbeat) {
		return fmt11.Errorf("Heartbeat this(%v) Not Equal that(%v)", this.Heartbeat, that1.Heartbeat)
	}
	if !this.ReconcileTasks.Equal(that1.ReconcileTasks) {
		return fmt11.Errorf("ReconcileTasks this(%v) Not Equal that(%v)", this.ReconcileTasks, that1.ReconcileTasks)
	}
	if !this.LostSlave.Equal(that1.LostSlave) {
		return fmt11.Errorf("LostSlave this(%v) Not Equal that(%v)", this.LostSlave, that1.LostSlave)
	}
	if !bytes2.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return fmt11.Errorf("XXX_unrecognized this(%v) Not Equal that(%v)", this.XXX_unrecognized, that1.XXX_unrecognized)
	}
	return nil
}
func (this *Envelope) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*Envelope)
	if !ok {
		return false
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if this.SenderId != nil && that1.SenderId != nil {
		if *this.SenderId != *that1.SenderId {
			return false
		}
	} else if this.SenderId != nil {
		return false
	} else if that1.SenderId != nil {
		return false
	}
	if this.DestinationId != nil && that1.DestinationId != nil {
		if *this.DestinationId != *that1.DestinationId {
			return false
		}
	} else if this.DestinationId != nil {
		return false
	} else if that1.DestinationId != nil {
		return false
	}
	if !this.ResourceOffer.Equal(that1.ResourceOffer) {
		return false
	}
	if !this.TaskInfo.Equal(that1.TaskInfo) {
		return false
	}
	if !this.TaskStatus.Equal(that1.TaskStatus) {
		return false
	}
	if !this.MasterInfo.Equal(that1.MasterInfo) {
		return false
	}
	if !this.SlaveInfo.Equal(that1.SlaveInfo) {
		return false
	}
	if !this.Request.Equal(that1.Request) {
		return false
	}
	if !this.RunTask.Equal(that1.RunTask) {
		return false
	}
	if !this.KillTask.Equal(that1.KillTask) {
		return false
	}
	if !this.RegisterSlave.Equal(that1.RegisterSlave) {
		return false
	}
	if !this.ReRegisterSlave.Equal(that1.ReRegisterSlave) {
		return false
	}
	if !this.SlaveReRegistered.Equal(that1.SlaveReRegistered) {
		return false
	}
	if !this.UnregisterSlave.Equal(that1.UnregisterSlave) {
		return false
	}
	if !this.Heartbeat.Equal(that1.Heartbeat) {
		return false
	}
	if !this.ReconcileTasks.Equal(that1.ReconcileTasks) {
		return false
	}
	if !this.LostSlave.Equal(that1.LostSlave) {
		return false
	}
	if !bytes2.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
