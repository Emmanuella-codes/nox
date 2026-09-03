package pipes

import (
	"testing"

	"github.com/emmanuella-codes/nox/models"
)

func TestStoryItemDurationUsesFixedImageLength(t *testing.T) {
	imageAsset := &models.MediaAsset{
		MediaKind:        models.ImageMediaKind,
		DurationSeconds:  storyItemImageDurationSeconds,
		ProcessingStatus: models.ReadyMediaStatus,
	}
	videoAsset := &models.MediaAsset{
		MediaKind:        models.VideoMediaKind,
		DurationSeconds:  42,
		ProcessingStatus: models.ReadyMediaStatus,
	}
	if !validStoryMedia(imageAsset) {
		t.Fatal("expected ready image asset to be valid for stories")
	}
	if storyItemDuration(imageAsset) != storyItemImageDurationSeconds {
		t.Fatalf("expected image duration %d, got %d", storyItemImageDurationSeconds, storyItemDuration(imageAsset))
	}
	if !validStoryMedia(videoAsset) {
		t.Fatal("expected ready video asset to be valid for stories")
	}
	if storyItemDuration(videoAsset) != 42 {
		t.Fatalf("expected video duration 42, got %d", storyItemDuration(videoAsset))
	}
}
