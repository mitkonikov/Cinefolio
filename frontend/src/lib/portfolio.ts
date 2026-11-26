export interface Thumbnail {
    url: string;
    crop: boolean;
}

export interface Film {
    id: number;
    title: string;
    description: string;
    published: string;
    link: string;
    videoUrl: string;
    thumbnail: Thumbnail | null;
}

export interface Trailer {
    url: string;
    film: Film;
}

export interface Playlist {
    id: number;
    name: string;
    films: Film[];
}

export interface Portfolio {
    trailer: Trailer | null;
    playlists: Playlist[];
}

export function thumbnailURL(film: Film): string {
    if (film.thumbnail) return film.thumbnail.url;
    if (!film.link) return '';
    const url = new URL(film.link);
    const host = url.hostname.replace(/^www\./, '');
    const parts = url.pathname.split('/').filter(Boolean);
    if (host === 'vimeo.com' || host === 'player.vimeo.com') {
        const id = parts.find((part) => /^\d+$/.test(part));
        return id ? `https://vumbnail.com/${id}.jpg` : '';
    }
    const id = host === 'youtu.be' ? parts[0]
        : host === 'youtube.com' ? url.searchParams.get('v') ?? parts[1] : null;
    return id && /^[\w-]{11}$/.test(id) ? `https://img.youtube.com/vi/${id}/sddefault.jpg` : '';
}
