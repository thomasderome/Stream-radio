function update_control_menu(name_station = "", img = "", status = "") {
  document.getElementById('button_control').style.display = "block";
  let button = document.getElementById('pause_resume');
  console.log(status)
  if (status == true) button.textContent = "⏸️";
  else button.textContent = "▶️";

  const img_play = document.getElementById('now-playing-img');
  img_play.src = img;
  img_play.style.display = "block"; 
  document.getElementById('now-playing-title').textContent = name_station;
}

function renderControl(data) {
  const controls = Object.values(data.data);
  const status = data.status;

  controls.forEach(control => {
    update_control_menu(control.name, control.img, status);
  });
}
function control_menu () {
  fetch(`${URL}/control_menu`, {method: "GET"})
  .then(response => response.json())
  .then(data => {
    renderControl(data);
  });
}

function setup() {
  control_menu();

  pause_resume = document.getElementById('pause_resume');
  pause_resume.addEventListener('click', () => {
    fetch(`${URL}/pause_resume`, {method: "PUT"})
    .then(response => response.json())
    .then(data => {
      control_menu();
    });
  });
}

document.addEventListener("DOMContentLoaded", () => {
  setup();
});