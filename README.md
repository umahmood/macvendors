# MAC Vendors

MAC vendors is a Go library which uses the [macvendors.co](https://macvendors.co/) 
lookup service, to get vendor information of a network device using its MAC address. 

# Installation

> $ go get github.com/umahmood/macvendors

# Usage

To get detailed information about a MAC address, use 'Lookup':
```
package main

import (
    "fmt"
    
    github.com/umahmood/macvendors"
)

func main() {
    vendor := macvendors.New()
    mac, err := vendor.Lookup("28:18:78:6D:64:42")
    if err != nil {
        //...
    }
    fmt.Println(mac.Address)
    fmt.Println(mac.Company)
    fmt.Println(mac.Country)
    fmt.Println(mac.Type)
    fmt.Println(mac.MacPrefix)
    fmt.Println(mac.StartHex)
    fmt.Println(mac.EndHex)
}
```
Output:
```
One Microsoft Way,Redmond  Washington  98052-6399,US
Microsoft Corporation
US
MA-L
28:18:78
281878000000
281878FFFFFF
```
If you are only interested in the vendor name, use the 'Name' method:
```
name, err := vendor.Name("28:18:78:6D:64:42")
if err != nil {
    //...
}
fmt.Println(name)
```
Output:
```
Microsoft Corporation
```

# Testing

- go version go1.9 darwin/amd64

# Documentation

> http://godoc.org/github.com/umahmood/macvendors

# License

See the [LICENSE](LICENSE.md) file for license rights and limitations (MIT).
