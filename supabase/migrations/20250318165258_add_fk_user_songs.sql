ALTER TABLE songs ADD COLUMN user_id UUID NOT NULL;

ALTER TABLE songs 
ADD CONSTRAINT fk_user 
FOREIGN KEY (user_id) REFERENCES users(id) 
ON DELETE CASCADE;