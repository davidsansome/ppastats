package launchpad

import (
  "encoding/json"
  "io/ioutil"
  "net/http"
  "net/url"
  "strconv"
  "strings"
)

type Archive struct {
  Repository string
  Id string
}

type Binary struct {
  Archive *Archive
  DisplayName string
  Id string
}

func Url(path_parts []string, op string) string {
  path := "/devel/" + strings.Join(path_parts, "/")

  ret := url.URL{
    Scheme: "https",
    Host:   "api.launchpad.net",
    Path:   path,
  }

  if op != "" {
    q := ret.Query()
    q.Set("ws.op", op)
    ret.RawQuery = q.Encode()
  }

  return ret.String()
}

func (archive *Archive) GetPublishedBinaries() ([]*Binary, error) {
  // Create the URL
  url := Url([]string{archive.Repository, "+archive", archive.Id},
             "getPublishedBinaries")

  // Fetch the URL
  resp, err := http.Get(url)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()

  // Parse the JSON
  reply := struct {
    Entries []struct {
      Display_name string
      Self_link string
    }
  }{}
  if err = json.NewDecoder(resp.Body).Decode(&reply); err != nil {
    return nil, err
  }

  // Fill the return structure
  ret := make([]*Binary, len(reply.Entries))
  for i, entry := range reply.Entries {
    ret[i] = &Binary{
      Archive:     archive,
      DisplayName: entry.Display_name,
    }

    if link_parts := strings.Split(entry.Self_link, "/"); len(link_parts) > 0 {
      ret[i].Id = link_parts[len(link_parts)-1]
    }
  }

  return ret, nil
}

func (binary *Binary) GetDownloadCount() (int64, error) {
  // Create the URL
  url := Url(
    []string{
      binary.Archive.Repository,
      "+archive", binary.Archive.Id,
      "+binarypub", binary.Id,
    }, "getDownloadCount")

  // Fetch the URL
  resp, err := http.Get(url)
  if err != nil {
    return 0, err
  }
  defer resp.Body.Close()

  // Read data from the stream
  data, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return 0, err
  }
  
  // Convert it to an integer and return
  ret, err := strconv.ParseInt(string(data), 10, 0)
  if err != nil {
    return 0, err
  }
  return ret, nil
}