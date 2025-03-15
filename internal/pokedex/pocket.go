package pokedex

import (
	"sync"
)

type StatInfo struct {
	Name string `json:"name"`
}

type Stat struct {
	BaseStat int      `json:"base_stat"`
	StatInfo StatInfo `json:"stat"`
}

type TypeInfo struct {
	Name string `json:"name"`
}

type Type struct {
	TypeInfo TypeInfo `json:"type"`
}

type PokemonInfo struct {
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []Stat `json:"stats"`
	Types          []Type `json:"types"`
}

type Pocket struct {
	mx     sync.Mutex
	Pocket map[string]PokemonInfo
}

func NewPocket() *Pocket {
	return &Pocket{
		Pocket: make(map[string]PokemonInfo),
	}
}

func (p *Pocket) AddPokemon(key string, pokemonInfo PokemonInfo) {
	p.mx.Lock()
	defer p.mx.Unlock()
	p.Pocket[key] = pokemonInfo
}

func (p *Pocket) GetPokemon(key string) (PokemonInfo, bool) {
	p.mx.Lock()
	defer p.mx.Unlock()
	pokemon, exists := p.Pocket[key]
	return pokemon, exists
}
