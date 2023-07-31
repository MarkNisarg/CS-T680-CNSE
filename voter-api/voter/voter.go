package voter

import (
	"errors"
	"time"
)

// voterPoll represents the voting information for a specific poll.
type voterPoll struct {
	PollID   uint      `json:"pollId"`
	VoteDate time.Time `json:"voteDate"`
}

// Voter represents a voter with a unique ID and voting history.
type Voter struct {
	VoterID     uint        `json:"voterId"`
	FirstName   string      `json:"firstName"`
	LastName    string      `json:"lastName"`
	VoteHistory []voterPoll `json:"voteHistory"`
}

// VoterList is a collection of voters.
type VoterList struct {
	Voters map[uint]Voter
}

// Create a new VoterList instance and initializes the Voters map.
func NewVoterList() *VoterList {
	voterList := &VoterList{
		Voters: make(map[uint]Voter),
	}

	return voterList
}

// Create a new Voter instance with the provided details.
func NewVoter(id uint, firstName string, lastName string) Voter {
	voter := Voter{
		VoterID:     id,
		FirstName:   firstName,
		LastName:    lastName,
		VoteHistory: make([]voterPoll, 0),
	}

	return voter
}

// Return a slice of all voters in the VoterList.
func (vl *VoterList) GetAllVoters() []Voter {
	var voters []Voter

	for _, voter := range vl.Voters {
		voters = append(voters, voter)
	}

	return voters
}

// Retrieve a single voter from the VoterList by voterID.
func (vl *VoterList) GetVoter(voterID uint) (Voter, error) {
	voter, exists := vl.Voters[voterID]
	if !exists {
		return Voter{}, errors.New("voter does not exist")
	}

	return voter, nil
}

// Add a new voter to the VoterList.
func (vl *VoterList) AddVoter(voter Voter) error {
	if _, exists := vl.Voters[voter.VoterID]; exists {
		return errors.New("voter already exists")
	}

	vl.Voters[voter.VoterID] = voter

	return nil
}

// Update an existing voter in the VoterList.
func (vl *VoterList) UpdateVoter(voter Voter) error {
	existingVoter, exists := vl.Voters[voter.VoterID]

	if !exists {
		return errors.New("voter does not exist")
	}

	existingVoter.FirstName = voter.FirstName
	existingVoter.LastName = voter.LastName

	vl.Voters[voter.VoterID] = existingVoter

	return nil
}

// Delete all voters from the VoterList.
func (vl *VoterList) DeleteAllVoters() error {
	vl.Voters = make(map[uint]Voter)

	return nil
}

// Delete a single voter from the VoterList by voterID.
func (vl *VoterList) DeleteVoter(voterID uint) error {
	if _, exists := vl.Voters[voterID]; !exists {
		return errors.New("voter does not exist")
	}

	delete(vl.Voters, voterID)

	return nil
}

// Retrieve the vote history of a voter by voterID.
func (vl *VoterList) GetVoterHistory(voterID uint) ([]voterPoll, error) {
	voter, exists := vl.Voters[voterID]
	if !exists {
		return nil, errors.New("voter does not exist")
	}

	return voter.VoteHistory, nil
}

// Retrieve a specific voter poll by voterID and pollID.
func (vl *VoterList) GetVoterPoll(voterID, pollID uint) (voterPoll, error) {
	voter, exists := vl.Voters[voterID]
	if !exists {
		return voterPoll{}, errors.New("voter does not exist")
	}

	for _, poll := range voter.VoteHistory {
		if poll.PollID == pollID {
			return poll, nil
		}
	}

	return voterPoll{}, errors.New("voter poll not found")
}

// Add a new voter poll to the vote history of a voter.
func (vl *VoterList) AddVoterPoll(voterID, pollID uint) (voterPoll, error) {
	voter, exists := vl.Voters[voterID]
	if !exists {
		return voterPoll{}, errors.New("voter does not exist")
	}

	for _, poll := range voter.VoteHistory {
		if poll.PollID == pollID {
			return voterPoll{}, errors.New("voter has already voted in this poll")
		}
	}

	newVoterPoll := voterPoll{
		PollID:   pollID,
		VoteDate: time.Now(),
	}

	voter.VoteHistory = append(voter.VoteHistory, newVoterPoll)
	vl.Voters[voter.VoterID] = voter

	return newVoterPoll, nil
}

// Update an existing voter poll in the vote history of a voter.
func (vl *VoterList) UpdateVoterPoll(voterID, pollID uint) (voterPoll, error) {
	voter, exists := vl.Voters[voterID]
	if !exists {
		return voterPoll{}, errors.New("voter does not exist")
	}

	var updatedVoterPoll voterPoll
	for i, poll := range voter.VoteHistory {
		if poll.PollID == pollID {
			updatedVoterPoll = voterPoll{
				PollID:   pollID,
				VoteDate: time.Now(),
			}
			voter.VoteHistory[i] = updatedVoterPoll
			break
		}
	}

	if updatedVoterPoll.PollID == 0 {
		return voterPoll{}, errors.New("voter poll not found")
	}

	vl.Voters[voter.VoterID] = voter

	return updatedVoterPoll, nil
}

// Remove a specific voter poll from the vote history of a voter.
func (vl *VoterList) DeleteVoterPoll(voterID, pollID uint) error {
	voter, exists := vl.Voters[voterID]
	if !exists {
		return errors.New("voter does not exist")
	}

	var updatedVoteHistory []voterPoll
	for _, poll := range voter.VoteHistory {
		if poll.PollID != pollID {
			updatedVoteHistory = append(updatedVoteHistory, poll)
		}
	}

	if len(updatedVoteHistory) == len(voter.VoteHistory) {
		return errors.New("voter poll not found")
	}

	voter.VoteHistory = updatedVoteHistory
	vl.Voters[voter.VoterID] = voter

	return nil
}
