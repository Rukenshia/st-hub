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
    font-weight: 500;
    color: #da3c3c;
  }

  .rate {
  }
</style>

<div class="flex items-center">
  <div class="w-4 h-auto mr-2">
    <TeamSvg />
  </div>
  <span class="rate" class:bad={$rate >= 50}>{$rate}%</span>
  <div>&nbsp;in division</div>
</div>
