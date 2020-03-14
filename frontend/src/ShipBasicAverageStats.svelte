<script>
  import { onMount } from 'svelte';
  import { MDCTextField } from '@material/textfield';
  export let battles;
  export let uid;

  const val = v => v.Corrected ? v.Corrected : v.Value;

  $: averageDamage = Math.round(battles.reduce((p, b) => p + (b.Results ? b.Results.Damage.Sum : val(b.Statistics.Damage)), 0) / battles.length);
  $: averageKills = Math.round(battles.reduce((p, b) => p + val(b.Statistics.Kills), 0) / battles.length * 100) / 100;
  $: survivalRate = Math.round(battles.reduce((p, b) => p + b.Statistics.Survived, 0) / battles.length * 100);
  $: winRate = Math.round(battles.reduce((p, b) => p + b.Statistics.Win, 0) / battles.length * 100);

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

<div>
  <div class="mdc-text-field stat-text-field">
    <input type="text" id="{uid}averageDamage" class="mdc-text-field__input" disabled value={averageDamage ? averageDamage : 'n/a'}>
    <label class="mdc-floating-label" for="{uid}averageDamage">Average Damage</label>
    <div class="mdc-line-ripple"></div>
  </div>
  <div class="mdc-text-field stat-text-field">
    <input type="text" id="{uid}averageKills" class="mdc-text-field__input" disabled value={averageKills ? averageKills : 'n/a'}>
    <label class="mdc-floating-label" for="{uid}averageKills">Average Kills</label>
    <div class="mdc-line-ripple"></div>
  </div>
  <div class="mdc-text-field stat-text-field">
    <input type="text" id="{uid}survivalRate" class="mdc-text-field__input" disabled value={survivalRate ? (survivalRate + '%') : 'n/a'}>
    <label class="mdc-floating-label" for="{uid}survivalRate">Survival Rate</label>
    <div class="mdc-line-ripple"></div>
  </div>
  <div class="mdc-text-field stat-text-field">
    <input type="text" id="{uid}winRate" class="mdc-text-field__input" disabled value={winRate ? (winRate + '%') : 'n/a'}>
    <label class="mdc-floating-label" for="{uid}winRate">Win Rate</label>
    <div class="mdc-line-ripple"></div>
  </div>
</div>
