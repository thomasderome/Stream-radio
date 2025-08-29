from api.web_radio import WEB_radio
from api.spot_credential import Spotify_PKCE
from lib.ffmplay import Stream_station
from lib.stream_spotify import Spotifyd
from lib.data import Data_manager as DB
from os.path import exists

class Core_system():
    def __init__(self):
        self.web_radio_system: WEB_radio = WEB_radio()
        self.spotify_credential: Spotify_PKCE = Spotify_PKCE()
        self.stream_station: Stream_station = Stream_station()
        self.stream_spotify: Spotifyd = Spotifyd()


        self.config_path = "data/config.json"
        if not exists(self.config_path):
            self.db: DB = DB(self.config_path)
            self.db.write({"is_radio": True,
                            "station": {
                                "url": "",
                                "name": ""
                            },
                            "volume": 50})
        else:
            self.db: DB = DB(self.config_path)

        self.config = self.db.read()
        self.is_radio = self.config['is_radio']

        if self.is_radio and self.config['station']['url']:
            self.stream_station.play(self.config['station']['url'], self.config['station']['name'])
            self.stream_station.set_volume(self.config['volume'])
        elif not self.is_radio:
            self.stream_spotify.start_spotifyd()

    def play_radio(self, url, name):
        if self.is_radio:
            if not url == self.config['station']:
                self.db.write({"station": {"url": url, "name": name}})
        
            return self.stream_station.play(url, name)
        
    def set_volume(self, volume: int):
        self.stream_station.set_volume(volume)
        self.db.write({"volume": volume})
        return volume

    def plateforme(self):
        self.is_radio = not self.is_radio
        if self.is_radio:
            self.stream_spotify.stop_spotifyd()
        else:
            self.stream_station.stop()
            self.stream_spotify.start_spotifyd()

        self.db.write({"is_radio": self.is_radio})
          
