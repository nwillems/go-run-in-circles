// Barcode scanning for lap recording
function setupBarcodeScanner() {
    const video = document.getElementById('barcode-video');
    const scanBtn = document.getElementById('barcode-scan-btn');
    const status = document.getElementById('barcode-status');
    const bibInput = document.querySelector('#lap-form input[name="bib_number"]');
    let stream = null;
    let scanning = false;

    if (!video || !scanBtn || typeof ZXing === 'undefined') return;

    scanBtn.addEventListener('click', async function () {
        if (scanning) return;
        scanning = true;
        status.textContent = 'Starting camera...';
        try {
            stream = await navigator.mediaDevices.getUserMedia({ video: { facingMode: 'environment' } });
            video.srcObject = stream;
            status.textContent = 'Scanning for barcode...';
            const codeReader = new ZXing.BrowserBarcodeReader();
            codeReader.decodeFromVideoDevice(null, video, (result, err) => {
                if (result) {
                    bibInput.value = result.text;
                    status.textContent = 'Barcode detected: ' + result.text;
                    codeReader.reset();
                    if (stream) {
                        stream.getTracks().forEach(track => track.stop());
                    }
                    video.srcObject = null;
                    scanning = false;
                } else if (err && !(err instanceof ZXing.NotFoundException)) {
                    status.textContent = 'Scan error: ' + err;
                }
            });
        } catch (err) {
            status.textContent = 'Camera error: ' + err;
            scanning = false;
        }
    });
}



// Hide forms until event is configured
function setFormsVisible(visible) {
    document.getElementById('lap-form').parentElement.style.display = visible ? '' : 'none';
}

// Check if event is configured
async function checkEventConfigured() {
    const res = await fetch('/event');
    if (!res.ok) return false;
    const event = await res.json();
    return !!event && !!event.name;
}


function setupEventListeners() {
    // Webcam setup
    const video = document.getElementById('webcam');
    const canvas = document.getElementById('photo-canvas');
    const captureBtn = document.getElementById('capture-btn');
    const photoInput = document.getElementById('photo-data');
    const photoPreview = document.getElementById('photo-preview');
    if (video && navigator.mediaDevices && navigator.mediaDevices.getUserMedia) {
        navigator.mediaDevices.getUserMedia({ video: true })
            .then(stream => {
                video.srcObject = stream;
            })
            .catch(err => {
                video.style.display = 'none';
                captureBtn.disabled = true;
                photoPreview.innerText = 'Webcam unavailable';
            });
    }
    if (captureBtn) {
        captureBtn.addEventListener('click', function (e) {
            e.preventDefault();
            canvas.getContext('2d').drawImage(video, 0, 0, canvas.width, canvas.height);
            const dataUrl = canvas.toDataURL('image/jpeg');
            photoInput.value = dataUrl.split(',')[1];
            photoPreview.innerHTML = `<img src="${dataUrl}" width="80"/>`;
        });
    }

    // Event setup logic
    document.getElementById('event-setup-form').addEventListener('submit', async function (e) {
        e.preventDefault();
        const form = e.target;
        const name = form.name.value.trim();
        const start = form.start.value;
        const lap_distance = parseInt(form.lap_distance.value, 10);
        if (!name || !start || !lap_distance) return;
        // Convert datetime-local (YYYY-MM-DDTHH:mm) to RFC3339
        const startTimestamp = new Date(start).toISOString();
        const res = await fetch('/event', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, start_timestamp: startTimestamp, lap_distance })
        });
        document.getElementById('event-setup-status').textContent = res.ok ? 'Event created!' : 'Error';
        if (res.ok) {
            document.getElementById('event-setup-section').style.display = 'none';
            setFormsVisible(true);
        }
    });

    // Record lap
    document.getElementById('lap-form').addEventListener('submit', async function (e) {
        e.preventDefault();
        const form = e.target;
        const bib_number = form.bib_number.value.trim();
        if (!bib_number) return;
        const res = await fetch('/laps', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ bib_number })
        });
        document.getElementById('lap-status').textContent = res.ok ? 'Lap recorded!' : 'Error';
        if (res.ok) form.reset();
    });
}

// On load, check event and show/hide forms, then register event listeners
window.onload = async function () {
    const configured = await checkEventConfigured();
    if (configured) {
        document.getElementById('event-setup-section').style.display = 'none';
        setFormsVisible(true);
    } else {
        setFormsVisible(false);
    }
    setupEventListeners();
    setupBarcodeScanner();
};
