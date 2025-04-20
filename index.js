// Configurable constants for better maintainability and clarity

const ALGOLIA = {
	APP_ID: "A4F5Q6U2GA",
	API_KEY: "a36a833f174d244d25e077857e4f5969",
	INDEX_NAME: "restaurants",
};
const DEFAULT_USER_LOCATION = { lat: 40.7127281, lng: -74.0060152 }; // New York City
const MAX_VALUES_PER_FACET = 7; // limit the values to not overload the UI

// technically it's an oblate spheroid (so this isn't necessarily exact)
const EARTH_RADIUS_MILES = 3958.8;

// Configure facets

const renderFacet = (content, element, facet) => {
	element.innerHTML = content
		.getFacetValues(facet, { sortBy: ["count:desc"] })
		.map(generateFilterItemHTML)
		.join("");
};

const FACETS = [
	{
		element: document.getElementById("cuisine-filter"),
		facet: "food_type",
		handler: (helper, facet, facetValue) => {
			helper.toggleFacetRefinement(facet, facetValue).search();
		},
		render: renderFacet,
	},
	{
		element: document.getElementById("rating-filter"),
		facet: "stars_count",
		handler: (helper, facet, facetValue) => {
			if (helper.hasRefinements(facet)) {
				helper.clearRefinements(facet).search();
			} else {
				helper.addNumericRefinement(facet, ">=", facetValue).search();
			}
		},
		render: (content, element, facet) => {
			// the mock up displays from 0 to 5 but that seems a questionable UI choice
			// as the user probably wants to filter by the highest rated restaurants
			let vals = [5, 4, 3, 2, 1];

			// if it was already defined, just render the refinement
			let values = helper.getRefinements(facet);
			if (values.length > 0) {
				vals = values[0].value;
			}

			element.innerHTML = generateRatingFilterHTML(vals);
		},
	},
	{
		element: document.getElementById("payment-filter"),
		facet: "payment_options",
		handler: (helper, facet, facetValue) => {
			helper.toggleFacetRefinement(facet, facetValue).search();
		},
		render: renderFacet,
	},
];

// Initialize Algolia search client and helper
const client = algoliasearch(ALGOLIA.APP_ID, ALGOLIA.API_KEY);
const helper = algoliasearchHelper(client, ALGOLIA.INDEX_NAME, {
	facets: FACETS.map((facetConfig) => facetConfig.facet),
	maxValuesPerFacet: MAX_VALUES_PER_FACET,
});

let userGeoloc = { ...DEFAULT_USER_LOCATION };
let geolocationSuccess = false;

// Utility Functions

// Converts degrees to radians
const toRadians = (degrees) => (degrees * Math.PI) / 180;

// Calculates the distance using the Haversine formula
const haversineDistance = (pos1, pos2) => {
	const dLat = toRadians(pos2.lat - pos1.lat);
	const dLon = toRadians(pos2.lng - pos1.lng);
	const a =
		Math.sin(dLat / 2) ** 2 +
		Math.cos(toRadians(pos1.lat)) *
			Math.cos(toRadians(pos2.lat)) *
			Math.sin(dLon / 2) ** 2;
	return EARTH_RADIUS_MILES * 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
};

// Formats distance in miles for display
const formatDistance = (miles) =>
	miles < 10 ? `${miles.toFixed(1)} mi` : `${Math.round(miles)} mi`;

// Converts a rating to a star-based representation
const ratingToStars = (rating) => {
	// round to nearest 0.5 increment
	const roundedRating = Math.round(rating * 2) / 2;
	const fullStars = Math.floor(roundedRating);
	const emptyStars = Math.floor(5 - roundedRating);
	const halfStars = 5 - fullStars - emptyStars;
	return `
		<span class="stars" role="img" aria-label="${rating} stars">
			${'<span class="star"></span>'.repeat(fullStars)}
			${'<span class="star half"></span>'.repeat(halfStars)}
			${'<span class="star empty"></span>'.repeat(emptyStars)}
		</span>
	`;
};

// Geolocation Logic

// create a search with the user's location
const searchWithLocation = () => {
	helper
		.setQueryParameter("aroundLatLng", `${userGeoloc.lat},${userGeoloc.lng}`)
		.search();
};

const getUserLocation = () => {
	if (navigator.geolocation) {
		navigator.geolocation.getCurrentPosition(
			(geoPosition) => {
				userGeoloc = {
					lat: geoPosition.coords.latitude,
					lng: geoPosition.coords.longitude,
				};

				// mark success to avoid override
				geolocationSuccess = true;
				searchWithLocation();
			},
			(error) => console.warn("Geolocation error:", error),
			{ timeout: 10000 }
		);
	}

	// fallback to IP-based geolocation
	fetch("https://ipapi.co/json/")
		.then((response) => response.json())
		.then((data) => {
			if (!geolocationSuccess && data.latitude && data.longitude) {
				userGeoloc = { lat: data.latitude, lng: data.longitude };
				searchWithLocation();
			}
		});
};

// Template Functions

const generateHitHTML = (hit) => {
	const distance = haversineDistance(userGeoloc, hit._geoloc);
	return `
		<div class="results__item ${
			hit.__position === 1 ? "results__item--first" : ""
		}">
			<div class="result">
				<div class="result__image-container">
					<img class="result__image" src="${hit.image_url}" alt="${hit.name}" />
				</div>
				<div class="result__text-container">
					<h1 class="result__title">${hit._highlightResult.name.value}</h1>
					<p class="result__rating">${hit.stars_count} ${ratingToStars(
		hit.stars_count
	)} (${hit.reviews_count} reviews)</p>
					<p class="result__summary">
						${hit._highlightResult.food_type.value} | ${
		hit._highlightResult.neighborhood.value
	} | ${hit._highlightResult.price_range.value} | ${formatDistance(distance)}
					</p>
				</div>
			</div>
		</div>
	`;
};

const generateNoResultsHTML = (query) => `
	<div id="no-results-message">
		<p>We didn't find any results for the search <em>"${query}"</em>.</p>
		<a href="." class="clear-all">Clear search</a>
	</div>
`;

const generateResultsStatsHTML = (nbHits, serverTimeMS) => `
	<div class="results__stats-bar">
		<span class="results__count-text">
			${nbHits} results found <span class="results__time-text">in ${
	serverTimeMS / 1000
} seconds</span>
		</span>
	</div>
`;

const generateFilterItemHTML = (value) => `
	<div class="filter__label ${
		value.isRefined ? "filter__label--active" : ""
	}" data-facet="${value.name}">
		<span class="filter__label-text">${value.name}</span>
		<span class="filter__label-number">${value.count}</span>
	</div>
`;

const generateRatingFilterHTML = (starValues) =>
	starValues
		.map(
			(numStars) => `
			<div class="filter__label" aria-label="${numStars} & up" data-facet="${numStars}">
				${ratingToStars(numStars)}
			</div>
		`
		)
		.join("");

// Rendering and Event Binding

const resultsContainer = document.getElementById("results");
let showMoreButton;

const renderHits = (res, hasMore) => {
	if (res.hits.length === 0) {
		resultsContainer.innerHTML = generateNoResultsHTML(res.query);
		return;
	}

	resultsContainer.innerHTML = `${generateResultsStatsHTML(
		res.nbHits,
		res.serverTimeMS
	)}<div class="results__items">${res.hits
		.map(generateHitHTML)
		.join("")}</div>`;

	if (!hasMore) {
		return;
	}

	resultsContainer.innerHTML += `<div class="button"><a id="show-more" href="#" class="button__link">Show More</a></div>`;

	showMoreButton = document.getElementById("show-more");
	showMoreButton.addEventListener("click", () => {
		helper.nextPage().search();
	});
};

const appendHits = (res, hasMore) => {
	showMoreButton.parentElement.insertAdjacentHTML(
		"beforebegin",
		res.hits.map(generateHitHTML).join("")
	);

	if (!hasMore) {
		showMoreButton.parentElement.remove();
		showMoreButton = null;
	}
};

const renderFacetLists = (content) => {
	FACETS.forEach(({ element, facet, render }) => {
		render(content, element, facet);
	});
};

const bindFacetEvents = () => {
	FACETS.forEach(({ element, handler, facet }) => {
		element.addEventListener("click", (event) => {
			const target = event.target.closest(".filter__label");
			if (!target) return;
			handler(helper, facet, target.dataset.facet);
		});
	});
};

// Initialize the application

const init = () => {
	getUserLocation();

	document.getElementById("search-box").addEventListener("keyup", (event) => {
		helper.setQuery(event.target.value).search();
	});

	bindFacetEvents();

	helper.on("result", (e) => {
		const hasMore = e.results.page < e.results.nbPages - 1;
		if (e.results.page > 0) {
			appendHits(e.results, hasMore);
		} else {
            renderFacetLists(e.results);
            renderHits(e.results, hasMore);
        }
	});

	searchWithLocation();
};

init();
