package services

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "time"

    "github.com/mouradev1/buscacepsgolang/internal/models"
)

type providerFunc func(string) (*models.Cep, error)

var externalProviders = []struct {
    name string
    fn   providerFunc
}{
    {"viacep", fetchFromViaCep},
    {"brasilapi", fetchFromBrasilApi},
    {"apicep", fetchFromApiCep},
}

func FetchFromExternalApis(cep string) (*models.Cep, string, error) {
    for _, p := range externalProviders {
        result, err := p.fn(cep)
        if err == nil && result != nil {
            return result, p.name, nil
        }
    }
    return nil, "", errors.New("CEP não encontrado em nenhum provedor")
}

func fetchFromViaCep(cep string) (*models.Cep, error) {
    url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
    client := http.Client{Timeout: 5 * time.Second}
    resp, err := client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    var data struct {
        Cep        string `json:"cep"`
        Logradouro string `json:"logradouro"`
        Bairro     string `json:"bairro"`
        Localidade string `json:"localidade"`
        Uf         string `json:"uf"`
        Ibge       string `json:"ibge"`
        Ddd        string `json:"ddd"`
        Erro       bool   `json:"erro"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil || data.Erro {
        return nil, errors.New("não encontrado")
    }
    return &models.Cep{
        Cep:        data.Cep,
        Logradouro: data.Logradouro,
        Bairro:     data.Bairro,
        Localidade: data.Localidade,
        Uf:         data.Uf,
        Ibge:       data.Ibge,
        Ddd:        data.Ddd,
    }, nil
}

func fetchFromBrasilApi(cep string) (*models.Cep, error) {
    url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
    client := http.Client{Timeout: 5 * time.Second}
    resp, err := client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    if resp.StatusCode != 200 {
        return nil, errors.New("não encontrado")
    }
    var data struct {
        Cep         string `json:"cep"`
        Street      string `json:"street"`
        Neighborhood string `json:"neighborhood"`
        City        string `json:"city"`
        State       string `json:"state"`
        Ibge        string `json:"ibge"`
        Ddd         string `json:"ddd"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return nil, err
    }
    return &models.Cep{
        Cep:        data.Cep,
        Logradouro: data.Street,
        Bairro:     data.Neighborhood,
        Localidade: data.City,
        Uf:         data.State,
        Ibge:       data.Ibge,
        Ddd:        data.Ddd,
    }, nil
}

func fetchFromApiCep(cep string) (*models.Cep, error) {
    url := fmt.Sprintf("https://cdn.apicep.com/file/apicep/%s.json", cep)
    client := http.Client{Timeout: 5 * time.Second}
    resp, err := client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    if resp.StatusCode != 200 {
        return nil, errors.New("não encontrado")
    }
    var data struct {
        Code    string `json:"code"`
        Ok      bool   `json:"ok"`
        Status  int    `json:"status"`
        State   string `json:"state"`
        City    string `json:"city"`
        Address string `json:"address"`
        District string `json:"district"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return nil, err
    }
    if !data.Ok || data.Status != 200 {
        return nil, errors.New("não encontrado")
    }
    return &models.Cep{
        Cep:        data.Code,
        Logradouro: data.Address,
        Bairro:     data.District,
        Localidade: data.City,
        Uf:         data.State,
        Ibge:       "",
        Ddd:        "",
    }, nil
}