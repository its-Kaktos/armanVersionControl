# File structures

## Tree
A Tree file structure in high level will look like:
* Tree header
* Tree entries

### Tree Header structure total size is N:
* signature and version: uint16
* null byte (\u0000)

### Tree entry size N:
* EntryKind: int32
* A pointer to Blob or Tree: Size unkown