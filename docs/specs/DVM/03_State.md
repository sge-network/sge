# **State**

## **KeyVault**

```proto
// KeyVault is the information of important keys stored in dvm state.
message KeyVault {
  // public_keys contains allowed public keys.
  repeated string public_keys = 1 [ (gogoproto.nullable) = false ];
}
```

The DVM module has only key in the module state that contains a list of strings which are the trusted public keys.
