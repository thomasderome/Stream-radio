class Radio_list {
    constructor() {
        this.stationContainer = document.getElementById("station_container");
        this.lazyscroll_div = document.getElementById("lazy-loading");
        this.in_update = false;
        this.no_more_data = false;
        this.isVisible = false;

        this.offset = 0;
        this.search = "";
        this.tags = "";

        this.base_url = "http://127.0.0.1:3000"

        this.searchInput = document.getElementById("search");
        this.tagFilter = document.getElementById("tag-filter");

        let searchTimeout = undefined;
        this.searchInput.addEventListener("input", () => {
            clearTimeout(searchTimeout);
            searchTimeout = setTimeout(() => {
                this.search = this.searchInput.value.trim();
                this.reset();
            }, 300);
        });

        this.tagFilter.addEventListener("change", () => {
            this.tags = this.tagFilter.value;
            this.reset();
        });

        document.addEventListener("favorites:changed", (e) => this.updateFavoriteButton(e.detail));

        let options = {
            root: null,
            rootMargin: "300px 0px",
            threshold: 0
        };

        const observer = new IntersectionObserver((entries) => this.handleObs(entries), options);
        observer.observe(this.lazyscroll_div);
    }

    reset() {
        this.offset = 0;
        this.no_more_data = false;
        this.stationContainer.innerHTML = "";
        this.stationContainer.appendChild(this.lazyscroll_div);
        this.updateFront();
    }

    handleObs(entries) {
        this.isVisible = entries[0].isIntersecting;
        if (this.isVisible) {
            this.updateFront()
        }
    }

    updateFront() {
        if (this.in_update || this.no_more_data) return;

        this.in_update = true;
        fetch(this.base_url+`/stations/station?offset=${this.offset}${this.search ? "&search="+encodeURIComponent(this.search) : ""}${this.tags ? "&tags="+this.tags : ""}`, {
            method: "GET"
        })
            .then((resp) => resp.json())
            .then((resp) => {
                if (!Array.isArray(resp) || resp.length === 0) {
                    this.no_more_data = true;
                } else {
                    this.renderData(resp);
                    this.offset += resp.length;
                }
                this.in_update = false;

                if (!this.no_more_data && this.isVisible) {
                    this.updateFront();
                }
            })
            .catch((e) => {
                this.in_update = false;
                console.error(e);
                alert("Une erreur est survenue au moment de la récupération des radios !")
            })
    }

    renderData(data) {
        data.forEach((radio) => {
            const new_station = document.createElement("div");

            new_station.innerHTML = `
                <div class="radio bg-neutral-700 hover:bg-neutral-600 transition-colors rounded-2xl p-4 flex items-center gap-4 shadow-lg">
                    <img class="icon_radio w-16 h-16 rounded-xl object-cover flex-shrink-0" src="${radio.favicon}" alt="image" />
                    <div class="flex-1 min-w-0">
                        <span class="radio_name block font-semibold truncate">${radio.name}</span>
                        <div class="flex gap-1 mt-1 flex-wrap tags">
                            ${radio.tags.map((item) => {
                if (!item || item.length >= 10) return "";
                return `<span class="tag text-[10px] uppercase tracking-wide bg-emerald-500/20 text-emerald-300 px-2 py-0.5 rounded-full">${item}</span>`
            }).join("")}
                        </div>
                    </div>
                    <button class="favorite_button w-9 h-9 rounded-full text-yellow-400 hover:text-yellow-300 transition-colors flex items-center justify-center flex-shrink-0">
                        ${STAR_ICON(radio.is_favorite)}
                    </button>
                    <button data-stationId="${radio.station_id}" class="radio_button bg-indigo-500 hover:bg-indigo-400 transition-colors text-sm font-medium px-4 py-2 rounded-full flex-shrink-0 flex items-center justify-center">
                        <span>Play</span>
                    </button>
                </div>
            `;
            new_station.dataset.stationId = radio.station_id;
            new_station.dataset.isFavorite = radio.is_favorite;
            new_station.querySelector(".radio_button").addEventListener("click", () => {
                window.controller?.playStation(radio.station_id);
            });
            new_station.querySelector(".favorite_button").addEventListener("click", () => {
                window.favorites?.toggleFavorite(radio.station_id, new_station.dataset.isFavorite === "true");
            });
            this.stationContainer.insertBefore(new_station, this.lazyscroll_div);
        });
        this.in_update = false;
    }

    updateFavoriteButton(data) {
        if (!data || data.station_id === undefined) return;

        const station = this.stationContainer.querySelector(`[data-station-id="${data.station_id}"]`);
        if (!station) return;

        station.dataset.isFavorite = data.is_favorite;
        station.querySelector(".favorite_button").innerHTML = STAR_ICON(data.is_favorite);
    }
}

function setup() {
    new Radio_list();
}

document.addEventListener("DOMContentLoaded", setup)
