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
var ajaxRequest;
var plotlist;
var plotlayers=[];

function initmap() {
    // Karte initialisieren
    map = new L.Map('map');

    var osmUrl='//{s}.tile.openstreetmap.org/{z}/{x}/{y}.png';
    var osmAttrib='Karte von <a href="https://openstreetmap.org">OpenStreetMap</a>';
    var osm = new L.TileLayer(osmUrl, {minZoom: 12, maxZoom: 28, attribution: osmAttrib});

    // OSM anzeigen
    map.addLayer(osm);

    // bei Ortungserfolg Standort anzeigen:
    map.on('locationfound', function(e) {
        var radius = e.accuracy / 2;

        L.marker(e.latlng).addTo(map)
            .bindPopup("Du bist ungefähr hier.").openPopup();

        L.circle(e.latlng, radius).addTo(map);
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
          var json = JSON.parse(data);
          for (var i = 0; i < json.length; i++) {
              var plakat = json[i];
              var plakatlatlng = new L.LatLng(plakat.Latitude,plakat.Longitude)
              var marker = new L.Marker(plakatlatlng, {draggable:false})
.bindPopup("<input type='button' value='Plakat löschen' data-id='"+plakat.ID+"' class='marker-delete-button'/>");

              marker.on("popupopen", onPopupOpen);
              map.addLayer(marker);
          }
      }
    );

    // neue Marker setzen:
    map.on('click', function(e) {
        if (confirm("Möchtest du hier ein neues Plakat melden?")) {
            var marker = new L.Marker(e.latlng, {draggable:false});
            marker.on("popupopen", onPopupOpen);
            map.addLayer(marker);

            $.post(
                "/neuesplakat",
                {
                    lat: e.latlng.lat,
                    lon: e.latlng.lng
                },
                function(data) {
                    alert(data);
                }
            );
        }

        return false;
    });

    // Ortung versuchen:
    map.locate({ setView: true, maxZoom: 28 });
}

function onPopupOpen() {
    var tempMarker = this;

    $(".marker-delete-button:visible").click(function () {
        $.post(
          "/delpost",
          { id: $(this).attr("data-id") },
          function(data) {
              // noop
          }
        );
        map.removeLayer(tempMarker);
    });
}

$(document).ready(function() {
    initmap();
});
