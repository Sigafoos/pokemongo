package pokemongo

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

// Gamemaster is the parsed gamemaster.json file. The file itself contains cups (Silph Arena), moves,
// pokemon and settings, but currently only the pokemon are supported by this type.
type Gamemaster struct {
	Pokemon   []Pokemon `json:"pokemon"`
	Shadow    []string  `json:"shadowPokemon"`
	shadowMap map[string]bool
	dexMap    map[int]*Pokemon
	nameMap   map[string]*Pokemon
	idMap     map[string]*Pokemon
}

// NewGamemaster instantiates a Gamemaster object. It expects the data to be marshalled in json.
// It does not handle closing the Reader: that is left to the calling code.
func NewGamemaster(fp io.Reader) (*Gamemaster, error) {
	body, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, fmt.Errorf("error reading gamemaster: %s", err.Error())
	}

	var gm Gamemaster
	err = json.Unmarshal(body, &gm)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling gamemaster: %s", err.Error())
	}

	gm.shadowMap = make(map[string]bool)
	for _, p := range gm.Shadow {
		gm.shadowMap[p] = true
	}

	gm.dexMap = make(map[int]*Pokemon)
	gm.nameMap = make(map[string]*Pokemon)
	gm.idMap = make(map[string]*Pokemon)

	// iterate over the list of Pokemon once and populate the lookup maps.
	for k := range gm.Pokemon {
		p := gm.Pokemon[k]
		if _, ok := gm.shadowMap[p.ID]; ok {
			p.Shadow = true
		}
		gm.dexMap[p.Dex] = &p
		gm.nameMap[p.Name] = &p
		gm.idMap[p.ID] = &p
	}

	return &gm, nil
}

// PokemonByNumber accepts a Pokedex number and returns the corresponding Pokemon.
func (gm *Gamemaster) PokemonByNumber(number int) *Pokemon {
	p, ok := gm.dexMap[number]
	if !ok {
		return nil
	}
	return p
}

// PokemonByName accepts a capitalized name -- ie Bulbasaur -- and returns the corresponding Pokemon.
func (gm *Gamemaster) PokemonByName(name string) *Pokemon {
	p, ok := gm.nameMap[name]
	if !ok {
		return nil
	}
	return p
}

// PokemonByName accepts an id -- ie mewtwo_armored -- and returns the corresponding Pokemon.
func (gm *Gamemaster) PokemonByID(ID string) *Pokemon {
	p, ok := gm.idMap[ID]
	if !ok {
		return nil
	}
	return p
}
