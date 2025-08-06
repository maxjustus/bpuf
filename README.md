# BPUF

A Go library for Union-find (Disjoint Set) data structures with support for generic values and bipartite graphs, plus ClickHouse UDF integration.

## Features

- Union-find operations with path compression and union by rank
- Bipartite Union-find for finding common roots between two disjoint sets
- Generic types support for any comparable type
- ClickHouse executable UDFs for both Union-find and Bipartite Union-find with XML UDF configuration gen

## Installation

```bash
go get github.com/maxjustus/bpuf
```


## Quick Start

### Union-find

```go
package main

import (
    "fmt"
    "github.com/maxjustus/bpuf"
)

func main() {
    uf := bpuf.NewUnionFind(10)
    
    // Union elements
    uf.Union(1, 2)
    uf.Union(3, 4)
    uf.Union(2, 3) // Now 1,2,3,4 are connected
    
    // Check if elements are in same set
    root1 := uf.Find(1)
    root4 := uf.Find(4)
    fmt.Printf("1 and 4 connected: %t\n", root1 == root4) // true
}
```

### Union-find with Values

Work with string keys instead of integer indices:

```go
package main

import (
    "fmt"
    "github.com/maxjustus/bpuf"
)

func main() {
    uf := bpuf.NewUnionFindWithValues[string](100)
    
    // Union string values
    uf.Union("alice", "bob")
    uf.Union("charlie", "dave")
    uf.Union("bob", "charlie")
    
    // Find root representative
    root := uf.FindReturningValue("alice")
    fmt.Printf("Alice's group representative: %s\n", root)
    
    // Check if two values are connected
    aliceRoot := uf.FindReturningValue("alice")
    daveRoot := uf.FindReturningValue("dave")
    fmt.Printf("Alice and Dave connected: %t\n", aliceRoot == daveRoot) // true
}
```

### Bipartite Union-find

Handle relationships between two disjoint sets (U and V):

```go
package main

import (
    "fmt"
    "github.com/maxjustus/bpuf"
)

func main() {
    buf := bpuf.NewBipartiteUnionFind(100)
    
    // Connect elements from set U to set V - think of U as users and V as tennis clubs
    // Multiple connections from same U element create transitive relationships
    buf.Union(1, 10) // U[1] -> V[10]
    buf.Union(1, 11) // U[1] -> V[11], creates V[10] <-> V[11]
    buf.Union(2, 11) // U[2] -> V[11], connects U[2] to existing group
    
    // Find associated root in V for element in U
    if root, exists := buf.FindAssociatedRoot(1); exists {
        fmt.Printf("U[1] maps to V[%d]\n", root)
    }
    // Prints:
    // U[1] maps to V[10]
    // U[2] maps to V[10]
}
```

### Bipartite Union-find with Values

Combine bipartite functionality with generic types:

```go
package main

import (
    "fmt"
    "github.com/maxjustus/bpuf"
)

func main() {
    buf := bpuf.NewBipartiteUnionFindWithValues[string, string](100)
    
    // Connect users (U) to groups (V)
    // This example may be a bit contrived - but imagine users connected via organizations
    buf.Union("Stanley McFred", "Rad Corp") // Creates cluster 100 <-> "Rad Corp"
    buf.Union("Stanley McFred", "Milquetoast Inc") // Connects user1 to another organization, creating transitive relationship
    buf.Union("Jebson Dougalthorpe", "Milquetoast Inc") // Connects another user to the same organization
    buf.Union("Jebson Dougalthorpe", "McNotDonalds")
    buf.Union("Tom", "McNotDonalds")
    
    // Find which "root" organization a user belongs to
    if org, exists := buf.FindVRootForU("Stanley McFred"); exists {
        fmt.Printf("Stanley belongs to root org: %s\n", org)
        // prints: Stanley belongs to root org: Rad Corp
    }
    
    // Check if users are in same "collection" - IE: they are connected through any organization
    org1, _ := buf.FindVRootForU("Stanley McFred")
    // Even though Jebson is not a direct member of Rad Corp, he is connected to it indirectly through Milquetoast Inc
    org2, _ := buf.FindVRootForU("Jebson Dougalthorpe")
    // Even though Tom is not a direct member of Rad Corp, the shared corps form links in a "chain"
    // or bipartite graph, so the root org for all 3 of these people is Rad Corp
    org3, _ := buf.FindVRootForU("Tom")
    fmt.Printf("Stanley and Jebson are connected: %t\n", org1 == org2) // true
    fmt.Printf("Stanley and Tom are connected: %t\n", org1 == org3) // true
}
```

## ClickHouse UDF Integration

The library includes two ClickHouse User Defined Functions for processing Union-find operations using JSONEachRow format.

### Union-find UDF

Processes symmetric relationships where A connects to B.

**Input Format (JSONEachRow batch):**
```json
[{"a": "user1", "b": "user2"}, {"a": "user2", "b": "user3"}, {"a": "user4", "b": "user5"}]```

**Output Format (all unique values with roots):**
```json
[{"value": "user3", "root": "user1"},{"value": "user4", "root": "user4"},{"value": "user5", "root": "user4"},{"value": "user1", "root": "user1"},{"value": "user2", "root": "user1"}]
```

**Build and Usage:**
```bash
# Build the UDF binary
go build -o bin/bpuf-clickhouse ./cmd/bpuf-clickhouse

# Test directly
echo '{"a": "user1", "b": "user2"}
{"a": "user2", "b": "user3"}' | ./bin/bpuf-clickhouse --mode=unionfind

# Generate ClickHouse XML configuration for both UDFs
./bin/bpuf-clickhouse --udf-xml > udfs.xml
```

### Bipartite Union-find UDF

Tracks transitive relationships to find common roots in V for elements in U between two disjoint sets U and V.

**Input Format (JSONEachRow batch):**
```json
[{"u": "entity1", "v": "group100"}, {"u": "entity1", "v": "group101"}, {"u": "entity2", "v": "group101"}]```

**Output Format (unique U values with V roots):**
```json
[{"u": "entity1", "v_root": "group100"}, {"u": "entity2", "v_root": "group100"}]
```

**Build and Usage:**
```bash
# Test bipartite UDF
echo '[{"u": "entity1", "v": "group100"},{"u": "entity1", "v": "group101"},{"u": "entity2", "v": "group101"}]' | ./bin/bpuf-clickhouse --mode=bipartite

# Generate ClickHouse XML configuration (same file contains both UDFs)
./bin/bpuf-clickhouse --udf-xml > udfs.xml
```

### ClickHouse Integration

Both UDFs work with ClickHouse arrays of tuples:

```sql
-- Standard Union-find
SELECT unionFind([('user1', 'user2'), ('user2', 'user3'), ('user4', 'user5')]) as result
-- Returns: [('user1','user1'), ('user2','user1'), ('user3','user1'), ('user4','user4'), ('user5','user4')]

-- Bipartite Union-find  
SELECT bipartiteUnionFind([('entity1', 'group100'), ('entity1', 'group101'), ('entity2', 'group101')]) as result
-- Returns: [('entity1','group100'), ('entity2','group100')]
```

## Development

```bash
# Run tests
go test -v ./...

# Build UDF binary
go build -o bin/bpuf-clickhouse ./cmd/bpuf-clickhouse
```

## License

MIT
