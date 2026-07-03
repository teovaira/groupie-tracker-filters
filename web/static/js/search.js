'use strict';

/**
 * debounce returns a function that delays invoking fn until after delay ms
 * have elapsed since the last invocation. Repeated calls within the delay
 * window reset the timer, so fn is only called once after the burst stops.
 *
 * @param {Function} fn - the function to debounce
 * @param {number} delay - delay in milliseconds
 * @returns {Function} debounced version of fn that forwards all arguments
 */
function debounce(fn, delay) {
    let timer;
    return function () {
        const args = arguments;
        clearTimeout(timer);
        timer = setTimeout(function () {
            fn.apply(null, args);
        }, delay);
    };
}

/**
 * renderCards builds an HTML string of artist card links from an array of
 * artist objects. Each object must have id, name, image, and creationDate.
 *
 * @param {Array<{id: number, name: string, image: string, creationDate: number}>} artists
 * @returns {string} HTML string of artist cards, or an empty string if the array is empty
 */
function renderCards(artists) {
    if (!artists || artists.length === 0) {
        return '';
    }
    let html = '';
    for (let i = 0; i < artists.length; i++) {
        const a = artists[i];
        html += '<a href="/artist/' + a.id + '" class="artist-card">' +
            '<img src="' + a.image + '" alt="' + a.name + '">' +
            '<div class="artist-card-info">' +
            '<h2>' + a.name + '</h2>' +
            '<p>Since ' + a.creationDate + '</p>' +
            '</div>' +
            '</a>';
    }
    return html;
}

// Expose pure functions globally so search.test.js can import them.
window.debounce = debounce;
window.renderCards = renderCards;

// init wires the search input to the live search API.
function init() {
    const input      = document.getElementById('search-input');
    const loading    = document.getElementById('loading');
    const noResults  = document.getElementById('no-results');
    const results    = document.getElementById('search-results');

    if (!input || !loading || !noResults || !results) {
        return;
    }

    // Save the original server-rendered cards to restore on empty query.
    const originalHTML = results.innerHTML;

    function showLoading()    { loading.classList.remove('hidden'); }
    function hideLoading()    { loading.classList.add('hidden'); }
    function showNoResults()  { noResults.classList.remove('hidden'); }
    function hideNoResults()  { noResults.classList.add('hidden'); }

    function handleInput() {
        const q = input.value.trim();

        // Empty query — restore original cards, clear any search state.
        if (q === '') {
            hideLoading();
            hideNoResults();
            results.innerHTML = originalHTML;
            return;
        }

        showLoading();
        hideNoResults();

        fetch('/api/search?q=' + encodeURIComponent(q))
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

    input.addEventListener('input', debounce(handleInput, 300));
}

document.addEventListener('DOMContentLoaded', init);
