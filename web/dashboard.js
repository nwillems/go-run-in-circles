// Dashboard auto-refresh logic
async function loadDashboard() {
  const res = await fetch('/dashboard');
  if (!res.ok) return;
  const data = await res.json();
  const tbody = document.querySelector('#dashboard-table tbody');
  tbody.innerHTML = '';
  data.forEach(r => {
    const tr = document.createElement('tr');
    tr.innerHTML = `
            <td>${r.photo ? `<img src="data:image/jpeg;base64,${r.photo}" width="48"/>` : ''}</td>
            <td>${r.bib_number}</td>
            <td>${r.name}</td>
            <td>${r.total_laps || 0}</td>
            <td>${r.last_lap_time || ''}</td>
            <td>${r.average_lap_time || ''}</td>
            <td>${r.fastest_lap || ''}</td>
            <td>${r.slowest_lap || ''}</td>
            <td>${r.avg_pace || ''}</td>
            <td>${r.total_distance || ''}</td>
        `;
    tbody.appendChild(tr);
  });
}
setInterval(loadDashboard, 5000);
window.onload = loadDashboard;
