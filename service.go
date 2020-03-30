package discovery

import (
	"encoding/json"
	"path"
)

type Service struct {
	Name    string
	Version string
	nodes   []*Node
}

type Node struct {
	Id   string            `json:"id"`
	Addr string            `json:"addr"`
	Meta map[string]string `json:"meta"`
}

func (n *Node) Decode(data []byte) error {
	return json.Unmarshal(data, n)
}

func (n *Node) Encode() []byte {
	data, _ := json.Marshal(n)
	return data
}

func (n *Node) String() string {
	return string(n.Encode())
}

func nodePath(name, version, addr string) string {
	return path.Join("/", name, version, addr)
}

func nodePrefix(name, version string) string {
	return path.Join("/", name, version)
}
