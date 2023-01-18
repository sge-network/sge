# **State**

## **PublicKeys**

```proto
message PublicKeys {
  repeated string list = 1;
}
```

The DVM module has only key in the module state that contains a list of strings which are the trusted public keys.
