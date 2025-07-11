package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v10/modules/core/exported"
)

var (
	_ codectypes.UnpackInterfacesMessage = (*QueryConnectionClientStateResponse)(nil)
	_ codectypes.UnpackInterfacesMessage = (*QueryConnectionConsensusStateResponse)(nil)
)

// NewQueryConnectionResponse creates a new QueryConnectionResponse instance
func NewQueryConnectionResponse(
	connection ConnectionEnd, proof []byte, height clienttypes.Height,
) *QueryConnectionResponse {
	return &QueryConnectionResponse{
		Connection:  &connection,
		Proof:       proof,
		ProofHeight: height,
	}
}

// NewQueryClientConnectionsResponse creates a new ConnectionPaths instance
func NewQueryClientConnectionsResponse(
	connectionPaths []string, proof []byte, height clienttypes.Height,
) *QueryClientConnectionsResponse {
	return &QueryClientConnectionsResponse{
		ConnectionPaths: connectionPaths,
		Proof:           proof,
		ProofHeight:     height,
	}
}

// NewQueryClientConnectionsRequest creates a new QueryClientConnectionsRequest instance
func NewQueryClientConnectionsRequest(clientID string) *QueryClientConnectionsRequest {
	return &QueryClientConnectionsRequest{
		ClientId: clientID,
	}
}

// NewQueryConnectionClientStateResponse creates a newQueryConnectionClientStateResponse instance
func NewQueryConnectionClientStateResponse(identifiedClientState clienttypes.IdentifiedClientState, proof []byte, height clienttypes.Height) *QueryConnectionClientStateResponse {
	return &QueryConnectionClientStateResponse{
		IdentifiedClientState: &identifiedClientState,
		Proof:                 proof,
		ProofHeight:           height,
	}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (resp QueryConnectionClientStateResponse) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return resp.IdentifiedClientState.UnpackInterfaces(unpacker)
}

// NewQueryConnectionConsensusStateResponse creates a newQueryConnectionConsensusStateResponse instance
func NewQueryConnectionConsensusStateResponse(clientID string, anyConsensusState *codectypes.Any, consensusStateHeight exported.Height, proof []byte, height clienttypes.Height) *QueryConnectionConsensusStateResponse {
	return &QueryConnectionConsensusStateResponse{
		ConsensusState: anyConsensusState,
		ClientId:       clientID,
		Proof:          proof,
		ProofHeight:    height,
	}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (resp QueryConnectionConsensusStateResponse) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return unpacker.UnpackAny(resp.ConsensusState, new(exported.ConsensusState))
}
