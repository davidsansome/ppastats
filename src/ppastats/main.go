package main

import (
  "fmt"
  "flag"
  "launchpad"
  "os"
  "time"
)

func main() {
  repo    := flag.String("repo", "", "Launchpad repository name")
  archive := flag.String("archive", "", "Launchpad archive name")
  flag.Parse()

  if *repo == "" || *archive == "" {
    fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "-repo <repo> -archive <archive>")
    return
  }

  a := launchpad.Archive{
    Repository: *repo,
    Id:         *archive,
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