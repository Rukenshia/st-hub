<script>
  import { onMount } from 'svelte';
  import { MDCTextField } from '@material/textfield';

  export let label;
  export let value;

  onMount(() => {
    document.querySelectorAll('.mdc-text-field').forEach(t => new MDCTextField(t));
  });
</script>

<style lang="scss">
@import '@material/textfield/mdc-text-field';

.mdc-text-field input {
  border: none;
}
</style>

<div class="mdc-text-field correctable-text-field mdc-text-field--with-trailing-icon">
  <input type="text" id="field" class="mdc-text-field__input" disabled {value}>

  <div style="position:absolute; left:0; right:0; top:0; bottom:0;"></div> <!-- hack to make clicking work on all platforms-->

  <label class="mdc-floating-label" for="field">{label}</label>
  <div class="mdc-line-ripple"></div>
</div>