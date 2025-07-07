import json

class Data_manager():
    def __init__(self, file):
        self.file = file
        try:
            with open(f'data/{self.file}.json','r') as f:
                f.close()
        except Exception:
            self._create_json()

    def _create_json(self):
        try:
            with open(f'data/{self.file}.json', 'a') as f:
                f.write(json.dumps({},indent=4))
                f.close()
            return True
        except:
            return False
        
    def write(self, data):
        with open(f'data/{self.file}.json', 'r') as f:
            json_file = json.load(f)
            json_file.update(data)
            
        with open(f'data/{self.file}.json', 'w') as f:
            json.dump(json_file, f, indent=4)
            
    def read(self): 
        with open(f'data/{self.file}.json', 'r') as f:
            return json.load(f)