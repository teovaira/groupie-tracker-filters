// search.test.js — unit tests for debounce and renderCards.
// Run with: node web/static/js/search.test.js

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

function assertContains(str, substr, msg) {
    if (!str.includes(substr)) {
        throw new Error(msg || `expected string to contain ${JSON.stringify(substr)}\ngot: ${str}`);
    }
}

// ---------------------------------------------------------------------------
// Load the module under test.
// search.js must export { debounce, renderCards } when run under Node.
// ---------------------------------------------------------------------------

const { debounce, renderCards } = require('./search.js');

// ---------------------------------------------------------------------------
// debounce tests
// ---------------------------------------------------------------------------

console.log('\ndebounce');

test('calls fn after delay', (done) => {
    let callCount = 0;
    const fn = () => { callCount++; };
    const debounced = debounce(fn, 50);

    debounced();
    assertEqual(callCount, 0, 'should not call immediately');

    setTimeout(() => {
        assertEqual(callCount, 1, 'should call once after delay');
    }, 100);
});

test('resets timer on repeated calls', () => {
    let callCount = 0;
    const fn = () => { callCount++; };
    const debounced = debounce(fn, 50);

    debounced();
    debounced();
    debounced();

    assertEqual(callCount, 0, 'should not have called fn yet');

    setTimeout(() => {
        assertEqual(callCount, 1, 'should call fn exactly once after burst');
    }, 100);
});

test('passes arguments to fn', () => {
    let received = null;
    const fn = (val) => { received = val; };
    const debounced = debounce(fn, 0);
    debounced('hello');
    setTimeout(() => {
        assertEqual(received, 'hello', 'should pass argument through');
    }, 10);
});

// ---------------------------------------------------------------------------
// renderCards tests
// ---------------------------------------------------------------------------

console.log('\nrenderCards');

test('renders one card with correct content', () => {
    const artists = [
        { id: 1, name: 'Queen', image: 'http://img/queen.jpg', creationDate: 1970 }
    ];
    const html = renderCards(artists);
    assertContains(html, '/artist/1');
    assertContains(html, 'Queen');
    assertContains(html, 'http://img/queen.jpg');
    assertContains(html, '1970');
});

test('renders multiple cards', () => {
    const artists = [
        { id: 1, name: 'Queen', image: 'http://img/1.jpg', creationDate: 1970 },
        { id: 2, name: 'Foo Fighters', image: 'http://img/2.jpg', creationDate: 1994 },
    ];
    const html = renderCards(artists);
    assertContains(html, '/artist/1');
    assertContains(html, '/artist/2');
    assertContains(html, 'Queen');
    assertContains(html, 'Foo Fighters');
});

test('returns empty string for empty array', () => {
    const html = renderCards([]);
    assertEqual(html, '', 'empty array should produce empty string');
});

test('escapes special characters in name', () => {
    const artists = [
        { id: 3, name: 'AC/DC', image: 'http://img/acdc.jpg', creationDate: 1973 }
    ];
    const html = renderCards(artists);
    assertContains(html, '/artist/3');
    // name appears in alt and h2 — both must be present
    assertContains(html, 'AC/DC');
});

// ---------------------------------------------------------------------------
// Summary
// ---------------------------------------------------------------------------

setTimeout(() => {
    console.log(`\n${passed + failed} tests: ${passed} passed, ${failed} failed`);
    if (failed > 0) process.exit(1);
}, 200);
