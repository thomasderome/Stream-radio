from api.web_radio import test_link_deco
import subprocess

class Stream_station:
    def __init__(self):
        self.stream_url = ""
        self.name = ""

        self.process = None

    @test_link_deco
    def play(self, url = "", name=""):
        if self.process is None and url != "":
            self.stream_url = url
            self.name = name
            try: 
                self.start_process()
                print(f"Lecture du stream : {self.stream_url}")
                return True
            
            except Exception as e:
                print(e)
                return False
            
        elif not self.process is None and url != "":
            try:
                self.stop()
                self.stream_url = url
                self.name = name
                self.play(self.stream_url, name)
                return True
                
            except Exception:
                return False
        
        elif self.stream_url != "" and self.process is None:
            self.start_process()
            return True

    def start_process(self):
        self.process = subprocess.Popen(
            ['ffplay', '-nodisp', '-autoexit' , '-i',self.stream_url], ## 
            stdin=subprocess.PIPE,
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL
        )
    
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
            
    def pause_resume(self):
        if self.process:
            self.stop()
            return True
        else:
            self.play()
            return True