// map.test.js — unit tests for initMap.
// Run with: node web/static/js/map.test.js

'use strict';

// ---------------------------------------------------------------------------
// Minimal test harness — no external dependencies
// ---------------------------------------------------------------------------

let passed = 0;
let failed = 0;

function test(name, fn) {
    try {
        fn();
        console.log(`  PASS  ${name}`);
        passed++;
    } catch (e) {
        console.log(`  FAIL  ${name}`);
        console.log(`        ${e.message}`);
        failed++;
    }
}

function assertEqual(actual, expected, msg) {
    if (actual !== expected) {
        throw new Error(msg || `expected ${JSON.stringify(expected)}, got ${JSON.stringify(actual)}`);
    }
}

// ---------------------------------------------------------------------------
// DOM stub
// ---------------------------------------------------------------------------

function makeDOM(markersJSON) {
    const elements = {};

    elements['map'] = { _id: 'map' };
    elements['artist-markers'] = {
        _id: 'artist-markers',
        textContent: markersJSON,
    };

    global.document = {
        getElementById: (id) => elements[id] || null,
        addEventListener: () => {},
    };
}

// ---------------------------------------------------------------------------
// Leaflet stub — records calls so tests can assert on them
// ---------------------------------------------------------------------------

function makeLeaflet() {
    const calls = {
        mapCreated: false,
        setView: [],
        tileLayer: false,
        markers: [],
        fitBounds: null,
    };

    const fakeMap = {
        setView(center, zoom) { calls.setView.push({ center, zoom }); return fakeMap; },
        fitBounds(bounds, opts) { calls.fitBounds = { bounds, opts }; return fakeMap; },
    };

    const fakeMarker = {
        addTo() { return fakeMarker; },
        bindPopup(name) { calls.markers.push(name); return fakeMarker; },
    };

    const fakeTileLayer = { addTo() {} };

    global.L = {
        map(el) { calls.mapCreated = true; return fakeMap; },
        tileLayer(url, opts) { calls.tileLayer = true; return fakeTileLayer; },
        marker(pos) { return fakeMarker; },
    };

    return calls;
}

// ---------------------------------------------------------------------------
// Load module under test — strips the DOMContentLoaded listener
// ---------------------------------------------------------------------------

function loadInitMap() {
    // Re-evaluate map.js source each test so state is fresh
    const fs = require('fs');
    const src = fs.readFileSync('./web/static/js/map.js', 'utf8')
        .replace("document.addEventListener('DOMContentLoaded', initMap);", '');
    const fn = new Function(src + '\nreturn initMap;');
    return fn();
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

console.log('\ninitMap');

test('does nothing when #map element is missing', () => {
    global.document = { getElementById: () => null, addEventListener: () => {} };
    global.L = { map() { throw new Error('L.map should not be called'); } };
    const initMap = loadInitMap();
    initMap(); // should not throw
});

test('does nothing when artist-markers element is missing', () => {
    global.document = {
        getElementById: (id) => id === 'map' ? { _id: 'map' } : null,
        addEventListener: () => {},
    };
    global.L = { map() { throw new Error('L.map should not be called'); } };
    const initMap = loadInitMap();
    initMap(); // should not throw
});

test('initialises map with default world view', () => {
    makeDOM('[]');
    const calls = makeLeaflet();
    const initMap = loadInitMap();
    initMap();
    assertEqual(calls.mapCreated, true, 'L.map should be called');
    assertEqual(calls.setView.length, 1, 'setView should be called once');
    assertEqual(calls.setView[0].center[0], 20, 'default lat should be 20');
    assertEqual(calls.setView[0].center[1], 0, 'default lng should be 0');
    assertEqual(calls.setView[0].zoom, 2, 'default zoom should be 2');
});

test('adds tile layer', () => {
    makeDOM('[]');
    const calls = makeLeaflet();
    const initMap = loadInitMap();
    initMap();
    assertEqual(calls.tileLayer, true, 'tileLayer should be added');
});

test('places a marker for each location', () => {
    const markers = [
        { Name: 'london-uk', Lat: 51.5074, Lng: -0.1278 },
        { Name: 'paris-france', Lat: 48.8566, Lng: 2.3522 },
    ];
    makeDOM(JSON.stringify(markers));
    const calls = makeLeaflet();
    const initMap = loadInitMap();
    initMap();
    assertEqual(calls.markers.length, 2, 'should place 2 markers');
    assertEqual(calls.markers[0], 'london-uk');
    assertEqual(calls.markers[1], 'paris-france');
});

test('calls fitBounds when more than one marker', () => {
    const markers = [
        { Name: 'london-uk', Lat: 51.5074, Lng: -0.1278 },
        { Name: 'paris-france', Lat: 48.8566, Lng: 2.3522 },
    ];
    makeDOM(JSON.stringify(markers));
    const calls = makeLeaflet();
    const initMap = loadInitMap();
    initMap();
    if (!calls.fitBounds) throw new Error('fitBounds should be called for multiple markers');
    assertEqual(calls.fitBounds.bounds.length, 2, 'fitBounds should receive 2 positions');
});

test('calls setView with zoom 8 for a single marker', () => {
    const markers = [{ Name: 'london-uk', Lat: 51.5074, Lng: -0.1278 }];
    makeDOM(JSON.stringify(markers));
    const calls = makeLeaflet();
    const initMap = loadInitMap();
    initMap();
    const zoomedIn = calls.setView.find(c => c.zoom === 8);
    if (!zoomedIn) throw new Error('setView with zoom 8 should be called for single marker');
});

test('returns early without placing markers for empty array', () => {
    makeDOM('[]');
    const calls = makeLeaflet();
    const initMap = loadInitMap();
    initMap();
    assertEqual(calls.markers.length, 0, 'no markers should be placed for empty array');
    assertEqual(calls.fitBounds, null, 'fitBounds should not be called');
});

// ---------------------------------------------------------------------------
// Summary
// ---------------------------------------------------------------------------

console.log(`\n${passed + failed} tests: ${passed} passed, ${failed} failed`);
if (failed > 0) process.exit(1);
