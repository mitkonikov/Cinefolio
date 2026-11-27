import { request } from './api';

export interface Asset {
    id: number;
    filename: string;
    contentType: string;
    size: number;
}
interface RevisionedRecord {
    id: number;
    revision: number;
}

export interface AdminFilm extends RevisionedRecord {
    title: string;
    description: string;
    published: string;
    link: string;
    asset?: Asset;
    thumbnailId: number | null;
    playlists?: AdminPlaylist[];
}
export interface AdminPlaylist extends RevisionedRecord {
    name: string;
    position: number;
    visible: boolean;
    films: AdminFilm[] | null;
}
export interface AdminThumbnail extends RevisionedRecord {
    title: string;
    crop: boolean;
    asset?: Asset;
}
export interface AdminTrailer extends RevisionedRecord {
    filmId: number;
    film?: AdminFilm;
    asset?: Asset;
    featured: boolean;
}

export interface Resources {
    films: AdminFilm;
    playlists: AdminPlaylist;
    thumbnails: AdminThumbnail;
    trailers: AdminTrailer;
}
export type Resource = keyof Resources;

const base = '/api/admin';

export function mediaURL(asset: Asset) {
    return `${base}/media/${asset.id}`;
}

export function json(method: string, input: unknown): RequestInit {
    return {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(input)
    };
}

export function list<K extends Resource>(resource: K, offset = 0, q = '', limit = 25) {
    return request<{ items: Resources[K][]; total: number }>(
        `${base}/${resource}?${new URLSearchParams({ q, offset: String(offset), limit: String(limit) })}`
    );
}

export async function choices<K extends Resource>(resource: K): Promise<Resources[K][]> {
    const all: Resources[K][] = [];
    for (;;) {
        const page = await list(resource, all.length, '', 100);
        all.push(...page.items);
        if (all.length >= page.total || page.items.length === 0) return all;
    }
}

export function save(resource: Resource, id: number | undefined, input: unknown, file: File | undefined) {
    const method = id ? 'PUT' : 'POST';
    const url = `${base}/${resource}${id ? '/' + id : ''}`;
    if (!file) return request(url, json(method, input));
    const body = new FormData();
    body.set('input', JSON.stringify(input));
    body.set('file', file);
    return request(url, { method, body });
}

export function deleteRecord(resource: Resource, record: RevisionedRecord) {
    return request(`${base}/${resource}/${record.id}?revision=${record.revision}`, { method: 'DELETE' });
}
