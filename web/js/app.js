(function() {
    const API = window.FOLIO_API || '';
    let token = localStorage.getItem('folio_token');
    let user = null;

    const $ = (sel, ctx) => (ctx || document).querySelector(sel);
    const $$ = (sel, ctx) => [...(ctx || document).querySelectorAll(sel)];
    const app = () => $('#app');

    async function api(method, path, body) {
        const opts = { method, headers: {} };
        if (token) opts.headers['Authorization'] = 'Bearer ' + token;
        if (body && !(body instanceof FormData)) {
            opts.headers['Content-Type'] = 'application/json';
            opts.body = JSON.stringify(body);
        } else if (body) {
            opts.body = body;
        }
        const res = await fetch(API + path, opts);
        if (res.status === 401) { logout(); throw new Error('Unauthorized'); }
        if (!res.ok) { const e = await res.json().catch(() => ({})); throw new Error(e.error || res.statusText); }
        if (res.headers.get('content-type')?.includes('application/json')) return res.json();
        return res;
    }

    function toast(msg) {
        const t = document.createElement('div');
        t.className = 'toast';
        t.textContent = msg;
        document.body.appendChild(t);
        setTimeout(() => t.remove(), 3000);
    }

    function route(path) {
        history.pushState(null, '', path);
        render();
    }

    window.addEventListener('popstate', render);

    function render() {
        const p = location.pathname;
        if (!token) return renderAuth();
        if (p === '/reader') return renderReader();
        if (p.startsWith('/series/')) return renderSeriesDetail();
        if (p.startsWith('/library/')) return renderLibrary();
        renderHome();
    }

    // --- Auth ---
    function renderAuth() {
        let isLogin = true;
        function draw() {
            app().innerHTML = `
                <div class="auth-page">
                    <div class="auth-card">
                        <h2>Folio</h2>
                        <form id="authForm">
                            ${!isLogin ? '<input name="username" placeholder="Username" required>' : ''}
                            <input name="email" type="email" placeholder="Email" required>
                            <input name="password" type="password" placeholder="Password" required>
                            <button type="submit" class="primary">${isLogin ? 'Sign In' : 'Register'}</button>
                        </form>
                        <div class="toggle">
                            <a href="#" id="toggleAuth">${isLogin ? 'Create an account' : 'Sign in instead'}</a>
                        </div>
                    </div>
                </div>`;
            $('#toggleAuth').onclick = e => { e.preventDefault(); isLogin = !isLogin; draw(); };
            $('#authForm').onsubmit = async e => {
                e.preventDefault();
                const fd = new FormData(e.target);
                const data = Object.fromEntries(fd);
                try {
                    const ep = isLogin ? '/api/auth/login' : '/api/auth/register';
                    const res = await api('POST', ep, data);
                    token = res.token;
                    localStorage.setItem('folio_token', token);
                    user = res.user;
                    render();
                } catch (err) { toast(err.message); }
            };
        }
        draw();
    }

    function logout() {
        token = null;
        user = null;
        localStorage.removeItem('folio_token');
        render();
    }

    // --- Topbar ---
    function topbar(title) {
        return `
            <div class="topbar">
                <h1 onclick="location.href='/'">${title || 'Folio'}</h1>
                <nav>
                    <button onclick="location.href='/'">Home</button>
                    <button id="importBtn">Import</button>
                    <button id="logoutBtn">Logout</button>
                </nav>
            </div>`;
    }

    function bindTopbar() {
        const lb = $('#logoutBtn');
        if (lb) lb.onclick = logout;
        const ib = $('#importBtn');
        if (ib) ib.onclick = showImportModal;
    }

    // --- Home ---
    async function renderHome() {
        app().innerHTML = topbar('Folio') + '<div class="container"><div class="empty">Loading...</div></div>';
        bindTopbar();
        try {
            const libs = await api('GET', '/api/libraries');
            let html = topbar('Folio') + '<div class="container">';
            html += '<div class="section-header"><h2>Libraries</h2><button class="secondary" id="addLib">+ New Library</button></div>';
            if (libs.length === 0) {
                html += '<div class="empty">No libraries yet. Create one to get started.</div>';
            } else {
                html += '<div class="library-grid">';
                for (const lib of libs) {
                    html += `<div class="library-card" data-id="${lib.id}">
                        <h3>${esc(lib.name)}</h3>
                        <p>${lib.type || 'Comic'}</p>
                    </div>`;
                }
                html += '</div>';
            }
            html += '</div>';
            app().innerHTML = html;
            bindTopbar();
            $$('.library-card').forEach(c => {
                c.onclick = () => route('/library/' + c.dataset.id);
            });
            const ab = $('#addLib');
            if (ab) ab.onclick = showNewLibraryModal;
        } catch (err) { toast(err.message); }
    }

    // --- Library ---
    async function renderLibrary() {
        const id = location.pathname.split('/').pop();
        app().innerHTML = topbar('Folio') + '<div class="container"><div class="empty">Loading...</div></div>';
        bindTopbar();
        try {
            const series = await api('GET', '/api/series?library_id=' + id);
            const libs = await api('GET', '/api/libraries');
            const lib = libs.find(l => l.id == id);
            let html = topbar('Folio') + '<div class="container">';
            html += `<div class="section-header"><h2>${esc(lib ? lib.name : 'Library')}</h2>
                <div style="display:flex;gap:8px">
                    <button class="secondary" id="scanBtn">Scan</button>
                    <button class="secondary" id="importBtn2">Import</button>
                </div></div>`;
            if (series.length === 0) {
                html += '<div class="empty">No series found. Import or scan a library.</div>';
            } else {
                html += '<div class="series-grid">';
                for (const s of series) {
                    const cover = s.id ? `/api/image/series-cover?seriesId=${s.id}` : '';
                    html += `<div class="series-card" data-id="${s.id}">
                        <div class="cover">${cover ? '<img src="'+cover+'" alt="" onerror="this.parentElement.innerHTML=\'??\'">' : '??'}</div>
                        <div class="info"><h3>${esc(s.name)}</h3></div>
                    </div>`;
                }
                html += '</div>';
            }
            html += '</div>';
            app().innerHTML = html;
            bindTopbar();
            $$('.series-card').forEach(c => {
                c.onclick = () => route('/series/' + c.dataset.id);
            });
            const sb = $('#scanBtn');
            if (sb) sb.onclick = async () => {
                try { await api('POST', '/api/libraries/' + id + '/scan'); toast('Scan started'); }
                catch (err) { toast(err.message); }
            };
            const ib2 = $('#importBtn2');
            if (ib2) ib2.onclick = showImportModal;
        } catch (err) { toast(err.message); }
    }

    // --- Series Detail ---
    async function renderSeriesDetail() {
        const id = location.pathname.split('/').pop();
        app().innerHTML = topbar('Folio') + '<div class="container"><div class="empty">Loading...</div></div>';
        bindTopbar();
        try {
            const ser = await api('GET', '/api/series/' + id);
            let html = topbar('Folio') + '<div class="container">';
            html += `<div class="section-header"><h2>${esc(ser.name)}</h2></div>`;
            html += '<div class="volume-list">';
            for (const vol of (ser.volumes || [])) {
                html += `<div class="volume-item"><h3>${esc(vol.name)}</h3>`;
                for (const ch of (vol.chapters || [])) {
                    html += `<div class="chapter-row" data-chapter="${ch.id}" data-pages="${ch.page_count}" data-filetype="${ch.file_type}" style="padding:8px 0;cursor:pointer;color:var(--text-secondary)">
                        ${esc(ch.title)} <span style="float:right">${ch.page_count} pages</span></div>`;
                }
                html += '</div>';
            }
            if (!ser.volumes || ser.volumes.length === 0) {
                html += '<div class="empty">No volumes found.</div>';
            }
            html += '</div></div>';
            app().innerHTML = html;
            bindTopbar();
            $$('.chapter-row').forEach(r => {
                r.onclick = () => {
                    window._readerChapter = { id: r.dataset.chapter, pages: parseInt(r.dataset.pages), fileType: r.dataset.filetype };
                    route('/reader');
                };
            });
        } catch (err) { toast(err.message); }
    }

    // --- Reader ---
    let readerPage = 0;
    function renderReader() {
        const ch = window._readerChapter;
        if (!ch) { route('/'); return; }
        const totalPages = ch.pages || 1;
        app().innerHTML = `
            <div class="reader-container" id="readerView">
                <div class="empty">Loading page...</div>
            </div>
            <div class="reader-controls">
                <button id="prevPage">&laquo;</button>
                <span class="page-info">${readerPage + 1} / ${totalPages}</span>
                <button id="nextPage">&raquo;</button>
            </div>`;
        loadPage(ch.id, readerPage);
        $('#prevPage').onclick = () => { if (readerPage > 0) { readerPage--; renderReader(); } };
        $('#nextPage').onclick = () => { if (readerPage < totalPages - 1) { readerPage++; renderReader(); } };
        document.onkeydown = e => {
            if (e.key === 'ArrowLeft' && readerPage > 0) { readerPage--; renderReader(); }
            if (e.key === 'ArrowRight' && readerPage < totalPages - 1) { readerPage++; renderReader(); }
        };
    }

    async function loadPage(chapterId, num) {
        const view = $('#readerView');
        try {
            const res = await fetch(API + `/api/reader/chapter/${chapterId}/page/${num}`, {
                headers: token ? { 'Authorization': 'Bearer ' + token } : {}
            });
            const ct = res.headers.get('content-type') || '';
            if (ct.includes('image/')) {
                const blob = await res.blob();
                view.innerHTML = `<img src="${URL.createObjectURL(blob)}" alt="Page ${num}">`;
            } else {
                const text = await res.text();
                view.innerHTML = `<div style="max-width:700px;line-height:1.8;padding:20px;white-space:pre-wrap">${esc(text)}</div>`;
            }
        } catch { view.innerHTML = '<div class="empty">Failed to load page</div>'; }
    }

    // --- Modals ---
    function showModal(html) {
        const overlay = document.createElement('div');
        overlay.className = 'modal-overlay';
        overlay.innerHTML = `<div class="modal">${html}</div>`;
        overlay.onclick = e => { if (e.target === overlay) overlay.remove(); };
        document.body.appendChild(overlay);
        return overlay;
    }

    function showNewLibraryModal() {
        const m = showModal(`
            <h3>New Library</h3>
            <form id="newLibForm">
                <input name="name" placeholder="Library Name" required>
                <select name="type"><option value="comic">Comic</option><option value="book">Book</option></select>
                <input name="folder_paths" placeholder="Folder path (e.g. /mnt/comics)">
                <div class="modal-actions">
                    <button type="button" class="secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                    <button type="submit" class="primary">Create</button>
                </div>
            </form>`);
        m.querySelector('#newLibForm').onsubmit = async e => {
            e.preventDefault();
            const fd = new FormData(e.target);
            const data = {
                name: fd.get('name'),
                type: fd.get('type'),
                folder_paths: [fd.get('folder_paths')].filter(Boolean),
            };
            try {
                await api('POST', '/api/libraries', data);
                m.remove();
                toast('Library created');
                renderHome();
            } catch (err) { toast(err.message); }
        };
    }

    function showImportModal() {
        const m = showModal(`
            <h3>Import Book</h3>
            <form id="importForm" enctype="multipart/form-data">
                <select name="library_id" id="importLibSelect"><option value="">Select library...</option></select>
                <input name="series_name" placeholder="Series name (optional)">
                <div class="upload-zone" id="dropZone">
                    <p>Drop file here or click to browse</p>
                    <input name="file" type="file" accept=".cbz,.cbr,.epub,.zip,.rar" style="display:none" id="fileInput">
                </div>
                <div class="modal-actions">
                    <button type="button" class="secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                    <button type="submit" class="primary">Import</button>
                </div>
            </form>`);
        const dz = m.querySelector('#dropZone');
        const fi = m.querySelector('#fileInput');
        dz.onclick = () => fi.click();
        dz.ondragover = e => { e.preventDefault(); dz.classList.add('dragover'); };
        dz.ondragleave = () => dz.classList.remove('dragover');
        dz.ondrop = e => { e.preventDefault(); dz.classList.remove('dragover'); fi.files = e.dataTransfer.files; dz.querySelector('p').textContent = fi.files[0]?.name || 'Drop file'; };
        fi.onchange = () => { dz.querySelector('p').textContent = fi.files[0]?.name || 'Drop file'; };

        api('GET', '/api/libraries').then(libs => {
            const sel = m.querySelector('#importLibSelect');
            libs.forEach(l => { const o = document.createElement('option'); o.value = l.id; o.textContent = l.name; sel.appendChild(o); });
        }).catch(() => {});

        m.querySelector('#importForm').onsubmit = async e => {
            e.preventDefault();
            const fd = new FormData(e.target);
            try {
                await api('POST', '/api/import', fd);
                m.remove();
                toast('Imported successfully');
                render();
            } catch (err) { toast(err.message); }
        };
    }

    function esc(s) {
        if (!s) return '';
        const d = document.createElement('div');
        d.textContent = s;
        return d.innerHTML;
    }

    // --- Init ---
    if (token) {
        api('GET', '/api/auth/me').then(u => { user = u; render(); }).catch(() => { logout(); });
    } else {
        render();
    }
})();
