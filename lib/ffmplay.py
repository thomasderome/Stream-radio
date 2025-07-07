from api.web_radio import test_link_deco
import subprocess

class Stream_station:
    def __init__(self):
        self.stream_url = ""
        self.process = None

    @test_link_deco
    def play(self, url = ""):
        if self.process is None and url != "":
            self.stream_url = url
            try: 
                self.process = subprocess.Popen(
                    ['ffplay', '-nodisp', '-autoexit' , '-i',self.stream_url], ## 
                    stdin=subprocess.PIPE,
                    stdout=subprocess.DEVNULL,
                    stderr=subprocess.DEVNULL
                )
                
                print(f"Lecture du stream : {self.stream_url}")
                return True
            
            except Exception as e:
                print(e)
                return False
            
        elif not self.process is None and url != "":
            try:
                self.stop()
                self.stream_url = url
                self.play(self.stream_url)
                return True
                
            except Exception:
                return False
        
        elif not self.stream_url == "":
            return False

    def stop(self):
        if self.process:
            self.process.terminate()
            self.process.wait()
            self.process = None
            print("Arrêt du stream...")
            return True
        else:
            print("Aucun stream à arrêter.")
            return False
            
    def pause(self):
        if self.process:
            self.stop()
            return True
        else:
            return False
    
    def resume(self):
        if not self.process:
            self.play(self.stream_url)
            return True
        else:
            return False