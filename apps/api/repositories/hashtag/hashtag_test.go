package hashtag

import (
	"reflect"
	"testing"
)

func TestExtractTagsNormalizesAndDeduplicates(t *testing.T) {
	got := ExtractTags("new #Amapiano night with #amapiano and #Alt-RnB_2 plus #too-long-tag-name")
	want := []string{"amapiano", "alt-rnb_2", "too-long-tag-name"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestNormalizeTagTrimsHashAndSeparators(t *testing.T) {
	got := NormalizeTag("  ##Amapiano-Night__  ")
	want := "amapiano-night"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestExtractTagsRejectsInvalidTags(t *testing.T) {
	got := ExtractTags("#-bad #_good #valid-tag #aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	want := []string{"good", "valid-tag"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
