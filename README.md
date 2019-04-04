# Plakate-App

Dies ist eine Webanwendung, die eine anklickbare Straßenkarte anzeigt, über die (zum Beispiel) die Position von Plakaten im Wahlkampf bestimmt werden kann.

## Hä? Haste mal ein Beispiel?

Klar. Nehmen wir an, ihr wollt für irgendeinen Verein eure Stadt vollplakatieren, wollt aber den Überblick behalten, wo überall Plakate hängen, um sie später wieder entfernen zu können. Genau diesem Zweck dient diese Anwendung.

## Technik

Ihr braucht auf eurem Server nur Go und ein paar Abhängigkeiten, alles Weitere passiert automatisch:

    git clone https://github.com/dertuxmalwieder/plakateapp ; cd plakateapp/src
    go get github.com/mattn/go-sqlite3
    go get github.com/jmoiron/sqlx
    go get github.com/gorilla/mux
    go build ./plakateapp.go
    ./plakateapp

Die Karte ist anschließend über den Port 6090 (einstellbar direkt in der Datei `plakateapp.go`) erreichbar. Unter `euerserver:6090/manageplakate` gibt es auch eine einfache Liste aller eingetragenen Plakate zum schnellen Löschen. Der Großteil des UIs wurde mit [Leafjet.js](http://leafletjs.com/) programmiert.

## Urheberrecht? Quatsch.

Die Plakateapp wurde ursprünglich für den Kommunalwahlkampf in Niedersachsen 2016 von [@tux0r](https://twitter.com/tux0r) hektisch (also eher zweckmäßig als gut) für die Piratenpartei Braunschweig programmiert (weshalb der Standard für die Position mitten in Braunschweig liegt, aber das könnt ihr im Javascript-Code ändern). 2019 wurde sie in Go neu implementiert. All dies hier steht unter der [WTFPL v2](http://www.wtfpl.net/txt/copying/), ihr dürft also gern damit wegrennen und es teuer verscherbeln. Viel Spaß!
