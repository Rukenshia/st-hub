<script>
  import { onMount } from 'svelte';
  import { MDCTextField } from '@material/textfield';
  export let ship;
  export let battles;

  $: averageDamage = Math.round(battles.reduce((p, b) => p + b.Statistics.Damage.Value, 0) / battles.length);
  $: averageKills = Math.round(battles.reduce((p, b) => p + b.Statistics.Kills.Value, 0) / battles.length * 100) / 100;
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
    <input type="text" id="averageDamage" class="mdc-text-field__input" disabled value={averageDamage ? averageDamage : ''}>
    <label class="mdc-floating-label" for="averageDamage">Average Damage</label>
    <div class="mdc-line-ripple"></div>
  </div>
  <div class="mdc-text-field stat-text-field">
    <input type="text" id="averageKills" class="mdc-text-field__input" disabled value={averageKills ? averageKills : ''}>
    <label class="mdc-floating-label" for="averageKills">Average Kills</label>
    <div class="mdc-line-ripple"></div>
  </div>
  <div class="mdc-text-field stat-text-field">
    <input type="text" id="survivalRate" class="mdc-text-field__input" disabled value={survivalRate ? (survivalRate + '%') : ''}>
    <label class="mdc-floating-label" for="survivalRate">Survival Rate</label>
    <div class="mdc-line-ripple"></div>
  </div>
  <div class="mdc-text-field stat-text-field">
    <input type="text" id="winRate" class="mdc-text-field__input" disabled value={winRate ? (winRate + '%') : ''}>
    <label class="mdc-floating-label" for="winRate">Win Rate</label>
    <div class="mdc-line-ripple"></div>
  </div>
</div>
