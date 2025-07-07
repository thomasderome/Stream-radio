from flask import Flask, render_template, request
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
    number = request.args.get('num')
    return radio.get_radio(number=int(number))

@app.route('/search', methods=['GET'])
def search():
    query = request.args.get('q') 
    return radio.search(query=query)

start(app, '0.0.0.0', 80)