let bugs = [];

function addBug() {

    const input = document.getElementById("bugInput");
    const priority = document.getElementById("priority").value;

    const text = input.value.trim();

    if (text === "") {
        alert("Введите баг!");
        return;
    }

    const bug = {
        text: text,
        priority: priority,
        status: "Open"
    };

    bugs.push(bug);

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
            <div class="bug-info">
                <h3>${bug.text}</h3>

                <p>
                    Priority:
                    <span class="priority ${bug.priority}">
                        ${bug.priority}
                    </span>
                </p>

                <p>Status: ${bug.status}</p>
            </div>

            <button class="delete-btn"
                onclick="deleteBug(${index})">
                ❌
            </button>
        `;

        list.appendChild(li);
    });
}