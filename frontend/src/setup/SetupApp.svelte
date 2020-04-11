<script>
  import {onMount} from 'svelte';

  import Banner from './Banner.svelte';
  import GameDirectorySelection from './GameDirectorySelection.svelte';
  import Step from './Step.svelte';

  let step = 1;

  onMount(() => {
    document.addEventListener('astilectron-ready', function() {
      astilectron.sendMessage("connect", function() {
        setTimeout(() => step = 2, 500);
      });
    });
  });

  function selectGameDirDone() {
    step = 3;

    installGameMod();
  }

  function installGameMod() {
    astilectron.sendMessage("install", function() {
      setTimeout(() => step = 4, 500);
    });
  }

  function redirect() {
    window.location.href = "https://sthub.in.fkn.space";
  }
</script>

<style global lang="scss">
@tailwind base;
@tailwind utilities;
@tailwind components;
</style>

<Banner />

<div class="container p-8">
  <div class="flex items-center mb-4">
    <div class="flex-initial md:w-1/6 sm:w-1/12"></div>
    <div class="flex-initial w-full">
      <Step done={step > 1} active={step === 1} number=1 title="Connecting to the API"></Step>
      <Step done={step > 2} active={step === 2} number=2 title="Select game directory"></Step>
      {#if step === 2}
        <GameDirectorySelection on:done={selectGameDirDone} />
      {/if}
      <Step done={step > 3} active={step === 3} number=3 title="Install game mod"></Step>
      {#if step === 4}
        <div class="mt-4 pl-12">
          <h1 class="text-xl text-gray-800">You are all set 🎉</h1>

          <button on:click={redirect} class="mt-4 rounded-sm px-3 py-2 mx-auto border-2 uppercase text-medium border-teal-300 text-teal-300 hover:border-teal-600 hover:text-teal-600">
            Let's go
          </button>
        </div>
      {/if}
    </div>
  </div>
</div>
