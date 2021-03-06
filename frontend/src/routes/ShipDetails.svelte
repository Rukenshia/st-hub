<script>
  import { iteration, battles } from '../stores';
  import { derived } from 'svelte/store';
  import { onMount } from 'svelte';
  import { MDCTextField } from '@material/textfield';

  import ShipBasicAverageStats from '../ShipBasicAverageStats.svelte';
  import DamageBreakdownGraph from '../DamageBreakdownGraph.svelte';
  import ShipName from '../ShipName.svelte';

  import { getShipInfo } from '../wows';

  export let id;
  export let location;

  function getAverage(arr, valFn) {
    const data = arr.reduce((o, v) => {
      if (!v.Results) {
        return o;
      }
      return [o[0] + valFn(v), o[1] + 1];
    }, [0, 0]);

    return Math.round(data[0] / data[1]);
  }

  let ship = derived(iteration, it => it ? it.Ships.find(s => `${s.id}` === id) : { name: '...' });
  let shipBattles = derived(battles, newBattles => newBattles.filter(b => b.ShipID === $ship.id));

  let averageExp = derived(shipBattles, newBattles => getAverage(newBattles, v => v.Results.Economics.BaseExp));
  let averageCredits = derived(shipBattles, newBattles => getAverage(newBattles, v => v.Results.Economics.Credits));
  let averageLifetime = derived(shipBattles, newBattles => getAverage(newBattles, v => v.Results.LifeTime));
  let averageBattleDuration = derived(shipBattles, newBattles => getAverage(newBattles, v => v.Results.Duration));
  let averagePlanesKilled = derived(shipBattles, newBattles => getAverage(newBattles, v => v.Results.PlanesKilled));
  let averageFloodsCaused = derived(shipBattles, newBattles => getAverage(newBattles, v => v.Results.FloodsCaused));

  let shipInfo;

  onMount(async () => {
    document.querySelectorAll('.mdc-text-field').forEach(t => new MDCTextField(t));

    shipInfo = await getShipInfo(id);
  });
</script>

<style lang="scss">
@import '@material/textfield/mdc-text-field';

.stat-text-field.mdc-text-field input {
  border: none;
}

.background {
  position: absolute;
  top: 0;
  z-index: -1;

  background: url(/bg.jpg) no-repeat fixed;
  background-size: cover;
  width: 100%;
  height: 400px;
  opacity: 0.3;
}
</style>

<div class="background"></div>

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
    <div class="pl-2 mb-32 w-full">
      {#if shipInfo}
      <div class="flex justify-between pb-2">
        <div class="w-auto max-w-xs mx-auto">
          <img class="w-auto" alt="Ship image" src={shipInfo.images.large} />
          <div class="flex justify-between">
            <div class="mx-auto">
            <ShipName name={shipInfo.name} tier={shipInfo.tier} type={shipInfo.type} nation={shipInfo.nation} />
            </div>
          </div>
        </div>
      </div>
      {/if}

      <div class="mx-auto w-full md:w-3/4 xl:w-1/2 border rounded-sm shadow-right-lg border-gray-900 bg-gray-900 p-4">
        <DamageBreakdownGraph battles={shipBattles} />

        <div class="mt-4">
          <h1 class="text-2xl">Overall Statistics</h1>

          <div class="p-4">
            <h2 class="text-xl">Basic info</h2>
            <div class="p-2">
              <div class="mdc-text-field stat-text-field">
                <input type="text" id="shipBattles" class="mdc-text-field__input" disabled value={$shipBattles.length}>
                <label class="absolute left-4 top-1 text-cool-gray-400 text-sm" for="shipBattles">Battles played</label>
                <div class="mdc-line-ripple"></div>
              </div>
              <div class="mdc-text-field stat-text-field">
                <input type="text" id="averageExp" class="mdc-text-field__input" disabled value={$averageExp ? $averageExp : 'n/a'}>
                <label class="absolute left-4 top-1 text-cool-gray-400 text-sm" for="averageExp">Average Base EXP</label>
                <div class="mdc-line-ripple"></div>
              </div>
              <div class="mdc-text-field stat-text-field">
                <input type="text" id="averageCredits" class="mdc-text-field__input" disabled value={$averageCredits ? $averageCredits : 'n/a'}>
                <label class="absolute left-4 top-1 text-cool-gray-400 text-sm" for="averageCredits">Average Credits</label>
                <div class="mdc-line-ripple"></div>
              </div>
            </div>

            <h2 class="text-xl">Effectiveness</h2>
            <div class="p-2">
              <ShipBasicAverageStats uid={ship.Name} battles={$shipBattles} />
            </div>

            <h2 class="text-xl">Survivability</h2>
            <div class="p-2">
              <div class="mdc-text-field stat-text-field">
                <input type="text" id="averageLifetime" class="mdc-text-field__input" disabled value={$averageLifetime ? $averageLifetime : '0'}>
                <label class="absolute left-4 top-1 text-cool-gray-400 text-sm" for="averageLifetime">Average Lifetime (seconds)</label>
                <div class="mdc-line-ripple"></div>
              </div>
              <div class="mdc-text-field stat-text-field w-1/2">
                <input type="text" id="averageBattleDuration" class="mdc-text-field__input" disabled value={$averageBattleDuration ? $averageBattleDuration : '0'}>
                <label class="absolute left-4 top-1 text-cool-gray-400 text-sm" for="averageBattleDuration">Average Battle duration (seconds)</label>
                <div class="mdc-line-ripple"></div>
              </div>
            </div>

            <h2 class="text-xl">Anti aircraft</h2>
            <div class="p-2">
              <div class="mdc-text-field stat-text-field">
                <input type="text" id="averagePlanesKilled" class="mdc-text-field__input" disabled value={$averagePlanesKilled ? $averagePlanesKilled : 'n/a'}>
                <label class="absolute left-4 top-1 text-cool-gray-400 text-sm" for="averagePlanesKilled">Average Planes Killed</label>
                <div class="mdc-line-ripple"></div>
              </div>
            </div>

            <h2 class="text-xl">Other data</h2>
            <div class="p-2">
              <div class="mdc-text-field stat-text-field">
                <input type="text" id="averageFloodsCaused" class="mdc-text-field__input" disabled value={$averageFloodsCaused ? $averageFloodsCaused : '0'}>
                <label class="absolute left-4 top-1 text-cool-gray-400 text-sm" for="averageFloodsCaused">Average Floods caused</label>
                <div class="mdc-line-ripple"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  {/if}
</div>
