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
  .good {
    color: #5ab15a;
  }

  .bad {
    color: #da3c3c;
  }

  .rate {
    font-weight: 500;
  }

  td {
    padding-top: 8px;
    padding-right: 16px;
  }

  td + td {
    text-align: right;
  }
</style>

<div class="flex items-center">
  <div class="w-4 h-auto mr-2">
    <TeamSvg />
  </div>
  <span class="rate" class:good={$rate < 50} class:bad={$rate >= 50}>{$rate}%</span>
  <div class="pl-2">in division</div>
</div>
