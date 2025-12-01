// Copyright (c) 2025 Trilitech Ltd.
// Author: tzstats@trili.tech

package rpc

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trilitech/tzgo/tezos"
)

func TestParseAttestationsAggregate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chains/main/blocks/BMABzWp5Y3iSJRaCkWVwsPKXVZ1iCwB94dB7GfKsigahQ3v5Czc/operations" {
			t.Errorf("Expected to request '/chains/main/blocks/BMABzWp5Y3iSJRaCkWVwsPKXVZ1iCwB94dB7GfKsigahQ3v5Czc/operations', got: %s", r.URL.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[[{"protocol":"PtSeouLouXkxhg39oWzjxDWaCydNfR3RxCUrNe4Q9Ro8BTehcbh","chain_id":"NetXd56aBs1aeW3","hash":"ooqgVxC8XDYHXSnznWSdkXktgUaA2PZdUs4azRbco9i6fhiY1ui","branch":"BLwNFCbHuF21bF4S9KybZd52si5QQG6k29mQTshbY7fVnehrxbh","contents":[{"kind":"attestations_aggregate","consensus_content":{"level":109771,"round":0,"block_payload_hash":"vh3RiisbNp7QBvLkJ6HAyoYfMh4AoZLcJFVE5M34LLKnwAWTRv6J"},"committee":[{"slot":0,"dal_attestation":"0"},{"slot":4,"dal_attestation":"0"},{"slot":15},{"slot":21,"dal_attestation":"0"},{"slot":203}],"metadata":{"committee":[{"delegate":"tz1NNT9EERmcKekRq2vdv6e8TL3WQpY8AXSF","consensus_pkh":"tz4X5GCfEHQCUnrBy9Qo1PSsmExYHXxiEkvp","consensus_power":1435},{"delegate":"tz1Zt8QQ9aBznYNk5LUBjtME9DuExomw9YRs","consensus_pkh":"tz4XbGtqxNZDq6CJNVbxktoqSTu9Db6aXQHL","consensus_power":1465},{"delegate":"tz4MvCEiLgK6vYRSkug9Nz64UNTbT4McNq8m","consensus_pkh":"tz4MvCEiLgK6vYRSkug9Nz64UNTbT4McNq8m","consensus_power":707},{"delegate":"tz3PgFHdYvEGEbUo1pUJmuNH8fgc8cwARKfC","consensus_pkh":"tz4EWkmNN93yE7HrjaRR6mGh22rUgSYJG1Sj","consensus_power":730},{"delegate":"tz1LmrwzCKUDibk7xaGC5RxvTbmUbCAtCA4a","consensus_pkh":"tz4QqpgVTG4CCETKMgV6YaUHyowE4Gkqwdfi","consensus_power":1}],"total_consensus_power":4338}}],"signature":"BLsigB1uDeiiuW1NPKWsZ6WRAKL3aSGTXKmtsDH2xxWgWqM3QBr3mFW8QqzH6VWGsvPGsrii3VVw7KA9CvC9LjC3VxH3MgHSvcWVK6Z7rBbEY79sKXi4XrbbfY8QJpE38B4u6mteGKHnVj"}],[],[],[]]`))
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	value, e := c.GetBlockOperations(context.TODO(), tezos.MustParseBlockHash("BMABzWp5Y3iSJRaCkWVwsPKXVZ1iCwB94dB7GfKsigahQ3v5Czc"))
	assert.Nil(t, e)
	assert.Len(t, value, 4)
	assert.Len(t, value[0], 1)
	assert.Equal(t, value[0][0].Contents.Len(), 1)
	op := value[0][0].Contents.N(0).(*AttestationsAggregate)
	assert.Equal(t, tezos.OpTypeAttestationsAggregate, op.Kind())
	assert.Equal(t, 5, len(op.Committee))
	assert.Equal(t, ConsensusContent{Level: 109771, Round: 0, PayloadHash: tezos.MustParsePayloadHash("vh3RiisbNp7QBvLkJ6HAyoYfMh4AoZLcJFVE5M34LLKnwAWTRv6J")}, op.ConsensusContent)

	// Verify committee metadata is parsed correctly
	metadata := op.Meta()
	assert.Len(t, metadata.CommitteeMetadata, 5, "should have 5 committee members in metadata")
	v, err := metadata.TotalConsensusPower.AsV023Value()
	assert.Nil(t, err)
	assert.Equal(t, 4338, v, "total consensus power should be 4338")

	// Verify first committee member
	assert.Equal(t, tezos.MustParseAddress("tz1NNT9EERmcKekRq2vdv6e8TL3WQpY8AXSF"), metadata.CommitteeMetadata[0].Delegate)
	assert.Equal(t, tezos.MustParseAddress("tz4X5GCfEHQCUnrBy9Qo1PSsmExYHXxiEkvp"), metadata.CommitteeMetadata[0].ConsensusPkh)
	v, err = metadata.CommitteeMetadata[0].ConsensusPower.AsV023Value()
	assert.Nil(t, err)
	assert.Equal(t, 1435, v)

	// Verify second committee member
	assert.Equal(t, tezos.MustParseAddress("tz1Zt8QQ9aBznYNk5LUBjtME9DuExomw9YRs"), metadata.CommitteeMetadata[1].Delegate)
	assert.Equal(t, tezos.MustParseAddress("tz4XbGtqxNZDq6CJNVbxktoqSTu9Db6aXQHL"), metadata.CommitteeMetadata[1].ConsensusPkh)
	v, err = metadata.CommitteeMetadata[1].ConsensusPower.AsV023Value()
	assert.Nil(t, err)
	assert.Equal(t, 1465, v)

	// Verify tz4 delegate (matching delegate and consensus_pkh)
	assert.Equal(t, tezos.MustParseAddress("tz4MvCEiLgK6vYRSkug9Nz64UNTbT4McNq8m"), metadata.CommitteeMetadata[2].Delegate)
	assert.Equal(t, tezos.MustParseAddress("tz4MvCEiLgK6vYRSkug9Nz64UNTbT4McNq8m"), metadata.CommitteeMetadata[2].ConsensusPkh)
	v, err = metadata.CommitteeMetadata[2].ConsensusPower.AsV023Value()
	assert.Nil(t, err)
	assert.Equal(t, 707, v)
}

func TestParseAttestationsAggregateV024(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chains/main/blocks/BKpbfCvh777DQHnXjU2sqHvVUNZ7dBAdqEfKkdw8EGSkD9LSYXb/operations" {
			t.Errorf("Expected to request '/chains/main/blocks/BKpbfCvh777DQHnXjU2sqHvVUNZ7dBAdqEfKkdw8EGSkD9LSYXb/operations', got: %s", r.URL.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[[{"protocol":"PtTALLiNtPec7mE7yY4m3k26J8Qukef3E3ehzhfXgFZKGtDdAXu","chain_id":"NetXd56aBs1aeW3","hash":"ooqgVxC8XDYHXSnznWSdkXktgUaA2PZdUs4azRbco9i6fhiY1ui","branch":"BKpbfCvh777DQHnXjU2sqHvVUNZ7dBAdqEfKkdw8EGSkD9LSYXb","contents":[{"kind":"attestations_aggregate","consensus_content":{"level":109771,"round":0,"block_payload_hash":"vh3RiisbNp7QBvLkJ6HAyoYfMh4AoZLcJFVE5M34LLKnwAWTRv6J"},"committee":[{"slot":0,"dal_attestation":"0"},{"slot":4,"dal_attestation":"0"},{"slot":15},{"slot":21,"dal_attestation":"0"},{"slot":203}],"metadata":{"committee":[{"delegate":"tz1NNT9EERmcKekRq2vdv6e8TL3WQpY8AXSF","consensus_pkh":"tz4X5GCfEHQCUnrBy9Qo1PSsmExYHXxiEkvp","consensus_power":{"slots":1575,"baking_power":"60261130120715"}},{"delegate":"tz1Zt8QQ9aBznYNk5LUBjtME9DuExomw9YRs","consensus_pkh":"tz4XbGtqxNZDq6CJNVbxktoqSTu9Db6aXQHL","consensus_power":{"slots":47,"baking_power":"2038150029478"}},{"delegate":"tz4MvCEiLgK6vYRSkug9Nz64UNTbT4McNq8m","consensus_pkh":"tz4MvCEiLgK6vYRSkug9Nz64UNTbT4McNq8m","consensus_power":{"slots":9,"baking_power":"492207228099"}},{"delegate":"tz3PgFHdYvEGEbUo1pUJmuNH8fgc8cwARKfC","consensus_pkh":"tz4EWkmNN93yE7HrjaRR6mGh22rUgSYJG1Sj","consensus_power":{"slots":447,"baking_power":"18010566072999"}},{"delegate":"tz1LmrwzCKUDibk7xaGC5RxvTbmUbCAtCA4a","consensus_pkh":"tz4QqpgVTG4CCETKMgV6YaUHyowE4Gkqwdfi","consensus_power":{"slots":1,"baking_power":"6083478248"}}],"total_consensus_power":{"slots":2079,"baking_power":"80808136929539"}}}],"signature":"BLsigB1uDeiiuW1NPKWsZ6WRAKL3aSGTXKmtsDH2xxWgWqM3QBr3mFW8QqzH6VWGsvPGsrii3VVw7KA9CvC9LjC3VxH3MgHSvcWVK6Z7rBbEY79sKXi4XrbbfY8QJpE38B4u6mteGKHnVj"}],[],[],[]]`))
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	value, e := c.GetBlockOperations(context.TODO(), tezos.MustParseBlockHash("BKpbfCvh777DQHnXjU2sqHvVUNZ7dBAdqEfKkdw8EGSkD9LSYXb"))
	assert.Nil(t, e)
	assert.Len(t, value, 4)
	assert.Len(t, value[0], 1)
	assert.Equal(t, value[0][0].Contents.Len(), 1)
	op := value[0][0].Contents.N(0).(*AttestationsAggregate)
	assert.Equal(t, tezos.OpTypeAttestationsAggregate, op.Kind())
	assert.Equal(t, 5, len(op.Committee))
	assert.Equal(t, ConsensusContent{Level: 109771, Round: 0, PayloadHash: tezos.MustParsePayloadHash("vh3RiisbNp7QBvLkJ6HAyoYfMh4AoZLcJFVE5M34LLKnwAWTRv6J")}, op.ConsensusContent)

	// Verify committee metadata is parsed correctly
	metadata := op.Meta()
	assert.Len(t, metadata.CommitteeMetadata, 5, "should have 5 committee members in metadata")
	v, err := metadata.TotalConsensusPower.AsV024Value()
	assert.Nil(t, err)
	assert.Equal(t, 2079, v.Slots)
	assert.Equal(t, int64(80808136929539), v.BakingPower)

	// Verify first committee member
	assert.Equal(t, tezos.MustParseAddress("tz1NNT9EERmcKekRq2vdv6e8TL3WQpY8AXSF"), metadata.CommitteeMetadata[0].Delegate)
	assert.Equal(t, tezos.MustParseAddress("tz4X5GCfEHQCUnrBy9Qo1PSsmExYHXxiEkvp"), metadata.CommitteeMetadata[0].ConsensusPkh)
	v, err = metadata.CommitteeMetadata[0].ConsensusPower.AsV024Value()
	assert.Nil(t, err)
	assert.Equal(t, 1575, v.Slots)
	assert.Equal(t, int64(60261130120715), v.BakingPower)

	// Verify second committee member
	assert.Equal(t, tezos.MustParseAddress("tz1Zt8QQ9aBznYNk5LUBjtME9DuExomw9YRs"), metadata.CommitteeMetadata[1].Delegate)
	assert.Equal(t, tezos.MustParseAddress("tz4XbGtqxNZDq6CJNVbxktoqSTu9Db6aXQHL"), metadata.CommitteeMetadata[1].ConsensusPkh)
	v, err = metadata.CommitteeMetadata[1].ConsensusPower.AsV024Value()
	assert.Nil(t, err)
	assert.Equal(t, 47, v.Slots)
	assert.Equal(t, int64(2038150029478), v.BakingPower)

	// Verify tz4 delegate (matching delegate and consensus_pkh)
	assert.Equal(t, tezos.MustParseAddress("tz4MvCEiLgK6vYRSkug9Nz64UNTbT4McNq8m"), metadata.CommitteeMetadata[2].Delegate)
	assert.Equal(t, tezos.MustParseAddress("tz4MvCEiLgK6vYRSkug9Nz64UNTbT4McNq8m"), metadata.CommitteeMetadata[2].ConsensusPkh)
	v, err = metadata.CommitteeMetadata[2].ConsensusPower.AsV024Value()
	assert.Nil(t, err)
	assert.Equal(t, 9, v.Slots)
	assert.Equal(t, int64(492207228099), v.BakingPower)
}

func TestParsePreattestationsAggregate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chains/main/blocks/BMABzWp5Y3iSJRaCkWVwsPKXVZ1iCwB94dB7GfKsigahQ3v5Czc/operations" {
			t.Errorf("Expected to request '/chains/main/blocks/BMABzWp5Y3iSJRaCkWVwsPKXVZ1iCwB94dB7GfKsigahQ3v5Czc/operations', got: %s", r.URL.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[[{"protocol":"PtSeouLouXkxhg39oWzjxDWaCydNfR3RxCUrNe4Q9Ro8BTehcbh","chain_id":"NetXd56aBs1aeW3","hash":"op1BDoVeq9WxXcJeaammBZwBb7snpGCQ2WQAyU7kKRYgivtezPX","branch":"BLYpjw24Ad8Y1LmvhXPy4ac8171sqbHzKNiJ2xLdmktEt5RNuWN","contents":[{"kind":"preattestations_aggregate","consensus_content":{"level":109772,"round":0,"block_payload_hash":"vh1jC8jsjWdJVETCNHiUdk9oto8aE88zGSz8aJCcF7PGxWVv3TgH"},"committee":[0,3,4,10],"metadata":{"committee":[{"delegate":"tz1NNT9EERmcKekRq2vdv6e8TL3WQpY8AXSF","consensus_pkh":"tz4X5GCfEHQCUnrBy9Qo1PSsmExYHXxiEkvp","consensus_power":1485},{"delegate":"tz4MvCEiLgK6vYRSkug9Nz64UNTbT4McNq8m","consensus_pkh":"tz4MvCEiLgK6vYRSkug9Nz64UNTbT4McNq8m","consensus_power":688},{"delegate":"tz1Zt8QQ9aBznYNk5LUBjtME9DuExomw9YRs","consensus_pkh":"tz4XbGtqxNZDq6CJNVbxktoqSTu9Db6aXQHL","consensus_power":1498},{"delegate":"tz3PgFHdYvEGEbUo1pUJmuNH8fgc8cwARKfC","consensus_pkh":"tz4EWkmNN93yE7HrjaRR6mGh22rUgSYJG1Sj","consensus_power":700}],"total_consensus_power":4371}}],"signature":"BLsigAFfUNms9mHTHxXkXH1YWGSjtuZQYtUfN2iW3bLUP6MnSZBQpxFL8CCzqi3K9DkejeDqbBpUAxnnWPQ96dd62jWfwUfUnQhuJ34wcDdG8ZJQiXmCiEuDyRDxWXQfnEBjpXbtdWpwuz"}],[],[],[]]`))
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	value, e := c.GetBlockOperations(context.TODO(), tezos.MustParseBlockHash("BMABzWp5Y3iSJRaCkWVwsPKXVZ1iCwB94dB7GfKsigahQ3v5Czc"))
	assert.Nil(t, e)
	assert.Len(t, value, 4)
	assert.Len(t, value[0], 1)
	assert.Equal(t, value[0][0].Contents.Len(), 1)
	op := value[0][0].Contents.N(0).(*PreattestationsAggregate)
	assert.Equal(t, tezos.OpTypePreattestationsAggregate, op.Kind())
	assert.Equal(t, []int{0, 3, 4, 10}, op.Committee)
	assert.Equal(t, ConsensusContent{Level: 109772, Round: 0, PayloadHash: tezos.MustParsePayloadHash("vh1jC8jsjWdJVETCNHiUdk9oto8aE88zGSz8aJCcF7PGxWVv3TgH")}, op.ConsensusContent)

	// Verify committee metadata is parsed correctly
	metadata := op.Meta()
	assert.Len(t, metadata.CommitteeMetadata, 4, "should have 4 committee members in metadata")
	v, err := metadata.TotalConsensusPower.AsV023Value()
	assert.Nil(t, err)
	assert.Equal(t, 4371, v, "total consensus power should be 4371")

	// Verify first committee member
	assert.Equal(t, tezos.MustParseAddress("tz1NNT9EERmcKekRq2vdv6e8TL3WQpY8AXSF"), metadata.CommitteeMetadata[0].Delegate)
	assert.Equal(t, tezos.MustParseAddress("tz4X5GCfEHQCUnrBy9Qo1PSsmExYHXxiEkvp"), metadata.CommitteeMetadata[0].ConsensusPkh)
	v, err = metadata.CommitteeMetadata[0].ConsensusPower.AsV023Value()
	assert.Nil(t, err)
	assert.Equal(t, 1485, v)

	// Verify tz4 delegate (matching delegate and consensus_pkh)
	assert.Equal(t, tezos.MustParseAddress("tz4MvCEiLgK6vYRSkug9Nz64UNTbT4McNq8m"), metadata.CommitteeMetadata[1].Delegate)
	assert.Equal(t, tezos.MustParseAddress("tz4MvCEiLgK6vYRSkug9Nz64UNTbT4McNq8m"), metadata.CommitteeMetadata[1].ConsensusPkh)
	v, err = metadata.CommitteeMetadata[1].ConsensusPower.AsV023Value()
	assert.Nil(t, err)
	assert.Equal(t, 688, v)
}

// TestAttestationsAggregateEmptyMetadata tests parsing when metadata fields are absent
func TestAttestationsAggregateEmptyMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// Response without metadata committee field
		w.Write([]byte(`[[{"protocol":"PtSeouLouXkxhg39oWzjxDWaCydNfR3RxCUrNe4Q9Ro8BTehcbh","chain_id":"NetXd56aBs1aeW3","hash":"ooqgVxC8XDYHXSnznWSdkXktgUaA2PZdUs4azRbco9i6fhiY1ui","branch":"BLwNFCbHuF21bF4S9KybZd52si5QQG6k29mQTshbY7fVnehrxbh","contents":[{"kind":"attestations_aggregate","consensus_content":{"level":109771,"round":0,"block_payload_hash":"vh3RiisbNp7QBvLkJ6HAyoYfMh4AoZLcJFVE5M34LLKnwAWTRv6J"},"committee":[{"slot":0,"dal_attestation":"0"}],"metadata":{}}],"signature":"BLsigB1uDeiiuW1NPKWsZ6WRAKL3aSGTXKmtsDH2xxWgWqM3QBr3mFW8QqzH6VWGsvPGsrii3VVw7KA9CvC9LjC3VxH3MgHSvcWVK6Z7rBbEY79sKXi4XrbbfY8QJpE38B4u6mteGKHnVj"}],[],[],[]]`))
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	value, e := c.GetBlockOperations(context.TODO(), tezos.MustParseBlockHash("BMABzWp5Y3iSJRaCkWVwsPKXVZ1iCwB94dB7GfKsigahQ3v5Czc"))
	assert.Nil(t, e)
	assert.Len(t, value, 4)
	assert.Len(t, value[0], 1)

	op := value[0][0].Contents.N(0).(*AttestationsAggregate)
	metadata := op.Meta()

	// Should have empty metadata fields when not provided
	assert.Len(t, metadata.CommitteeMetadata, 0, "should have no committee metadata")
	v, err := metadata.TotalConsensusPower.AsV023Value()
	assert.Nil(t, err)
	assert.Equal(t, 0, v, "total consensus power should be 0")
}

// TestAttestationsAggregateEmptyMetadataV024 tests parsing when metadata fields are absent
func TestAttestationsAggregateEmptyMetadataV024(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// Response without metadata committee field
		w.Write([]byte(`[[{"protocol":"PtTALLiNtPec7mE7yY4m3k26J8Qukef3E3ehzhfXgFZKGtDdAXu","chain_id":"NetXd56aBs1aeW3","hash":"ooqgVxC8XDYHXSnznWSdkXktgUaA2PZdUs4azRbco9i6fhiY1ui","branch":"BLwNFCbHuF21bF4S9KybZd52si5QQG6k29mQTshbY7fVnehrxbh","contents":[{"kind":"attestations_aggregate","consensus_content":{"level":109771,"round":0,"block_payload_hash":"vh3RiisbNp7QBvLkJ6HAyoYfMh4AoZLcJFVE5M34LLKnwAWTRv6J"},"committee":[{"slot":0,"dal_attestation":"0"}],"metadata":{}}],"signature":"BLsigB1uDeiiuW1NPKWsZ6WRAKL3aSGTXKmtsDH2xxWgWqM3QBr3mFW8QqzH6VWGsvPGsrii3VVw7KA9CvC9LjC3VxH3MgHSvcWVK6Z7rBbEY79sKXi4XrbbfY8QJpE38B4u6mteGKHnVj"}],[],[],[]]`))
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	value, e := c.GetBlockOperations(context.TODO(), tezos.MustParseBlockHash("BMABzWp5Y3iSJRaCkWVwsPKXVZ1iCwB94dB7GfKsigahQ3v5Czc"))
	assert.Nil(t, e)
	assert.Len(t, value, 4)
	assert.Len(t, value[0], 1)

	op := value[0][0].Contents.N(0).(*AttestationsAggregate)
	metadata := op.Meta()

	// Should have empty metadata fields when not provided
	assert.Len(t, metadata.CommitteeMetadata, 0, "should have no committee metadata")
	v, err := metadata.TotalConsensusPower.AsV024Value()
	assert.Nil(t, err)
	assert.Equal(t, 0, v.Slots)
	assert.Equal(t, int64(0), v.BakingPower)
}

// TestAttestationsAggregateSingleCommitteeMember tests parsing with one committee member
func TestAttestationsAggregateSingleCommitteeMember(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[[{"protocol":"PtSeouLouXkxhg39oWzjxDWaCydNfR3RxCUrNe4Q9Ro8BTehcbh","chain_id":"NetXd56aBs1aeW3","hash":"ooqgVxC8XDYHXSnznWSdkXktgUaA2PZdUs4azRbco9i6fhiY1ui","branch":"BLwNFCbHuF21bF4S9KybZd52si5QQG6k29mQTshbY7fVnehrxbh","contents":[{"kind":"attestations_aggregate","consensus_content":{"level":109771,"round":0,"block_payload_hash":"vh3RiisbNp7QBvLkJ6HAyoYfMh4AoZLcJFVE5M34LLKnwAWTRv6J"},"committee":[{"slot":0,"dal_attestation":"0"}],"metadata":{"committee":[{"delegate":"tz4MvCEiLgK6vYRSkug9Nz64UNTbT4McNq8m","consensus_pkh":"tz4MvCEiLgK6vYRSkug9Nz64UNTbT4McNq8m","consensus_power":5000}],"total_consensus_power":5000}}],"signature":"BLsigB1uDeiiuW1NPKWsZ6WRAKL3aSGTXKmtsDH2xxWgWqM3QBr3mFW8QqzH6VWGsvPGsrii3VVw7KA9CvC9LjC3VxH3MgHSvcWVK6Z7rBbEY79sKXi4XrbbfY8QJpE38B4u6mteGKHnVj"}],[],[],[]]`))
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	value, e := c.GetBlockOperations(context.TODO(), tezos.MustParseBlockHash("BMABzWp5Y3iSJRaCkWVwsPKXVZ1iCwB94dB7GfKsigahQ3v5Czc"))
	assert.Nil(t, e)

	op := value[0][0].Contents.N(0).(*AttestationsAggregate)
	metadata := op.Meta()

	assert.Len(t, metadata.CommitteeMetadata, 1, "should have 1 committee member")
	v, err := metadata.TotalConsensusPower.AsV023Value()
	assert.Nil(t, err)
	assert.Equal(t, 5000, v)
	assert.Equal(t, tezos.MustParseAddress("tz4MvCEiLgK6vYRSkug9Nz64UNTbT4McNq8m"), metadata.CommitteeMetadata[0].Delegate)
	assert.Equal(t, tezos.MustParseAddress("tz4MvCEiLgK6vYRSkug9Nz64UNTbT4McNq8m"), metadata.CommitteeMetadata[0].ConsensusPkh)
	v, err = metadata.CommitteeMetadata[0].ConsensusPower.AsV023Value()
	assert.Nil(t, err)
	assert.Equal(t, 5000, v)
}

// TestPreattestationsAggregateMetadataConsistency verifies metadata consistency between operations
func TestPreattestationsAggregateMetadataConsistency(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[[{"protocol":"PtSeouLouXkxhg39oWzjxDWaCydNfR3RxCUrNe4Q9Ro8BTehcbh","chain_id":"NetXd56aBs1aeW3","hash":"op1BDoVeq9WxXcJeaammBZwBb7snpGCQ2WQAyU7kKRYgivtezPX","branch":"BLYpjw24Ad8Y1LmvhXPy4ac8171sqbHzKNiJ2xLdmktEt5RNuWN","contents":[{"kind":"preattestations_aggregate","consensus_content":{"level":109772,"round":0,"block_payload_hash":"vh1jC8jsjWdJVETCNHiUdk9oto8aE88zGSz8aJCcF7PGxWVv3TgH"},"committee":[0,1],"metadata":{"committee":[{"delegate":"tz1NNT9EERmcKekRq2vdv6e8TL3WQpY8AXSF","consensus_pkh":"tz4X5GCfEHQCUnrBy9Qo1PSsmExYHXxiEkvp","consensus_power":2000},{"delegate":"tz1Zt8QQ9aBznYNk5LUBjtME9DuExomw9YRs","consensus_pkh":"tz4XbGtqxNZDq6CJNVbxktoqSTu9Db6aXQHL","consensus_power":3000}],"total_consensus_power":5000}}],"signature":"BLsigAFfUNms9mHTHxXkXH1YWGSjtuZQYtUfN2iW3bLUP6MnSZBQpxFL8CCzqi3K9DkejeDqbBpUAxnnWPQ96dd62jWfwUfUnQhuJ34wcDdG8ZJQiXmCiEuDyRDxWXQfnEBjpXbtdWpwuz"}],[],[],[]]`))
	}))
	defer server.Close()

	c, _ := NewClient(server.URL, nil)
	value, e := c.GetBlockOperations(context.TODO(), tezos.MustParseBlockHash("BMABzWp5Y3iSJRaCkWVwsPKXVZ1iCwB94dB7GfKsigahQ3v5Czc"))
	assert.Nil(t, e)

	op := value[0][0].Contents.N(0).(*PreattestationsAggregate)
	metadata := op.Meta()

	// Verify the sum of individual consensus powers equals total
	sum := 0
	for _, member := range metadata.CommitteeMetadata {
		v, err := member.ConsensusPower.AsV023Value()
		assert.Nil(t, err, "should not get error while extracting consensus power")
		sum += v
	}
	v, err := metadata.TotalConsensusPower.AsV023Value()
	assert.Nil(t, err)
	assert.Equal(t, v, sum, "sum of individual powers should equal total consensus power")

	// Verify all addresses are valid Tezos addresses
	for i, member := range metadata.CommitteeMetadata {
		assert.True(t, member.Delegate.IsValid(), "delegate address at index %d should be valid", i)
		assert.True(t, member.ConsensusPkh.IsValid(), "consensus_pkh at index %d should be valid", i)
		v, err := member.ConsensusPower.AsV023Value()
		assert.Nil(t, err)
		assert.Greater(t, v, 0, "consensus power at index %d should be positive", i)
	}
}
