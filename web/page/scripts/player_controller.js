class Player_controller {
    constructor() {
        this.base_url = "http://127.0.0.1:3000"
        this.state_play = false;

        this.buttonPlay = document.getElementById("button_play");
        this.voluleSlider = document.getElementById("volume");
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
    }

    updateGui(data) {
        //buttonPlay = resp.is_playing
        this.voluleSlider.value = data.volume;
    }

}

function setup() {
    new Player_controller()
}

document.addEventListener("DOMContentLoaded", setup)