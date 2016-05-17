package niconico

import "testing"

func TestToVideoID(t *testing.T) {
  v := ToVideoID("http://www.nicovideo.jp/watch/sm25715842?playlist_type=mylist&group_id=46356904&mylist_sort=1&ref=mylist_s1_p1_n52")
  if v != "sm25715842" {
    t.Error("Expected sm25715842, got ", v)
  }
}

func TestGetThumbInfo(t *testing.T) {
}
func TestGetHistory(t *testing.T) {
}
func TestGetFlv(t *testing.T) {
}
func TestDownloadVideoSource(t *testing.T) {
}
