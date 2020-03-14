<script>
  import { derived } from 'svelte/store';
  import TeamSvg from './svg/team.svelte';

  export let battles;

  const rate = derived(battles, newBattles => {
    const inDiv = newBattles.filter(b => b.Statistics.InDivision.Corrected ? true : b.Statistics.InDivision.Value === true).length;

    const rate = Math.round(inDiv / newBattles.length * 100);
    return isNaN(rate) ? 0 : rate;
  });
</script>

<style>
  .bad {
    @apply text-red-600;

    font-weight: 500;
  }

  .rate {
  }
</style>

<div class="flex items-center text-sm text-gray-500">
  <span class="rate" class:bad={$rate >= 50}>{$rate}%</span>
  <div>&nbsp;in division</div>
</div>
