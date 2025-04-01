let data = [];

async function fetchData() {
  await fetch("http://localhost:8080/get_all", {
    method: "GET",
  })
    .then((response) => response.json())
    .then((resData) => {
      data = resData;
    });
  renderTasks();
}
fetchData();

function renderTasks() {
  const taskList = document.getElementById("taskList");
  taskList.innerHTML = "";

  for (let i = 0; i < data.length; i++) {
    const task = data[i];
    const li = document.createElement("li");

    li.className = `task-item ${task.status ? "completed" : ""}`;
    li.innerHTML = `
                <div class="checkbox ${
                  task.status ? "checked" : ""
                }" data-index="${task.Id}">
                    ${task.status ? '<i class="fas fa-check"></i>' : ""}
                </div>
                <span class="task-text">${task.title}</span>
                <button class="delete-btn" data-index="${task.Id}">
                      <i class="fas fa-trash"></i>
                </button>
            `;
    taskList.appendChild(li);

    const deleteButtons = document.getElementsByClassName("delete-btn");
    const checkboxes = document.getElementsByClassName("checkbox");
    for (let i = 0; i < deleteButtons.length; i++) {
      deleteButtons[i].addEventListener("click", deleteTask);
    }
    for (let i = 0; i < checkboxes.length; i++) {
      checkboxes[i].addEventListener("click", toggleTask);
    }
  }
}

const addTaskButton = document.getElementById("addButton");
addTaskButton.addEventListener("click", addTask);
async function addTask() {
  const taskInput = document.getElementById("taskInput");
  await fetch("http://localhost:8080/new_task", {
    method: "POST",
    body: JSON.stringify({
      title: taskInput.value,
    }),
  });
  fetchData();
}

async function deleteTask(event) {
  await fetch(
    `http://localhost:8080/delete/${event.target.attributes["data-index"].value}`,
    {
      method: "DELETE",
    }
  );
  fetchData();
}

async function toggleTask(event) {
  await fetch(
    `http://localhost:8080/change_status/${event.target.attributes["data-index"].value}`,
    {
      method: "POST",
    }
  );
  fetchData();
}
