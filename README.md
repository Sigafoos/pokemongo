# pokemon-go
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v1.4%20adopted-ff69b4.svg)](code-of-conduct.md)
[![GoDoc](https://godoc.org/github.com/Sigafoos/pokemongo?status.svg)](https://godoc.org/github.com/Sigafoos/pokemongo)

A library for dealing with Pokemon Go CP and IVs.

It requires you to have a `gamemaster.json` file in the same format as what PVPoke uses. The easiest way would be to... use PVPoke's:

```
curl -O https://raw.githubusercontent.com/pvpoke/pvpoke/master/src/data/gamemaster.json
```

## Usage
Here's an example `main.go`. It:

* loads the gamemaster file
* looks up Wobbuffet by its Pokedex number
* sets its level and IVs
* calculates its CP and stat product
* prints the object

```go
package main

import (
	"fmt"
	"os"

	"github.com/Sigafoos/pokemongo"
)

func main() {
	fp, err := os.Open("gamemaster.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fp.Close()

	gm, err := pokemongo.NewGamemaster(fp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	wob := gm.PokemonByNumber(202)
	wob.Level = 23.5
	wob.IVs = pokemongo.Stats{
		Attack:  10,
		Defense: 13,
		HP:      12,
	}
	wob.Calculate()
	fmt.Printf("%+v\n", wob)
}
```
