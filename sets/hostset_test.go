package sets

import "testing"

func TestHostSet(t *testing.T) {
	set := NewHostSet([]string{"*.google.com", "*.ya.ru"})
	set.Add("*.vk.com")

	for _, h := range []string{"translate.google.com", "ya.ru", "music.vk.com"} {
		if !set.Contains(h){
			t.Errorf("not contains %s", h)
		}
	}
}
