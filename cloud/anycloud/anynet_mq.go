package anycloud

import (
	"fmt"
	"goutil/basic/gerrors"
	"goutil/basic/glog"
	"goutil/container/gqueue"
	"goutil/dsa/guuid"
	"goutil/net/gmsg"
	"net"
	"time"
)

type (
	Request struct {
		Func   string
		Params map[string]any
	}

	Reply map[string]any

	ParamChecker struct {
		InRequirements  map[string]map[string]any
		OutRequirements map[string]map[string]any
	}

	MQ struct {
		anyNet         *AnyNet
		msgConn        gmsg.Conn
		targetNodeID   string
		targetNodeConn net.Conn
		cachedPop      map[string]*gqueue.Queue // map[queue]*gqueue.Queue
		cachedRequest  map[string]chan *Msg     // map[service|requestID]chan *Msg
		cachedSub      map[string]SubCallback
		cachedReply    map[string]ReplyCallback
		checker        ParamChecker
	}

	SubCallback   func(subj string, b []byte)
	ReplyCallback func(service string, b []byte) ([]byte, error)
)

const (
	mqMaxCached = 1024
)

func (mq *MQ) Push(queue string, data []byte) error {
	return mq.sndMsg(CommCmdPush, queue, data)
}

func (mq *MQ) Pop(queue string, waitTimeout *time.Duration) ([]byte, error) {
	if mq.cachedPop[queue] == nil {
		mq.cachedPop[queue] = gqueue.NewQueue()
	}
	msg := mq.cachedPop[queue].Pop(waitTimeout)
	return msg.(*Msg).Data, nil
}

func (mq *MQ) Pub(subj string, data []byte) error {
	return mq.sndMsg(CommCmdPub, subj, data)
}

func (mq *MQ) Sub(subj string, callback SubCallback) {
	mq.cachedSub[subj] = callback
}

func (mq *MQ) UnSub(subj string) {
	delete(mq.cachedSub, subj)
}

func (mq *MQ) Request(service string, data []byte, waitTimeout *time.Duration) ([]byte, error) {
	requestID := guuid.NewString(false, false)
	if mq.cachedRequest[requestID] == nil {
		mq.cachedRequest[requestID] = make(chan *Msg, 1)
	}
	if err := mq.sndMsg(CommCmdRequest, fmt.Sprintf("%s|%s", service, requestID), data); err != nil {
		delete(mq.cachedReply, requestID)
		return nil, err
	}
	if waitTimeout == nil {
		reply := <-mq.cachedRequest[requestID]
		return reply.Data, nil
	} else {
		ticker := time.NewTicker(*waitTimeout)
		select {
		case <-ticker.C:
			return nil, gerrors.ErrTimeout
		case reply := <-mq.cachedRequest[requestID]:
			return reply.Data, nil
		}
	}
}

func (mq *MQ) Reply(service string, callback ReplyCallback) {
	mq.cachedReply[service] = callback
}

func (mq *MQ) UnReply(service string) {
	delete(mq.cachedReply, service)
}

func (mq *MQ) SetParamChecker(checker ParamChecker) {
	mq.checker = checker
}

func (mq *MQ) Call(name string, args Request, reply *Reply) error {
	return nil
}

func (mq *MQ) Close() {
}

func (mq *MQ) sndMsg(mode CommCmd, subj string, data []byte) error {
	msg, err := NewMsg(mode, subj, data)
	if err != nil {
		return err
	}
	msgBuf, err := EncodeMsg(msg)
	if err != nil {
		return err
	}
	return mq.msgConn.Write(msgBuf)
}

func (mq *MQ) rcvMsg() error {
	for {
		msgBuf, err := mq.msgConn.Read()
		if err != nil {
			return err
		}

		msg, err := DecodeMsg(msgBuf)
		if err != nil {
			return err
		}

		deleteOldest := false
		switch msg.Cmd {
		case CommCmdPop:
			if mq.cachedPop[msg.Subj] == nil {
				mq.cachedPop[msg.Subj] = gqueue.NewQueue()
				mq.cachedPop[msg.Subj].SetLimit(mqMaxCached)
			}
			deleteOldest = mq.cachedPop[msg.Subj].Push(msg)
		case CommCmdSub:
			if mq.cachedSub != nil && mq.cachedSub[msg.Subj] != nil {
				go mq.cachedSub[msg.Subj](msg.Subj, msg.Data)
			} else {
				glog.Warnf("CommCmd(%d) subj(%s) ignored message because of nil SubCallback.", msg.Cmd, msg.Subj)
			}
		case CommCmdRequest:
			if mq.cachedReply != nil && mq.cachedReply[msg.Subj] != nil {
				go func() {
					if replyBuf, err := mq.cachedReply[msg.Subj](msg.Subj, msg.Data); err != nil {
						glog.Erro(err)
					} else {
						if err := mq.sndMsg(CommCmdReply, msg.Subj, replyBuf); err != nil {
							glog.Erro(err)
						}
					}
				}()
			} else {
				glog.Warnf("CommCmd(%d) subj(%s) ignored message because of nil ReplyCallback.", msg.Cmd, msg.Subj)
			}
		case CommCmdReply:
			if mq.cachedRequest[msg.Subj] != nil {
				if len(mq.cachedRequest[msg.Subj]) == 0 {
					mq.cachedRequest[msg.Subj] <- msg
				} else {
					return gerrors.New("duplicated reply service(%s) found", msg.Subj)
				}
			} else {
				return gerrors.New("invalid reply service(%s) found because of invalid requestID", msg.Subj)
			}
		default:
			return gerrors.New("unknown CommCmd %d", msg.Cmd)
		}
		if deleteOldest {
			glog.Warnf("CommCmd(%d) subj(%s) deleted oldest message because of max limit reached.", msg.Cmd, msg.Subj)
		}
	}
}
