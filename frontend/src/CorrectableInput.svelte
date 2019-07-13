<script>
  import { onMount } from 'svelte';
  import { MDCTextField } from '@material/textfield';
  import {MDCTextFieldHelperText} from '@material/textfield/helper-text';

  export let value;
  export let label;
  export let helptext;

  onMount(() => {
    document.querySelectorAll('.mdc-text-field').forEach(t => new MDCTextField(t));
    document.querySelectorAll('.mdc-text-field-helper-text').forEach(t => new MDCTextFieldHelperText(t));
  });
</script>

<style lang="scss">
@import '@material/textfield/mdc-text-field';
@import "@material/textfield/helper-text/mdc-text-field-helper-text";

.correctable-text-field.mdc-text-field input {
  border: none;
}
</style>

<div class="battle-card__content">
  <div class="mdc-text-field correctable-text-field">
    <input type="text" id="field" class="mdc-text-field__input" disabled value={value.Corrected ? value.Corrected : value.Value}>
    <label class="mdc-floating-label" for="field">{label}</label>
    <div class="mdc-line-ripple"></div>
  </div>
  {#if helptext}
  <div class="mdc-text-field-helper-line">
    <div id="field-helper-text" class="mdc-text-field-helper-text mdc-text-field-helper-text--persistent" aria-hidden="false">
      {helptext}
    </div>
  </div>
  {/if}
</div>