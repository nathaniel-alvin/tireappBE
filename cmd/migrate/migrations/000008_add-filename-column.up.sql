ALTER TABLE image 
ADD COLUMN IF NOT EXISTS filename varchar(255);
