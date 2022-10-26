# **State**

## **PublicKeys**

```proto
syntax = "proto3";
package sgenetwork.sge.dvm;

option go_package = "github.com/sge-network/sge/x/dvm/types";

message PublicKeys {
  repeated string list = 1;
}
```

The DVM module has only one state table consisting of a list of strings which are the trusted public keys.
