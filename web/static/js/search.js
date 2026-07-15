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

// Expose pure functions for search.test.js to import under Node, and for
// filter.js to reuse in the browser, where it drives #search-input directly
// as part of the combined search+filter request in /api/filter.
if (typeof module !== 'undefined' && module.exports) {
    module.exports = { debounce, renderCards };
} else {
    window.debounce = debounce;
    window.renderCards = renderCards;
}
