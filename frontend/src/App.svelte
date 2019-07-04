<script>
	import { battles, iteration } from './stores';
	import DivisionStatistics from './DivisionStatistics.svelte';
	import ShipsList from './ShipsList.svelte';
	import AppBar from './AppBar.svelte';
	import Battles from './Battles.svelte';
	import axios from 'axios';
	import { onMount } from 'svelte';

	const fetchIntegration = () => {
		return axios.get('http://100.115.92.205:1323/iterations/current');
	};

	const fetchBattles = () => {
		return axios.get(`http://100.115.92.205:1323/iterations/${$iteration.ClientVersion}/${$iteration.IterationName}/battles`);
	};

	onMount(async () => {
		const res = await fetchIntegration();
		$iteration = res.data;
		const resBattles = await fetchBattles();
		$battles = resBattles.data;

		setInterval(async () => {
			const resBattles = await fetchBattles();
			$battles = resBattles.data;
		}, 2500);
	});
</script>

<style global lang="scss">
@import '@material/typography/mdc-typography';

body {
	@include mdc-typography-base();
	width: 100%;
}

header {
	margin: -8px;
}
</style>

<AppBar iteration={$iteration} />

<DivisionStatistics/>

<ShipsList ships={$iteration.Ships} />

<Battles />

<footer>
	<div class="mdc-typography--subtitle2" style="text-align: center;">
		<div>Icons made by <a href="https://www.freepik.com/" title="Freepik">Freepik</a> from <a href="https://www.flaticon.com/"                 title="Flaticon">www.flaticon.com</a> is licensed by <a href="http://creativecommons.org/licenses/by/3.0/" title="Creative Commons BY 3.0" target="_blank">CC 3.0 BY</a></div>
	</div>
</footer>
