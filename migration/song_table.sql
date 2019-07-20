CREATE TABLE IF NOT EXISTS songs(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT,
    link TEXT,
	 author TEXT,
	 album TEXT
);
