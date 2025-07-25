from bs4 import BeautifulSoup
from lib.data import Data_manager as db 
import requests, time, random

class WEB_radio():
    def __init__(self):
        self.url = 'https://onlineradiobox.com/fr'

        self.max_page = 92
        self.db = db("cache_radio")
        self.radio = self._cache_loading()
        if self.radio == {}:
            self.radio = self._scrap_radio()
        
    def _cache_loading(self):
        try:
            return self.db.read()
        except Exception as e:
            print(f'ERREUR: {e}')
            return False
        
    def _scrap_radio(self):
        url_radio = {}
        
        for page in range(0, self.max_page + 1):
            response = requests.get(f'{self.url}/?cs=fr.nrjfrance&p={page}')
            
            if self._response(response):
                soup = BeautifulSoup(response.content, 'lxml')
                radio_scrap = soup.find("ul", id="stations")

                image = radio_scrap.find_all('img', class_="station__title__logo")
                names = radio_scrap.find_all('figcaption', class_="station__title__name")
                radios = radio_scrap.find_all('button', class_="station_play")

                for name, radio, image in zip(names, radios, image):
                    url_radio[name.get_text()] = {"name": name.get_text(), "url": radio.get('stream'), "img": image.get('src')}
                    print(name.text)

            time.sleep(random.randint(1,3))
            print(f"radio: {len(url_radio)}")    
            print(f"page: {page}")
            
        print(len(url_radio))    
        self.db.write({"station": url_radio})  

    def get_radio(self, line_sel, num_sel):
        cache = list(self.radio['station'].items())[line_sel:line_sel + num_sel]
        return dict(cache)
    
    def search(self, query):
        result = {}
        for cle, valeur in self.radio['station'].items():
            if query.lower() in cle.lower():
                result[cle] = valeur
        
        return dict(result.items())

    def _response(self, response):
        if response.status_code == 200:
            return True
        else:
            return False


def test_link_deco(func):
    def wrapper(self, url = "", name_staiton = ""):
        if url == "":
            func(self)
        else:
            headers = {
                "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0",
                "Accept": "*/*",
                "Accept-Language": "fr,fr-FR;q=0.8,en-US;q=0.5,en;q=0.3",
                "Sec-GPC": "1",
                "Sec-Fetch-Dest": "empty",
                "Sec-Fetch-Mode": "cors",
                "Sec-Fetch-Site": "cross-site",
                "Priority": "u=0",
                "Referer": "https://onlineradiobox.com/"
            }
            
            response = requests.get(url, headers=headers, stream=True)
    
            if response.status_code == 200:
                func(self, url, name_staiton)
            else:
                response=requests.get(f'https://onlineradiobox.com/search?q={name_staiton}&c=fr')
                soup = BeautifulSoup(response.content, 'lxml')
                radio_scrap = soup.find("ul", id="stations")
                radio_scrap = radio_scrap.find_all('button', class_="station_play")[0]
                
                func(self, radio_scrap.get('stream'), name_staiton)
                
                
    return wrapper