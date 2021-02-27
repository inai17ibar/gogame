package model

type Character struct {
	ID          string `json:"userCharacterID"`
	CharacterID string `json:"characterID"`
	Name        string `json:"name"`
}

type CharactersListResponse struct {
	Characters []Character `json:"characters"`
}
