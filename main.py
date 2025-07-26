from flask import Flask, render_template, request, jsonify, redirect
from api import web_radio, spot_credential
from lib.ffmplay import Stream_station
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

radio = web_radio.WEB_radio()
play_radio = Stream_station()
sp_pkce = spot_credential.Spotify_PKCE()

def start(app: Flask, ip, port):
    return app.run(ip, port)

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
@app.route('/', methods=['GET'])
def page_web():
    return render_template("index.html")

@app.route('/radio', methods=['GET'])
def radio_list():
    line_sel = request.args.get('line_sel')
    num_sel = request.args.get('num_sel')

    return jsonify(radio.get_radio(line_sel=int(line_sel), num_sel=int(num_sel)))

@app.route('/search', methods=['GET'])
def search():
    query = request.args.get('q') 
    return radio.search(query=query)

@app.route('/play_radio', methods=['PUT'])
def play_audio():
    try:
        stream_url = request.args.get('stream_url')
        name_station = request.args.get('name_station')
        
        play_radio.play(stream_url, name_station)
        
        return jsonify({"response": "True"})
    
    except Exception:
        return jsonify({"response": "False"})
    
@app.route('/pause_resume', methods=['PUT'])
def pause_resume():
    return jsonify(play_radio.pause_resume())

@app.route('/control_menu', methods=['GET'])
def control_menu():
    if play_radio.name != "":
        return {"data": radio.search(query=play_radio.name), "status": True if play_radio.process else False}
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
    return redirect(sp_pkce.gen_url())

@app.route('/login')
def login():
    print('---------------------------------------------')
    print(request.args.get('code'))
    print('---------------------------------------------')
    print(sp_pkce.get_token_acces(request.args.get('code')))
    return jsonify('test')


start(app, '0.0.0.0', 8080)     