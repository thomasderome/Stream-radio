const URL = `${window.location.protocol}//${window.location.host}`;
const range = 2;
let counter = 0;

function createRadioElement(radio) {
    console.log(radio)
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
    button.dataset.img = radio[1].img;
    button.dataset.name = radio[0];

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
    console.log("fetch some JSON data");
    fetch(`${URL}/radio?line_sel=${counter*range}&num_sel=${range}`)
        .then(response => response.json())
        .then(data => {
            console.log(data)
            renderRadios(data)
    });
}
