package model

type KeygenRequest struct {
	Nodes      []string `json:"nodes"`
	ElectionId uint64   `json:"election_id"`
	RequestId  string   `json:"request_id"`
	Threshold  int      `json:"threshold"`
	Timestamp  int64    `json:"timestamp"`
}

type KeygenResponse struct {
	ClusterPublicKey string `json:"cluster_public_key"`
}

type TransactionSignRequest struct {
	MessageHash string `json:"message_hash"`
	ElectionId  uint64 `json:"election_id"`
	RequestId   string `json:"request_id"`
}
