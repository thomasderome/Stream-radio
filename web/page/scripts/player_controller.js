const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

class Player_controller {
    constructor() {
        this.base_url = "http://127.0.0.1:3000"
        this.state_play = false;

        this.buttonPlay = document.getElementById("button_play");

        this.volumeSlider = document.getElementById("volume");
        this.volumeLevel = document.getElementById("volume_value");
        this.allowUpdateVolume = true;
        this.volumeSlider.addEventListener("input", () => this.updateVolumeNumber());

        this.imgRadioRead = document.getElementById("name_radio_read");
        this.nameRadioRead = document.getElementById("img_radio_read");

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

    updateGui(data) {
        //buttonPlay = resp.is_playing
        this.volumeSlider.value = data.volume * 100;
        this.volumeLevel.textContent = `${this.volumeSlider.value}%`;
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
                    method: "PUT"
                })
                    .then((resp) => resp.json())
                    .then((resp) => {
                        this.volumeSlider.value = resp.volume;
                        this.volumeLevel.textContent = `${this.volumeSlider.value}%`;
                    })
                    .catch((err) => {
                        alert("Une erreur c'est produite !");
                        console.error(err);
                    })
                this.allowUpdateVolume = true
            }, 200);
        }
    }
}

function setup() {
    new Player_controller()
}

document.addEventListener("DOMContentLoaded", setup)