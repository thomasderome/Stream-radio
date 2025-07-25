from flask import Flask, render_template, request, jsonify
from api.web_radio import WEB_radio
from lib.ffmplay import Stream_station
import os

template_dir = os.path.join(os.path.dirname(__file__), 'web', 'templates')
static_dir = os.path.join(os.path.dirname(__file__), 'web', 'static')
app = Flask("Radio Garage", template_folder=template_dir, static_folder=static_dir)

radio = WEB_radio()
play_radio = Stream_station()

def start(app: Flask, ip, port):
    return app.run(ip, port)

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
    return play_radio.pause_resume()

start(app, '0.0.0.0', 80)