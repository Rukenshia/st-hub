<script>
  import { iteration, battles } from '../stores';
  import { derived } from 'svelte/store';

  import ShipStatistics from '../ShipStatistics.svelte';
  import DamageBreakdownGraph from '../DamageBreakdownGraph.svelte';

  export let id;
  export let location;

  let ship = derived(iteration, it => it ? it.Ships.find(s => `${s.ID}` === id) : { Name: '...' });
  let shipBattles = derived(battles, newBattles => newBattles.filter(b => b.ShipID === $ship.ID));
</script>

<div class="p-4">
  {#if !$ship}
    Loading...
  {:else}
    <a href="/" class="text-teal-400">
      <div class="flex items-center">
        <div class="flex -mr-1">
          <i class="material-icons">arrow_left</i>
        </div>
        <div class="flex">
          Back
        </div>
      </div>
    </a>
    <div class="pl-2 mt-4 mb-32">
      <div class="text-2xl">Detailed Ship Statistics for {$ship.Name}</div>

      <div class="flex">
        <DamageBreakdownGraph battles={shipBattles} />
      </div>

      <div class="flex mt-4 ml-4 mb-4">
        <div class="">
          <strong>Battles:</strong> {$shipBattles.length}<br />
          <strong>Battles counted for graphs:</strong> {$shipBattles.filter(b => b.Results !== undefined).length}<br />
          <strong>Battles (in division):</strong> {$shipBattles.filter(b => b.Statistics.InDivision.Value === true).length}<br />
          <strong>Battles (in division):</strong> {$shipBattles.filter(b => b.Statistics.InDivision.Value === true).length}<br />
        </div>
      </div>

      <ShipStatistics battles={$shipBattles} />

      <div class="flex">
        <div class="w-3/4 p-4 mt-4 ml-4 rounded-sm bg-gray-900 text-gray-100">Watch this space! Additional information will follow soon™️</div>
      </div>
    </div>
  {/if}
</div>
