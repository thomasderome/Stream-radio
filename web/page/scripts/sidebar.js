class Sidebar {
    constructor(default_state = true) {
        this.state = default_state;

        this.container = document.getElementById("container");
        this.sidebar_button = document.getElementById("sidebar_button");

        this.init();
    }

    init() {
        this.sidebar_button.addEventListener("click", () => this.toggle());
        if (this.state) this.container.classList.remove("close");
        else this.container.classList.add("close");
    }

    toggle() {
        this.render();
        this.state = !this.state;
    }

    render() {
        if (this.state) this.container.classList.add("close");
        else this.container.classList.remove("close");
    }
}

function run() { const sidebar = new Sidebar(); }
document.addEventListener("DOMContentLoaded", run);