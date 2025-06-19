#
# Copyright (c) 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

import jsonpickle
import os
import quart
import requests
from quart_cors import cors
from typing import List
from coherence import NamedMap, Session, Options
from dataclasses import dataclass
from coherence.serialization import proxy

from quart import Quart, request, redirect, jsonify, send_from_directory

WEATHER_API_URL_TEMPLATE = "https://wttr.in/{city}?format=j1"
CACHE_TTL = 30000  # 30 seconds

# ---- init ------------

app: Quart = cors(
    Quart(__name__, static_url_path='', static_folder='./'),
    allow_origin="*"
)

# the Session with the gRPC proxy
session: Session

weather: NamedMap[str, any]


@app.before_serving
async def init():

    # initialize the session using the default localhost:1408 or the value of COHERENCE_SERVER_ADDRESS
    global session
    session = await Session.create()

    global weather
    weather = await session.get_map('weather')

def fetch_weather(city):
    """Fetch full weather data from wttr.in JSON API"""
    url = WEATHER_API_URL_TEMPLATE.format(city=city)
    response = requests.get(url)
    response.raise_for_status()
    return response.json()

# ----- routes --------------------------------------------------------------
@app.route("/")
async def index():
    return await send_from_directory(os.path.join(app.static_folder), "index.html")

@app.route("/weather/")
async def weather():
    city = request.args.get("city")
    if not city:
        return jsonify({"error": "Missing 'city' query parameter"}), 40

    # Optional TTL override (in seconds)
    ttl_param = request.args.get("ttl")
    try:
        ttl_seconds = int(ttl_param) if ttl_param is not None else 30  # default 30s
    except ValueError:
        return jsonify({"error": "Invalid 'ttl' parameter; must be an integer"}), 400

    ttl_ms = ttl_seconds * 1000

    key = city.lower()
    cache_value = await weather.get(key)
    source = "Cache"

    if cache_value == None:
        try:
           cache_value = fetch_weather(city)
           source = "API Call"
        except Exception as e:
            return jsonify({"error": f"Failed to fetch weather: {str(e)}"}), 500
        await weather.put(key, cache_value, ttl_ms)

    return jsonify({"source": source, "data": cache_value})

# ----- main ----------------------------------------------------------------

if __name__ == '__main__':
    # run the application on port 8080
    app.run(host='0.0.0.0', port=8080)