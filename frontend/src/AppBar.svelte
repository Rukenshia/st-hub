<script>
  import { onMount } from 'svelte';
  import { MDCTopAppBar } from '@material/top-app-bar';
  import { MDCSwitch } from '@material/switch';
  import { darkMode } from './stores';

  onMount(() => {
    // Initialise app bar
    const topAppBarElement = document.querySelector('.mdc-top-app-bar');
    const topAppBar = new MDCTopAppBar(topAppBarElement);
    const switchControl = new MDCSwitch(document.querySelector('.mdc-switch'));
    switchControl.getDefaultFoundation().setChecked($darkMode);
  });

  export let iteration = {
    ClientVersion: 'loading',
    IterationName: 'loading',
  };

  export let apiError = false;
</script>

<style lang="scss">
@import '@material/top-app-bar/mdc-top-app-bar';
@import '@material/switch/mdc-switch';

header.mdc-top-app-bar {
  position: relative;
  width: 100%;
  margin-right: 0;

  .mdc-top-app-bar__section {
    span {
      margin-right: 8px;
    }
  }

  .mdc-switch.toggle {
    --mdc-theme-secondary: #cecece;
  }

  label {
    margin-left: 12px;
  }

  &.is-error {
    --mdc-theme-primary: #B00020;
  }
}
</style>

<header class="mdc-top-app-bar header" class:is-error={apiError}>
  <div class="mdc-top-app-bar__row">
    <section class="mdc-top-app-bar__section mdc-top-app-bar__section--align-start">
      <span class="mdc-top-app-bar__title">StHub</span>
      <section class="mdc-top-app-bar__section">
        {#if apiError}
          <span><strong>Unable to connect to the API. Are you running sthub locally?</strong></span>
        {:else}
          <span>Client version: <strong>{iteration.ClientVersion}</strong></span>
          <span>Iteration: <strong>{iteration.IterationName}</strong></span>
        {/if}
      </section>
      <div class="toggle mdc-switch" class:mdc-switch--checked={$darkMode}>
        <div class="mdc-switch__track"></div>
        <div class="mdc-switch__thumb-underlay">
          <div class="mdc-switch__thumb">
              <input type="checkbox" id="mode-switch" on:change={() => $darkMode = !$darkMode} class="mdc-switch__native-control" role="switch">
          </div>
        </div>
      </div>
      <label for="mode-switch" class="mdc-typography--subtitle1"><strong>Dark mode</strong></label>
    </section>
  </div>
</header>
