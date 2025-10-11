

# Distributed Key-Value Store — System Requirements & Architecture Specification

## 1. Overview

This document describes the system requirements, architecture, and operational design for a **distributed key-value store** implemented in **Go**, supporting:

* **Concurrent access** via goroutines
* **Fault tolerance** through replication and leader election
* **Strong consistency** using a **(TODO)**
* **Scalable frontend proxies** and a **registry-based node discovery** mechanism
* **gRPC-based APIs** for client, inter-node, and control plane communication


## 2. High-Level Architecture

### 2.1 Components Overview

```
Clients
  └─ gRPC ─▶ Frontend Proxy(s)  <─(periodic refresh / health)─▶ Registry
                        │
     ┌──────────────────┼──────────────────┐
     │                  │                  │
     ▼                  ▼                  ▼
  Node A (Primary) <─ Replication RPCs ─> Node B (Follower)
  (gRPC client port)                       (gRPC peer port)
      ↑                                        ↑
      │                                        │
  WAL + Snapshot                         WAL + Snapshot
  Persistent store (BadgerDB/BoltDB)     Persistent store
```

### 2.2 Components Description

| Component             | Description                                                                                                                          |
| --------------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| **Frontend Proxy**    | Stateless service that handles client requests and routes them to the appropriate cluster node (usually the leader).                 |
| **Registry**          | Central membership and discovery service where nodes register and renew heartbeats. Provides leader and cluster metadata to proxies. |
| **KV Nodes**          | Core data replicas that maintain the key-value state, replicate data via Raft, and expose client and peer gRPC APIs.                 |
| **Persistence Layer** | Each node maintains Write-Ahead Logs (WAL) and periodic snapshots for durability and recovery.                                       |

---

## 3. Functional Flows

### 3.1 Node Startup and Registration

1. Node boots and starts its gRPC servers.
2. Node registers with the **Registry** using `RegisterNode(node_id, rpc_addr, peer_addr)`.
3. Registry stores node info with a TTL.
4. Node sends periodic **heartbeats** to renew registration.

### 3.2 Leader Election

* Implemented via **Raft** (or a custom election algorithm).
* The elected leader updates the cluster state and optionally notifies the Registry.

### 3.3 Client Write (Put)

1. Client → Proxy → Leader node.
2. Leader appends entry to WAL and replicates to followers.
3. On quorum commit, leader responds to Proxy.
4. Proxy returns success to Client.

### 3.4 Client Read (Get)

* **Strong consistency:** Proxy routes read to the leader.
* **Eventual consistency:** Proxy can route to any follower for low latency.

### 3.5 Leader Change

* New leader is elected on failure.
* Proxies detect leader change via:

  * `NotLeader` gRPC error responses, or
  * Periodic refresh from Registry.

---

## 4. Registry Design

### 4.1 Responsibilities

* Maintain node registrations with TTL and heartbeats.
* Provide `GetClusterNodes` and `GetLeader` APIs.
* Serve as a service discovery point for Frontend Proxies.

### 4.2 gRPC Service (Example)

```proto
service Registry {
  rpc RegisterNode(NodeInfo) returns (RegisterResponse);
  rpc Heartbeat(NodeHeartbeat) returns (HeartbeatResponse);
  rpc GetClusterNodes(Empty) returns (ClusterNodesResponse);
  rpc GetLeader(Empty) returns (LeaderResponse);
}

message NodeInfo { string node_id = 1; string client_addr = 2; string peer_addr = 3; }
message NodeHeartbeat { string node_id = 1; int64 ts = 2; }
message ClusterNodesResponse { repeated NodeInfo nodes = 1; }
message LeaderResponse { string leader_id = 1; string leader_addr = 2; }
```

---

## 5. Frontend Proxy Design

### 5.1 Responsibilities

* Route **writes** to the leader node.
* Route **reads** based on read policy (leader/follower).
* Handle retries, leader redirection, and caching.

### 5.2 Behavior

1. Cache leader address with TTL.
2. On write:

   * Send to cached leader.
   * Retry with exponential backoff on failure.
3. On read:

   * Route per consistency mode (leader or follower).
4. On `NotLeader` response:

   * Update leader cache using new info or Registry query.

### 5.3 Availability

* Stateless — multiple proxies behind a load balancer.
* No persistent data — easily horizontally scalable.

---

## 6. gRPC APIs

### 6.1 Client-Facing (on nodes)

```proto
service KVStore {
  rpc Get(GetRequest) returns (GetResponse);
  rpc Put(PutRequest) returns (PutResponse);
  rpc Delete(DeleteRequest) returns (DeleteResponse);
}
```

### 6.2 Peer-Facing (for Raft)

```proto
service Replication {
  rpc AppendEntries(AppendRequest) returns (AppendResponse);
  rpc RequestVote(VoteRequest) returns (VoteResponse);
  rpc InstallSnapshot(SnapshotRequest) returns (SnapshotResponse);
}
```

---

## 7. File & Project Structure

```
distrokv/
├── api/
│   ├── kvstore.proto
│   ├── replication.proto
│   └── registry.proto
├── cmd/
│   ├── node/
│   │   └── main.go
│   ├── proxy/
│   │   └── main.go
│   └── registry/
│       └── main.go
├── internal/
│   ├── node/
│   │   ├── server.go
│   │   ├── raft_wrapper.go
│   │   ├── store.go
│   │   └── persistence/
│   ├── replication/
│   ├── registry/
│   ├── proxy/
│   ├── config/
│   ├── telemetry/
│   └── utils/
├── pkg/
│   ├── proto/
│   └── hash/
├── deploy/
│   ├── k8s/
│   └── docker/
├── test/
│   ├── integration/
│   └── chaos/
├── tools/
│   └── bench/
├── scripts/
│   ├── build.sh
│   └── run-local.sh
├── Makefile
├── Dockerfile.node
├── Dockerfile.proxy
├── README.md
└── go.mod
```

---

## 8. Production Considerations

| Area              | Recommendation                                                  |
| ----------------- | --------------------------------------------------------------- |
| **Persistence**   | Use BadgerDB or BoltDB; always enable WAL + periodic snapshots. |
| **Consensus**     | Implement Raft (or use HashiCorp Raft library).                 |
| **Registry**      | Use etcd or replicate Registry using Raft to avoid SPOF.        |
| **Security**      | mTLS between all components; gRPC interceptors for auth.        |
| **Observability** | Use Prometheus + OpenTelemetry; structured logging via zap.     |
| **Testing**       | Integration tests for failover, partition, recovery, and chaos. |
| **Rate Limiting** | Implement proxy-level concurrency control.                      |
| **Scaling**       | Stateless proxies; node count = odd (3 or 5) for quorum safety. |

---

## 9. Reliability, Availability, and Consistency

| Attribute                   | Description                                                                          |
| --------------------------- | ------------------------------------------------------------------------------------ |
| **Consistency**         | Strong (linearizable) with Raft-based majority commit.                               |
| **Availability**        | High for reads; reduced for writes when quorum lost.                                 |
| **Partition Tolerance** | Achieved via Raft; system sacrifices availability for consistency during partitions. |
| **Reliability**             | WAL + snapshots ensure durability; automatic leader election ensures continuity.     |

### Cluster behavior:

* **3-node cluster** tolerates 1 node failure.
* **5-node cluster** tolerates 2 node failures.
* Frontend proxies remain operational (stateless).

---

## 10. Default Configuration

| Parameter          | Recommended Default                    |
| ------------------ | -------------------------------------- |
| Cluster size       | 3 (production minimum)                 |
| Consensus          | Raft (majority quorum)                 |
| Read mode          | Leader (linearizable)                  |
| Registry           | etcd or Raft-based                     |
| Storage            | BadgerDB                               |
| Transport security | mTLS                                   |
| Logging            | zap structured logs                    |
| Monitoring         | Prometheus metrics                     |
| Snapshots          | Periodic + on WAL compaction threshold |

---

## 11. Operational Flow Summary

1. **Startup:** Nodes register to Registry → Leader elected via Raft.
2. **Client Requests:** Proxy routes client calls to appropriate node.
3. **Replication:** Leader replicates logs to followers (majority commit).
4. **Failover:** Automatic leader election on failure.
5. **Recovery:** WAL replay + snapshot restore.
6. **Scaling:** Add proxies freely; add nodes via registration.

---

## 12. CAP Summary

| Property                | Achieved                   | Notes                                      |
| ----------------------- | -------------------------- | ------------------------------------------ |
| **Consistency**         | Strong (CP system)       | via Raft majority commit                   |
| **Availability**        | Limited under partition | when quorum unavailable                    |
| **Partition Tolerance** |                           | continues consistent operation with quorum |

---

## 13. Future Extensions

* Sharding / partitioning (multiple Raft groups)
* Cross-region replication
* Pluggable storage backends
* Async replication modes (tunable consistency)
* Dynamic rebalancing
* Multi-tenant isolation
* Snapshot streaming

