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

/**
 * sliderRangeState converts a dual-handle slider's current thumb positions
 * into the min/max strings to send to /api/filter. A thumb resting at its
 * data edge (min thumb at dataMin, max thumb at dataMax) imposes no real
 * constraint, so that side is returned as '' and omitted from the query —
 * leaving untouched sliders equivalent to "no filter".
 *
 * @param {number} minVal - current position of the min thumb
 * @param {number} maxVal - current position of the max thumb
 * @param {number} dataMin - the slider's lowest possible value
 * @param {number} dataMax - the slider's highest possible value
 * @returns {{min: string, max: string}} values to send, '' where at the edge
 */
function sliderRangeState(minVal, maxVal, dataMin, dataMax) {
    return {
        min: minVal > dataMin ? String(minVal) : '',
        max: maxVal < dataMax ? String(maxVal) : '',
    };
}

// Expose pure functions for filter.test.js to import under Node, and
// globally in the browser where filter.test.js runs them directly.
if (typeof module !== 'undefined' && module.exports) {
    module.exports = { buildFilterQuery, sliderRangeState };
} else {
    window.buildFilterQuery = buildFilterQuery;
    window.sliderRangeState = sliderRangeState;
}

// rangeSlider wraps a single dual-handle range fieldset, exposing helpers to
// keep the two thumbs from crossing, paint the filled segment, refresh the
// value labels, and read the current min/max as query values.
function rangeSlider(fieldset) {
    const minThumb  = fieldset.querySelector('.range-thumb--min');
    const maxThumb  = fieldset.querySelector('.range-thumb--max');
    const fill      = fieldset.querySelector('[data-fill]');
    const bubbleMin = fieldset.querySelector('[data-bubble-min]');
    const bubbleMax = fieldset.querySelector('[data-bubble-max]');
    const dataMin   = Number(minThumb.min);
    const dataMax   = Number(maxThumb.max);

    function clamp() {
        // Thumbs share a track; stop them from crossing so min never exceeds max.
        if (Number(minThumb.value) > Number(maxThumb.value)) {
            const mid = Math.min(Number(minThumb.value), Number(maxThumb.value));
            minThumb.value = mid;
            maxThumb.value = mid;
        }
    }

    function percent(value) {
        const span = dataMax - dataMin || 1;
        return ((Number(value) - dataMin) / span) * 100;
    }

    function paint() {
        const left = percent(minThumb.value);
        const right = percent(maxThumb.value);
        if (fill) {
            fill.style.left = left + '%';
            fill.style.width = (right - left) + '%';
        }
        // Bubbles ride above each thumb so the current value is visible under
        // the finger while dragging; the head readout mirrors the same values.
        if (bubbleMin) {
            bubbleMin.textContent = minThumb.value;
            bubbleMin.style.left = left + '%';
        }
        if (bubbleMax) {
            bubbleMax.textContent = maxThumb.value;
            bubbleMax.style.left = right + '%';
        }

        // When the min thumb is dragged past the midpoint it overlaps the max
        // thumb; raise the max thumb above it so it never becomes untappable.
        if (left > 50) {
            minThumb.style.zIndex = '3';
            maxThumb.style.zIndex = '4';
        } else {
            minThumb.style.zIndex = '4';
            maxThumb.style.zIndex = '3';
        }
    }

    function refresh() {
        clamp();
        paint();
    }

    function state() {
        return sliderRangeState(Number(minThumb.value), Number(maxThumb.value), dataMin, dataMax);
    }

    function reset() {
        minThumb.value = dataMin;
        maxThumb.value = dataMax;
        paint();
    }

    return { minThumb, maxThumb, refresh, state, reset };
}

// init wires the search box and filter panel inputs to the combined
// /api/filter endpoint. It replaces search.js's own request handling — the
// search box's input event is read here too, so search and filters combine
// into a single request rather than each overwriting the other's results.
function init() {
    const searchInput   = document.getElementById('search-input');
    const loading       = document.getElementById('loading');
    const noResults     = document.getElementById('no-results');
    const results       = document.getElementById('search-results');
    const resetButton   = document.getElementById('filter-reset');
    const locationsRoot = document.getElementById('filter-locations');

    if (!searchInput || !loading || !noResults || !results) {
        return;
    }

    const sliders = {};
    document.querySelectorAll('.filter-range[data-range]').forEach(function (fs) {
        sliders[fs.getAttribute('data-range')] = rangeSlider(fs);
    });

    // Paint the initial fill/labels so the sliders reflect their starting span.
    Object.keys(sliders).forEach(function (k) { sliders[k].refresh(); });

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
        const creation = sliders['creation'] ? sliders['creation'].state() : { min: '', max: '' };
        const album = sliders['first-album'] ? sliders['first-album'].state() : { min: '', max: '' };
        const members = sliders['members'] ? sliders['members'].state() : { min: '', max: '' };
        return {
            q: searchInput.value,
            creationMin: creation.min,
            creationMax: creation.max,
            firstAlbumMin: album.min,
            firstAlbumMax: album.max,
            membersMin: members.min,
            membersMax: members.max,
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

    // Repaint the slider fill immediately on every drag (snappy visual), but
    // debounce the actual network request so dragging doesn't spam /api/filter.
    function onSliderInput(slider) {
        return function () {
            slider.refresh();
            debouncedChange();
        };
    }

    searchInput.addEventListener('input', debouncedChange);
    Object.keys(sliders).forEach(function (k) {
        const s = sliders[k];
        s.minThumb.addEventListener('input', onSliderInput(s));
        s.maxThumb.addEventListener('input', onSliderInput(s));
    });
    if (locationsRoot) {
        locationsRoot.addEventListener('change', debouncedChange);
    }

    if (resetButton) {
        resetButton.addEventListener('click', function () {
            searchInput.value = '';
            Object.keys(sliders).forEach(function (k) { sliders[k].reset(); });
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
