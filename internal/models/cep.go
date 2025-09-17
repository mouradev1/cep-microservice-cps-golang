package models

import "time"

type Cep struct {
    Cep         string    `bson:"cep" json:"cep"`
    Logradouro  string    `bson:"logradouro" json:"logradouro"`
    Bairro      string    `bson:"bairro" json:"bairro"`
    Localidade  string    `bson:"localidade" json:"localidade"`
    Uf          string    `bson:"uf" json:"uf"`
    Ibge        string    `bson:"ibge" json:"ibge"`
    Ddd         string    `bson:"ddd" json:"ddd"`
    LastProvider string   `bson:"lastProvider" json:"lastProvider"`
    CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
}