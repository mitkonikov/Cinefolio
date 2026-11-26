// Run with mongosh --quiet <local Cinefolio database URI> --file tools/export-legacy.js.
const collections = {
    sites: ['name'],
    files: ['filepath', 'link', 'title', 'description', 'published', 'thumbnail'],
    playlists: ['name', 'site', 'visible', 'order'],
    filetoplaylists: ['fileId', 'playlistId'],
    thumbnails: ['title', 'filepath', 'crop'],
    trailers: ['film', 'filepath', 'visible']
};

const output = {};
const available = new Set(db.getCollectionNames());
for (const [name, fields] of Object.entries(collections)) {
    if (!available.has(name)) throw new Error(`Missing collection ${name}; check the database in the URI`);
    const projection = Object.fromEntries(fields.map(field => [field, 1]));
    output[name] = db.getCollection(name).find({}, projection).sort({ _id: 1 }).toArray().map(record =>
        Object.fromEntries(Object.entries(record).map(([field, value]) => [
            field, typeof value?.toHexString === 'function' ? value.toHexString() : value
        ]))
    );
}

print(JSON.stringify(output));
