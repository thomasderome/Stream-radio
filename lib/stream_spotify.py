import subprocess

class Spotifyd():
    def __init__(self):
        self.process = None

    def start_spotifyd(self):
        if self.process is None:
            self.process = subprocess.Popen(
                ["spotifyd", "--no-daemon", "--disable-discovery"],
                stdin=subprocess.PIPE,
                stdout=subprocess.DEVNULL,
                stderr=subprocess.DEVNULL
            )

    def stop_spotifyd(self):
        if self.process:
            self.process.kill()