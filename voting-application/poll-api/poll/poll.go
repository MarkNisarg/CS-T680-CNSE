package poll

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
	RedisKeyPrefix       = "poll:"
)

// pollOptions represents the poll information for a specific poll.
type pollOption struct {
	PollOptionID   uint   `json:"pollOptionId"`
	PollOptionText string `json:"pollOptionText"`
}

// Poll represents a poll with a unique ID and poll information.
type Poll struct {
	PollID       uint         `json:"pollId"`
	PollTitle    string       `json:"pollTitle"`
	PollQuestion string       `json:"pollQuestion"`
	PollOptions  []pollOption `json:"pollOptions"`
}

type cache struct {
	cacheClient *redis.Client
	jsonHelper  *rejson.Handler
	context     context.Context
}

// The reference to a cache object.
type PollCache struct {
	cache
}

// The constructor function that returns a pointer to a new PollCache.
// It uses the default Redis URL with the companion constructor newPollCacheInstance.
func NewPollCache() (*PollCache, error) {
	redisUrl := os.Getenv("REDIS_URL")

	if redisUrl == "" {
		redisUrl = RedisDefaultLocation
	}

	return newPollCacheInstance(redisUrl)
}

// The constructor function that returns a pointer to a new PollCache.
func newPollCacheInstance(url string) (*PollCache, error) {
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

	return &PollCache{
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

// Create a new Poll instance with the provided details.
func NewPoll(pollId uint, pollTitle string, pollQuestion string) Poll {
	poll := Poll{
		PollID:       pollId,
		PollTitle:    pollTitle,
		PollQuestion: pollQuestion,
		PollOptions:  make([]pollOption, 0),
	}

	return poll
}

// Helper to return a Poll from redis provided a key.
func (pc *PollCache) getItemFromRedis(key string, item *Poll) error {
	itemObject, err := pc.jsonHelper.JSONGet(key, ".")
	if err != nil {
		return err
	}

	err = json.Unmarshal(itemObject.([]byte), item)
	if err != nil {
		return err
	}

	return nil
}

// Return a slice of all polls from the PollCache.
func (pc *PollCache) GetAllPolls() ([]Poll, error) {
	var polls []Poll

	pattern := fmt.Sprintf("%s*", RedisKeyPrefix)
	keys, err := pc.cacheClient.Keys(pc.context, pattern).Result()

	if err != nil {
		return polls, err
	}

	for _, key := range keys {
		var poll Poll

		err := pc.getItemFromRedis(key, &poll)
		if err != nil {
			return polls, err
		}

		polls = append(polls, poll)
	}

	return polls, nil
}

// Retrieve a single poll from the PollCache by pollId.
func (pc *PollCache) GetPoll(pollID uint) (Poll, error) {
	var poll Poll

	redisKey := redisKeyFromId(pollID)
	if err := pc.getItemFromRedis(redisKey, &poll); err != nil {
		return Poll{}, errors.New("poll does not exist")
	}

	return poll, nil
}

// Add a poll to the PollCache.
func (pc *PollCache) AddPoll(poll Poll) error {
	if _, err := pc.GetPoll(poll.PollID); err == nil {
		return errors.New("poll already exists")
	}

	redisKey := redisKeyFromId(poll.PollID)
	if _, setErr := pc.jsonHelper.JSONSet(redisKey, ".", poll); setErr != nil {
		return setErr
	}

	return nil
}

// Delete all polls from the PollCache.
func (pc *PollCache) DeleteAllPolls() error {
	pattern := fmt.Sprintf("%s*", RedisKeyPrefix)
	keys, err := pc.cacheClient.Keys(pc.context, pattern).Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		if _, deleteErr := pc.cacheClient.Del(pc.context, key).Result(); deleteErr != nil {
			return deleteErr
		}
	}

	return nil
}

// Delete a single poll from the PollCache by pollID.
func (pc *PollCache) DeletePoll(pollID uint) error {
	if _, err := pc.GetPoll(pollID); err != nil {
		return errors.New("poll does not exist")
	}

	redisKey := redisKeyFromId(pollID)
	if _, deleteErr := pc.jsonHelper.JSONDel(redisKey, "."); deleteErr != nil {
		return deleteErr
	}

	return nil
}

// Retrieve the poll options of a poll by pollID.
func (pc *PollCache) GetPollOptions(pollID uint) ([]pollOption, error) {
	poll, err := pc.GetPoll(pollID)
	if err != nil {
		return nil, errors.New("poll does not exist")
	}

	return poll.PollOptions, nil
}

// Retrieve a specific poll option by pollID and pollOptionID.
func (pc *PollCache) GetPollOption(pollID, pollOptionID uint) (pollOption, error) {
	poll, err := pc.GetPoll(pollID)
	if err != nil {
		return pollOption{}, errors.New("poll does not exist")
	}

	for _, option := range poll.PollOptions {
		if option.PollOptionID == pollOptionID {
			return option, nil
		}
	}

	return pollOption{}, errors.New("poll option not found")
}

// Add a new poll option to the poll options of a poll.
func (pc *PollCache) AddPollOption(pollID, pollOptionID uint, pollOptionText string) (pollOption, error) {
	poll, err := pc.GetPoll(pollID)
	if err != nil {
		return pollOption{}, errors.New("poll does not exist")
	}

	for _, option := range poll.PollOptions {
		if option.PollOptionID == pollOptionID {
			return pollOption{}, errors.New("poll option has already in poll")
		}
	}

	newPollOption := pollOption{
		PollOptionID:   pollOptionID,
		PollOptionText: pollOptionText,
	}

	poll.PollOptions = append(poll.PollOptions, newPollOption)

	redisKey := redisKeyFromId(poll.PollID)
	if _, setErr := pc.jsonHelper.JSONSet(redisKey, ".", poll); setErr != nil {
		return pollOption{}, setErr
	}

	return newPollOption, nil
}

// Remove a specific poll option from the poll options of a poll.
func (pc *PollCache) DeletePollOption(pollID, pollOptionID uint) error {
	poll, err := pc.GetPoll(pollID)
	if err != nil {
		return errors.New("poll does not exist")
	}

	var updatedPollOptions []pollOption
	for _, pollOpt := range poll.PollOptions {
		if pollOpt.PollOptionID != pollOptionID {
			updatedPollOptions = append(updatedPollOptions, pollOpt)
		}
	}

	if len(updatedPollOptions) == len(poll.PollOptions) {
		return errors.New("poll option not found")
	}

	poll.PollOptions = updatedPollOptions
	redisKey := redisKeyFromId(poll.PollID)
	if _, setErr := pc.jsonHelper.JSONSet(redisKey, ".", poll); setErr != nil {
		return setErr
	}

	return nil
}
