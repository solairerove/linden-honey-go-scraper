CREATE TABLE IF NOT EXISTS verses(
   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
   ordinal integer,
   data TEXT,
	 song_id UUID REFERENCES songs (id)
);
