package main

import (
  "fmt"
  "launchpad"
  "time"
)

func main() {
  a := launchpad.Archive{
    Repository: "~me-davidsansome",
    Id: "clementine",
  }

  if ret, err := a.GetPublishedBinaries(); err != nil {
    fmt.Println(err)
  } else {
    for _, binary := range ret {
      time.Sleep(time.Second)

      downloads, err := binary.GetDownloadCount()
      if err != nil {
        fmt.Println(err)
      } else {
        fmt.Println(binary.DisplayName, downloads)
      }
    }
  }
}