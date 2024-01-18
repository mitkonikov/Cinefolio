import { initializeApp } from 'firebase/app';
import 'firebase/firestore';
import { getFirestore } from 'firebase/firestore';

const firebaseConfig = {
    apiKey: "AIzaSyAFnrbghv1TvfzlVne5ytuviLo3UbqQ9Ww",
    authDomain: "cinefolio-f4969.firebaseapp.com",
    databaseURL: "https://cinefolio-f4969.firebaseio.com",
    projectId: "cinefolio-f4969",
    storageBucket: "cinefolio-f4969.appspot.com",
    messagingSenderId: "1027599997430",
    appId: "1:1027599997430:web:cb42942b07c66e3bc40fd7",
    measurementId: "G-YXTNJL8GTJ"
};

const app = initializeApp(firebaseConfig);
export const db = getFirestore(app);

