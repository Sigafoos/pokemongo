package pokemongo

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockReadWriter implements the io.Reader and io.Writer interfaces.
type mockReadWriter struct {
	mock.Mock
}

// Write is a mocked implementation of reading bytes from a file.
func (m *mockReadWriter) Read(p []byte) (int, error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

// Write is a mocked implementation of writing bytes to a file.
func (m *mockReadWriter) Write(p []byte) (int, error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

const wigglytuffYaml = `id: wigglytuff
dex: 40
name: Wigglytuff
basestats:
    attack: 156
    defense: 90
    hp: 295
ivs:
    attack: 10
    defense: 15
    hp: 12
calculatedstats:
    attack: 117.3427772
    defense: 74.222841
    hp: 217
level: 28
cp: 1489
moves:
    fast: CHARM
    charge:
      - ICE_BEAM
      - PLAY_ROUGH
`

var wigglytuffStruct = Pokemon{
	ID:   "wigglytuff",
	Dex:  40,
	Name: "Wigglytuff",
	BaseStats: Stats{
		Attack:  156.0,
		Defense: 90.0,
		HP:      295.0,
	},
	IVs: Stats{
		Attack:  10.0,
		Defense: 15.0,
		HP:      12.0,
	},
	CalculatedStats: Stats{
		Attack:  117.3427772,
		Defense: 74.222841,
		HP:      217,
	},
	Level: 28.0,
	CP:    1489,
	Moves: Moves{
		Fast: "CHARM",
		Charge: [2]string{
			"ICE_BEAM",
			"PLAY_ROUGH",
		},
	},
}

// TestLoadHappyPath ensures that when all goes well we can unmarshal a Pokemon
// from yaml.
func TestLoadHappyPath(t *testing.T) {
	fp := strings.NewReader(wigglytuffYaml)

	p, err := Load(fp)

	assert.Equal(t, &wigglytuffStruct, p)
	assert.Nil(t, err)
}

// TestLoadBadYaml ensures that invalid yaml will be treated as such.
func TestLoadBadYaml(t *testing.T) {
	fp := strings.NewReader("!@#$%^&*()")

	p, err := Load(fp)

	assert.Nil(t, p)
	assert.Error(t, err)
}

// TestLoadCannotReadAll tests that an issue with the io.Reader will return the error.
func TestLoadCannotReadAll(t *testing.T) {
	expectedErr := fmt.Errorf("sorry, something's gone wrong")
	fp := new(mockReadWriter)
	fp.On("Read", mock.AnythingOfType("[]uint8")).
		Return(0, expectedErr)

	p, err := Load(fp)

	assert.Nil(t, p)
	assert.Error(t, err)
}

// TestSaveHappyPath ensures that when all goes well we can marshal a Pokemon
// into yaml.
func TestSaveHappyPath(t *testing.T) {
	fp := new(mockReadWriter)
	fp.On("Write", []uint8(wigglytuffYaml)).
		Return(23, nil)

	err := wigglytuffStruct.Save(fp)

	assert.Nil(t, err)
}

// TestSaveErrorWriting tests that an error while writing to the io.Writer will
// return an error
func TestSaveErrorWriting(t *testing.T) {
	expectedErr := fmt.Errorf("sorry, something's gone wrong")
	fp := new(mockReadWriter)
	fp.On("Write", []uint8(wigglytuffYaml)).
		Return(0, expectedErr)

	err := wigglytuffStruct.Save(fp)

	assert.Error(t, err)
}

// TestCalculate tests a series of Pokemon stats and ensures that the correct CP/etc
// will be generated.
func TestCalculate(t *testing.T) {
	tests := []struct {
		name       string
		atkIV      float64
		defIV      float64
		hpIV       float64
		level      float64
		expectedCP int
	}{
		{
			name:       "wigglytuff",
			atkIV:      0.0,
			defIV:      0.0,
			hpIV:       0.0,
			level:      36.0,
			expectedCP: 1496,
		},
		{
			name:       "wigglytuff",
			atkIV:      15.0,
			defIV:      15.0,
			hpIV:       15.0,
			level:      27.0,
			expectedCP: 1486,
		},
		{
			name:       "wigglytuff",
			atkIV:      0.0,
			defIV:      0.0,
			hpIV:       0.0,
			level:      1.0,
			expectedCP: 22,
		},
		{
			name:       "shedinja",
			atkIV:      0.0,
			defIV:      0.0,
			hpIV:       0.0,
			level:      4.5,
			expectedCP: 10,
		},
	}

	pokes := make(map[string]*Pokemon)
	pokes["wigglytuff"] = &Pokemon{
		BaseStats: Stats{
			Attack:  156.0,
			Defense: 90.0,
			HP:      295.0,
		},
	}
	pokes["shedinja"] = &Pokemon{
		BaseStats: Stats{
			Attack:  153.0,
			Defense: 73.0,
			HP:      1.0,
		},
	}

	for _, test := range tests {
		p, ok := pokes[test.name]
		assert.True(t, ok)

		p.Level = test.level
		p.IVs = Stats{
			Attack:  test.atkIV,
			Defense: test.defIV,
			HP:      test.hpIV,
		}
		p.Calculate()
		assert.Equal(t, test.expectedCP, p.CP)
	}
}
