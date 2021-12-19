import adapter from '@sveltejs/adapter-static';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	kit: {
		adapter: adapter(),
    trailingSlash: "never",
    // If you are not using a .nojekyll file, change your `appDir` to
    // something not starting with an underscore. For example, instead
    // of `_app`, use `app_`, `internal`, etc.
		appDir: 'internal',
		// Hydrate the `<div id="svelte"></div>` element in `src/app.html`.
		target: '#svelte'
	}
};

export default config;
