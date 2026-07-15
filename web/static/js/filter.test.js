// filter.test.js — unit tests for buildFilterQuery.
// Run with: node web/static/js/filter.test.js

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
// Load the module under test.
// filter.js must export { buildFilterQuery } when run under Node.
// ---------------------------------------------------------------------------

const { buildFilterQuery } = require('./filter.js');

// ---------------------------------------------------------------------------
// buildFilterQuery tests
// ---------------------------------------------------------------------------

console.log('\nbuildFilterQuery');

test('empty state produces empty query string', () => {
    const qs = buildFilterQuery({});
    assertEqual(qs, '');
});

test('includes q when present', () => {
    const qs = buildFilterQuery({ q: 'queen' });
    assertEqual(qs, 'q=queen');
});

test('trims and omits empty q', () => {
    const qs = buildFilterQuery({ q: '   ' });
    assertEqual(qs, '');
});

test('includes creation date bounds', () => {
    const qs = buildFilterQuery({ creationMin: '1990', creationMax: '2000' });
    assertEqual(qs, 'creation_min=1990&creation_max=2000');
});

test('includes first album bounds', () => {
    const qs = buildFilterQuery({ firstAlbumMin: '1980', firstAlbumMax: '1990' });
    assertEqual(qs, 'first_album_min=1980&first_album_max=1990');
});

test('includes members bounds', () => {
    const qs = buildFilterQuery({ membersMin: '1', membersMax: '4' });
    assertEqual(qs, 'members_min=1&members_max=4');
});

test('omits blank numeric fields', () => {
    const qs = buildFilterQuery({ creationMin: '', creationMax: '2000' });
    assertEqual(qs, 'creation_max=2000');
});

test('includes one locations param per selected location', () => {
    const qs = buildFilterQuery({ locations: ['texas-usa', 'washington-usa'] });
    assertEqual(qs, 'locations=texas-usa&locations=washington-usa');
});

test('url-encodes location values', () => {
    const qs = buildFilterQuery({ locations: ['san francisco-usa'] });
    assertEqual(qs, 'locations=san%20francisco-usa');
});

test('omits empty locations array', () => {
    const qs = buildFilterQuery({ locations: [] });
    assertEqual(qs, '');
});

test('combines all fields in a stable order', () => {
    const qs = buildFilterQuery({
        q: 'queen',
        creationMin: '1970', creationMax: '1980',
        firstAlbumMin: '1973', firstAlbumMax: '1980',
        membersMin: '4', membersMax: '4',
        locations: ['london-uk'],
    });
    assertEqual(
        qs,
        'q=queen&creation_min=1970&creation_max=1980&first_album_min=1973&first_album_max=1980&members_min=4&members_max=4&locations=london-uk'
    );
});

// ---------------------------------------------------------------------------
// Summary
// ---------------------------------------------------------------------------

setTimeout(() => {
    console.log(`\n${passed + failed} tests: ${passed} passed, ${failed} failed`);
    if (failed > 0) process.exit(1);
}, 200);
