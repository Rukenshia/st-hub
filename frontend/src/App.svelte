<script>
	import { battles } from './stores';
	import Navigation from './Navigation.svelte';
	import DivisionStatistics from './DivisionStatistics.svelte';
	import Battles from './Battles.svelte';
	import Feed from './Feed.svelte';

	$battles = [
		{id: 4711, statistics: { ship: 'Colbert', win: true, division: true, damage: Math.round(Math.random() * 100000), kills: 2, hits: 250, survived: true }},
		{id: 4712, statistics: { ship: 'Colbert', win: true, division: false, damage: Math.round(Math.random() * 100000), kills: 0, hits: 250, survived: false }},
		{id: 4715, statistics: { ship: 'Le Fantasque', win: true, division: true, damage: Math.round(Math.random() * 100000), kills: 2, hits: 250, survived: true }},
		{id: 4713, statistics: { ship: 'Colbert', win: false, division: false, damage: Math.round(Math.random() * 100000), kills: 2, hits: 250, survived: false }},
		{id: 4714, statistics: { ship: 'Colbert', win: false, division: false, damage: Math.round(Math.random() * 100000), kills: 1, hits: 250, survived: true }},
		{id: 4715, statistics: { ship: 'Colbert', win: true, division: false, damage: Math.round(Math.random() * 100000), kills: 0, hits: 250, survived: false }},
		{id: 4714, statistics: { ship: 'Le Fantasque', win: false, division: false, damage: Math.round(Math.random() * 100000), kills: 2, hits: 250, survived: false }},
		{id: 4716, statistics: { ship: 'Le Fantasque', win: false, division: true, damage: Math.round(Math.random() * 100000), kills: 2, hits: 250, survived: true }},
		{id: 4716, statistics: { ship: 'Colbert', win: false, division: true, damage: Math.round(Math.random() * 100000), kills: 2, hits: 250, survived: true }},
		{id: 4717, statistics: { ship: 'Colbert', win: true, division: true, damage: Math.round(Math.random() * 100000), kills: 2, hits: 250, survived: false }},
		{id: 4718, statistics: { ship: 'Le Fantasque', win: true, division: true, damage: Math.round(Math.random() * 100000), kills: 2, hits: 250, survived: true }},
		{id: 4719, statistics: { ship: 'Slava', win: true, division: true, damage: Math.round(Math.random() * 100000), kills: 2, hits: 250, survived: true }},
		{id: 4720, statistics: { ship: 'Slava', win: true, division: false, damage: Math.round(Math.random() * 100000), kills: 2, hits: 250, survived: false }},
		{id: 4721, statistics: { ship: 'Slava', win: true, division: true, damage: Math.round(Math.random() * 100000), kills: 2, hits: 250, survived: false }},
		{id: 4720, statistics: { ship: 'Slava', win: true, division: false, damage: Math.round(Math.random() * 100000), kills: 2, hits: 250, survived: true }},
		{id: 4721, statistics: { ship: 'Slava', win: true, division: true, damage: Math.round(Math.random() * 100000), kills: 2, hits: 250, survived: false }},
		{id: 4722, statistics: { ship: 'Slava', win: true, division: true, damage: Math.round(Math.random() * 100000), kills: 2, hits: 250, survived: true }},
		{id: 4723, statistics: { ship: 'Slava', win: true, division: true, damage: Math.round(Math.random() * 100000), kills: 2, hits: 250, survived: true }},
	];

	function addBattle() {
		$battles = [...$battles, {id: 4723, statistics: { ship: 'Slava', win: true, division: true, damage: Math.round(Math.random() * 100000), kills: 2, hits: 250 }}];
	}
</script>

<style global lang="scss">
@import '@material/typography/mdc-typography';

body {
	@include mdc-typography-base();
}
</style>

<DivisionStatistics/>

<button on:click={addBattle}>Add Slava Battle</button>

<Navigation links={["battles", "feed"]}>
	<span slot="battles">
		<Battles />
	</span>

	<span slot="feed">
		<Feed/>
	</span>
</Navigation>

<footer>
	<div class="mdc-typography--subtitle2" style="text-align: center;">
		<div>Icons made by <a href="https://www.freepik.com/" title="Freepik">Freepik</a> from <a href="https://www.flaticon.com/"                 title="Flaticon">www.flaticon.com</a> is licensed by <a href="http://creativecommons.org/licenses/by/3.0/" title="Creative Commons BY 3.0" target="_blank">CC 3.0 BY</a></div>
	</div>
</footer>