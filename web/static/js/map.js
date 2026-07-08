'use strict';
function initMap(){
    const mapElement = document.getElementById('map');
    if (!mapElement) return;

    const markersElement = document.getElementById('artist-markers');
    if (!markersElement) return;

    const markers = JSON.parse(markersElement.textContent);

    const map = L.map(mapElement).setView([20, 0], 2);

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        maxZoom: 19,
        attribution: '&copy; OpenStreetMap contributors'
    }).addTo(map);

    if (!markers || markers.length === 0) return;

    const bounds = [];

    markers.forEach(function (marker) {
        const position = [marker.Lat, marker.Lng];
        L.marker(position).addTo(map).bindPopup(marker.Name);
        bounds.push(position);
    });

    if (bounds.length > 1) {
        map.fitBounds(bounds, { padding: [30, 30] });
    } else {
        map.setView(bounds[0], 8);
    }
}

document.addEventListener('DOMContentLoaded', initMap);
