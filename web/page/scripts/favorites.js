const STAR_ICON = (isFavorite) => `
    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 block mx-auto" viewBox="0 0 24 24" fill="${isFavorite ? "currentColor" : "none"}" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" />
    </svg>
`;

class Favorites {
    constructor() {
        this.base_url = "http://127.0.0.1:3000"
        this.favoritesList = document.getElementById("favorites-list");

        window.favorites = this;
        document.addEventListener("favorites:changed", () => this.loadFavorites());

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

    toggleFavorite(stationId, isFavorite) {
        const request = isFavorite
            ? fetch(this.base_url+`/stations/favorite/${stationId}`, { method: "DELETE" })
            : fetch(this.base_url+"/stations/favorite", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    "station_id": Number(stationId)
                })
            });

        return request
            .then((resp) => resp.json())
            .then((resp) => {
                if (resp.error) throw new Error(resp.error);

                document.dispatchEvent(new CustomEvent("favorites:changed", { detail: resp }));
                return resp;
            })
            .catch((err) => {
                console.error(err);
                alert("Une erreur est survenue au moment de la modification des favoris !")
            })
    }

    renderFavorites(favorites) {
        this.favoritesList.innerHTML = "";
        favorites.forEach((station) => {
            const item = document.createElement("div");
            const volume = station.volume !== undefined ? `${Math.round(station.volume * 100)}%` : "";

            item.innerHTML = `
                <div class="radio flex items-center gap-3 bg-neutral-800 hover:bg-neutral-700 transition-colors rounded-xl p-2 cursor-pointer">
                    <img class="icon_radio w-10 h-10 rounded-lg object-cover" src="${station.favicon}" alt="image" />
                    <div class="flex-1 min-w-0">
                        <span class="radio_name block text-sm font-medium truncate">${station.name}</span>
                        <span class="text-xs text-neutral-400">${volume ? "Volume : "+volume : ""}</span>
                    </div>
                    <button class="favorite_button w-8 h-8 rounded-full text-yellow-400 hover:text-yellow-300 transition-colors flex items-center justify-center flex-shrink-0">
                        ${STAR_ICON(true)}
                    </button>
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
            item.querySelector(".favorite_button").addEventListener("click", () => {
                this.toggleFavorite(station.station_id, true);
            });
            this.favoritesList.appendChild(item);
        });
    }
}

function setup() {
    new Favorites();
}

document.addEventListener("DOMContentLoaded", setup)
