package services

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mouradev1/buscacepsgolang/internal/config"
	"github.com/mouradev1/buscacepsgolang/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func isCepDataComplete(data *models.Cep) bool {
    if data == nil {
        return false
    }
    onlyDigits := strings.Map(func(r rune) rune {
        if r >= '0' && r <= '9' {
            return r
        }
        return -1
    }, data.Cep)
    return len(onlyDigits) == 8 &&
        strings.TrimSpace(data.Logradouro) != "" &&
        strings.TrimSpace(data.Bairro) != "" &&
        strings.TrimSpace(data.Localidade) != "" &&
        strings.TrimSpace(data.Uf) != ""
}

func GetCepDataService(c *fiber.Ctx, cep string) (interface{}, int, error) {
	cleanCep := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, cep)

	if len(cleanCep) != 8 {
		return nil, fiber.StatusBadRequest, errors.New("CEP inválido")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go LogRequestByIp(c.IP())

	cached, err := GetCepCache(ctx, cleanCep)
	if err == nil {
		if cached != nil && isCepDataComplete(cached) {
			return cached, fiber.StatusOK, nil
		}
		if cached == nil {
			return nil, fiber.StatusNotFound, errors.New("CEP não encontrado")
		}
	} else if err.Error() != "redis: nil" {
		log.Printf("Erro ao acessar o cache: %v", err)
	}

	cepCollection := config.GetCollection("ceps")

	var cepDoc models.Cep
	err = cepCollection.FindOne(ctx, bson.M{"cep": cleanCep}).Decode(&cepDoc)
	if err == nil && isCepDataComplete(&cepDoc) {
		_ = SetCepCache(ctx, cleanCep, &cepDoc, time.Hour)
		return cepDoc, fiber.StatusOK, nil
	}

	lastProvider := ""
	if err == nil {
		lastProvider = cepDoc.LastProvider
	}
	fetchedCep, provider, err := FetchFromExternalApisWithRotation(cleanCep, lastProvider)
	if err != nil || !isCepDataComplete(fetchedCep) {
		// Salva "null" no cache por 4 horas
		_ = SetCepNotFoundCache(ctx, cleanCep, 4*time.Hour)
		return nil, fiber.StatusNotFound, errors.New("CEP não encontrado")
	}

	fetchedCep.Cep = cleanCep
	fetchedCep.LastProvider = provider
	fetchedCep.CreatedAt = time.Now()

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	var updated models.Cep
	err = cepCollection.FindOneAndUpdate(
		ctx,
		bson.M{"cep": cleanCep},
		bson.M{"$set": fetchedCep},
		opts,
	).Decode(&updated)
	if err != nil {
		_, _ = cepCollection.InsertOne(ctx, fetchedCep)
		_ = SetCepCache(ctx, cleanCep, fetchedCep, time.Hour)
		return fetchedCep, fiber.StatusOK, nil
	}
	_ = SetCepCache(ctx, cleanCep, &updated, time.Hour)
	return updated, fiber.StatusOK, nil
}

func FetchFromExternalApisWithRotation(cep string, lastProvider string) (*models.Cep, string, error) {
	providers := []struct {
		name string
		fn   providerFunc
	}{
		{"viacep", fetchFromViaCep},
		{"brasilapi", fetchFromBrasilApi},
		{"apicep", fetchFromApiCep},
	}

	startIdx := 0
	if lastProvider != "" {
		for i, p := range providers {
			if p.name == lastProvider {
				startIdx = (i + 1) % len(providers)
				break
			}
		}
	}

	ordered := append(providers[startIdx:], providers[:startIdx]...)
	for _, p := range ordered {
		log.Printf("Tentando buscar no provider: %s", p.name)
		result, err := p.fn(cep)
		if err == nil && result != nil {
			return result, p.name, nil
		}
		if err != nil {
			log.Printf("Erro ao buscar no provider %s: %v", p.name, err)
		}
	}
	return nil, "", errors.New("CEP não encontrado em nenhum provedor")
}

func LogRequestByIp(ip string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	now := time.Now()
	logCollection := config.GetCollection("requestlogs")
	logCollection.UpdateOne(
		ctx,
		bson.M{"ip": ip},
		bson.M{
			"$inc": bson.M{"count": 1},
			"$set": bson.M{"lastRequest": now, "createdAt": now},
		},
		options.Update().SetUpsert(true),
	)
}
