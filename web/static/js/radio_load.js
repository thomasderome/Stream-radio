const radioList = document.getElementById('radio-list');
const searchInput = document.getElementById('search');
const URL = `${window.location.protocol}//${window.location.host}`;
const range = 20;
let counter = 0;

function createRadioElement(radio) {
    console.log(radio)
    const div = document.createElement('div');
    div.className = 'radio';

    const img = document.createElement('img');
    img.src = radio.img;
    img.alt = radio.name;

    const info = document.createElement('div');
    info.className = 'radio-info';
    info.textContent = radio.name;

    const button = document.createElement('button');
    button.textContent = 'Play';
    button.dataset.url = radio.url;
    button.dataset.img = radio.img;
    button.dataset.name = radio.name;

    button.addEventListener('click', () => {
        fetch(`${URL}/play_radio?stream_url=${encodeURIComponent(button.dataset.url)}&name_station=${encodeURIComponent(button.dataset.name)}`, {method: "PUT"})
        .then(response => response.json())
        .then(data => {
        if (data.response == "True") {
            update_control_menu(button.dataset.name, button.dataset.img);
        }
        });
    });

    div.appendChild(img);
    div.appendChild(info);
    div.appendChild(button);

    return div;
}

function renderRadios(radios) {
    const radioArray = Object.values(radios);
    radioArray.forEach(radio => {
        radioList.appendChild(createRadioElement(radio));
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
            radioList.innerHTML = "";
            counter = 0;
            renderRadios(data);
        });
    }
});

document.addEventListener("DOMContentLoaded", async () => {
    let options = {
        root: null,
        rootMargins: "0px",
        threshold: 0.5
    };

    const observer = new IntersectionObserver(handleIntersect, options);
    observer.observe(document.querySelector("footer"));
});

function handleIntersect(entries) {
    if (entries[0].isIntersecting) {
        getData();
        counter++;
    }
}

function getData() {
    if (counter == 0) radioList.innerHTML = ""
    fetch(`${URL}/radio?line_sel=${counter*range}&num_sel=${range}`)
        .then(response => response.json())
        .then(data => {
            renderRadios(data)
    });
}
