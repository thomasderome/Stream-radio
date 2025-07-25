const radioList = document.getElementById('radio-list');
const searchInput = document.getElementById('search');
let currentAudio = null;

function update_control_menu(name_station = "", img = "") {
  if (name_station == "" && img == "") {
    document.getElementById('button_control').style.display = "none";
    const img = document.getElementById('now-playing-img');
    img.src = "";
    img.style.display = "none";

    document.getElementById('now-playing-title').textContent = "Aucune lecture en cour";
  } else {
    document.getElementById('button_control').style.display = "block";
    const img_play = document.getElementById('now-playing-img');
    img_play.src = img;
    img_play.style.display = "block";

    document.getElementById('now-playing-title').textContent = name_station;
  }
}

function setup() {
  pause_resume = document.getElementById('pause_resume');
  pause_resume.addEventListener('click', () => {
    fetch(`${URL}/pause_resume`, {method: "PUT"})
    .then(response => response.json())
    .then(data => {
      if (data.response == 'pause') {
        pause_resume.textContent = "▶️"
      } else if (data.response == 'resume') {
        pause_resume.textContent = "⏸️";
      }
    });
  });
}