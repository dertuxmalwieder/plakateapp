# Plakate-App

Dies ist eine Webanwendung, die eine anklickbare Straßenkarte anzeigt, über die (zum Beispiel) die Position von Plakaten im Wahlkampf bestimmt werden kann. Eure Position wird automatisch ermittelt, sofern möglich.

## Hä? Haste mal ein Beispiel?

Klar. Nehmen wir an, ihr wollt für irgendeinen Verein eure Stadt vollplakatieren, wollt aber den Überblick behalten, wo überall Plakate hängen, um sie später wieder entfernen zu können. Genau diesem Zweck dient diese Anwendung.

## Technik

Ihr braucht auf eurem Server nur Fossil und Go, alles Weitere passiert automatisch:

    % fossil clone https://code.rosaelefanten.org/plakateapp plakateapp.fossil ; fossil open plakateapp.fossil
    % cd src
    % go build

Falls ihr Fossil nicht mögt: Es gibt auch einen [GitHub-Mirror](https://github.com/dertuxmalwieder/plakateapp).

Das Ausführen ist dann leicht:

    % ./plakateapp

Unter unixoiden Systemen legt die Plakateapp eine PID-Datei an, läuft also im Hintergrund. Zum Beenden sollte der `kill`-Befehl funktionieren:

    % kill `cat plakateapp.pid`

Unter Windows ist das zurzeit noch nicht vorgesehen, die Plakateapp läuft immer im Vordergrund.

Die Karte ist anschließend über den Port 6090 (einstellbar direkt in der Datei `plakateapp.go`) erreichbar. Unter `euerserver:6090/manageplakate` gibt es auch eine einfache Liste aller eingetragenen Plakate zum schnellen Löschen. Der Großteil des UIs wurde mit [Leafjet.js](http://leafletjs.com/) programmiert.

## Urheberrecht? Quatsch.

Die Plakateapp wurde ursprünglich für den Kommunalwahlkampf in Niedersachsen 2016 von [@tux0r](https://twitter.com/tux0r) hektisch (also eher zweckmäßig als gut) für die Piratenpartei Braunschweig programmiert (weshalb der Standard für die Position mitten in Braunschweig liegt, aber das könnt ihr im Javascript-Code ändern). 2019 wurde sie in Go neu implementiert. Seit Ende 2020 steht die Plakateapp unter der [CDDL](LICENSE), somit ist sie endlich auch offiziell Freie Software. Ihr dürft also damit wegrennen und sie teuer verscherbeln. Viel Spaß!
