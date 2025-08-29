from flask import Flask, render_template, request, jsonify, redirect
from api import web_radio, spot_credential
from lib.ffmplay import Stream_station
from lib.core import Core_system
import os

"""
$$$$$$\           $$\   $$\                                           $$\
\_$$  _|          \__|  $$ |                                          $$ |
  $$ |  $$$$$$$\  $$\ $$$$$$\          $$$$$$$\ $$\   $$\  $$$$$$$\ $$$$$$\    $$$$$$\  $$$$$$\$$$$\
  $$ |  $$  __$$\ $$ |\_$$  _|        $$  _____|$$ |  $$ |$$  _____|\_$$  _|  $$  __$$\ $$  _$$  _$$\
  $$ |  $$ |  $$ |$$ |  $$ |          \$$$$$$\  $$ |  $$ |\$$$$$$\    $$ |    $$$$$$$$ |$$ / $$ / $$ |
  $$ |  $$ |  $$ |$$ |  $$ |$$\        \____$$\ $$ |  $$ | \____$$\   $$ |$$\ $$   ____|$$ | $$ | $$ |
$$$$$$\ $$ |  $$ |$$ |  \$$$$  |      $$$$$$$  |\$$$$$$$ |$$$$$$$  |  \$$$$  |\$$$$$$$\ $$ | $$ | $$ |
\______|\__|  \__|\__|   \____/       \_______/  \____$$ |\_______/    \____/  \_______|\__| \__| \__|
                                                $$\   $$ |
                                                \$$$$$$  |
                                                 \______/
"""
template_dir = os.path.join(os.path.dirname(__file__), 'web', 'templates')
static_dir = os.path.join(os.path.dirname(__file__), 'web', 'static')
app = Flask("Radio Garage", template_folder=template_dir, static_folder=static_dir)

Core = Core_system()

#radio = web_radio.WEB_radio()
#play_radio = Stream_station()
#sp_pkce = spot_credential.Spotify_PKCE()

def start(app: Flask, ip, port):
    return app.run(ip, port, ssl_context=('ssl/cert.pem', 'ssl/key.pem'), debug=False)

"""
 $$$$$$\                                       $$$$$$\                        $$\
$$  __$$\                                     $$  __$$\                       $$ |
$$ /  \__| $$$$$$\   $$$$$$\   $$$$$$\        $$ /  \__|$$\   $$\  $$$$$$$\ $$$$$$\    $$$$$$\  $$$$$$\$$$$\
$$ |      $$  __$$\ $$  __$$\ $$  __$$\       \$$$$$$\  $$ |  $$ |$$  _____|\_$$  _|  $$  __$$\ $$  _$$  _$$\
$$ |      $$ /  $$ |$$ |  \__|$$$$$$$$ |       \____$$\ $$ |  $$ |\$$$$$$\    $$ |    $$$$$$$$ |$$ / $$ / $$ |
$$ |  $$\ $$ |  $$ |$$ |      $$   ____|      $$\   $$ |$$ |  $$ | \____$$\   $$ |$$\ $$   ____|$$ | $$ | $$ |
\$$$$$$  |\$$$$$$  |$$ |      \$$$$$$$\       \$$$$$$  |\$$$$$$$ |$$$$$$$  |  \$$$$  |\$$$$$$$\ $$ | $$ | $$ |
 \______/  \______/ \__|       \_______|       \______/  \____$$ |\_______/    \____/  \_______|\__| \__| \__|
                                                        $$\   $$ |
                                                        \$$$$$$  |
                                                         \______/
"""
@app.route('/', methods=['GET'])
def page_web():
    return render_template("index.html")

@app.route('/plateform', methods=['PUT'])
def plateforme():
    data = request.args.get('set_plateforme')
    


"""
$$$$$$$\                  $$\ $$\                                                 $$\
$$  __$$\                 $$ |\__|                                                $$ |
$$ |  $$ | $$$$$$\   $$$$$$$ |$$\  $$$$$$\         $$$$$$$\ $$\   $$\  $$$$$$$\ $$$$$$\    $$$$$$\  $$$$$$\$$$$\
$$$$$$$  | \____$$\ $$  __$$ |$$ |$$  __$$\       $$  _____|$$ |  $$ |$$  _____|\_$$  _|  $$  __$$\ $$  _$$  _$$\
$$  __$$<  $$$$$$$ |$$ /  $$ |$$ |$$ /  $$ |      \$$$$$$\  $$ |  $$ |\$$$$$$\    $$ |    $$$$$$$$ |$$ / $$ / $$ |
$$ |  $$ |$$  __$$ |$$ |  $$ |$$ |$$ |  $$ |       \____$$\ $$ |  $$ | \____$$\   $$ |$$\ $$   ____|$$ | $$ | $$ |
$$ |  $$ |\$$$$$$$ |\$$$$$$$ |$$ |\$$$$$$  |      $$$$$$$  |\$$$$$$$ |$$$$$$$  |  \$$$$  |\$$$$$$$\ $$ | $$ | $$ |
\__|  \__| \_______| \_______|\__| \______/       \_______/  \____$$ |\_______/    \____/  \_______|\__| \__| \__|
                                                            $$\   $$ |
                                                            \$$$$$$  |
                                                             \______/
"""

@app.route('/radio', methods=['GET'])
def radio_list():
    line_sel = request.args.get('line_sel')
    num_sel = request.args.get('num_sel')

    return jsonify(Core.web_radio_system.get_radio(line_sel=int(line_sel), num_sel=int(num_sel)))

@app.route('/search', methods=['GET'])
def search():
    query = request.args.get('q') 
    return Core.web_radio_system.search(query=query)

@app.route('/play_radio', methods=['PUT'])
def play_audio():
    try:
        stream_url = request.args.get('stream_url')
        name_station = request.args.get('name_station')
        
        Core.play_radio(stream_url, name_station)
        
        return jsonify({"response": "True"})
    
    except Exception:
        return jsonify({"response": "False"})
    
@app.route('/pause_resume', methods=['PUT'])
def pause_resume():
    return jsonify(Core.stream_station.pause_resume())

@app.route('/set_volume', methods=['PUT'])
def set_volume():
    volume = request.args.get('volume')
    Core.set_volume(int(volume))
    return volume

@app.route('/control_menu', methods=['GET'])
def control_menu():
    if Core.stream_station.name != "":
        return {"data": Core.web_radio_system.search(query=Core.stream_station.name), "status": Core.stream_station.is_playing(), "volume":Core.stream_station.get_volume()}
    else:
        return jsonify("nothing")
    
"""
 $$$$$$\                       $$\     $$\  $$$$$$\                                                  $$\
$$  __$$\                      $$ |    \__|$$  __$$\                                                 $$ |
$$ /  \__| $$$$$$\   $$$$$$\ $$$$$$\   $$\ $$ /  \__|$$\   $$\        $$$$$$$\ $$\   $$\  $$$$$$$\ $$$$$$\    $$$$$$\  $$$$$$\$$$$\
\$$$$$$\  $$  __$$\ $$  __$$\\_$$  _|  $$ |$$$$\     $$ |  $$ |      $$  _____|$$ |  $$ |$$  _____|\_$$  _|  $$  __$$\ $$  _$$  _$$\
 \____$$\ $$ /  $$ |$$ /  $$ | $$ |    $$ |$$  _|    $$ |  $$ |      \$$$$$$\  $$ |  $$ |\$$$$$$\    $$ |    $$$$$$$$ |$$ / $$ / $$ |
$$\   $$ |$$ |  $$ |$$ |  $$ | $$ |$$\ $$ |$$ |      $$ |  $$ |       \____$$\ $$ |  $$ | \____$$\   $$ |$$\ $$   ____|$$ | $$ | $$ |
\$$$$$$  |$$$$$$$  |\$$$$$$  | \$$$$  |$$ |$$ |      \$$$$$$$ |      $$$$$$$  |\$$$$$$$ |$$$$$$$  |  \$$$$  |\$$$$$$$\ $$ | $$ | $$ |
 \______/ $$  ____/  \______/   \____/ \__|\__|       \____$$ |      \_______/  \____$$ |\_______/    \____/  \_______|\__| \__| \__|
          $$ |                                       $$\   $$ |                $$\   $$ |
          $$ |                                       \$$$$$$  |                \$$$$$$  |
          \__|                                        \______/                  \______/
"""
@app.route('/spotify_pkce')
def spotify_pkce():
    return redirect(Core.spotify_credential.gen_url())

@app.route('/login')
def login():
    print('---------------------------------------------')
    print(request.args.get('code'))
    print('---------------------------------------------')
    print(Core.spotify_credential.get_token_acces(request.args.get('code')))
    return jsonify('test')

start(app, '0.0.0.0', 8080)     