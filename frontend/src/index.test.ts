import { expect, test } from 'vitest';
import { playbackURL, thumbnailURL, type Film } from './lib/portfolio';

const film: Film = {
    id: 1, title: 'Short film', description: '', published: '2024-10-12',
    link: '', videoUrl: '', thumbnail: null
};

test('an uploaded thumbnail takes priority over a provider thumbnail', () => {
    expect(thumbnailURL({
        ...film, link: 'https://youtu.be/aaaaaaaaaaa', thumbnail: { url: '/media/42', crop: true }
    })).toBe('/media/42');
});

test('unusable legacy links fall back to uploaded media', () => {
    for (const link of ['https://youtu.be:99999/aaaaaaaaaaa', 'javascript:alert(1)', 'http://[invalid']) {
        const legacy = { ...film, link, videoUrl: '/media/12' };
        expect(thumbnailURL(legacy)).toBe('');
        expect(playbackURL(legacy)).toBe('/media/12');
    }
});

test('local and external films have the expected thumbnail fallbacks', () => {
    expect(thumbnailURL(film)).toBe('');
    expect(thumbnailURL({ ...film, link: 'https://www.youtube.com/watch?v=aaaaaaaaaaa' }))
        .toBe('https://img.youtube.com/vi/aaaaaaaaaaa/sddefault.jpg');
    expect(thumbnailURL({ ...film, link: 'https://player.vimeo.com/video/123456' }))
        .toBe('https://vumbnail.com/123456.jpg');
    expect(thumbnailURL({ ...film, link: 'https://example.com/aaaaaaaaaaa' })).toBe('');
});
