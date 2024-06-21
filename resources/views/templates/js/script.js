const inputBox = document.getElementById("input-box");
const dueDateInput = document.getElementById("due-date");
const listContainer = document.getElementById("list-container");

function addTask() {
    if (inputBox.value === '') {
        alert("You must write something!");
    } else {
        let li = document.createElement("li");
        
        // Membuat elemen untuk teks tugas
        let taskText = document.createElement("div");
        taskText.className = "task-details";
        taskText.innerHTML = inputBox.value;
        li.appendChild(taskText);
        
        // Membuat elemen untuk tanggal
        let taskDates = document.createElement("div");
        taskDates.className = "task-dates";
        const currentDate = new Date().toISOString().split('T')[0];
        taskDates.innerHTML = `Created: ${currentDate} | Due: ${dueDateInput.value || 'No due date'}`;
        li.appendChild(taskDates);
        
        // Membuat elemen tombol hapus
        let span = document.createElement("span");
        span.innerHTML = "\u00D7";
        span.className = "delete-btn";
        li.appendChild(span);
        
        listContainer.appendChild(li);
    }
    inputBox.value = "";
    dueDateInput.value = "";
}

listContainer.addEventListener("click", function(e) {
    if (e.target.tagName === "LI") {
        e.target.classList.toggle("checked");
    } else if (e.target.className === "delete-btn") {
        e.target.parentElement.remove();
    }
}, false);
