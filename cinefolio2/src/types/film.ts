import type { Timestamp } from "firebase/firestore";

export interface IFilm {
    thumb : string;
    alt : string;
    title : string;
    content : string;
    link : string;
	timestamp : Timestamp;
	embedlink : string;
};
