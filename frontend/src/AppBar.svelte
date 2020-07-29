<script>
  import { onMount } from 'svelte';
  import { darkMode } from './stores';

  onMount(() => {

    // Load cookie acceptance
    cookiesAccepted = window.localStorage.getItem('STHUB_COOKIES_ACCEPTED') !== null;
  });

  export let iteration = {
    ClientVersion: 'loading',
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

</style>

<div class="bg-gray-900">
  <nav class="bg-gray-900">
    <div class="max-w-7xl mx-auto sm:px-2 lg:px-4">
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center">
              <div class="ml-10 flex items-baseline">
                <a href="#" class="ml-4 px-2 py-2 rounded-md text-sm font-medium text-gray-300"><span class="text-gray-400">StHub</span> {version}</a>
                <a href="#" class="ml-4 px-2 py-2 rounded-md text-sm font-medium text-gray-300"><span class="text-gray-400">WoWS</span> {iteration.ClientVersion}</a>
              </div>
          </div>
  </nav>
</div>

{#if apiError}
<div class="mt-8 flex justify-around items-center mx-auto">
  <div class="rounded-md bg-red-900 p-4">
    <div class="flex">
      <div class="flex-shrink-0">
        <svg class="h-5 w-5 text-red-50" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
        </svg>
      </div>
      <div class="ml-3">
        <h3 class="text-sm leading-5 font-medium text-red-200">
          I'm having trouble to reach the app
        </h3>
        <div class="mt-2 text-sm leading-5 text-red-300">
          <ul class="list-disc pl-5">
            <li>
              It looks like StHub.exe is not running properly
            </li>
            <li class="mt-1">
              Browsing from the web? You'll need to download the App from <strong><a href="https://github.com/Rukenshia/st-hub/releases">here</a></strong>.
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</div>
{/if}

{#if !cookiesAccepted}
<div class="fixed z-50 bottom-0 p-4 bg-gray-500 text-gray-900">
  This website uses cookies to gain anonymised usage statistics (namely, using Google Analytics to find out if people actually use StHub).

  <button on:click={accept} class="rounded-sm bg-gray-800 text-gray-200 ml-2 px-2 py-1 hover:bg-gray-900">Alright</button>
</div>
{/if}
