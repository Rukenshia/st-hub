<script>
  import { onMount } from 'svelte';
  import { MDCTextField } from '@material/textfield';
  import {MDCTextFieldHelperText} from '@material/textfield/helper-text';
  import { MDCDialog } from '@material/dialog';
  import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

  export let value;
  export let label;
  export let helptext = '';
  export let noEdit = true;

  let editing = false;
  let buffer;

  onMount(() => {
    document.querySelectorAll('.mdc-text-field').forEach(t => new MDCTextField(t));
    document.querySelectorAll('.mdc-text-field-helper-text').forEach(t => new MDCTextFieldHelperText(t));
  });

  function startEditing() {
    if (noEdit) { return }

    editing = true;
    buffer = value.Corrected ? value.Corrected : value.Value;
    setTimeout(() => {
      new MDCDialog(document.querySelector('.mdc-dialog')).open();
      document.querySelectorAll('.mdc-text-field').forEach(t => new MDCTextField(t));
    }, 1);
  }

  function save() {
    if (isNaN(parseInt(value, 10))) {
      value = value.Corrected ? value.Corrected : value.Value;
    }
    dispatch('save', {
      value: buffer,
    });
  }
</script>

<style lang="scss">
@import '@material/textfield/mdc-text-field';
@import '@material/textfield/helper-text/mdc-text-field-helper-text';
@import '@material/dialog/mdc-dialog';

.correctable-text-field.mdc-text-field input {
  border: none;
}
</style>

<div class="battle-card__content">
  {#if editing}
  <div class="mdc-dialog mdc-dialog--open input-dialog"
      role="alertdialog"
      aria-modal="true"
      aria-labelledby="my-dialog-title"
      aria-describedby="my-dialog-content">
    <div class="mdc-dialog__container">
      <div class="mdc-dialog__surface">
        <h2 class="mdc-dialog__title" id="my-dialog-title">
          Manual Update
        </h2>
        <div class="mdc-dialog__content" id="my-dialog-content">
          <div class="mdc-text-field correctable-text-field mdc-text-field--with-trailing-icon">
            <input type="text" id="field" class="mdc-text-field__input" bind:value={buffer}>
            <label class="absolute left-2 top-1 text-cool-gray-400 text-sm" for="field">{label}</label>
            <div class="mdc-line-ripple"></div>
          </div>
        </div>
        <footer class="mdc-dialog__actions">
          <button type="button" class="mdc-button mdc-dialog__button" data-mdc-dialog-action="cancel" on:click={() => editing = false}>
            <span class="mdc-button__label">Cancel</span>
          </button>
          <button type="button" class="mdc-button mdc-dialog__button" data-mdc-dialog-action="save" on:click={save}>
            <span class="mdc-button__label">Save</span>
          </button>
        </footer>
      </div>
    </div>
    <div class="mdc-dialog__scrim"></div>
  </div>
  {/if}


  <div class="mdc-text-field correctable-text-field mdc-text-field--with-trailing-icon" on:click={startEditing}>
    <input type="text" id="field" class="mdc-text-field__input" disabled value={value.Corrected ? value.Corrected : value.Value}>

    <div style="position:absolute; left:0; right:0; top:0; bottom:0;"></div> <!-- hack to make clicking work on all platforms-->

    <label class="absolute left-2 top-1 text-cool-gray-400 text-sm" for="field">{label}</label>
    {#if !noEdit}
    <i class="material-icons mdc-text-field__icon">edit</i>
    {/if}
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
