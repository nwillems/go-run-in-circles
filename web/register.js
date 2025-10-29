function setFormsVisible(visible) {
  document.getElementById('register-form').parentElement.style.display = visible ? '' : 'none';
}

function setupRegisterListeners() {
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

  // Register runner
  document.getElementById('register-form').addEventListener('submit', async function (e) {
    e.preventDefault();
    const form = e.target;
    const bib_number = form.bib_number.value.trim();
    const name = form.name.value.trim();
    const photoBase64 = photoInput.value;
    if (!bib_number || !name || !photoBase64) return;
    const res = await fetch('/runners', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ bib_number, name, photo: photoBase64 })
    });
    document.getElementById('register-status').textContent = res.ok ? 'Registered!' : 'Error';
    if (res.ok) {
      form.reset();
      photoPreview.innerHTML = '';
      photoInput.value = '';
    }
  });
}

window.onload = function () {
  setFormsVisible(true);
  setupRegisterListeners();
};
