<script>
  import { onMount } from 'svelte';
  import { MDCTopAppBar } from '@material/top-app-bar';

  onMount(() => {
    // Initialise app bar
    const topAppBarElement = document.querySelector('.mdc-top-app-bar');
    const topAppBar = new MDCTopAppBar(topAppBarElement);
  });

  export let iteration = {
    ClientVersion: 'loading',
    IterationName: 'loading',
  };
</script>

<style lang="scss">
@import '@material/top-app-bar/mdc-top-app-bar';

header.mdc-top-app-bar {
  position: relative;
  width: 100%;
  margin-right: 0;

  .mdc-top-app-bar__section {
    span {
      margin-right: 8px;
    }
  }
}
</style>

<header class="mdc-top-app-bar">
  <div class="mdc-top-app-bar__row">
    <section class="mdc-top-app-bar__section mdc-top-app-bar__section--align-start">
      <span class="mdc-top-app-bar__title">StHub</span>
      <section class="mdc-top-app-bar__section">
        <span>Client version: <strong>{iteration.ClientVersion}</strong></span>
        <span>Iteration: <strong>{iteration.IterationName}</strong></span>
      </section>
    </section>
  </div>
</header>