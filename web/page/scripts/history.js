class History {
    constructor() {
        this.base_url = "http://127.0.0.1:3000"
        this.historyList = document.getElementById("history-list");
        this.clearButton = document.getElementById("clear-history");

        this.clearButton.addEventListener("click", () => this.clearHistory());
        document.addEventListener("station:played", () => this.loadHistory());
        document.addEventListener("favorites:changed", () => this.loadHistory());

        this.loadHistory();
    }

    loadHistory() {
        fetch(this.base_url+"/stations/history", {
            method: "GET"
        })
            .then((resp) => resp.json())
            .then((resp) => {
                this.renderHistory(resp);
            })
            .catch((err) => {
                console.error(err);
                alert("Une erreur est survenue au moment de la récupération de l'historique !")
            })
    }

    clearHistory() {
        fetch(this.base_url+"/stations/history", {
            method: "DELETE"
        })
            .then((resp) => resp.json())
            .then((resp) => {
                this.renderHistory(resp);
            })
            .catch((err) => {
                console.error(err);
                alert("Une erreur est survenue au moment de la suppression de l'historique !")
            })
    }

    renderHistory(history) {
        this.historyList.innerHTML = "";

        if (history.length === 0) {
            this.historyList.innerHTML = `<span class="text-xs text-neutral-500 px-1">Aucune écoute pour le moment</span>`;
            return;
        }

        history.forEach((station) => {
            const item = document.createElement("div");

            item.innerHTML = `
                <div class="radio flex items-center gap-3 bg-neutral-800 hover:bg-neutral-700 transition-colors rounded-xl p-2 cursor-pointer">
                    <img class="icon_radio w-10 h-10 rounded-lg object-cover" src="${station.favicon}" alt="image" />
                    <span class="radio_name flex-1 text-sm font-medium truncate">${station.name}</span>
                    <button class="favorite_button w-8 h-8 rounded-full text-yellow-400 hover:text-yellow-300 transition-colors flex items-center justify-center flex-shrink-0">
                        ${STAR_ICON(station.is_favorite)}
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
                window.favorites?.toggleFavorite(station.station_id, station.is_favorite);
            });
            this.historyList.appendChild(item);
        });
    }
}

function setup() {
    new History();
}

document.addEventListener("DOMContentLoaded", setup)
