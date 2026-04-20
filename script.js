let bugs = [];

function addBug() {
    const input = document.getElementById("bugInput");
    const text = input.value;

    if (text === "") {
        alert("Введите баг!");
        return;
    }

    bugs.push(text);
    input.value = "";

    renderBugs();
}

function deleteBug(index) {
    bugs.splice(index, 1);
    renderBugs();
}

function renderBugs() {
    const list = document.getElementById("bugList");
    list.innerHTML = "";

    bugs.forEach((bug, index) => {
        const li = document.createElement("li");

        li.innerHTML = `
            ${bug}
            <button onclick="deleteBug(${index})">❌</button>
        `;

        list.appendChild(li);
    });
}