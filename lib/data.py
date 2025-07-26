import json

class Data_manager():
    def __init__(self, file):
        self.file_path = file

        try:
            with open(f'{self.file_path}','r') as f:
                f.close()
        except Exception:
            self._create_json()

    def _create_json(self):
        try:
            with open(f'{self.file_path}', 'a') as f:
                f.write(json.dumps({},indent=4))
                f.close()
            return True
        except:
            return False
        
    def write(self, data: dict):
        with open(f'{self.file_path}', 'r') as f:
            json_file = json.load(f)
            json_file.update(data)
            
        with open(f'{self.file_path}', 'w') as f:
            json.dump(json_file, f, indent=4)
            
    def read(self): 
        with open(f'{self.file_path}', 'r') as f:
            return json.load(f)
        
    def erase(self, data: dict):
        with open(f'{self.file_path}', "w") as f:
            json.dump(data, f, indent=4)