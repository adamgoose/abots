package patreon

import (
	"github.com/adamgoose/abots/lib/stash"
	"github.com/adamgoose/abots/lib/structure"
	"github.com/charmbracelet/log"
	"github.com/defval/di"
)

type Organizer struct {
	di.Inject
	Log   *log.Logger
	DB    *structure.DB
	Repo  *Repo
	Stash *stash.Stash
}

func (o *Organizer) OrganizeCampaign(id string) error {
	return o.DB.Update(func(tx *structure.Tx) error {
		return o.Repo.ForEachCampaignPostMedia(tx, id, func(campaign Entity[Campaign], post Entity[Post], media Entity[Media]) error {
			// See if it's already organized
			if downloaded, _ := o.Repo.MediaIsOrganized(tx, media.ID); downloaded {
				return nil
			}

			// TODO: Do the organizing

			return nil // Not ready to mark things as organized
			// Mark it as organized
			return o.Repo.MarkMediaAsOrganized(tx, media.ID)
		})
	})
}
