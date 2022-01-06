package tree_test

import (
	"fmt"
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

func TestMapTree(t *testing.T) {
	var zoo = Zoo{
		Name: "Beasts & Monsters",
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

	json := encoding.NewJSON()
	zooJSON, err := json.Marshal(zoo)
	if err != nil {
		t.Fatal(err)
	}
	var data map[string]interface{}
	err = json.Unmarshal(zooJSON, &data)
	if err != nil {
		t.Fatal(err)
	}
	mapTree := tree.NewMapTree()
	mapTree.SetData(data)

	v := mapTree.Get(magic.UnixChain("animals.0.species_name"))
	if s, ok := v.(string); !ok {
		t.Fatalf("expected a string, but got %v", v)
	} else if s != "Elephant" {
		t.Fatalf("expected a string of `Elephant`, but got %s", s)
	}

	mapTree.Set(magic.UnixChain("name"), "Sea World")
	mapTree.Set(magic.UnixChain("animals.0.species_name"), "Dolphin")
	mapTree.Set(magic.UnixChain("keepers.3"), map[string]interface{}{"Name": "Mike", "Age": 29})
	mapTree.Set(magic.UnixChain("animals.4.1.3"), "a mistake")

	if mapTree.Get(magic.UnixChain("name")) != "Sea World" {
		t.Fatalf("expected zoo name `Sea World`, but got %s", mapTree.Get(magic.UnixChain("sea")))
	}
	if mapTree.Get(magic.UnixChain("animals.0.species_name")) != "Dolphin" {
		t.Fatalf("expected species name `Dolphin`, but got %s", mapTree.Get(magic.UnixChain("animals.0.species_name")))
	}
	if mapTree.Get(magic.UnixChain("keepers.3.name")) != "Mike" {
		t.Fatalf("expected keeper name `Mike`, but got %s", mapTree.Get(magic.UnixChain("keepers.3.name")))
	}

	mapTree.Remove(magic.UnixChain("animals.4.1.3"))
	mapTree.Remove(magic.UnixChain("Keepers.1.Name"))
	mapTree.Remove(magic.UnixChain("Keepers.1.Age"))
	zooJSON, err = json.Marshal(mapTree.Data())
	if err != nil {
		t.Fatal(err)
	}
	zoo = Zoo{}
	err = json.Unmarshal(zooJSON, &zoo)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", zoo)
	if zoo.Name != "Sea World" {
		t.Fatalf("expected species name `Sea World`, but got %s", zoo.Name)
	}
	if zoo.Animals[0].SpeciesName != "Dolphin" {
		t.Fatalf("expected species name `Dolphin`, but got %s", zoo.Animals[0].SpeciesName)
	}
	if zoo.Keepers[3].Name != "Mike" {
		t.Fatalf("expected keeper name `Mike`, but got %s", zoo.Keepers[3].Name)
	}

	mapTree2 := mapTree.Dulplicate()

	mapTree2.Remove(magic.UnixChain("Keepers.3.Name"))
	mapTree2.Remove(magic.UnixChain("Keepers.3.Age"))
	mapTree2.Set(magic.UnixChain("Name"), "Universe World")
	mapTree2.Set(magic.UnixChain("Keepers.0.Age"), 100)
	zooJSON, err = json.Marshal(mapTree2.Data())
	if err != nil {
		t.Fatal(err)
	}
	zoo = Zoo{}
	err = json.Unmarshal(zooJSON, &zoo)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", zoo)
	if len(zoo.Keepers) != 1 {
		t.Fatalf("expected 1 keeper, but got %d", len(zoo.Keepers))
	}

	mapTree.Merge(mapTree2)
	zooJSON, err = json.Marshal(mapTree.Data())
	if err != nil {
		t.Fatal(err)
	}
	zoo = Zoo{}
	err = json.Unmarshal(zooJSON, &zoo)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", zoo)
	if zoo.Name != "Universe World" {
		t.Fatalf("expected name `Universe World`, but got %s", zoo.Name)
	}

	hash, err := encoding.NewHash().Marshal(mapTree)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(hash))

	mapTree3 := tree.NewMapTree()
	if err := encoding.NewHash().Unmarshal(hash, mapTree3); err != nil {
		t.Fatal(err)
	}

	// DeepEqual doesnot match
	if fmt.Sprintf("%v", mapTree3) != fmt.Sprintf("%v", mapTree) {
		t.Fatalf("expected `mapTree3` equals to `mapTree`, mapTree3: %+v, mapTree: %+v", mapTree3, mapTree)
	}
}
