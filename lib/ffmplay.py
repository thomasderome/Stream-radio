from api.web_radio import test_link_deco
import vlc

class Stream_station:
    def __init__(self):
        self.stream_url = ""
        self.name = ""

        self.instance = vlc.Instance("--aout=alsa", "--audio-filter=stereo")
        self.player = self.instance.media_player_new()

    @test_link_deco
    def play(self, url = "", name=""):
        if self.is_playing() is False and url != "":
            self.stream_url = url
            self.name = name
            try: 
                self.start_process()
                print(f"Lecture du stream : {self.stream_url}")
                return True
            
            except Exception as e:
                print(e)
                return False
            
        elif self.is_playing() and url != "":
            try:
                self.stop()
                self.stream_url = url
                self.name = name
                self.play(self.stream_url, name)
                return True
                
            except Exception:
                return False
        
        elif self.stream_url != "" and self.is_playing() is False:
            self.start_process()
            return True

    def start_process(self):
        media = self.instance.media_new(self.stream_url)
        self.player.set_media(media)
        self.player.play()

    def stop(self):
        if self.is_playing():
            self.player.stop()
            print("Arrêt du stream...")
            return True
        else:
            print("Aucun stream à arrêter.")
            return False
            
    def pause_resume(self):
        if self.is_playing():
            self.stop()
            return True
        else:
            self.play()
            return True
        
    def set_volume(self, set):
        self.player.audio_set_volume(set)
        return self.get_volume()
        
    def get_volume(self):
        return self.player.audio_get_volume()
        
    def is_playing(self):
        return True if self.player.is_playing() == 1 else False