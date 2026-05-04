import { useState, useEffect } from 'react';

function App() {
    const [sites, setSites] = useState([]);
    const [url, setUrl] = useState('');
    const [name, setName] = useState('');
    const [editId, setEditId] = useState(null);

    const API = '/api/sites';

    const fetchSites = () => {
        fetch(API).then(r => r.json()).then(data => setSites(data || []));
    };

    useEffect(() => { fetchSites(); }, []);

    const addSite = () => {
        if (editId) {
            fetch(API, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ id: editId, url, name }),
            }).then(() => { fetchSites(); setUrl(''); setName(''); setEditId(null); });
        } else {
            fetch(API, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ url, name }),
            }).then(() => { fetchSites(); setUrl(''); setName(''); });
        }
    };

    const deleteSite = (id) => {
        fetch(API + '?id=' + id, { method: 'DELETE' }).then(fetchSites);
    };

    const editSite = (site) => {
        setEditId(site.id);
        setUrl(site.url);
        setName(site.name);
    };

    const formatDate = (d) => {
        if (!d) return '—';
        return new Date(d).toLocaleString();
    };

    return (
        <main className="container">
            <h1>Uptime Monitor</h1>

            <div role="group">
                <input placeholder="URL" value={url} onChange={e => setUrl(e.target.value)} disabled={!!editId} />
                <input placeholder="Name" value={name} onChange={e => setName(e.target.value)} />
                <button onClick={addSite}>{editId ? 'Save' : 'Add'}</button>
                {editId && <button className="outline secondary" onClick={() => { setEditId(null); setUrl(''); setName(''); }}>Cancel</button>}
            </div>

            <table>
                <thead>
                <tr><th>Name</th><th>URL</th><th>Status</th><th>Uptime</th><th>Checked</th><th>Added</th><th></th></tr>
                </thead>
                <tbody>
                {sites.map(s => (
                    <tr key={s.id}>
                        <td>{s.name}</td>
                        <td>{s.url}</td>
                        <td style={{ color: s.status ? '#22c55e' : '#ef4444' }}>{s.status ? 'UP' : 'DOWN'}</td>
                        <td>{s.uptime != null ? s.uptime.toFixed(2) + '%' : '—'}</td>
                        <td>{formatDate(s.checked_at)}</td>
                        <td>{formatDate(s.created_at)}</td>
                        <td><button className="outline" onClick={() => editSite(s)}>Edit</button></td>
                        <td><button className="outline secondary" onClick={() => deleteSite(s.id)}>Delete</button></td>
                    </tr>
                ))}
                </tbody>
            </table>
        </main>
    );
}

export default App;