import sqlite3
import json
from flask import Flask, render_template, request, send_from_directory, redirect
from gevent.wsgi import WSGIServer
from gevent import monkey

monkey.patch_all()

app = Flask(__name__)

PORT = 6090

@app.route('/')
def hauptseite():
    return render_template("index.htm")


@app.route("/manageplakate")
def manage_plakate():
    try:
        conn = sqlite3.connect("plakate.db")
        with conn:
            c = conn.cursor()
            c.execute("SELECT * FROM plakate")
            plakate = [dict(id=row[0], lat=row[1], lon=row[2]) for row in c.fetchall()]
            return render_template("delete.htm",plakate=plakate)
    except:
        return "err"


@app.route("/listplakate", methods=['POST'])
def list_plakate():
    try:
        conn = sqlite3.connect("plakate.db")
        with conn:
            c = conn.cursor()
            c.execute("SELECT * FROM plakate")
            plakate = [dict(id=row[0], lat=row[1], lon=row[2]) for row in c.fetchall()]
            return json.dumps(plakate)
    except:
        return "err"


@app.route("/neuesplakat", methods=['POST'])
def neues_plakat():
    try:
        latitude = request.form["lat"]
        longitude = request.form["lon"]

        conn = sqlite3.connect("plakate.db", detect_types=sqlite3.PARSE_DECLTYPES|sqlite3.PARSE_COLNAMES)
        with conn:
            c = conn.cursor()
            c.execute("INSERT INTO plakate (lat, lon) VALUES ({}, {});".format(float(latitude), float(longitude)))
            conn.commit()

        return "Plakat erfolgreich eingetragen!"
    except sqlite3.Error as e:
        return "Fehler! :-( {}".format(e.message)


@app.route("/del/<int:plakatid>", strict_slashes = False)
def del_plakat(plakatid):
    try:
        conn = sqlite3.connect("plakate.db", detect_types=sqlite3.PARSE_DECLTYPES|sqlite3.PARSE_COLNAMES)
        with conn:
            c = conn.cursor()
            c.execute("DELETE FROM plakate WHERE id={}".format(plakatid))
            conn.commit()
        return redirect("/manageplakate")
    except sqlite3.Error as e:
        return "Fehler! :-( {}".format(e.message)


@app.route("/delpost", methods=['POST'])
def delpost():
    try:
        plakatid = request.form["id"]
        conn = sqlite3.connect("plakate.db", detect_types=sqlite3.PARSE_DECLTYPES|sqlite3.PARSE_COLNAMES)
        with conn:
            c = conn.cursor()
            c.execute("DELETE FROM plakate WHERE id={}".format(plakatid))
            conn.commit()
        return "Geklappt."
    except sqlite3.Error as e:
        return "Fehler! :-( {}".format(e.message)


@app.route('/templates/<path:path>', strict_slashes = False)
def web_static(path):
    # Statische Templatedateien ausliefern
    return send_from_directory('templates', path)

http_server = WSGIServer(('0.0.0.0', PORT), app)
try:
    http_server.serve_forever()
except KeyboardInterrupt:
    print("Beende die Plakate-App...")
    if http_server.started:
        http_server.stop()