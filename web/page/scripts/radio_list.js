class Radio_list {
    constructor() {
        this.lazyscroll_div = document.getElementById("lazy-loading");
        this.in_update = false;

        this.stationContainer = document.getElementById("station_container");

        this.offset = 0;
        this.search = "";
        this.tags = "";

        this.base_url = "http://127.0.0.1:3000"


        let options = {
            root: null,
            rootMargins: "0px",
            threshold: 0.5
        };

        const observer = new IntersectionObserver((entries) => this.handleObs(entries), options);
        observer.observe(this.lazyscroll_div);
    }

    handleObs(entries) {
        if (entries[0].isIntersecting && !this.in_update) {
            this.updateFront()
        }
    }

    updateFront() {
        if (this.in_update) return;

        this.in_update = true;
        fetch(this.base_url+`/stations/station?offset=${this.offset}${this.search ? "&search="+this.search : ""}${this.tags ? "&tags="+this.tags : ""}`, {
            method: "GET"
        })
            .then((resp) => resp.json())
            .then((resp) => {
                console.log(resp);
                this.renderData(resp);

            })
            .catch((e) => {
                console.error(e);
                alert("Une erreur est survenue au moment de la récupération des radios !")
            })
    }

    renderData(data) {
        console.log(data.tags)
        data.forEach((radio) => {
            const new_station = document.createElement("div");

            new_station.innerHTML = `
                <div class="radio bg-neutral-700 hover:bg-neutral-600 transition-colors rounded-2xl p-4 flex items-center gap-4 shadow-lg">
                    <img class="icon_radio w-16 h-16 rounded-xl object-cover flex-shrink-0" src="${radio.favicon}" alt="image" />
                    <div class="flex-1 min-w-0">
                        <span class="radio_name block font-semibold truncate">${radio.name}</span>
                        <div class="flex gap-1 mt-1 flex-wrap tags">
                            ${radio.tags.map((item) => {
                if (item.length >= 10) return "";
                return `<span className="tag text-[10px] uppercase tracking-wide bg-emerald-500/20 text-emerald-300 px-2 py-0.5 rounded-full">${item}</span>`
            }).join("")}
                        </div>
                    </div>
                    <button data-stationId="${radio.station_id}" class="radio_button bg-indigo-500 hover:bg-indigo-400 transition-colors text-sm font-medium px-4 py-2 rounded-full flex-shrink-0 flex items-center justify-center">
                        <span>Play</span>
                    </button>
                </div>
            `;
            this.stationContainer.appendChild(new_station);
        });
        this.in_update = false;
    }
}

function setup() {
    new Radio_list();
}

document.addEventListener("DOMContentLoaded", setup)