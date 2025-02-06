let currentText = '';
let startTime = null;

async function init() {
    await loadNewText();
    setupEventListeners();
}

async function loadNewText() {
    const response = await fetch('/api/text');
    const data = await response.json();
    currentText = data.text;
    document.getElementById('textDisplay').textContent = currentText;
    document.getElementById('inputField').value = '';
    document.getElementById('results').innerHTML = '';
    startTime = Date.now();
}

async function checkText() {
    const input = document.getElementById('inputField').value;
    const response = await fetch('/api/check', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({
            input: input,
            original: currentText
        })
    });
    
    const result = await response.json();
    updateUI(result);
}

function updateUI(data) {
    const timer = document.getElementById('timer');
    const time = Math.floor((Date.now() - startTime) / 1000);
    timer.textContent = `${Math.floor(time/60)}:${time%60 < 10 ? '0' : ''}${time%60}`;
    
    document.getElementById('results').innerHTML = `
        <p>Скорость: ${data.wpm} слов/мин</p>
        <p>Точность: ${data.accuracy.toFixed(1)}%</p>
    `;
}

async function saveResults() {
    await fetch('/api/save', {
        method: 'POST',
        body: JSON.stringify({
            wpm: document.getElementById('results').querySelector('p').textContent.match(/\d+/)[0]
        })
    });
}

function setupEventListeners() {
    document.getElementById('inputField').addEventListener('input', checkText);
    document.getElementById('nextBtn').addEventListener('click', loadNewText);
    document.getElementById('saveBtn').addEventListener('click', saveResults);
}

// Инициализация
init();