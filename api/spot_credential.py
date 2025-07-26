import string, random, hashlib, secrets, base64, requests, os
from lib.data import Data_manager as db
from urllib.parse import urlencode

################################
#   CLASS CREATED BY DOC: https://developer.spotify.com/documentation/web-api/tutorials/code-pkce-flow
################################
class Spotify_PKCE(db):
    def __init__(self):
        self.url_token ="https://accounts.spotify.com/api/token"

        self.base_url = "https://accounts.spotify.com/authorize?"
        self.client_id = "65b708073fc0480ea92a077233ca87bd"
        self.redirect_url = "http://127.0.0.1:8080/login"

        self.state = self.gen_state()
        self.code_verif = self.generate_code_verifier()
        self.code_challenge = self.generate_code_challenge(self.code_verif)

        self.scope = [        
            "playlist-modify",
            "playlist-modify-private",
            "playlist-modify-public",
            "playlist-read",
            "playlist-read-collaborative",
            "playlist-read-private",
            "streaming",
            "ugc-image-upload",
            "user-follow-modify",
            "user-follow-read",
            "user-library-modify",
            "user-library-read",
            "user-modify",
            "user-modify-playback-state",
            "user-modify-private",
            "user-personalized",
            "user-read-birthdate",
            "user-read-currently-playing",
            "user-read-email",
            "user-read-play-history",
            "user-read-playback-position",
            "user-read-playback-state",
            "app-remote-control",
            "user-read-recently-played",
            "user-top-read",
            "user-read-private",
        ]

        self.param = {
            'response_type': "code",
            'client_id': self.client_id,
            'code_challenge_method': 'S256',
            'code_challenge': self.code_challenge,
            'redirect_uri': self.redirect_url
        }

        path = os.path.expanduser('~/.cache/spotifyd/oauth')
        os.makedirs(path, exist_ok=True)
        self.data = os.path.join(path, 'credentials.json')
        super().__init__(self.data)

    def gen_state(self, length = 16):
        return ''.join(random.choices(string.ascii_letters + string.digits, k=length))
    
    def generate_code_verifier(self, length=64):
        allowed_chars = string.ascii_letters + string.digits + "-._~"
        return ''.join(secrets.choice(allowed_chars) for _ in range(length))
    
    def generate_code_challenge(self, code_verifier):
        sha256_digest = hashlib.sha256(code_verifier.encode('utf-8')).digest()
        challenge = base64.urlsafe_b64encode(sha256_digest).decode('utf-8')
        return challenge.rstrip('=')
    
    def gen_url(self):
        query_url = urlencode(self.param)
        return self.base_url+query_url+f"&scope={"+".join(self.scope)}"
    
    def get_token_acces(self, code):
        response = requests.post(url=self.url_token, headers={'Content-Type': 'application/x-www-form-urlencoded'},
                      data={"client_id": self.client_id,
                            "grant_type": 'authorization_code',
                            "code": code,
                            "redirect_uri": self.redirect_url,
                            "code_verifier": self.code_verif,
                            }, allow_redirects=True)
        
        self.setup_spotifyd_credential(response.json())
        
        return response.status_code, response.url, response.content

    def setup_spotifyd_credential(self, data):
        self.erase({"username": data['username'],
                    "auth_type": 1,
                    "auth_data": data['access_token']
                    })
