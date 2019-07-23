// Package gamemaster allows for the reading and use of a gamemaster.json file in the format that
// PvPoke uses (https://github.com/pvpoke/pvpoke/).
package gamemaster

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"git.theconley.club/dconley/pokemon-go/pokemon"
)

// Gamemaster is the parsed gamemaster.json file. The file itself contains cups (Silph Arena), moves,
// pokemon and settings, but currently only the pokemon are supported by this type.
type Gamemaster struct {
	Pokemon []pokemon.Pokemon `json:"pokemon"`
	dexMap  map[int]*pokemon.Pokemon
	nameMap map[string]*pokemon.Pokemon
	idMap   map[string]*pokemon.Pokemon
}

// New instantiates a Gamemaster object. It expects the data to be marshalled in json. It does not
// handle closing the Reader: that is left to the calling code.
func New(fp io.Reader) (*Gamemaster, error) {
	body, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, fmt.Errorf("error reading gamemaster: %s", err.Error())
	}

	var gm Gamemaster
	err = json.Unmarshal(body, &gm)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling gamemaster: %s", err.Error())
	}

	gm.dexMap = make(map[int]*pokemon.Pokemon)
	gm.nameMap = make(map[string]*pokemon.Pokemon)
	gm.idMap = make(map[string]*pokemon.Pokemon)

	// iterate over the list of Pokemon once and populate the lookup maps.
	for k := range gm.Pokemon {
		p := gm.Pokemon[k]
		gm.dexMap[p.Dex] = &p
		gm.nameMap[p.Name] = &p
		gm.idMap[p.ID] = &p
	}

	return &gm, nil
}

// PokemonByNumber accepts a Pokedex number and returns the corresponding Pokemon.
func (gm *Gamemaster) PokemonByNumber(number int) *pokemon.Pokemon {
	p, ok := gm.dexMap[number]
	if !ok {
		return nil
	}
	return p
}

// PokemonByName accepts a capitalized name -- ie Bulbasaur -- and returns the corresponding Pokemon.
func (gm *Gamemaster) PokemonByName(name string) *pokemon.Pokemon {
	p, ok := gm.nameMap[name]
	if !ok {
		return nil
	}
	return p
}

// PokemonByName accepts an id -- ie mewtwo_armored -- and returns the corresponding Pokemon.
func (gm *Gamemaster) PokemonByID(ID string) *pokemon.Pokemon {
	p, ok := gm.idMap[ID]
	if !ok {
		return nil
	}
	return p
}
