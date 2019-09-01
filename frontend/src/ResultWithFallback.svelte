<script>
  import { onMount } from 'svelte';
  import { MDCTextField } from '@material/textfield';
  import { createEventDispatcher } from 'svelte';
  import CorrectableInput from './CorrectableInput.svelte';
  import Result from './Result.svelte';

	const dispatch = createEventDispatcher();

  export let label;
  export let value;
  export let correctable;

  onMount(() => {
    document.querySelectorAll('.mdc-text-field').forEach(t => new MDCTextField(t));
  });

  function save(value) {
    dispatch('save', {
      value,
    });
  }
</script>

<style lang="scss">
@import '@material/textfield/mdc-text-field';

.mdc-text-field input {
  border: none;
}
</style>

{#if value}
  <Result {value} {label} />
{:else}
  <CorrectableInput on:save={({ detail: {value} }) => save(value)} noEdit={false} value={correctable} {label} />
{/if}