package member

import (
	"fmt"

	"github.com/hashicorp/memberlist"
)

type Member struct {
}

func NewMember() *Member {
	list, err := memberlist.Create(memberlist.DefaultWANConfig())
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}
	
	list.Members()
}

func main() {
	list, err := memberlist.Create(memberlist.DefaultWANConfig())
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}

	// Join an existing cluster by specifying at least one known member.
	n, err := list.Join([]string{"1.2.3.4"})
	if err != nil {
		panic("Failed to join cluster: " + err.Error())
	}

	// Ask for members of the cluster
	for _, member := range list.Members() {
		fmt.Printf("Member: %s %s\n", member.Name, member.Addr)
	}
}
