import App from './App.svelte';
import SetupApp from './setup/SetupApp.svelte';

if (window.location.hash == '#setup') {
	var app = new SetupApp({
		target: document.body
	});
} else {
	var app = new App({
		target: document.body
	});
}