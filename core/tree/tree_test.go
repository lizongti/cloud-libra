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

func TestGetSet(t *testing.T) {
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
	mapTree := tree.NewMapTree(data)

	style := magic.NewChainStyle(magic.SeparatorPeriod, magic.SeparatorUnderscore)

	v := mapTree.Get(style.Chain("animals.0.species_name"))
	if s, ok := v.(string); !ok {
		t.Fatalf("expected a string, but got %v", v)
	} else if s != "Elephant" {
		t.Fatalf("expected a string of `Elephant`, but got %s", s)
	}

	mapTree.Set(style.Chain("name"), "Sea World")
	mapTree.Set(style.Chain("animals.0.species_name"), "Dolphin")
	mapTree.Set(style.Chain("keepers.3"), map[string]interface{}{"Name": "Mike", "Age": 29})
	mapTree.Set(style.Chain("animals.4.1.3"), "a mistake")

	if mapTree.Get(style.Chain("name")) != "Sea World" {
		t.Fatalf("expected zoo name `Sea World`, but got %s", mapTree.Get(style.Chain("sea")))
	}
	if mapTree.Get(style.Chain("animals.0.species_name")) != "Dolphin" {
		t.Fatalf("expected species name `Dolphin`, but got %s", mapTree.Get(style.Chain("animals.0.species_name")))
	}
	if mapTree.Get(style.Chain("keepers.3.name")) != "Mike" {
		t.Fatalf("expected keeper name `Mike`, but got %s", mapTree.Get(style.Chain("keepers.3.name")))
	}

	zooJSON, err = json.Marshal(mapTree.Data())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", string(zooJSON))

	mapTree.Remove(style.Chain("animals.4.1.3"))
	mapTree.Remove(style.Chain("Keepers.1.Name"))
	mapTree.Remove(style.Chain("Keepers.1.Age"))
	zooJSON, err = json.Marshal(mapTree.Data())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", string(zooJSON))
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

	mapTree.Remove(style.Chain("Keepers.3.Name"))
	mapTree.Remove(style.Chain("Keepers.3.Age"))
	zooJSON, err = json.Marshal(mapTree.Data())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", string(zooJSON))
	zoo = Zoo{}
	err = json.Unmarshal(zooJSON, &zoo)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", zoo)
	if len(zoo.Keepers) != 1 {
		t.Fatalf("expected 1 keeper, but got %d", len(zoo.Keepers))
	}
}
