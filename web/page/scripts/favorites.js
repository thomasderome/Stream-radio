class Favorites {
    constructor() {
        this.base_url = "http://127.0.0.1:3000"
        this.favoritesList = document.getElementById("favorites-list");
        this.loadFavorites();
    }

    loadFavorites() {
        fetch(this.base_url+"/stations/favorite", {
            method: "GET"
        })
            .then((resp) => resp.json())
            .then((resp) => {
                this.renderFavorites(resp);
            })
            .catch((err) => {
                console.error(err);
                alert("Une erreur est survenue au moment de la récupération des favoris !")
            })
    }

    renderFavorites(favorites) {
        this.favoritesList.innerHTML = "";
        favorites.forEach((station) => {
            const item = document.createElement("div");

            item.innerHTML = `
                <div class="radio flex items-center gap-3 bg-neutral-800 hover:bg-neutral-700 transition-colors rounded-xl p-2 cursor-pointer">
                    <img class="icon_radio w-10 h-10 rounded-lg object-cover" src="${station.favicon}" alt="image" />
                    <span class="radio_name flex-1 text-sm font-medium truncate">${station.name}</span>
                    <button data-stationId="${station.station_id}" class="radio_button w-9 h-9 rounded-full bg-indigo-500 hover:bg-indigo-400 transition-colors flex items-center justify-center flex-shrink-0">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-white block mx-auto" viewBox="0 0 24 24" fill="currentColor">
                            <path d="M8 5v14l11-7z" />
                        </svg>
                    </button>
                </div>
            `;
            item.querySelector(".radio_button").addEventListener("click", () => {
                window.controller?.playStation(station.station_id);
            });
            this.favoritesList.appendChild(item);
        });
    }
}

function setup() {
    new Favorites();
}

document.addEventListener("DOMContentLoaded", setup)
