'use strict';

/**
 * buildFilterQuery converts a plain filter-state object into a URL query
 * string for GET /api/filter. Every field is optional; absent, blank, or
 * whitespace-only values are omitted entirely rather than sent as empty
 * parameters, so the server only ever sees bounds it should actually apply.
 *
 * @param {{
 *   q?: string,
 *   creationMin?: string, creationMax?: string,
 *   firstAlbumMin?: string, firstAlbumMax?: string,
 *   membersMin?: string, membersMax?: string,
 *   locations?: string[]
 * }} state
 * @returns {string} query string with no leading '?', or '' if state has no active filters
 */
function buildFilterQuery(state) {
    const parts = [];

    function addIfPresent(name, value) {
        if (value !== undefined && value !== null && String(value).trim() !== '') {
            parts.push(name + '=' + encodeURIComponent(value));
        }
    }

    addIfPresent('q', state.q);
    addIfPresent('creation_min', state.creationMin);
    addIfPresent('creation_max', state.creationMax);
    addIfPresent('first_album_min', state.firstAlbumMin);
    addIfPresent('first_album_max', state.firstAlbumMax);
    addIfPresent('members_min', state.membersMin);
    addIfPresent('members_max', state.membersMax);

    if (state.locations && state.locations.length > 0) {
        for (let i = 0; i < state.locations.length; i++) {
            parts.push('locations=' + encodeURIComponent(state.locations[i]));
        }
    }

    return parts.join('&');
}

// Expose pure functions for filter.test.js to import under Node, and
// globally in the browser where filter.test.js runs them directly.
if (typeof module !== 'undefined' && module.exports) {
    module.exports = { buildFilterQuery };
} else {
    window.buildFilterQuery = buildFilterQuery;
}

// init wires the search box and filter panel inputs to the combined
// /api/filter endpoint. It replaces search.js's own request handling — the
// search box's input event is read here too, so search and filters combine
// into a single request rather than each overwriting the other's results.
function init() {
    const searchInput    = document.getElementById('search-input');
    const loading        = document.getElementById('loading');
    const noResults      = document.getElementById('no-results');
    const results        = document.getElementById('search-results');
    const resetButton    = document.getElementById('filter-reset');
    const creationMin    = document.getElementById('filter-creation-min');
    const creationMax    = document.getElementById('filter-creation-max');
    const firstAlbumMin  = document.getElementById('filter-first-album-min');
    const firstAlbumMax  = document.getElementById('filter-first-album-max');
    const membersMin     = document.getElementById('filter-members-min');
    const membersMax     = document.getElementById('filter-members-max');
    const locationsRoot  = document.getElementById('filter-locations');

    if (!searchInput || !loading || !noResults || !results) {
        return;
    }

    // Save the original server-rendered cards to restore when every filter is cleared.
    const originalHTML = results.innerHTML;

    function showLoading()    { loading.classList.remove('hidden'); }
    function hideLoading()    { loading.classList.add('hidden'); }
    function showNoResults()  { noResults.classList.remove('hidden'); }
    function hideNoResults()  { noResults.classList.add('hidden'); }

    function checkedLocations() {
        if (!locationsRoot) {
            return [];
        }
        const boxes = locationsRoot.querySelectorAll('input[type="checkbox"]:checked');
        const values = [];
        for (let i = 0; i < boxes.length; i++) {
            values.push(boxes[i].value);
        }
        return values;
    }

    function currentState() {
        return {
            q: searchInput.value,
            creationMin: creationMin ? creationMin.value : '',
            creationMax: creationMax ? creationMax.value : '',
            firstAlbumMin: firstAlbumMin ? firstAlbumMin.value : '',
            firstAlbumMax: firstAlbumMax ? firstAlbumMax.value : '',
            membersMin: membersMin ? membersMin.value : '',
            membersMax: membersMax ? membersMax.value : '',
            locations: checkedLocations(),
        };
    }

    function handleChange() {
        const query = buildFilterQuery(currentState());

        // No active filters or search term — restore the original page load.
        if (query === '') {
            hideLoading();
            hideNoResults();
            results.innerHTML = originalHTML;
            return;
        }

        showLoading();
        hideNoResults();

        fetch('/api/filter?' + query)
            .then(function (res) { return res.json(); })
            .then(function (artists) {
                hideLoading();
                if (!artists || artists.length === 0) {
                    results.innerHTML = '';
                    showNoResults();
                } else {
                    hideNoResults();
                    results.innerHTML = renderCards(artists);
                }
            })
            .catch(function () {
                hideLoading();
            });
    }

    const debouncedChange = debounce(handleChange, 300);

    searchInput.addEventListener('input', debouncedChange);
    [creationMin, creationMax, firstAlbumMin, firstAlbumMax, membersMin, membersMax].forEach(function (el) {
        if (el) {
            el.addEventListener('input', debouncedChange);
        }
    });
    if (locationsRoot) {
        locationsRoot.addEventListener('change', debouncedChange);
    }

    if (resetButton) {
        resetButton.addEventListener('click', function () {
            searchInput.value = '';
            [creationMin, creationMax, firstAlbumMin, firstAlbumMax, membersMin, membersMax].forEach(function (el) {
                if (el) {
                    el.value = '';
                }
            });
            if (locationsRoot) {
                const boxes = locationsRoot.querySelectorAll('input[type="checkbox"]:checked');
                for (let i = 0; i < boxes.length; i++) {
                    boxes[i].checked = false;
                }
            }
            hideLoading();
            hideNoResults();
            results.innerHTML = originalHTML;
        });
    }
}

if (typeof document !== 'undefined') {
    document.addEventListener('DOMContentLoaded', init);
}
