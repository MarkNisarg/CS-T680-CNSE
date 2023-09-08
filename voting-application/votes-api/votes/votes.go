package votes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

const (
	RedisNilError        = "redis: nil"
	RedisDefaultLocation = "redis:6379"
	RedisKeyPrefix       = "votes:"
)

// Vote represents a voter who voted in poll with vote value.
type Vote struct {
	VoteID    uint `json:"voteId"`
	VoterID   uint `json:"voterId"`
	PollID    uint `json:"pollId"`
	VoteValue uint `json:"voteValue"`
}

type cache struct {
	cacheClient *redis.Client
	jsonHelper  *rejson.Handler
	context     context.Context
}

// The reference to a cache object.
type VotesCache struct {
	cache
}

// The constructor function that returns a pointer to a new VotesCache.
// It uses the default Redis URL with the companion constructor newVotesCacheInstance.
func NewVotesCache() (*VotesCache, error) {
	redisUrl := os.Getenv("REDIS_URL")

	if redisUrl == "" {
		redisUrl = RedisDefaultLocation
	}

	return newVotesCacheInstance(redisUrl)
}

// The constructor function that returns a pointer to a new VotesCache.
func newVotesCacheInstance(url string) (*VotesCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: url,
	})

	ctx := context.Background()

	err := client.Ping(ctx).Err()
	if err != nil {
		log.Println("Error connecting to redis" + err.Error())
		return nil, err
	}

	jsonHelper := rejson.NewReJSONHandler()
	jsonHelper.SetGoRedisClientWithContext(ctx, client)

	return &VotesCache{
		cache: cache{
			cacheClient: client,
			jsonHelper:  jsonHelper,
			context:     ctx,
		},
	}, nil
}

// Get a string that can be used as a key in redis.
func redisKeyFromId(id uint) string {
	return fmt.Sprintf("%s%d", RedisKeyPrefix, id)
}

// Helper to return a Vote from redis provided a key.
func (vc *VotesCache) getItemFromRedis(key string, item *Vote) error {
	itemObject, err := vc.jsonHelper.JSONGet(key, ".")
	if err != nil {
		return err
	}

	err = json.Unmarshal(itemObject.([]byte), item)
	if err != nil {
		return err
	}

	return nil
}

// Return a slice of all votes from the VotesCache.
func (vc *VotesCache) GetAllVotes() ([]Vote, error) {
	var votes []Vote

	pattern := fmt.Sprintf("%s*", RedisKeyPrefix)
	keys, err := vc.cacheClient.Keys(vc.context, pattern).Result()

	if err != nil {
		return votes, err
	}

	for _, key := range keys {
		var vote Vote

		err := vc.getItemFromRedis(key, &vote)
		if err != nil {
			return votes, err
		}

		votes = append(votes, vote)
	}

	return votes, nil
}

// Retrieve a single vote from the VotesCache by voteID.
func (vc *VotesCache) GetVote(voteID uint) (Vote, error) {
	var vote Vote

	redisKey := redisKeyFromId(voteID)
	if err := vc.getItemFromRedis(redisKey, &vote); err != nil {
		return Vote{}, errors.New("vote does not exist")
	}

	return vote, nil
}

// Add a new vote to the VotesCache.
func (vc *VotesCache) AddVote(vote Vote) error {
	if _, err := vc.GetVote(vote.VoteID); err == nil {
		return errors.New("vote already exists")
	}

	redisKey := redisKeyFromId(vote.VoteID)
	if _, setErr := vc.jsonHelper.JSONSet(redisKey, ".", vote); setErr != nil {
		return setErr
	}

	return nil
}

// Delete a single vote from the VotesCache by voteID.
func (vc *VotesCache) DeleteVote(voteID uint) error {
	if _, err := vc.GetVote(voteID); err != nil {
		return errors.New("vote does not exist")
	}

	redisKey := redisKeyFromId(voteID)
	if _, deleteErr := vc.jsonHelper.JSONDel(redisKey, "."); deleteErr != nil {
		return deleteErr
	}

	return nil
}
