<script>
  export let availableVersion;
  export let version;
</script>

<style lang="scss">
@import '@material/top-app-bar/mdc-top-app-bar';

header.mdc-top-app-bar {
  @include mdc-top-app-bar-fill-color(rgb(244, 247, 118));
  @include mdc-top-app-bar-ink-color(#121212);
  position: relative;

  a {
    color: rgb(232, 126, 4);
    text-decoration: none;
  }
}
</style>

<header class="mdc-top-app-bar">
  <div class="mdc-top-app-bar__row">
    <section class="mdc-top-app-bar__section mdc-top-app-bar__section--align-start">
      <section class="mdc-top-app-bar__section">
        <span class="text-md">You are running version {version}, but {availableVersion} is available. Please update it <a target="_blank" href="https://github.com/Rukenshia/st-hub/releases/tag/{availableVersion}">here</a>.</span>
      </section>
    </section>
  </div>
</header>
