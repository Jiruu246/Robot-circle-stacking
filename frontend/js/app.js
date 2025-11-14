const GRID = document.getElementById("grid");
const EXPORT_BTN = document.getElementById("export-btn");
const MESSAGE = document.getElementById("message");
const HOLDING = document.getElementById("holding");

let _messageTimer = null;
const BASE_URL = "http://localhost:8080";
const END_POINTS = {
    state: `${BASE_URL}/state`,
    command: `${BASE_URL}/command`,
    export: `${BASE_URL}/export`
};

function showErrorMessage(text) {
    if (!MESSAGE) return;
    clearTimeout(_messageTimer);
    if (!text) {
        MESSAGE.className = 'message';
        MESSAGE.textContent = '';
        return;
    }
    MESSAGE.textContent = text;
    MESSAGE.className = `message show`;
    _messageTimer = setTimeout(() => {
        MESSAGE.className = 'message';
        MESSAGE.textContent = '';
    }, 4000);
}

function showWinMessage() {
    if (!MESSAGE) return;
    MESSAGE.textContent = 'Task successfully completed!';
    MESSAGE.className = `message show success`;
}

async function fetchInitialState() {
    const res = await fetch(END_POINTS.state);
    const data = await res.json();
    render(data);
}

async function sendCommand(action, direction = null) {
    const res = await fetch(END_POINTS.command, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ action, direction: direction || undefined })
    });

    if (!res.ok) {
        const msg = await res.json();
        showErrorMessage(`Command failed: ${msg.error}`);
        return;
    }

    const data = await res.json();
    await render(data);
}

async function render(state) {
    GRID.innerHTML = "";

    const gridData = state.grid;
    const robotX = state.position_x;
    const robotY = state.position_y;

    for (let y = 0; y < gridData.length; y++) {
        for (let x = 0; x < gridData[y].length; x++) {
            const cell = document.createElement("div");
            cell.className = "cell";

            const stack = gridData[x][y];

            stack.forEach(color => {
                const div = document.createElement("div");
                div.className = `circle ${color}`;
                cell.appendChild(div);
            });

            if (robotX === x && robotY === y) {
                const robot = document.createElement("div");
                robot.className = "robot";
                cell.appendChild(robot);
            }

            GRID.appendChild(cell);
        }
    }

    if (HOLDING) {
        HOLDING.innerHTML = '';
        if (state.holding) {
            HOLDING.classList.remove('empty');
            const dot = document.createElement('div');
            dot.className = `circle ${state.holding}`;
            HOLDING.appendChild(dot);
        } else {
            HOLDING.classList.add('empty');
            HOLDING.textContent = 'Empty';
        }
    }

    if (state.won) {
        showWinMessage();
    }
}


document.addEventListener("click", async (e) => {
    if (e.target.tagName !== "BUTTON") return;

    const action = e.target.dataset.action;
    if (!action) return;

    const direction = e.target.dataset.direction || null;
    await sendCommand(action, direction);
});

EXPORT_BTN.addEventListener("click", async () => {
    window.location.href = END_POINTS.export;
});

fetchInitialState();
