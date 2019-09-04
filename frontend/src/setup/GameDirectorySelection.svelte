<script>
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  let selectedDir = undefined;

  function selectGameDir() {
    astilectron.sendMessage("selectGameDir", function(dir) {
      selectedDir = dir;
      
      setTimeout(() => {
        dispatch('done', { value: 2 });
      }, 500);
    });
  }
</script>

<div class="pl-12 pb-4">
  StHub works by installing a game mod to your World of Warships installation.
  To continue, you will need to select the path your game is installed to.

  <div class="mt-4">
    <button
      disabled={selectedDir !== undefined}
      on:click={selectGameDir}
      class="border-2 rounded-sm px-3 py-2 uppercase font-medium"
      class:text-teal-400={!selectedDir}
      class:border-teal-400={!selectedDir}
      class:hover:border-teal-600={!selectedDir}
      class:hover:text-teal-600={!selectedDir}
      class:text-gray-300={selectedDir !== undefined}
      class:border-gray-300={selectedDir !== undefined}
    >
      Select game path
    </button>
  </div>
</div>