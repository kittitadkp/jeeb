package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	mongoAdapter "github.com/kittitadkp/jeeb-learning/internal/adapter/out/mongo"
	"github.com/kittitadkp/jeeb-learning/internal/config"
	"github.com/kittitadkp/jeeb-learning/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	client, err := mongoAdapter.NewClient(cfg.MongoDB.URI)
	if err != nil {
		slog.Error("failed to connect to mongodb", "error", err)
		os.Exit(1)
	}
	defer client.Disconnect(context.Background())

	db := client.Database(cfg.MongoDB.Database)
	ctx := context.Background()

	if err := ensureIndexes(ctx, db); err != nil {
		slog.Error("failed to create indexes", "error", err)
		os.Exit(1)
	}

	topicID, err := seedTopic(ctx, db)
	if err != nil {
		slog.Error("failed to seed topic", "error", err)
		os.Exit(1)
	}

	if err := seedItems(ctx, db, topicID); err != nil {
		slog.Error("failed to seed items", "error", err)
		os.Exit(1)
	}

	slog.Info("seed complete")
}

func ensureIndexes(ctx context.Context, db *mongo.Database) error {
	progressColl := db.Collection("progress")
	_, err := progressColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "item_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "topic_id", Value: 1}},
		},
	})
	return err
}

func seedTopic(ctx context.Context, db *mongo.Database) (string, error) {
	coll := db.Collection("topics")

	existing := &domain.Topic{}
	err := coll.FindOne(ctx, bson.M{"name": "IPA Phonetics"}).Decode(existing)
	if err == nil {
		slog.Info("topic already exists", "id", existing.ID)
		return existing.ID, nil
	}

	topic := domain.Topic{
		Name:        "IPA Phonetics",
		Description: "International Phonetic Alphabet — the 44 sounds of English",
		Category:    "Language",
		Icon:        "🔤",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	res, err := coll.InsertOne(ctx, topic)
	if err != nil {
		return "", err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("unexpected inserted id type")
	}
	id := oid.Hex()
	slog.Info("topic created", "id", id)
	return id, nil
}

func seedItems(ctx context.Context, db *mongo.Database, topicID string) error {
	coll := db.Collection("items")

	count, err := coll.CountDocuments(ctx, bson.M{"topic_id": topicID})
	if err != nil {
		return err
	}
	if count > 0 {
		slog.Info("items already seeded", "count", count)
		return nil
	}

	now := time.Now()
	items := buildIPAItems(topicID, now)

	docs := make([]interface{}, len(items))
	for i, item := range items {
		docs[i] = item
	}

	res, err := coll.InsertMany(ctx, docs)
	if err != nil {
		return err
	}
	slog.Info("items seeded", "count", len(res.InsertedIDs))
	return nil
}

func buildIPAItems(topicID string, now time.Time) []domain.Item {
	type raw struct {
		term     string
		meaning  string
		example  string
		category string
	}

	rows := []raw{
		// Plosives
		{"/p/", "voiceless bilabial plosive", "pit", "Plosives"},
		{"/b/", "voiced bilabial plosive", "bit", "Plosives"},
		{"/t/", "voiceless alveolar plosive", "tip", "Plosives"},
		{"/d/", "voiced alveolar plosive", "dip", "Plosives"},
		{"/k/", "voiceless velar plosive", "cat", "Plosives"},
		{"/ɡ/", "voiced velar plosive", "gap", "Plosives"},
		// Fricatives
		{"/f/", "voiceless labiodental fricative", "fat", "Fricatives"},
		{"/v/", "voiced labiodental fricative", "vat", "Fricatives"},
		{"/θ/", "voiceless dental fricative", "thin", "Fricatives"},
		{"/ð/", "voiced dental fricative", "this", "Fricatives"},
		{"/s/", "voiceless alveolar fricative", "sat", "Fricatives"},
		{"/z/", "voiced alveolar fricative", "zap", "Fricatives"},
		{"/ʃ/", "voiceless postalveolar fricative", "ship", "Fricatives"},
		{"/ʒ/", "voiced postalveolar fricative", "vision", "Fricatives"},
		{"/h/", "voiceless glottal fricative", "hat", "Fricatives"},
		// Affricates
		{"/tʃ/", "voiceless postalveolar affricate", "chip", "Affricates"},
		{"/dʒ/", "voiced postalveolar affricate", "jam", "Affricates"},
		// Nasals
		{"/m/", "voiced bilabial nasal", "map", "Nasals"},
		{"/n/", "voiced alveolar nasal", "nap", "Nasals"},
		{"/ŋ/", "voiced velar nasal", "sing", "Nasals"},
		// Approximants
		{"/l/", "voiced alveolar lateral approximant", "lip", "Approximants"},
		{"/r/", "voiced alveolar approximant", "rip", "Approximants"},
		{"/j/", "voiced palatal approximant", "yes", "Approximants"},
		{"/w/", "voiced bilabial approximant", "wet", "Approximants"},
		// Short Vowels
		{"/ɪ/", "near-close near-front unrounded vowel", "bit", "Short Vowels"},
		{"/e/", "close-mid front unrounded vowel", "bet", "Short Vowels"},
		{"/æ/", "near-open front unrounded vowel", "bat", "Short Vowels"},
		{"/ʌ/", "open-mid back unrounded vowel", "but", "Short Vowels"},
		{"/ɒ/", "open back rounded vowel", "bot", "Short Vowels"},
		{"/ʊ/", "near-close near-back rounded vowel", "book", "Short Vowels"},
		{"/ə/", "schwa — mid central vowel", "about", "Short Vowels"},
		// Long Vowels
		{"/iː/", "close front unrounded vowel (long)", "beat", "Long Vowels"},
		{"/ɑː/", "open back unrounded vowel (long)", "bar", "Long Vowels"},
		{"/ɔː/", "open-mid back rounded vowel (long)", "bore", "Long Vowels"},
		{"/uː/", "close back rounded vowel (long)", "boot", "Long Vowels"},
		{"/ɜː/", "open-mid central unrounded vowel (long)", "bird", "Long Vowels"},
		// Diphthongs
		{"/eɪ/", "closing diphthong from mid to close-front", "bait", "Diphthongs"},
		{"/aɪ/", "closing diphthong from open to close-front", "bite", "Diphthongs"},
		{"/ɔɪ/", "closing diphthong from open-mid back to close-front", "boy", "Diphthongs"},
		{"/əʊ/", "closing diphthong from mid to close-back", "boat", "Diphthongs"},
		{"/aʊ/", "closing diphthong from open to close-back", "bout", "Diphthongs"},
		{"/ɪə/", "centering diphthong to schwa", "beer", "Diphthongs"},
		{"/eə/", "centering diphthong from mid to schwa", "bear", "Diphthongs"},
		{"/ʊə/", "centering diphthong from near-back to schwa", "tour", "Diphthongs"},
	}

	items := make([]domain.Item, len(rows))
	for i, r := range rows {
		items[i] = domain.Item{
			TopicID:   topicID,
			Term:      r.term,
			Meaning:   r.meaning,
			Example:   r.example,
			Category:  r.category,
			SortOrder: i + 1,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	return items
}
