package tree_test

import (
	"testing"

	"github.com/aceaura/libra/core/encoding"
	"github.com/aceaura/libra/core/magic"
	"github.com/aceaura/libra/core/tree"
)

type Zoo struct {
	Name    string
	Keepers []Keeper
	Animals []Animal
}

type Keeper struct {
	Name string
	Age  int
}

type Animal struct {
	SpeciesName string
	Count       int
}

var zoo = Zoo{
	Name: "Beasts & Montsters",
	Keepers: []Keeper{
		{
			Name: "Alice",
			Age:  20,
		},
		{
			Name: "Tony",
			Age:  50,
		},
	},
	Animals: []Animal{
		{
			SpeciesName: "Elephant",
			Count:       3,
		},
		{
			SpeciesName: "Money",
			Count:       20,
		},
	},
}

func TestGet(t *testing.T) {
	zooJSON, err := encoding.Encode(encoding.NewJSON(), zoo)
	if err != nil {
		t.Fatal(err)
	}
	var data map[string]interface{}
	err = encoding.Decode(encoding.NewJSON(), zooJSON, &data)
	if err != nil {
		t.Fatal(err)
	}

	style := magic.NewChainStyle(magic.SeparatorPeriod, magic.SeparatorUnderscore)
	chain := style.Chain("animals.0.species_name")
	mapTree := tree.NewMapTree(data)
	t.Log(mapTree.Get(chain))
	v := mapTree.Get(chain)
	if s, ok := v.(string); !ok {
		t.Fatalf("expected a string, but got %v", v)
	} else if s != "Elephant" {
		t.Fatalf("expected a string of `Elephant`, but got %s", s)
	}
}
