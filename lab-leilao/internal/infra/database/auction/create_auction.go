package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection              *mongo.Collection
	auctionInterval         time.Duration
	closedAuctionsMap       map[string]bool
	closedAuctionsMapMutex  *sync.Mutex
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	repo := &AuctionRepository{
		Collection:             database.Collection("auctions"),
		auctionInterval:        getAuctionInterval(),
		closedAuctionsMap:      make(map[string]bool),
		closedAuctionsMapMutex: &sync.Mutex{},
	}

	// Inicia goroutine para monitorar e fechar leilões vencidos
	go repo.monitorExpiredAuctions(context.Background())

	return repo
}

func getAuctionInterval() time.Duration {
	auctionInterval := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(auctionInterval)
	if err != nil {
		return time.Minute * 5
	}
	return duration
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	return nil
}

// UpdateAuctionStatus atualiza o status de um leilão
func (ar *AuctionRepository) UpdateAuctionStatus(
	ctx context.Context,
	auctionId string,
	status auction_entity.AuctionStatus) *internal_error.InternalError {

	filter := bson.M{"_id": auctionId}
	update := bson.M{"$set": bson.M{"status": status}}

	_, err := ar.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Error trying to update auction status", err)
		return internal_error.NewInternalServerError("Error trying to update auction status")
	}

	return nil
}

// monitorExpiredAuctions monitora e fecha leilões que ultrapassaram o tempo definido
func (ar *AuctionRepository) monitorExpiredAuctions(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 10) // Verifica a cada 10 segundos
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ar.checkAndCloseExpiredAuctions(ctx)
		case <-ctx.Done():
			return
		}
	}
}

// checkAndCloseExpiredAuctions busca e fecha leilões vencidos
func (ar *AuctionRepository) checkAndCloseExpiredAuctions(ctx context.Context) {
	// Busca leilões ativos
	filter := bson.M{"status": auction_entity.Active}
	cursor, err := ar.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error("Error finding active auctions", err)
		return
	}
	defer cursor.Close(ctx)

	var auctions []AuctionEntityMongo
	if err := cursor.All(ctx, &auctions); err != nil {
		logger.Error("Error decoding auctions", err)
		return
	}

	now := time.Now()

	// Verifica cada leilão ativo
	for _, auction := range auctions {
		// Verifica se já foi processado
		ar.closedAuctionsMapMutex.Lock()
		alreadyClosed := ar.closedAuctionsMap[auction.Id]
		ar.closedAuctionsMapMutex.Unlock()

		if alreadyClosed {
			continue
		}

		// Calcula o tempo de expiração do leilão
		auctionEndTime := time.Unix(auction.Timestamp, 0).Add(ar.auctionInterval)

		// Se o leilão expirou, fecha ele
		if now.After(auctionEndTime) {
			err := ar.UpdateAuctionStatus(ctx, auction.Id, auction_entity.Completed)
			if err != nil {
				logger.Error("Error closing expired auction", err)
				continue
			}

			// Marca como fechado no mapa para evitar múltiplos updates
			ar.closedAuctionsMapMutex.Lock()
			ar.closedAuctionsMap[auction.Id] = true
			ar.closedAuctionsMapMutex.Unlock()

			logger.Info("Auction " + auction.Id + " closed automatically after expiration")
		}
	}
}
