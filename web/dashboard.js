// Dashboard auto-refresh logic
async function loadDashboard() {
  const res = await fetch('/dashboard');
  if (!res.ok) return;
  const data = await res.json();
  const tbody = document.querySelector('#dashboard-table tbody');
  tbody.innerHTML = '';

  const trs = data.map(r => {
    const tr = document.createElement('tr');
    var distance = r.total_distance || 0;
    distance = distance / 1000;
    tr.innerHTML = `
            <td>${r.photo ? `<img src="data:image/jpeg;base64,${r.photo}" width="48"/>` : ''}</td>
            <td>${r.bib_number}</td>
            <td>${r.name}</td>
            <td>${r.total_laps || 0}</td>
            <td>${r.last_lap_time || ''}</td>
            <td>${r.average_lap_time || ''}</td>
            <td>${r.fastest_lap || ''}</td>
            <td>${r.slowest_lap || ''}</td>
            <td>${distance.toFixed(2)}</td>
        `;
    return tr;
  });
  tbody.append(...trs);

  const tfoot = document.querySelector('#dashboard-table tfoot');
  const trFoot = document.createElement('tr');
  const totalDistance = data.reduce((sum, r) => sum + (r.total_distance || 0), 0) / 1000;
  trFoot.innerHTML = `
            <td colspan="8">Total Runners: ${data.length}</td>
            <td colspan="1">${totalDistance.toFixed(2)}</td>
        `;
  tfoot.innerHTML = '';
  tfoot.appendChild(trFoot);

}
setInterval(loadDashboard, 5000);
window.onload = loadDashboard;
