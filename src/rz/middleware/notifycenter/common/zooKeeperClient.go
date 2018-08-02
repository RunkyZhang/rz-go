package common

import (
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

func NewZooKeeperClient(addresses []string, timeout time.Duration, rootPath string) (*ZooKeeperClient) {
	var err error
	zooKeeperClient := &ZooKeeperClient{}
	zooKeeperClient.conn, _, err = zk.Connect(addresses, timeout, zk.WithEventCallback(zooKeeperClient.process))
	if nil != err {
		panic(err)
	}

	return zooKeeperClient
}

type zNode struct {
	zk.Stat

	Value string
}

type ZooKeeperClient struct {
	conn *zk.Conn
}

func (myself *ZooKeeperClient) CreatePersistent(path string, value string, acl ...[]zk.ACL) {
	myself.conn.Create(path, []byte(value), 0, myself.getDefaultACL(acl...))
}

func (myself *ZooKeeperClient) CreateEphemeral(path string, value string, acl ...[]zk.ACL) {
	myself.conn.Create(path, []byte(value), 1, myself.getDefaultACL(acl...))
}

func (myself *ZooKeeperClient) CreatePersistentSequence(path string, value string, acl ...[]zk.ACL) {
	myself.conn.Create(path, []byte(value), 2, myself.getDefaultACL(acl...))
}

func (myself *ZooKeeperClient) CreateEphemeralSequence(path string, value string, acl ...[]zk.ACL) {
	myself.conn.Create(path, []byte(value), 3, myself.getDefaultACL(acl...))
}

func (myself *ZooKeeperClient) Set(path string, value string, version ...int32) {
	myself.conn.Set(path, []byte(value), myself.getDefaultVersion(version...))
}

func (myself *ZooKeeperClient) Get(path string, watcher zk.EventCallback) (*zNode, error) {
	if nil == watcher {
		buffer, state, err := myself.conn.Get(path)
		if nil != err {
			return nil, err
		}

		zNode := &zNode{}
		zNode.Stat = *state
		zNode.Value = string(buffer)

		return zNode, nil
	} else {
		buffer, state, channel, err := myself.conn.GetW(path)
		if nil != err {
			return nil, err
		}

		zNode := &zNode{}
		zNode.Stat = *state
		zNode.Value = string(buffer)

		go func() {
			event := <-channel
			watcher(event)
		}()

		return zNode, nil
	}
}

func (myself *ZooKeeperClient) Exists(path string, watcher zk.EventCallback) (*zk.Stat, error) {
	if nil == watcher {
		_, state, err := myself.conn.Exists(path)
		if nil != err {
			return nil, err
		}

		return state, err
	} else {
		_, state, channel, err := myself.conn.ExistsW(path)
		if nil != err {
			return nil, err
		}

		if nil != state {
			go func() {
				event := <-channel
				watcher(event)
			}()
		}

		return state, err
	}
}

func (myself *ZooKeeperClient) Watch(path string, watcher zk.EventCallback, around ...bool) (error) {
	_, state, channel, err := myself.conn.ExistsW(path)
	if nil != err {
		return err
	}

	if nil != state {
		go func() {
			event := <-channel
			watcher(event)
			if myself.getDefaultAround(around...) {
				myself.Watch(path, watcher, around...)
			}
		}()
	}

	return nil
}

func (myself *ZooKeeperClient) WatchChildren(path string, watcher zk.EventCallback, around ...bool) (error) {
	_, state, channel, err := myself.conn.ChildrenW(path)
	if nil != err {
		return err
	}

	if nil != state {
		go func() {
			event := <-channel
			watcher(event)
			if myself.getDefaultAround(around...) {
				myself.WatchChildren(path, watcher, around...)
			}
		}()
	}

	return nil
}

func (myself *ZooKeeperClient) process(event zk.Event) {
	defer func() {

	}()
}

func (myself *ZooKeeperClient) getDefaultACL(acl ...[]zk.ACL) ([]zk.ACL) {
	if nil == acl || 0 == len(acl) {
		return zk.WorldACL(zk.PermAll)
	} else {
		return acl[0]
	}
}

func (myself *ZooKeeperClient) getDefaultVersion(version ...int32) (int32) {
	if nil == version || 0 == len(version) {
		return 0
	} else {
		return version[0]
	}
}

func (myself *ZooKeeperClient) getDefaultAround(around ...bool) (bool) {
	if nil == around || 0 == len(around) {
		return false
	} else {
		return around[0]
	}
}
