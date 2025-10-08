package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB(t *testing.T) (*mongo.Database, func()) {
	// Configuração do MongoDB de teste
	mongoURL := os.Getenv("MONGODB_URL")
	if mongoURL == "" {
		mongoURL = "mongodb://admin:admin@localhost:27017/auctions_test?authSource=admin"
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		t.Skipf("Skipping test: MongoDB not available: %v", err)
		return nil, func() {}
	}

	// Verifica conexão
	err = client.Ping(ctx, nil)
	if err != nil {
		t.Skipf("Skipping test: MongoDB not available: %v", err)
		return nil, func() {}
	}

	db := client.Database("auctions_test")

	cleanup := func() {
		// Limpa coleção após testes
		db.Collection("auctions").Drop(context.Background())
		client.Disconnect(context.Background())
	}

	return db, cleanup
}

func TestAuctionAutoClose(t *testing.T) {
	// Define intervalo curto para teste
	os.Setenv("AUCTION_INTERVAL", "3s")
	defer os.Unsetenv("AUCTION_INTERVAL")

	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	// Cria repository
	repo := NewAuctionRepository(db)

	// Cria um leilão de teste
	auction, err := auction_entity.CreateAuction(
		"Test Product",
		"Electronics",
		"Test Description for auction",
		auction_entity.New,
	)
	assert.Nil(t, err)

	ctx := context.Background()

	// Insere o leilão
	internalErr := repo.CreateAuction(ctx, auction)
	assert.Nil(t, internalErr)

	// Verifica que o leilão foi criado com status Active
	auctionCreated, internalErr := repo.FindAuctionById(ctx, auction.Id)
	assert.Nil(t, internalErr)
	assert.Equal(t, auction_entity.Active, auctionCreated.Status)

	// Aguarda mais do que o tempo do leilão + tempo de verificação (10s)
	time.Sleep(14 * time.Second)

	// Busca o leilão novamente
	auctionAfterExpiry, internalErr := repo.FindAuctionById(ctx, auction.Id)
	assert.Nil(t, internalErr)

	// Verifica que o status foi alterado para Completed
	assert.Equal(t, auction_entity.Completed, auctionAfterExpiry.Status, "Auction should be automatically closed after expiration")
}

func TestAuctionNotClosedBeforeExpiry(t *testing.T) {
	// Define intervalo longo para teste
	os.Setenv("AUCTION_INTERVAL", "30s")
	defer os.Unsetenv("AUCTION_INTERVAL")

	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	// Cria repository
	repo := NewAuctionRepository(db)

	// Cria um leilão de teste
	auction, err := auction_entity.CreateAuction(
		"Test Product 2",
		"Books",
		"Test Description for second auction",
		auction_entity.Used,
	)
	assert.Nil(t, err)

	ctx := context.Background()

	// Insere o leilão
	internalErr := repo.CreateAuction(ctx, auction)
	assert.Nil(t, internalErr)

	// Verifica que o leilão foi criado com status Active
	auctionCreated, internalErr := repo.FindAuctionById(ctx, auction.Id)
	assert.Nil(t, internalErr)
	assert.Equal(t, auction_entity.Active, auctionCreated.Status)

	// Aguarda menos do que o tempo do leilão
	time.Sleep(5 * time.Second)

	// Busca o leilão novamente
	auctionBeforeExpiry, internalErr := repo.FindAuctionById(ctx, auction.Id)
	assert.Nil(t, internalErr)

	// Verifica que o status ainda é Active
	assert.Equal(t, auction_entity.Active, auctionBeforeExpiry.Status, "Auction should still be active before expiration")
}

func TestUpdateAuctionStatus(t *testing.T) {
	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	repo := NewAuctionRepository(db)

	// Cria um leilão de teste
	auction, err := auction_entity.CreateAuction(
		"Test Product 3",
		"Clothing",
		"Test Description for third auction",
		auction_entity.Refurbished,
	)
	assert.Nil(t, err)

	ctx := context.Background()

	// Insere o leilão
	internalErr := repo.CreateAuction(ctx, auction)
	assert.Nil(t, internalErr)

	// Atualiza o status manualmente
	internalErr = repo.UpdateAuctionStatus(ctx, auction.Id, auction_entity.Completed)
	assert.Nil(t, internalErr)

	// Verifica que o status foi atualizado
	var auctionMongo AuctionEntityMongo
	err2 := repo.Collection.FindOne(ctx, bson.M{"_id": auction.Id}).Decode(&auctionMongo)
	assert.Nil(t, err2)
	assert.Equal(t, auction_entity.Completed, auctionMongo.Status)
}

func TestGetAuctionInterval(t *testing.T) {
	// Test com valor válido
	os.Setenv("AUCTION_INTERVAL", "10m")
	duration := getAuctionInterval()
	assert.Equal(t, 10*time.Minute, duration)
	os.Unsetenv("AUCTION_INTERVAL")

	// Test com valor inválido (deve retornar padrão)
	os.Setenv("AUCTION_INTERVAL", "invalid")
	duration = getAuctionInterval()
	assert.Equal(t, 5*time.Minute, duration)
	os.Unsetenv("AUCTION_INTERVAL")

	// Test sem variável de ambiente (deve retornar padrão)
	duration = getAuctionInterval()
	assert.Equal(t, 5*time.Minute, duration)
}
