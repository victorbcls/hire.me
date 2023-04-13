package shorter

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Struct para registrar as informações de URL e alias
type Register struct {
	Alias   string `json:"alias"`
	URL     string `json:"url"`
	Acessos int    `json:"acessos"`
}

// Struct para registrar o tempo gasto na rotina
type Statistics struct {
	TimeTaken string `json:"time_taken"`
}

// Struct para registrar a resposta a ser enviada ao cliente
type Response struct {
	Register
	Statistics `json:"statistics"`
}

// Struct para registrar mensagens de erro
type ErrorResponse struct {
	Alias       string `json:"alias"`
	ErrCode     string `json:"err_code"`
	Description string `json:"description"`
}

// Gera um alias curto para uma URL. Se o alias personalizado estiver vazio, um alias aleatório será gerado.
// Se o alias já existe no banco de dados, retorna um erro.
func GenerateShortURL(url, customAlias string, client *mongo.Client) ([]byte, error) {

	if isValidURL(url) == false {

		response := ErrorResponse{
			Alias:       url,
			ErrCode:     "003",
			Description: "INVALID URL",
		}
		return json.Marshal(response)
	}

	// Inicia a contagem do tempo
	start := time.Now()
	// Define a string de caracteres permitidos para o alias
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var alias string

	// Gera um alias aleatório se o customAlias estiver vazio
	if customAlias == "" {
		shortURL := make([]byte, 7)
		for i := range shortURL {
			shortURL[i] = letters[rand.Intn(len(letters))]
		}
		alias = string(shortURL)
	} else {
		alias = customAlias
	}

	// Verifica se o alias já existe no banco de dados. Se customAlias estiver definido, retorna um erro.
	for checkAliasExistence(alias, client) {
		if customAlias != "" {
			errResp := ErrorResponse{
				Alias:       alias,
				ErrCode:     "001",
				Description: "CUSTOM ALIAS ALREADY EXISTS",
			}
			return json.Marshal(errResp)
		}

		// Gera um novo alias aleatório se o alias personalizado já existir no banco de dados
		shortURL := make([]byte, 7)
		for i := range shortURL {
			shortURL[i] = letters[rand.Intn(len(letters))]
		}
		alias = string(shortURL)
	}

	// Registra a URL e o alias no banco de dados
	register := Register{
		Alias:   alias,
		URL:     url,
		Acessos: 0,
	}

	if err := saveRegister(register, client); err != nil {
		errResp := ErrorResponse{
			Alias:       alias,
			ErrCode:     "002",
			Description: "DATABASE ERROR: FAILED TO SAVE REGISTER TO DATABASE",
		}
		return json.Marshal(errResp)
	}

	// Finaliza a contagem do tempo e registra o tempo gasto em milissegundos
	end := time.Now()
	duration := end.Sub(start).Milliseconds()

	// Cria a resposta a ser enviada ao cliente
	resp := Response{
		Register:   register,
		Statistics: Statistics{TimeTaken: fmt.Sprintf("%dms", duration)},
	}

	return json.Marshal(resp)
}

// Salva a URL e o alias no banco de dados
func saveRegister(register Register, client *mongo.Client) error {
	collection := client.Database("mydb").Collection("shorter")
	if _, err := collection.InsertOne(context.Background(), register); err != nil {
		log.Fatal(err)
		return err
	}
	log.Println("Information inserted successfully")
	return nil
}

func checkAliasExistence(alias string, client *mongo.Client) bool {
	collection := client.Database("mydb").Collection("shorter")
	filter := bson.M{"alias": alias}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	return count > 0
}

func RetrieveUrl(shortURL string, client *mongo.Client) ([]byte, error) {
	// Busca o registro no banco de dados pelo alias encurtado
	collection := client.Database("mydb").Collection("shorter")
	filter := bson.M{"alias": shortURL}
	var result Register
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		errResp := ErrorResponse{
			Alias:       shortURL,
			ErrCode:     "002",
			Description: "SHORTENED URL NOT FOUND IN DATABASE",
		}
		return json.Marshal(errResp)
	}

	// Atualiza o valor de acessos e salva a alteração no banco de dados
	result.Acessos++
	update := bson.M{"$set": bson.M{"acessos": result.Acessos}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		errResp := ErrorResponse{
			Alias:       shortURL,
			ErrCode:     "002",
			Description: "FAILED TO UPDATE SHORTENED URL IN DATABASE",
		}
		return json.Marshal(errResp)
	}

	return json.Marshal(result)
}
func isValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}
func TopAccess(client *mongo.Client) ([]byte, error) {
	// Busca os registros no banco de dados, ordenando por acessos em ordem decrescente e limitando a 10
	collection := client.Database("mydb").Collection("shorter")
	options := options.Find().SetSort(bson.M{"acessos": -1}).SetLimit(10)
	cursor, err := collection.Find(context.Background(), bson.M{}, options)
	if err != nil {
		errResp := ErrorResponse{
			Alias:       "",
			ErrCode:     "004",
			Description: "DATABASE ERROR: FAILED TO RETRIEVE TOP 10 REGISTERS FROM DATABASE",
		}
		return json.Marshal(errResp)
	}
	defer cursor.Close(context.Background())

	// Cria um slice com os resultados da busca
	var results []Register
	if err = cursor.All(context.Background(), &results); err != nil {
		errResp := ErrorResponse{
			Alias:       "",
			ErrCode:     "004",
			Description: "DATABASE ERROR: FAILED TO RETRIEVE TOP 10 REGISTERS FROM DATABASE",
		}
		return json.Marshal(errResp)
	}

	return json.Marshal(results)
}
