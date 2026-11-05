class Player_controller {
    constructor() {
        this.base_url = "http://127.0.0.1:3000"
        this.isPlaying = undefined;

        window.controller = this;

        this.buttonPlay = document.getElementById("button_play");

        this.volumeSlider = document.getElementById("volume");
        this.volumeLevel = document.getElementById("volume_value");
        this.allowUpdateVolume = true;
        this.volumeSlider.addEventListener("input", () => this.updateVolumeNumber());

        this.imgRadioRead = document.getElementById("img_radio_read");
        this.nameRadioRead = document.getElementById("name_radio_read");
        this.idRadioRead = undefined;
        this.isFavorite = false;

        this.buttonFavorite = document.getElementById("button_favorite");
        this.buttonFavorite.addEventListener("click", () => {
            if (!this.idRadioRead) return;
            window.favorites?.toggleFavorite(this.idRadioRead, this.isFavorite);
        });
        document.addEventListener("favorites:changed", (e) => {
            if (e.detail && e.detail.station_id === this.idRadioRead) {
                this.updateGui({ is_favorite: e.detail.is_favorite });
            }
        });

        this.getPlayerState()
    }

    getPlayerState() {
        fetch(this.base_url+"/player/get_play_data")
            .then((resp) => resp.json())
            .then((resp) => {
                this.updateGui(resp)
            })
            .catch((err) => {
                alert("Une erreur c'est produite !");
                console.error(err);
            })
    }

    playStation(stationId) {
        fetch(this.base_url+"/player/set_play_station_id", {
            method: "PUT",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                "station_id": String(stationId)
            })
        })
            .then((resp) => resp.json())
            .then((resp) => {
                this.updateGui(resp);
            })
            .catch((e) => {
                alert("Une erreur c'est produite !");
                console.error(e);
            });
    }

    updateGui(data) {
        if (data.volume !== undefined && data.volume !== null) {
            const volumeTo100 = Math.round(data.volume * 100);
            this.volumeSlider.value = volumeTo100;
            this.volumeLevel.textContent = `${volumeTo100}%`;
        }

        if (data.is_playing !== undefined && data.is_playing !== null) {
            if (data.is_playing) {
                this.buttonPlay.innerHTML = `
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-white block mx-auto" viewBox="0 0 24 24" fill="currentColor">
                        <rect x="6" y="5" width="4" height="14" rx="1"></rect>
                        <rect x="14" y="5" width="4" height="14" rx="1"></rect>
                    </svg>
                `;
            } else {
                this.buttonPlay.innerHTML = `
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-white block mx-auto" viewBox="0 0 24 24" fill="currentColor">
                        <path d="M8 5v14l11-7z" />
                    </svg>
                `;
            }
            this.isPlaying = data.is_playing;
            this.buttonPlay.onclick = () => this.playResume(this);
        }
        if (data.station_name) {
            this.nameRadioRead.textContent = data.station_name
        }
        if (data.station_id) {
            this.idRadioRead = data.station_id;
        }
        if (data.station_img) {
            this.imgRadioRead.src = data.station_img;
        }
        if (data.is_favorite !== undefined && data.is_favorite !== null) {
            this.isFavorite = data.is_favorite;
            this.buttonFavorite.innerHTML = STAR_ICON(this.isFavorite);
        }
    }

    playResume() {
        this.buttonPlay.disabled = true;
        fetch(this.base_url+"/player/set_play_state", {
            method: "PUT",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                "play_state": !this.isPlaying
            })
        })
            .then((resp) => resp.json())
            .then((resp) => {
                this.updateGui(resp);
                this.buttonPlay.disabled = false;
            })
            .catch((e) => {
                alert("Une erreur c'est produite !");
                console.error(e);
                this.buttonPlay.disabled = false;
            });
    }

    updateVolumeNumber() {
        this.volumeLevel.textContent = `${this.volumeSlider.value}%`;

        if (this.allowUpdateVolume) {
            this.allowUpdateVolume = false

            setTimeout(() => {
                fetch(this.base_url+"/player/set_volume",  {
                    body: JSON.stringify({
                        volume: parseFloat(this.volumeSlider.value) / 100,
                    }),
                    method: "PUT",
                    headers: {
                        "Content-Type": "application/json"
                    }
                })
                    .then((resp) => resp.json())
                    .then((resp) => {
                        this.updateGui(resp);
                        this.allowUpdateVolume = true;

                        // Le volume est mémorisé côté API pour un favori : on rafraîchit la liste
                        if (this.isFavorite) {
                            window.favorites?.loadFavorites();
                        }
                    })
                    .catch((err) => {
                        alert("Une erreur c'est produite !");
                        console.error(err);
                        this.allowUpdateVolume = true;
                    })
            }, 200);
        }
    }
}

function setup() {
    controller = new Player_controller();
}

document.addEventListener("DOMContentLoaded", setup)