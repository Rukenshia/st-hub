<script>
  import { onMount } from 'svelte';
  import { MDCTopAppBar } from '@material/top-app-bar';
  // import { MDCSwitch } from '@material/switch';
  import { darkMode } from './stores';

  onMount(() => {
    // Initialise app bar
    const topAppBarElement = document.querySelector('.mdc-top-app-bar');
    const topAppBar = new MDCTopAppBar(topAppBarElement);
    // const switchControl = new MDCSwitch(document.querySelector('.mdc-switch'));
    // switchControl.getDefaultFoundation().setChecked($darkMode);

    // Load cookie acceptance
    cookiesAccepted = window.localStorage.getItem('STHUB_COOKIES_ACCEPTED') !== null;
  });

  export let iteration = {
    ClientVersion: 'loading',
    IterationName: 'loading',
    Ships: [],
  };

  export let version = null;

  export let apiError = false;

  let cookiesAccepted = true;

  function accept() {
    cookiesAccepted = true;
    window.localStorage.setItem('STHUB_COOKIES_ACCEPTED', true);
  }
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

  @apply bg-teal-600;
}

code {
  font-size: 1.375rem;
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
          <span><code>sthub</code> version: <strong>{version}</strong></span>
        {/if}
      </section>
      <!-- <div class="toggle mdc-switch" class:mdc-switch--checked={$darkMode}>
        <div class="mdc-switch__track"></div>
        <div class="mdc-switch__thumb-underlay">
          <div class="mdc-switch__thumb">
              <input type="checkbox" id="mode-switch" on:change={() => $darkMode = !$darkMode} class="mdc-switch__native-control" role="switch">
          </div>
        </div>
      </div>
      <label for="mode-switch" class="mdc-typography--subtitle1"><strong>Dark mode</strong></label> -->
    </section>
  </div>
</header>

{#if !cookiesAccepted}
<div class="fixed z-50 bottom-0 p-4 bg-gray-500 text-gray-900">
  This website uses cookies to gain anonymised usage statistics (namely, using Google Analytics to find out if people actually use StHub).

  <button on:click={accept} class="rounded-sm bg-gray-800 text-gray-200 ml-2 px-2 py-1 hover:bg-gray-900">Alright</button>
</div>
{/if}
