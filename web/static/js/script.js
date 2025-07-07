const radioList = document.getElementById('radio-list');
const searchInput = document.getElementById('search');
let currentAudio = null;

const URL = `${window.location.protocol}//${window.location.host}`;
const range = 20;
let counter = 1;

function createRadioElement(radio) {
  const div = document.createElement('div');
  div.className = 'radio';

  const img = document.createElement('img');
  img.src = radio[1].img;
  img.alt = radio[0];

  const info = document.createElement('div');
  info.className = 'radio-info';
  info.textContent = radio[0];

  const button = document.createElement('button');
  button.textContent = 'Play';
  button.dataset.url = radio[1].url;
  
  button.addEventListener('click', () => {
    const url = button.dataset.url;
    console.log('URL à lire :', url);
  });

  div.appendChild(img);
  div.appendChild(info);
  div.appendChild(button);

  return div;
}

function renderRadios(radios) {
  radioList.innerHTML = "";
  radios
    .forEach(r => {
      radioList.appendChild(createRadioElement(r));
    });
}

searchInput.addEventListener('input', (e) => {
  var search = e.target.value.trim();
  if (search.length < 3) {  
    document.querySelector("footer").style.display = "block";

    getData();
  } else if (search.length >= 3) {
    document.querySelector("footer").style.display = "none";

    fetch(`${URL}/search?q=${search}`)
    .then(response => response.json())
    .then(data => {
      renderRadios(data);
    });
  }
});


document.addEventListener("DOMContentLoaded", () => {
  let options = {
    root: null,
    rootMargins: "0px",
    threshold: 0.5
  };
  getData();
  
  const observer = new IntersectionObserver(handleIntersect, options);
  observer.observe(document.querySelector("footer"));
});

function handleIntersect(entries) {
  if (entries[0].isIntersecting) {
    console.warn("something is intersecting with the viewport");
    counter++;
    getData();
  }
}

function getData() {
  console.log("fetch some JSON data");
  fetch(`${URL}/radio?num=${counter*range}`)
    .then(response => response.json())
    .then(data => {
      renderRadios(data)
    });
}