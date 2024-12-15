# File structures

## Tree
A Tree file structure in high level will look like:
* Tree header
* Tree entries

### Tree Header structure:
* HeaderSize: int32 // This field is generated when storing the tree in a file
* signature and version: uint16
* null byte (\u0000)

### Tree entry:
* EntryKind: int32
* EntryHashSize: int32 // This field is generated when storing the Tree in a file
* EntryHash: string
* NameSize: int32 // This field is generated when storing the Tree in a file
* Name: string
* CreatedDateSize: int32 // This field is generated when storing the Tree in a file
* CreatedDate: string (representation of time.Time)
* ModifiedDateSize: int32 // This field is generated when storing the Tree in a file
* ModifiedDate: string