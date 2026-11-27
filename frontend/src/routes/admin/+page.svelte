<script lang="ts">
    import { onMount } from 'svelte';
    import { APIError, request } from '#lib/api.ts';
    import {
        choices, deleteRecord, json, list, mediaURL, save,
        type Asset, type AdminFilm, type AdminPlaylist, type AdminThumbnail, type AdminTrailer,
        type Resource, type Resources
    } from '#lib/admin.ts';
    import FrameCapture from '#lib/FrameCapture.svelte';
    import RecordChoices from '#lib/RecordChoices.svelte';

    type Row = Resources[Resource];
    const resources: { key: Resource; label: string; singular: string }[] = [
        { key: 'films', label: 'Films', singular: 'film' },
        { key: 'playlists', label: 'Playlists', singular: 'playlist' },
        { key: 'thumbnails', label: 'Thumbnails', singular: 'thumbnail' },
        { key: 'trailers', label: 'Trailers', singular: 'trailer' }
    ];
    const pageSize = 25;
    let section = $state(resources[0]);
    let email = $state('');
    let password = $state('');
    let owner = $state<{ email: string }>();
    let sessionLoading = $state(true);
    let loading = $state(false);
    let saving = $state(false);
    let error = $state('');
    let notice = $state('');
    let items = $state<Row[]>([]);
    let marked = $state<number[]>([]);
    let offset = $state(0);
    let total = $state(0);
    let query = $state('');
    let selection = $state<Row | 'new'>();
    let pendingDelete = $state<Row>();
    let file = $state<File>();
    let fileInput = $state<HTMLInputElement>();
    let preview = $state('');
    let films = $state<AdminFilm[]>([]);
    let playlists = $state<AdminPlaylist[]>([]);
    let thumbnails = $state<AdminThumbnail[]>([]);
    let draft = $state({
        text: '', description: '', published: '', link: '', thumbnailId: null as number | null,
        removeUpload: false, selected: [] as number[], visible: true, position: 0, crop: false,
        filmId: 0, featured: false
    });
    let formVersion = 0;
    let listVersion = 0;

    onMount(async () => {
        try {
            owner = await request('/api/admin/session');
            await refresh();
        } catch (cause) {
            if (!(cause instanceof APIError && cause.status === 401)) showError(cause);
        } finally {
            sessionLoading = false;
        }
    });

    $effect(() => {
        if (!file) {
            preview = '';
            return;
        }
        const url = URL.createObjectURL(file);
        preview = url;
        return () => URL.revokeObjectURL(url);
    });

    function showError(cause: unknown) {
        if (cause instanceof APIError && cause.status === 401 && owner) {
            owner = undefined;
            error = 'Your session expired. Sign in again.';
        } else {
            error = cause instanceof Error ? cause.message : 'The request failed.';
        }
    }

    async function refresh() {
        const version = ++listVersion;
        loading = true;
        marked = [];
        try {
            const result = await list(section.key, offset, query.trim(), pageSize);
            if (version !== listVersion) return;
            total = result.total;
            if (!result.items.length && offset > 0) {
                offset = Math.max(0, Math.floor((total - 1) / pageSize) * pageSize);
                await refresh();
                return;
            }
            items = result.items;
        } catch (cause) {
            if (version === listVersion) showError(cause);
        } finally {
            if (version === listVersion) loading = false;
        }
    }

    async function signIn(event: SubmitEvent) {
        event.preventDefault();
        saving = true;
        error = '';
        try {
            owner = await request('/api/admin/login', json('POST', { email, password }));
            password = '';
            await refresh();
        } catch (cause) {
            showError(cause);
        } finally {
            saving = false;
        }
    }

    async function signOut() {
        error = '';
        try {
            await request('/api/admin/logout', { method: 'POST' });
            owner = undefined;
            closeEditor();
            items = [];
        } catch (cause) { showError(cause); }
    }

    function switchSection(next: typeof section) {
        if (next.key === section.key) return;
        closeEditor();
        pendingDelete = undefined;
        section = next;
        query = '';
        offset = 0;
        items = [];
        total = 0;
        notice = '';
        error = '';
        void refresh();
    }

    function closeEditor() {
        formVersion++;
        selection = undefined;
        file = undefined;
    }

    async function edit(row: Row | 'new') {
        error = '';
        notice = '';
        pendingDelete = undefined;
        file = undefined;
        if (fileInput) fileInput.value = '';
        const record = row === 'new' ? undefined : row;
        draft = {
            text: record && 'name' in record ? record.name : record && 'title' in record ? record.title : '',
            description: record && 'description' in record ? record.description : '',
            published: record && 'published' in record ? record.published : '',
            link: record && 'link' in record ? record.link : '',
            thumbnailId: record && 'thumbnailId' in record ? record.thumbnailId : null,
            removeUpload: false,
            selected: record && 'films' in record ? (record.films ?? []).map((film) => film.id)
                : record && 'playlists' in record ? (record.playlists ?? []).map((item) => item.id) : [],
            visible: record && 'visible' in record ? record.visible : true,
            position: record && 'position' in record ? record.position : 0,
            crop: record && 'crop' in record ? record.crop : false,
            filmId: record && 'filmId' in record ? record.filmId : 0,
            featured: record && 'featured' in record ? record.featured : false
        };
        selection = row;
        const version = ++formVersion;
        saving = true;
        try {
            if (section.key === 'films') {
                const [allLists, allThumbnails] = await Promise.all([choices('playlists'), choices('thumbnails')]);
                if (version === formVersion) { playlists = allLists; thumbnails = allThumbnails; }
            } else if (section.key === 'playlists' || section.key === 'trailers') {
                const allFilms = await choices('films');
                if (version === formVersion) films = allFilms;
            }
        } catch (cause) {
            if (version === formVersion) { closeEditor(); showError(cause); }
        } finally { saving = false; }
    }

    function input() {
        switch (section.key) {
            case 'films': return {
                title: draft.text, description: draft.description, published: draft.published,
                link: draft.link, thumbnailId: draft.thumbnailId, removeUpload: draft.removeUpload,
                playlistIds: draft.selected
            };
            case 'playlists': return {
                name: draft.text, position: draft.position, visible: draft.visible, filmIds: draft.selected
            };
            case 'thumbnails': return { title: draft.text, crop: draft.crop };
            case 'trailers': return { filmId: draft.filmId, featured: draft.featured };
        }
    }

    async function submit(event: SubmitEvent) {
        event.preventDefault();
        if (!selection) return;
        saving = true;
        error = '';
        try {
            await save(section.key, selection === 'new' ? undefined : selection.id, input(), file);
            notice = 'Saved.';
            closeEditor();
            await refresh();
        } catch (cause) { showError(cause); }
        finally { saving = false; }
    }

    async function confirmDelete() {
        if (!pendingDelete) return;
        saving = true;
        error = '';
        try {
            await deleteRecord(section.key, pendingDelete.id);
            pendingDelete = undefined;
            notice = 'Deleted.';
            await refresh();
        } catch (cause) { showError(cause); }
        finally { saving = false; }
    }

    async function enableCrop() {
        saving = true;
        error = '';
        try {
            await request('/api/admin/thumbnails/enable-crop', json('POST', { ids: marked }));
            notice = 'Cropping enabled.';
            await refresh();
        } catch (cause) { showError(cause); }
        finally { saving = false; }
    }

    function name(row: Row): string {
        if ('name' in row) return row.name;
        if ('title' in row) return row.title || row.asset?.filename || `#${row.id}`;
        return row.film?.title || row.asset?.filename || `#${row.id}`;
    }

    function detail(row: Row): string {
        if ('featured' in row) return row.featured ? 'Featured' : 'Not featured';
        if ('position' in row) return `${row.visible ? 'Visible' : 'Hidden'} · Position ${row.position}`;
        if ('published' in row) return row.published;
        return row.crop ? 'Cropped' : '';
    }

    let currentAsset = $derived<Asset | undefined>(selection && selection !== 'new' && 'asset' in selection
        ? selection.asset : undefined);
    let currentPreview = $derived(preview || (currentAsset && !draft.removeUpload ? mediaURL(currentAsset) : ''));
</script>

<svelte:head>
    <title>Cinefolio Admin</title>
    <meta name="robots" content="noindex" />
</svelte:head>

<main class="admin">
    <header>
        <div><a href="/">View portfolio</a><h1>Cinefolio Admin</h1></div>
        {#if owner}<button class="secondary" onclick={signOut} disabled={saving}>Sign out</button>{/if}
    </header>

    {#if error}<p class="feedback error" role="alert">{error}</p>{/if}
    {#if notice}<p class="feedback" role="status">{notice}</p>{/if}

    {#if sessionLoading}
        <p>Checking session…</p>
    {:else if !owner}
        <form class="card login" onsubmit={signIn}>
            <label>Email <input type="email" autocomplete="username" required bind:value={email} /></label>
            <label>Password <input type="password" autocomplete="current-password" required bind:value={password} /></label>
            <button disabled={saving}>Sign in</button>
        </form>
    {:else}
        <nav aria-label="Catalog">
            {#each resources as resource (resource.key)}
                <button class:active={resource.key === section.key} aria-pressed={resource.key === section.key}
                    disabled={saving} onclick={() => switchSection(resource)}>{resource.label}</button>
            {/each}
        </nav>

        {#if selection}
            <form class="card editor" onsubmit={submit} aria-busy={saving}>
                <h2>{selection === 'new' ? 'Add' : 'Edit'} {section.singular}</h2>
                <fieldset class="fields" disabled={saving}>
                    {#if section.key === 'films'}
                        <label>Title <input required maxlength="200" bind:value={draft.text} /></label>
                        <label>Description <textarea rows="4" bind:value={draft.description}></textarea></label>
                        <label>Publication date <input type="date" required bind:value={draft.published} /></label>
                        <label>External video URL <input type="url" bind:value={draft.link} /></label>
                        <div class="field">
                            <label for="film-thumbnail">Thumbnail</label>
                            <select id="film-thumbnail" bind:value={draft.thumbnailId}>
                                <option value={null}>Automatic thumbnail</option>
                                {#each thumbnails as row (row.id)}<option value={row.id}>{name(row)}</option>{/each}
                            </select>
                        </div>
                        <RecordChoices title="Playlists" items={playlists.map((row) => ({ id: row.id, text: row.name }))}
                            bind:selected={draft.selected} />
                    {:else if section.key === 'playlists'}
                        <label>Name <input required maxlength="200" bind:value={draft.text} /></label>
                        <label>Display position <input type="number" required min="0" bind:value={draft.position} /></label>
                        <label class="check"><input type="checkbox" bind:checked={draft.visible} /> Visible on portfolio</label>
                        <RecordChoices title="Films" items={films.map((row) => ({ id: row.id, text: row.title }))}
                            bind:selected={draft.selected} />
                    {:else if section.key === 'thumbnails'}
                        <label>Title <input maxlength="200" bind:value={draft.text} /></label>
                        <label class="check"><input type="checkbox" bind:checked={draft.crop} /> Crop on portfolio</label>
                    {:else}
                        <div class="field">
                            <label for="trailer-film">Full film</label>
                            <select id="trailer-film" required bind:value={draft.filmId}>
                                <option value={0} disabled>Choose a film</option>
                                {#each films as row (row.id)}<option value={row.id}>{row.title}</option>{/each}
                            </select>
                        </div>
                        <label class="check"><input type="checkbox" bind:checked={draft.featured} /> Featured trailer</label>
                    {/if}

                    {#if section.key !== 'playlists'}
                        <label>{currentAsset ? 'Replacement upload' : 'Upload'}
                            <input type="file" bind:this={fileInput}
                                accept={section.key === 'thumbnails' ? 'image/jpeg' : section.key === 'trailers' ? 'video/mp4' : 'video/*'}
                                onchange={(event) => { file = event.currentTarget.files?.[0]; if (file) draft.removeUpload = false; }} />
                        </label>
                        {#if section.key === 'thumbnails'}
                            <FrameCapture onCapture={(captured) => { file = captured; if (fileInput) fileInput.value = ''; }} />
                        {:else if section.key === 'films' && currentAsset}
                            <label class="check"><input type="checkbox" bind:checked={draft.removeUpload} disabled={!!file} /> Remove current upload</label>
                        {/if}
                        {#if currentPreview}
                            {#if section.key === 'thumbnails'}
                                <img class="preview" src={currentPreview} alt="Thumbnail preview" />
                            {:else}
                                <!-- svelte-ignore a11y_media_has_caption -->
                                <video class="preview" src={currentPreview} controls playsinline preload="metadata"></video>
                            {/if}
                        {/if}
                        {#if file}<p class="filename">{file.name}</p>{/if}
                    {/if}
                </fieldset>
                <div class="actions"><button disabled={saving}>Save</button><button type="button" class="secondary"
                    disabled={saving} onclick={closeEditor}>Cancel</button></div>
            </form>
        {/if}

        {#if pendingDelete}
            <section class="card" aria-label="Confirm deletion">
                <p>Delete <strong>{name(pendingDelete)}</strong>? Remove any dependent references first.</p>
                <div class="actions"><button class="danger" disabled={saving} onclick={confirmDelete}>Delete</button>
                    <button class="secondary" disabled={saving} onclick={() => pendingDelete = undefined}>Cancel</button></div>
            </section>
        {/if}

        <section aria-label={section.label}>
            <div class="toolbar">
                {#if section.key !== 'trailers'}
                    <form class="actions" onsubmit={(event) => { event.preventDefault(); offset = 0; void refresh(); }}>
                        <input type="search" aria-label={`Search ${section.label.toLowerCase()}`} placeholder="Search…" bind:value={query}
                            oninput={() => { if (!query) { offset = 0; void refresh(); } }} />
                        <button class="secondary">Search</button>
                    </form>
                {/if}
                <button disabled={saving} onclick={() => edit('new')}>Add {section.singular}</button>
                {#if section.key === 'thumbnails' && marked.length}
                    <button class="secondary" disabled={saving} onclick={enableCrop}>Enable crop for {marked.length}</button>
                {/if}
            </div>
            <div class="listing" aria-busy={loading}>
                {#each items as row (row.id)}
                    <div class="row">
                        {#if section.key === 'thumbnails'}
                            <input type="checkbox" aria-label={`Select ${name(row)}`} value={row.id} bind:group={marked} />
                        {/if}
                        <div class="identity"><strong>{name(row)}</strong><span>{detail(row)}</span></div>
                        <button class="secondary" disabled={saving} onclick={() => edit(row)}>Edit</button>
                        <button class="secondary" disabled={saving} onclick={() => { closeEditor(); pendingDelete = row; error = ''; }}>Delete</button>
                    </div>
                {:else}
                    <p>{loading ? 'Loading…' : 'No records found.'}</p>
                {/each}
            </div>
            {#if total > pageSize}
                <div class="actions pagination">
                    <button class="secondary" disabled={loading || offset === 0}
                        onclick={() => { offset -= pageSize; void refresh(); }}>Previous</button>
                    <span>{offset + 1}–{Math.min(offset + pageSize, total)} of {total}</span>
                    <button class="secondary" disabled={loading || offset + pageSize >= total}
                        onclick={() => { offset += pageSize; void refresh(); }}>Next</button>
                </div>
            {/if}
        </section>
    {/if}
</main>

<style>
    .admin { margin: 0 auto; padding: 1.5rem 1rem 5rem; max-width: 62rem; width: 100%;
        box-sizing: border-box; font: 1rem/1.45 system-ui, sans-serif; color-scheme: dark; color: #e8e8ea; }
    header, nav, .actions, .toolbar, .row { display: flex; gap: 0.65rem; align-items: center; flex-wrap: wrap; }
    header { justify-content: space-between; margin-bottom: 1rem; }
    h1, h2 { margin: 0.4rem 0; }
    a { color: #acccff; }
    nav, .toolbar, .feedback, .card, .pagination { margin: 0.75rem 0; }
    .card, .listing { background: #222326; border: 1px solid #414349; border-radius: 0.55rem; padding: 0.85rem; }
    .login, .fields { display: grid; gap: 0.75rem; }
    .login { max-width: 25rem; }
    .login button { justify-self: start; }
    .fields { border: 0; padding: 0; margin: 0.75rem 0; min-inline-size: 0; }
    label, .field { display: grid; gap: 0.25rem; }
    label.check { display: flex; gap: 0.4rem; align-items: center; }
    input, button, select, textarea { box-sizing: border-box; font: inherit; color: inherit; }
    input:not([type='checkbox']), select, textarea { max-width: 100%; background: #161719; border: 1px solid #696c73;
        border-radius: 0.35rem; padding: 0.45rem; }
    button { border: 1px solid transparent; border-radius: 0.35rem; padding: 0.45rem 0.7rem; cursor: pointer;
        background: #bcd5fd; color: #101820; }
    button.secondary, nav button:not(.active) { background: #272930; color: #e9e9ee; border-color: #666b74; }
    button.danger { background: #ffa9a7; }
    button:disabled { opacity: 0.55; cursor: default; }
    .row { flex-wrap: nowrap; border-bottom: 1px solid #414349; padding: 0.55rem 0; }
    .row:last-child { border-bottom: 0; }
    .identity { display: grid; flex: 1 1 auto; overflow-wrap: anywhere; }
    .identity span, .filename { color: #bfc2c8; }
    .preview { display: block; max-width: 100%; max-height: 20rem; justify-self: start; }
    .error { color: #ffb4b6; }
</style>
