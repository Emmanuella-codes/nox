package pipes

import "github.com/emmanuella-codes/nox/models"

// attachmentsOrEmpty normalizes nil attachment slices into empty response arrays.
func attachmentsOrEmpty(attachments []*models.MediaAsset) []*models.MediaAsset {
	if attachments == nil {
		return []*models.MediaAsset{}
	}
	return attachments
}
