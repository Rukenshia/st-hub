<script>
  import { iteration, battles } from '../stores';
  import { derived } from 'svelte/store';
  import { onMount } from 'svelte';
  import { MDCTextField } from '@material/textfield';

  import ShipStatistics from '../ShipStatistics.svelte';
  import DamageBreakdownGraph from '../DamageBreakdownGraph.svelte';

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

  let ship = derived(iteration, it => it ? it.Ships.find(s => `${s.ID}` === id) : { Name: '...' });
  let shipBattles = derived(battles, newBattles => newBattles.filter(b => b.ShipID === $ship.ID));

  let averageExp = derived(battles, newBattles => getAverage(newBattles, v => v.Results.Economics.BaseExp));
  let averageCredits = derived(battles, newBattles => getAverage(newBattles, v => v.Results.Economics.Credits));

  onMount(() => {
    document.querySelectorAll('.mdc-text-field').forEach(t => new MDCTextField(t));
  });
</script>

<style lang="scss">
@import '@material/textfield/mdc-text-field';

.stat-text-field.mdc-text-field input {
  border: none;
}
</style>

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

      <div class="w-full md:w-1/2">
        <DamageBreakdownGraph battles={shipBattles} />
      </div>
      <div class="mt-4">
        <div class="mdc-text-field stat-text-field">
          <input type="text" id="shipBattles" class="mdc-text-field__input" disabled value={$shipBattles.length}>
          <label class="mdc-floating-label" for="shipBattles">Battles played</label>
          <div class="mdc-line-ripple"></div>
        </div>
        <div class="mdc-text-field stat-text-field">
          <input type="text" id="averageExp" class="mdc-text-field__input" disabled value={$averageExp ? $averageExp : 'n/a'}>
          <label class="mdc-floating-label" for="averageExp">Average Base EXP</label>
          <div class="mdc-line-ripple"></div>
        </div>
        <div class="mdc-text-field stat-text-field">
          <input type="text" id="averageCredits" class="mdc-text-field__input" disabled value={$averageCredits ? $averageCredits : 'n/a'}>
          <label class="mdc-floating-label" for="averageCredits">Average Credits</label>
          <div class="mdc-line-ripple"></div>
        </div>
      </div>

      <ShipStatistics battles={$shipBattles} />

      <div class="flex">
        <div class="w-3/4 p-4 mt-4 ml-4 rounded-sm bg-gray-900 text-gray-100">Watch this space! Additional information will follow soon™️</div>
      </div>
    </div>
  {/if}
</div>
