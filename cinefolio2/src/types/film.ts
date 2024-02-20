import type { IPlaylist } from "./playlist";

export interface IThumbnail {
    title: string;
    filepath: string;
    crop: boolean;
}

export interface IFilm {
    title: string;
    description: string;
    filepath: string;
    link: string | undefined;
    thumbnail: IThumbnail;
	published: string;
	embedlink: string;
};

export interface ITrailer {
    film: IFilm;
    filepath: string;
    visible: string;
};

export interface IFilmToPlaylist {
    fileId: IFilm;
    playlistId: IPlaylist;
}
