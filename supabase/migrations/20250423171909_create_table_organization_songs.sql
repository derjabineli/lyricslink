CREATE TABLE IF NOT EXISTS organizations_songs(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL,
    song_id UUID NOT NULL,
    CONSTRAINT fk_organization
        FOREIGN KEY (organization_id) REFERENCES planning_center_organizations(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_song
        FOREIGN KEY (song_id) REFERENCES songs(id)
        ON DELETE CASCADE,
    CONSTRAINT unique_org_song_pair
        UNIQUE(organization_id, song_id)
);