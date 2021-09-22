/*
 * The contents of this file are subject to the terms of the
 * Common Development and Distribution License, Version 1.0 only
 * (the "License").  You may not use this file except in compliance
 * with the License.
 *
 * See the file LICENSE in this distribution for details.
 * A copy of the CDDL is also available via the Internet at
 * http://www.opensource.org/licenses/cddl1.txt
 *
 * When distributing Covered Code, include this CDDL HEADER in each
 * file and include the contents of the LICENSE file from this
 * distribution.
 */

var map;
var umkreis;

function deleteconfirm_list(id) {
    // Löscht das Plakat <id> erst nach Bestätigung.
    // Leitet danach zur Plakatliste zurück.
    if (confirm("Möchten Sie dieses Plakat wirklich löschen?")) {
        location.href = "/del/" + id;
    }
}

function initmap() {
    L.Map.include({
        // Funktion zum Ermitteln einzelner Marker per ID.
        getMarkerById: function (id) {
            let marker = null;
            this.eachLayer(function (layer) {
                if (layer instanceof L.Marker) {
                    if (layer.options.id === id) {
                        marker = layer;
                    }
                }
            });
            return marker;
        }
    });

    // Karte initialisieren
    map = new L.Map('map', { tap: false });

    const osmUrl='//{s}.tile.openstreetmap.org/{z}/{x}/{y}.png';
    const osmAttrib='Karte von <a href="https://openstreetmap.org">OpenStreetMap</a>';
    const osm = new L.TileLayer(
        osmUrl, {
            minZoom: 12,
            maxZoom: 28,
            maxNativeZoom: 19,
            attribution: osmAttrib,
            useCache: true,
            crossOrigin: true }
    );

    // OSM anzeigen
    map.addLayer(osm);

    // bei Ortungserfolg Standort anzeigen:
    map.on('locationfound', function(e) {
        const radius = e.accuracy / 2;

        // Haben wir schon einen Marker? Dann verschieben, sonst hinzufügen:
        let posmarker = map.getMarkerById(1);
        if (posmarker == null) {
            let standortMarkerIcon = new L.Icon({
                iconUrl: "/static/images/marker-icon-green.png",
                shadowUrl: "/static/images/marker-shadow.png",
                iconSize: [25, 41],
                iconAnchor: [12, 41],
                popupAnchor: [1, -34],
                shadowSize: [41, 41]
            });

            L.marker(e.latlng, { icon: standortMarkerIcon, id: 1 })
                .addTo(map)
                .bindPopup("Du bist ungefähr hier.")
                .openPopup();
        }
        else {
            map.removeLayer(umkreis);
            posmarker.setLatLng(e.latlng);
        }

        umkreis = L.circle(e.latlng, { radius: radius }).addTo(map);
    });

    // bei Ortungsmisserfolg Karte mitten in Braunschweig setzen:
    map.on('locationerror', function() {
        map.setView(new L.LatLng(52.269167, 10.521111), 17);
    });

    // bestehende Marker laden:
    $.post(
        "/listplakate",
        {},
        function(data) {
            const json = JSON.parse(data);
            for (let i = 0; i < json.length; i++) {
                let plakat = json[i];
                let plakatlatlng = new L.LatLng(plakat.Latitude,plakat.Longitude);
                let marker = new L.Marker(plakatlatlng, {draggable:false})
                    .bindPopup("<input type='button' value='Plakat löschen' data-id='"+plakat.ID+"' class='marker-delete-button'/>");

                marker.on("popupopen", onPopupOpen);
                map.addLayer(marker);
            }
        }
    );

    // neue Marker setzen:
    map.on('click', function(e) {
        askForNewPlakat(e.latlng);
        return false;
    });

    // Ortung versuchen:
    map.locate({ setView: true });

    // Ort-Verfolgen-Button (ein/aus; Standard: aus):
    $("#tracebtn").on("click", function() {
        if ($(this).hasClass("strike")) {
            // Verfolgen ausmachen.
            map.stopLocate();
            map.locate({ setView: true, maxZoom: map.getZoom() });
        }
        else {
            // Verfolgen anmachen.
            map.stopLocate();
            map.locate({ setView: true, watch: true, maxZoom: map.getZoom() });
        }

        $(this).toggleClass("strike"); 
    });

    // Neues-Plakat-Button:
    $("#newplakatbtn").on("click", function() {
        // Position ermitteln:
        let posmarker = map.getMarkerById(1);
        if (posmarker != null) {
            // Normalerweise sollte dieser Marker existieren.
            // Ein else-Fall ergibt vermutlich keinen Sinn hier.
            askForNewPlakat(posmarker.getLatLng());
        }
        return false;
    });
}

function askForNewPlakat(latlng) {
    // Abfrage für ein neues Plakat an Position <latlng>.
    if (confirm("Möchtest du hier ein neues Plakat melden?")) {
        $.post(
            "/neuesplakat",
            {
                lat: latlng.lat,
                lon: latlng.lng
            },
            function(data) {
                Toastify({
                    text: data.text
                }).showToast();

                let marker = new L.Marker(latlng, {draggable:false})
                    .bindPopup("<input type='button' value='Plakat löschen' data-id='"+data.id+"' class='marker-delete-button'/>");
                marker.on("popupopen", onPopupOpen);
                map.addLayer(marker);
            }
        );
    }
}

function onPopupOpen() {
    let tempMarker = this;

    $(".marker-delete-button:visible").click(function () {
        $.post(
            "/delpost",
            { id: $(this).attr("data-id") },
            function(_data) {
                // noop
            }
        );
        map.removeLayer(tempMarker);
    });
}

if (window.jQuery) {
    $(document).ready(function() {
        initmap();
    });
}
